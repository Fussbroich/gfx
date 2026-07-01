package gfx

// fenster.go — Fensterverwaltung, Ebitengine-Game-Loop, Double-Buffering.

import (
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
)

// ===================== ebiten.Game =====================

func (g *gfxGame) Update() error {
	if !offen.Load() {
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
	fensterBereit = make(chan struct{})
	fensterDone = make(chan struct{})

	frontBuf = ebiten.NewImage(int(breite), int(hoehe))
	backBuf = ebiten.NewImage(int(breite), int(hoehe))
	drawTarget = frontBuf

	initEingabe()

	go func() {
		ebiten.SetWindowSize(int(breite), int(hoehe))
		ebiten.SetWindowTitle("gfx")
		ebiten.SetWindowClosingHandled(true)
		offen.Store(true)
		close(fensterBereit)

		_ = ebiten.RunGame(&gfxGame{})

		// Aufräumen nach Ende des Game-Loops
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
	if offen.Load() {
		offen.Store(false)
		// Warte, bis der Game-Loop beendet ist
		select {
		case <-fensterDone:
		case <-time.After(2 * time.Second):
		}
	}
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
