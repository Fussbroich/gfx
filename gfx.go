// Paket gfx stellt eine einfache 2D-Grafik-, Sound- und Eingabe-Bibliothek
// für Lernzwecke bereit. Die API orientiert sich an der ursprünglichen
// gfx/gfxw-Bibliothek von St. Schmidt (FU Berlin), verwendet intern
// jedoch Ebitengine statt SDL und benötigt keinen externen Server-Prozess.
//
// Alle Zeichenoperationen sind nebenläufig sicher.
package gfx

import "io"

// ===================== Fensterverwaltung =====================

// Fenster öffnet ein Grafikfenster mit breite x hoehe Pixeln.
// Der Ursprung (0,0) ist oben links. Die Stiftfarbe ist Schwarz.
func Fenster(breite, hoehe uint16) { fensterStarten(breite, hoehe) }

// FensterOffen liefert true, wenn das Grafikfenster geöffnet ist.
func FensterOffen() bool { return istFensterOffen() }

// FensterAus schließt das Grafikfenster.
func FensterAus() { fensterSchliessen() }

// Fenstertitel setzt den Titel in der Titelleiste des Fensters.
func Fenstertitel(s string) { setzeFenstertitel(s) }

// Grafikspalten liefert die Breite des Fensters in Pixeln.
func Grafikspalten() uint16 { return gibGrafikspalten() }

// Grafikzeilen liefert die Höhe des Fensters in Pixeln.
func Grafikzeilen() uint16 { return gibGrafikzeilen() }

// ===================== Zeichensteuerung =====================

// Cls füllt das gesamte Fenster mit der aktuellen Stiftfarbe.
func Cls() { clsBuf() }

// UpdateAus aktiviert Double-Buffering: Zeichenbefehle gehen
// in einen Hintergrundpuffer und sind noch nicht sichtbar.
func UpdateAus() { updateAus() }

// UpdateAn zeigt alle seit UpdateAus gepufferten Zeichenbefehle an.
func UpdateAn() { updateAn() }

// Archivieren sichert den gesamten aktuellen Fensterinhalt in einem
// Zwischenspeicher (siehe Restaurieren).
func Archivieren() { archivieren() }

// Restaurieren kopiert den rechteckigen Bereich (x,y) mit Breite b und Höhe h
// aus dem mit Archivieren gesicherten Zwischenspeicher zurück ins Fenster.
func Restaurieren(x, y, b, h uint16) { restaurieren(x, y, b, h) }

// Stiftfarbe setzt die Zeichenfarbe im RGB-Modell (je 0–255).
func Stiftfarbe(r, g, b uint8) { setzeStiftfarbe(r, g, b) }

// Transparenz stellt die Durchsichtigkeit ein.
// 0 = deckend (Standard), 255 = vollständig durchsichtig.
func Transparenz(t uint8) { setzeTransparenz(t) }

// Sperren blockiert andere Goroutinen am Zeichnen.
func Sperren() { zeichenSperre.Lock() }

// Entsperren gibt das Zeichnen wieder frei.
func Entsperren() { zeichenSperre.Unlock() }

// ===================== Zeichenprimitiven =====================

// Linie zeichnet eine 1px-Strecke von (x1,y1) nach (x2,y2).
func Linie(x1, y1, x2, y2 uint16) { linie(x1, y1, x2, y2) }

// Kreis zeichnet einen Kreisumriss um (x,y) mit Radius r.
func Kreis(x, y, r uint16) { kreis(x, y, r) }

// Vollkreis zeichnet einen ausgefüllten Kreis um (x,y) mit Radius r.
func Vollkreis(x, y, r uint16) { vollkreis(x, y, r) }

// Ellipse zeichnet einen Ellipsenumriss um (x,y) mit Halbachsen rx und ry.
func Ellipse(x, y, rx, ry uint16) { ellipse(x, y, rx, ry) }

// Vollellipse zeichnet eine ausgefüllte Ellipse um (x,y) mit Halbachsen rx und ry.
func Vollellipse(x, y, rx, ry uint16) { vollellipse(x, y, rx, ry) }

