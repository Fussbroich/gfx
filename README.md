# gfx (Go)

**Einfache 2D-Grafik, Sound und Eingabe für den Informatik-Unterricht — in Go,
auf Basis von [Ebitengine](https://ebitengine.org/).**

`gfx` ist die moderne Neufassung der klassischen *gfx*-/*gfxw*-Bibliothek von
Stefan Schmidt (FU Berlin). Die API bleibt gleich (deutsche Befehle,
blockierende Eingabe), aber intern läuft jetzt alles **rein in Go über
Ebitengine** — **kein SDL, keine DLLs, kein externer Server-Prozess**. Ein
`go get`, und es läuft unter Windows, Linux und macOS.

---

## Ein erstes Programm

```go
package main

import "github.com/Fussbroich/gfx"

func main() {
	gfx.Fenster(800, 600)
	gfx.Stiftfarbe(255, 0, 0)
	gfx.Vollkreis(400, 300, 50)
	gfx.TastaturLesen1()   // wartet auf einen Tastendruck
	gfx.FensterAus()
}
```

Ein roter Kreis, dann ein Tastendruck zum Beenden — mehr braucht es nicht.

---

## Das Besondere: geradliniges Programmieren

Für den Unterricht entscheidend: Man schreibt **normalen, sequentiellen Code**
— keine Spielschleife, keine Callbacks.

`gfx.Fenster(...)` startet den Ebitengine-Renderloop im Hintergrund (eigene
Goroutine). Das eigene Programm läuft einfach von oben nach unten weiter und
**blockiert** an `gfx.TastaturLesen1()` bzw. `gfx.MausLesen1()`, bis eine
Eingabe kommt. So können Anfängerinnen und Anfänger mit dem vertrauten
„erst dies, dann das"-Denken arbeiten.

Gezeichnet wird **bleibend** (retained mode): Was einmal gemalt ist, bleibt
stehen, bis man es überzeichnet oder mit `Cls()` löscht. Für flackerfreie
Animationen puffert man mit `UpdateAus()` / `UpdateAn()` doppelt.

---

## Funktionsumfang

- **Fenster** — `Fenster`, `FensterOffen`, `FensterAus`, `Fenstertitel`,
  `Grafikspalten`, `Grafikzeilen`.
- **Zeichensteuerung** — `Cls`, `Stiftfarbe`, `Transparenz`,
  `UpdateAus`/`UpdateAn` (Double-Buffering), `Archivieren`/`Restaurieren`
  (Bildbereich sichern & zurückholen), `Sperren`/`Entsperren`.
- **Formen** — `Linie`, `Kreis`/`Vollkreis`, `Ellipse`/`Vollellipse`,
  `Rechteck`/`Vollrechteck`, `Kreissektor`/`Vollkreissektor`, `Volldreieck`.
- **Text** — `Schreibe` (eingebauter Font, ohne Vorbereitung),
  `SetzeFont`/`SetzeFontDaten` + `SchreibeFont` (eigene TTF), `GibTextBreite`.
- **Eingabe (blockierend)** — `TastaturLesen1`, `MausLesen1`. Die Tastencodes
  sind wie früher SDL-kompatibel (z. B. Enter = 13, ESC = 27, Pfeile 273–276).
- **Sound** — `SpieleSound`, `SpieleSoundDaten`, `SpieleSoundStream`,
  `SetzeKlangparameter`.

---

## Installation

```
go get github.com/Fussbroich/gfx
```

Voraussetzung ist eine aktuelle **Go-Toolchain (1.25+)**. Ebitengine und die
übrigen Abhängigkeiten holt `go` automatisch. **Keine** SDL-Bibliotheken,
**kein** separates Serverprogramm.

---

## Unterschiede zur alten gfxw

Wer von der klassischen SDL-/Server-Fassung kommt:

- **Kein `gfxwserver.exe`** mehr und keine SDL-/Freetype-DLLs — alles steckt im
  Go-Modul.
- **Plattformübergreifend** (Windows/Linux/macOS) über Ebitengine, ohne C-Compiler.
- **API-kompatibel**: bestehende Lehr-Programme laufen weiter — gleiche
  Funktionsnamen, gleiches blockierendes `TastaturLesen1`, gleiche Tastencodes.
- Das Fenster lässt sich zur Laufzeit in der **Größe ändern** (erneuter
  `Fenster`-Aufruf), ohne den Renderloop neu zu starten.

---

## Verwandtes

Auf dieser `gfx` baut u. a. **robiw** auf — der 2D-Lernroboter „Robi" in einer
Kachelwelt, mit dem sich das Entwerfen von Algorithmen üben lässt.

---

## Credits

| Rolle | |
|---|---|
| **Ursprüngliche gfx-/gfxw-Bibliothek** — API-Entwurf & didaktisches Konzept | Stefan Schmidt, Lehrerweiterbildung Informatik, Freie Universität Berlin |
| **Go-/Ebitengine-Neufassung** | Thomas Schrader (GitHub: Fussbroich) |
| **Grundlage** | [Ebitengine](https://ebitengine.org/) |

Vielen Dank an Stefan Schmidt, dessen klar entworfene gfx-API dieser Bibliothek
zugrunde liegt.

© Thomas Schrader 2026
