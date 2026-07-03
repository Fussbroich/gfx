package gfx

// fenster.go — Fensterverwaltung, Ebitengine-Game-Loop, Double-Buffering.

import (
	"image"
	"image/color"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// gfxGame implementiert ebiten.Game und bildet die Brücke
// zwischen der gfx-API und dem Ebitengine-Renderloop.
type gfxGame struct{}

var (
	// Fensterzustand (atomar, da von mehreren Goroutinen gelesen)
	offen atomic.Bool

	fensterBreite uint16
	fensterHoehe  uint16

	// Signalkanäle für Lebenszyklussteuerung
	fensterBereit chan struct{} // wird geschlossen, wenn Fenster offen ist
	fensterDone   chan struct{} // wird geschlossen, wenn Fenster zu ist

	// Double-Buffering: frontBuf wird angezeigt, backBuf beschrieben.
	frontBuf   *ebiten.Image
	backBuf    *ebiten.Image
	drawTarget *ebiten.Image // aktuelles Zeichenziel
	bufMu      sync.Mutex    // schützt den Pointer-Swap

	// Stift-Zustand (nur aus der Zeichen-Goroutine beschrieben)
	stiftR, stiftG, stiftB uint8
	stiftAlpha             uint8 = 255

	// Für Sperren/Entsperren (benutzerseitig)
	zeichenSperre sync.Mutex

	// Ebitengine erlaubt pro Prozess nur ein RunGame. Deshalb bleibt der
	// Loop nach dem ersten Fenster() aktiv. FensterAus() schließt nur
	// logisch; ein erneutes Fenster() passt lediglich die Größe an.
	gameGestartet atomic.Bool

	// Zwischenspeicher für Archivieren/Restaurieren.
	archivBuf *ebiten.Image
)

// ===================== ebiten.Game =====================

func (g *gfxGame) Update() error {
	// Der Loop endet erst, wenn der Benutzer das Fenster wirklich schließt
	// (Klick auf das X). FensterAus() setzt nur -offen- auf false und beendet
	// den Loop NICHT, damit das Fenster später wieder geöffnet bzw. in der
	// Größe geändert werden kann (Ebitengine erlaubt kein zweites RunGame).
	if ebiten.IsWindowBeingClosed() {
		offen.Store(false)
		return ebiten.Termination
	}
	updateEingabe()
	return nil
}

func (g *gfxGame) Draw(screen *ebiten.Image) {
	bufMu.Lock()
	front := frontBuf
	bufMu.Unlock()
	if front != nil {
		screen.DrawImage(front, nil)
	}
}

func (g *gfxGame) Layout(_, _ int) (int, int) {
	return int(fensterBreite), int(fensterHoehe)
}

// ===================== Interne Funktionen =====================

func fensterStarten(breite, hoehe uint16) {
	fensterBreite = breite
	fensterHoehe = hoehe

	// Zeichenpuffer in der (evtl. neuen) Größe anlegen.
	frontBuf = ebiten.NewImage(int(breite), int(hoehe))
	backBuf = ebiten.NewImage(int(breite), int(hoehe))
	drawTarget = frontBuf
	archivBuf = nil

	if gameGestartet.Load() {
		// Der RunGame-Loop läuft bereits: nur die Fenstergröße anpassen.
		// (Layout() liefert dynamisch fensterBreite/fensterHoehe.)
		offen.Store(true)
		ebiten.SetWindowSize(int(breite), int(hoehe))
		return
	}

	gameGestartet.Store(true)
	fensterBereit = make(chan struct{})
	fensterDone = make(chan struct{})
	initEingabe()

	go func() {
		ebiten.SetWindowSize(int(breite), int(hoehe))
		ebiten.SetWindowTitle("gfx")
		ebiten.SetWindowClosingHandled(true)
		offen.Store(true)
		close(fensterBereit)

		_ = ebiten.RunGame(&gfxGame{})

		// Aufräumen nach Ende des Game-Loops (Fenster wurde geschlossen)
		offen.Store(false)
		close(fensterDone)
	}()

	<-fensterBereit
	// Kurz warten, damit der erste Frame gerendert wird
	time.Sleep(50 * time.Millisecond)
}

func istFensterOffen() bool {
	return offen.Load()
}

func fensterSchliessen() {
	// Nur logisch schließen: -offen- auf false setzen. Der RunGame-Loop läuft
	// weiter (er darf pro Prozess nur einmal gestartet werden). Das Fenster
	// verschwindet, wenn der Prozess endet (z. B. nach os.Exit) oder wird mit
	// einem erneuten Fenster()-Aufruf in der Größe angepasst.
	offen.Store(false)
}

func setzeFenstertitel(s string) {
	ebiten.SetWindowTitle(s)
}

// gibStiftfarbe liefert die aktuelle Zeichenfarbe inkl. Transparenz.
func gibStiftfarbe() color.NRGBA {
	return color.NRGBA{R: stiftR, G: stiftG, B: stiftB, A: stiftAlpha}
}

func setzeStiftfarbe(r, g, b uint8) {
	stiftR, stiftG, stiftB = r, g, b
}

func setzeTransparenz(t uint8) {
	stiftAlpha = 255 - t
}

func clsBuf() {
	drawTarget.Fill(gibStiftfarbe())
}

func updateAus() {
	drawTarget = backBuf
}

func updateAn() {
	bufMu.Lock()
	frontBuf, backBuf = backBuf, frontBuf
	bufMu.Unlock()
	drawTarget = frontBuf
}

func gibGrafikspalten() uint16 { return fensterBreite }

func gibGrafikzeilen() uint16 { return fensterHoehe }

// archivieren sichert den aktuellen Zeicheninhalt in einem Zwischenspeicher.
func archivieren() {
	if archivBuf == nil {
		archivBuf = ebiten.NewImage(int(fensterBreite), int(fensterHoehe))
	}
	archivBuf.Clear()
	archivBuf.DrawImage(drawTarget, nil)
}

// restaurieren kopiert den Bereich (x,y,b,h) aus dem Zwischenspeicher zurück.
func restaurieren(x, y, b, h uint16) {
	if archivBuf == nil {
		return
	}
	r := image.Rect(int(x), int(y), int(x)+int(b), int(y)+int(h))
	sub, ok := archivBuf.SubImage(r).(*ebiten.Image)
	if !ok {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	drawTarget.DrawImage(sub, op)
}
