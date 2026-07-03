package gfx

// text.go — Schriftartenausgabe mit TrueType-Fonts (TTF).
// Geladene Fonts werden gecacht, damit SetzeFont ohne Disk-I/O
// durchläuft, wenn derselbe Font erneut angefordert wird.

import (
	"bytes"
	"os"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"
)

var (
	fontQuelle  *text.GoTextFaceSource
	fontGroesse float64 = 12
	fontPfad    string

	// Cache: Pfad → bereits geladene Fontquelle
	fontCache = make(map[string]*text.GoTextFaceSource)
)

func setzeFont(pfad string, groesse int) bool {
	fontGroesse = float64(groesse)
	if pfad == fontPfad && fontQuelle != nil {
		return true
	}
	if cached, ok := fontCache[pfad]; ok {
		fontQuelle = cached
		fontPfad = pfad
		return true
	}
	daten, err := os.ReadFile(pfad)
	if err != nil {
		return false
	}
	return setzeFontAusDaten(daten, pfad)
}

func setzeFontDaten(daten []byte, name string, groesse int) bool {
	fontGroesse = float64(groesse)
	if name == fontPfad && fontQuelle != nil {
		return true
	}
	if cached, ok := fontCache[name]; ok {
		fontQuelle = cached
		fontPfad = name
		return true
	}
	return setzeFontAusDaten(daten, name)
}

func setzeFontAusDaten(daten []byte, name string) bool {
	quelle, err := text.NewGoTextFaceSource(bytes.NewReader(daten))
	if err != nil {
		return false
	}
	fontCache[name] = quelle
	fontQuelle = quelle
	fontPfad = name
	return true
}

func gibTextBreite(s string) float64 {
	if fontQuelle == nil {
		return 0
	}
	face := &text.GoTextFace{
		Source: fontQuelle,
		Size:   fontGroesse,
	}
	w, _ := text.Measure(s, face, 0)
	return w
}

func schreibeFont(x, y uint16, s string) {
	if fontQuelle == nil {
		return
	}
	face := &text.GoTextFace{
		Source: fontQuelle,
		Size:   fontGroesse,
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(gibStiftfarbe())
	text.Draw(drawTarget, s, face, op)
}

// stdFace ist ein eingebauter Bitmap-Font (7x13) für Schreibe(),
// damit auch ohne vorheriges SetzeFont() Text ausgegeben werden kann.
var stdFace *text.GoXFace

func schreibe(x, y uint16, s string) {
	if stdFace == nil {
		stdFace = text.NewGoXFace(basicfont.Face7x13)
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(gibStiftfarbe())
	text.Draw(drawTarget, s, stdFace, op)
}
