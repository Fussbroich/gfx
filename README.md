# gfx

Einfache 2D-Grafik-, Sound- und Eingabe-Bibliothek für Lernzwecke.

Die API orientiert sich an der ursprünglichen gfx/gfxw-Bibliothek von
St. Schmidt (FU Berlin), verwendet intern jedoch [Ebitengine](https://ebitengine.org/)
statt SDL und benötigt keinen externen Server-Prozess.

## Installation

```
go get github.com/Fussbroich/gfx
```

## Verwendung

```go
import "github.com/Fussbroich/gfx"

func main() {
	gfx.Fenster(800, 600)
	gfx.Stiftfarbe(255, 0, 0)
	gfx.Vollkreis(400, 300, 50)
	gfx.TastaturLesen1()
}
```