// Rechteck zeichnet einen Rechteckumriss bei (x1,y1) mit Breite b und Höhe h.
func Rechteck(x1, y1, b, h uint16) { rechteck(x1, y1, b, h) }

// Vollrechteck zeichnet ein ausgefülltes Rechteck.
func Vollrechteck(x1, y1, b, h uint16) { vollrechteck(x1, y1, b, h) }

// Kreissektor zeichnet einen Kreisbogen von Winkel w1 bis w2 (Grad, 0=Ost, gegen Uhrzeiger).
func Kreissektor(x, y, r, w1, w2 uint16) { kreissektor(x, y, r, w1, w2) }

// Vollkreissektor zeichnet einen ausgefüllten Kreissektor.
func Vollkreissektor(x, y, r, w1, w2 uint16) { vollkreissektor(x, y, r, w1, w2) }

// Volldreieck zeichnet ein ausgefülltes Dreieck mit Ecken (x1,y1), (x2,y2), (x3,y3).
func Volldreieck(x1, y1, x2, y2, x3, y3 uint16) { volldreieck(x1, y1, x2, y2, x3, y3) }

// ===================== Textausgabe =====================

// Schreibe schreibt Text an Position (x,y) mit einem eingebauten
// Standardfont — ohne dass zuvor SetzeFont aufgerufen werden muss.
// (x,y) ist die obere linke Ecke des Textes.
func Schreibe(x, y uint16, s string) { schreibe(x, y, s) }

// SetzeFont lädt eine TTF-Datei als aktuellen Zeichensatz.
// Liefert true bei Erfolg.
func SetzeFont(s string, groesse int) bool { return setzeFont(s, groesse) }

// SchreibeFont schreibt Text an Position (x,y) mit dem zuletzt gesetzten Font.
func SchreibeFont(x, y uint16, s string) { schreibeFont(x, y, s) }

// SetzeFontDaten lädt einen Font aus bereits gelesenen TTF-Daten.
// name dient als Cache-Schlüssel (z.B. "MonoBold").
func SetzeFontDaten(daten []byte, name string, groesse int) bool {
	return setzeFontDaten(daten, name, groesse)
}

// GibTextBreite liefert die Breite des Textes in Pixeln
// mit dem aktuell gesetzten Font und Schriftgröße.
func GibTextBreite(s string) float64 { return gibTextBreite(s) }

// ===================== Eingabe =====================

// TastaturLesen1 blockiert, bis eine Taste gedrückt oder losgelassen wird.
// taste: Tastencode, gedrueckt: 1=gedrückt/0=losgelassen, tiefe: Modifikatoren.
func TastaturLesen1() (taste uint16, gedrueckt uint8, tiefe uint16) {
	return tastaturLesen1()
}

// MausLesen1 blockiert, bis ein Mausereignis eintritt.
// taste: Maustaste (1=links, 4=ScrollHoch, 5=ScrollRunter),
// status: 1=gedrückt, -1=losgelassen, 0=Bewegung.
func MausLesen1() (taste uint8, status int8, mausX, mausY uint16) {
	return mausLesen1()
}

// ===================== Sound =====================

// SetzeKlangparameter konfiguriert die Audioausgabe.
// Bei Ebitengine intern verwaltet — wird für Kompatibilität akzeptiert.
func SetzeKlangparameter(rate uint32, aufloesung, kanaele, signal uint8, p float64) {
	setzeKlangparameter(rate, aufloesung, kanaele, signal, p)
}

// SpieleSound spielt eine WAV-Datei im Hintergrund ab.
func SpieleSound(s string) { spieleSound(s) }

// SpieleSoundDaten spielt WAV-Daten aus dem Speicher ab.
// name dient als Cache-Schlüssel.
func SpieleSoundDaten(daten []byte, name string) { spieleSoundDaten(daten, name) }

// SpieleSoundStream streamt eine WAV-Datei direkt aus einem Reader.
// Kein PCM-Cache, minimaler RAM-Verbrauch. Ideal für Musik.
func SpieleSoundStream(r io.ReadSeeker) { spieleSoundStream(r) }
