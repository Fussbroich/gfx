package gfx

// eingabe.go — Tastatur- und Mauseingabe.
// Die blockierenden Funktionen TastaturLesen1 und MausLesen1 werden
// über Channels bedient, die vom Ebitengine-Update-Loop befüllt werden.

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type tastaturEvent struct {
	taste    uint16
	gedrueckt uint8
	tiefe    uint16
}

type mausEvent struct {
	taste  uint8
	status int8
	mausX  uint16
	mausY  uint16
}

var (
	tastaturChan chan tastaturEvent
	mausChan     chan mausEvent
	letztesMausX int
	letztesMausY int
)

func initEingabe() {
	tastaturChan = make(chan tastaturEvent, 256)
	mausChan = make(chan mausEvent, 256)
}

// closeEingabe wird nicht mehr benötigt — fensterDone übernimmt das Signaling.

// ===================== Ebitengine → Channel =====================

// updateEingabe wird jeden Tick aus gfxGame.Update() aufgerufen
// und wandelt Ebitengine-Eingaben in Channel-Events um.
func updateEingabe() {
	// --- Tastatur ---
	keys := inpututil.AppendJustPressedKeys(nil)
	for _, k := range keys {
		if code := keyToCode(k); code != 0 {
			sendTastatur(tastaturEvent{taste: code, gedrueckt: 1})
		}
	}
	keysRel := inpututil.AppendJustReleasedKeys(nil)
	for _, k := range keysRel {
		if code := keyToCode(k); code != 0 {
			sendTastatur(tastaturEvent{taste: code, gedrueckt: 0})
		}
	}

	// --- Maus: Buttons ---
	mx, my := ebiten.CursorPosition()
	for _, btn := range []ebiten.MouseButton{
		ebiten.MouseButtonLeft, ebiten.MouseButtonMiddle, ebiten.MouseButtonRight,
	} {
		code := mausButtonCode(btn)
		if inpututil.IsMouseButtonJustPressed(btn) {
			sendMaus(mausEvent{taste: code, status: 1, mausX: uint16(mx), mausY: uint16(my)})
		}
		if inpututil.IsMouseButtonJustReleased(btn) {
			sendMaus(mausEvent{taste: code, status: -1, mausX: uint16(mx), mausY: uint16(my)})
		}
	}

	// --- Maus: Scrollrad ---
	_, scrollY := ebiten.Wheel()
	if scrollY > 0 {
		sendMaus(mausEvent{taste: 4, status: 1, mausX: uint16(mx), mausY: uint16(my)})
	} else if scrollY < 0 {
		sendMaus(mausEvent{taste: 5, status: 1, mausX: uint16(mx), mausY: uint16(my)})
	}

	// --- Maus: Bewegung ---
	if mx != letztesMausX || my != letztesMausY {
		sendMaus(mausEvent{taste: 0, status: 0, mausX: uint16(mx), mausY: uint16(my)})
		letztesMausX, letztesMausY = mx, my
	}
}

func sendTastatur(ev tastaturEvent) {
	select {
	case tastaturChan <- ev:
	default: // Puffer voll — Event verwerfen
	}
}

func sendMaus(ev mausEvent) {
	select {
	case mausChan <- ev:
	default:
	}
}

// ===================== Blockierende API =====================

func tastaturLesen1() (uint16, uint8, uint16) {
	select {
	case ev := <-tastaturChan:
		return ev.taste, ev.gedrueckt, ev.tiefe
	case <-fensterDone:
		panic("Das Grafikfenster wurde geschlossen! Programmabbruch!!")
	}
}

func mausLesen1() (uint8, int8, uint16, uint16) {
	select {
	case ev := <-mausChan:
		return ev.taste, ev.status, ev.mausX, ev.mausY
	case <-fensterDone:
		panic("Das Grafikfenster wurde geschlossen! Programmabbruch!!")
	}
}

// ===================== Tasten-Mapping =====================
// Bildet Ebitengine-Keys auf die gleichen Zeichencodes ab,
// die die alte gfx-Bibliothek (SDL-Keycodes) verwendete.

func keyToCode(k ebiten.Key) uint16 {
	switch k {
	case ebiten.KeyA: return 'a'
	case ebiten.KeyB: return 'b'
	case ebiten.KeyC: return 'c'
	case ebiten.KeyD: return 'd'
	case ebiten.KeyE: return 'e'
	case ebiten.KeyF: return 'f'
	case ebiten.KeyG: return 'g'
	case ebiten.KeyH: return 'h'
	case ebiten.KeyI: return 'i'
	case ebiten.KeyJ: return 'j'
	case ebiten.KeyK: return 'k'
	case ebiten.KeyL: return 'l'
	case ebiten.KeyM: return 'm'
	case ebiten.KeyN: return 'n'
	case ebiten.KeyO: return 'o'
	case ebiten.KeyP: return 'p'
	case ebiten.KeyQ: return 'q'
	case ebiten.KeyR: return 'r'
	case ebiten.KeyS: return 's'
	case ebiten.KeyT: return 't'
	case ebiten.KeyU: return 'u'
	case ebiten.KeyV: return 'v'
	case ebiten.KeyW: return 'w'
	case ebiten.KeyX: return 'x'
	case ebiten.KeyY: return 'y'
	case ebiten.KeyZ: return 'z'
	case ebiten.KeyDigit0: return '0'
	case ebiten.KeyDigit1: return '1'
	case ebiten.KeyDigit2: return '2'
	case ebiten.KeyDigit3: return '3'
	case ebiten.KeyDigit4: return '4'
	case ebiten.KeyDigit5: return '5'
	case ebiten.KeyDigit6: return '6'
	case ebiten.KeyDigit7: return '7'
	case ebiten.KeyDigit8: return '8'
	case ebiten.KeyDigit9: return '9'
	case ebiten.KeySpace:  return ' '
	case ebiten.KeyEnter:  return 13
	case ebiten.KeyEscape: return 27
	case ebiten.KeyArrowUp:    return 273
	case ebiten.KeyArrowDown:  return 274
	case ebiten.KeyArrowRight: return 275
	case ebiten.KeyArrowLeft:  return 276
	}
	return 0
}

func mausButtonCode(b ebiten.MouseButton) uint8 {
	switch b {
	case ebiten.MouseButtonLeft:   return 1
	case ebiten.MouseButtonMiddle: return 2
	case ebiten.MouseButtonRight:  return 3
	}
	return 0
}
