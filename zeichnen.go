package gfx

// zeichnen.go — Zeichenprimitiven mit wiederverwendbaren Puffern.
// ALLE Primitiven nutzen einen gemeinsamen vector.Path und Vertex/Index-Puffer,
// um pro Frame keine Heap-Allokationen zu erzeugen.

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var whitePixel *ebiten.Image

func init() {
	whitePixel = ebiten.NewImage(1, 1)
	whitePixel.Fill(color.White)
}

// Wiederverwendbare Puffer — einmal allokiert, dann recycled.
var (
	sp vector.Path
	sv []ebiten.Vertex
	si []uint16
)

// pfadFuellen zeichnet den aktuellen sharedPath gefüllt.
func pfadFuellen(clr color.NRGBA) {
	sv, si = sp.AppendVerticesAndIndicesForFilling(sv[:0], si[:0])
	farbeAufVertices(sv, clr)
	drawTarget.DrawTriangles(sv, si, whitePixel, nil)
}

// pfadUmriss zeichnet den aktuellen sharedPath als Umriss.
func pfadUmriss(clr color.NRGBA, breite float32) {
	sv, si = sp.AppendVerticesAndIndicesForStroke(sv[:0], si[:0], &vector.StrokeOptions{Width: breite})
	farbeAufVertices(sv, clr)
	drawTarget.DrawTriangles(sv, si, whitePixel, nil)
}

// ===================== Linien =====================

func linie(x1, y1, x2, y2 uint16) {
	sp.Reset()
	sp.MoveTo(float32(x1), float32(y1))
	sp.LineTo(float32(x2), float32(y2))
	pfadUmriss(gibStiftfarbe(), 1)
}

// ===================== Kreise =====================

const kreisSegmente = 32

func kreispfad(cx, cy, r float32) {
	sp.Reset()
	for i := 0; i <= kreisSegmente; i++ {
		a := float64(i) * 2 * math.Pi / kreisSegmente
		px := cx + r*float32(math.Cos(a))
		py := cy + r*float32(math.Sin(a))
		if i == 0 {
			sp.MoveTo(px, py)
		} else {
			sp.LineTo(px, py)
		}
	}
	sp.Close()
}

func kreis(x, y, r uint16) {
	kreispfad(float32(x), float32(y), float32(r))
	pfadUmriss(gibStiftfarbe(), 1)
}

func vollkreis(x, y, r uint16) {
	kreispfad(float32(x), float32(y), float32(r))
	pfadFuellen(gibStiftfarbe())
}

// ===================== Rechtecke =====================

func rechteckpfad(x, y, b, h float32) {
	sp.Reset()
	sp.MoveTo(x, y)
	sp.LineTo(x+b, y)
	sp.LineTo(x+b, y+h)
	sp.LineTo(x, y+h)
	sp.Close()
}

func rechteck(x1, y1, b, h uint16) {
	rechteckpfad(float32(x1), float32(y1), float32(b), float32(h))
	pfadUmriss(gibStiftfarbe(), 1)
}

func vollrechteck(x1, y1, b, h uint16) {
	rechteckpfad(float32(x1), float32(y1), float32(b), float32(h))
	pfadFuellen(gibStiftfarbe())
}

// ===================== Kreissektoren =====================
// Konvention: Winkel in Grad, 0° = Osten, gegen Uhrzeigersinn.

func gradAufBildschirm(grad float64) (float64, float64) {
	rad := grad * math.Pi / 180.0
	return math.Cos(rad), -math.Sin(rad)
}

func sektorpfad(cx, cy, r float32, w1, w2 uint16) {
	start := float64(w1)
	end := float64(w2)
	if end <= start {
		end += 360
	}
	segments := int((end-start)/5) + 1
	if segments < 4 {
		segments = 4
	}

	sp.Reset()
	sp.MoveTo(cx, cy)
	for i := 0; i <= segments; i++ {
		t := start + (end-start)*float64(i)/float64(segments)
		cosT, sinT := gradAufBildschirm(t)
		sp.LineTo(cx+r*float32(cosT), cy+r*float32(sinT))
	}
	sp.Close()
}

func kreissektor(x, y, r, w1, w2 uint16) {
	sektorpfad(float32(x), float32(y), float32(r), w1, w2)
	pfadUmriss(gibStiftfarbe(), 1)
}

func vollkreissektor(x, y, r, w1, w2 uint16) {
	sektorpfad(float32(x), float32(y), float32(r), w1, w2)
	pfadFuellen(gibStiftfarbe())
}

// ===================== Dreiecke =====================

func volldreieck(x1, y1, x2, y2, x3, y3 uint16) {
	sp.Reset()
	sp.MoveTo(float32(x1), float32(y1))
	sp.LineTo(float32(x2), float32(y2))
	sp.LineTo(float32(x3), float32(y3))
	sp.Close()
	pfadFuellen(gibStiftfarbe())
}

// ===================== Ellipsen =====================
// Mittelpunkt (cx,cy), horizontale Halbachse rx, vertikale Halbachse ry.

func ellipsepfad(cx, cy, rx, ry float32) {
	sp.Reset()
	for i := 0; i <= kreisSegmente; i++ {
		a := float64(i) * 2 * math.Pi / kreisSegmente
		px := cx + rx*float32(math.Cos(a))
		py := cy + ry*float32(math.Sin(a))
		if i == 0 {
			sp.MoveTo(px, py)
		} else {
			sp.LineTo(px, py)
		}
	}
	sp.Close()
}

func ellipse(x, y, rx, ry uint16) {
	ellipsepfad(float32(x), float32(y), float32(rx), float32(ry))
	pfadUmriss(gibStiftfarbe(), 1)
}

func vollellipse(x, y, rx, ry uint16) {
	ellipsepfad(float32(x), float32(y), float32(rx), float32(ry))
	pfadFuellen(gibStiftfarbe())
}

// ===================== Hilfsfunktion =====================

func farbeAufVertices(vs []ebiten.Vertex, clr color.NRGBA) {
	cr := float32(clr.R) / 255
	cg := float32(clr.G) / 255
	cb := float32(clr.B) / 255
	ca := float32(clr.A) / 255
	for i := range vs {
		vs[i].SrcX = 0
		vs[i].SrcY = 0
		vs[i].ColorR = cr * ca
		vs[i].ColorG = cg * ca
		vs[i].ColorB = cb * ca
		vs[i].ColorA = ca
	}
}
