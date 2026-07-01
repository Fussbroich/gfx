package gfx

// klang.go — WAV-Sound-Wiedergabe über Ebitengine-Audio.
// Geladene WAV-Dateien werden gecacht, damit wiederholtes Abspielen
// desselben Sounds ohne Disk-I/O auskommt.

import (
	"bytes"
	"io"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const audioSampleRate = 44100

var (
	audioCtx *audio.Context

	// Cache: Pfad → dekodierte PCM-Daten ([]byte)
	soundCache = make(map[string][]byte)
)

func initAudio() {
	if audioCtx == nil {
		audioCtx = audio.NewContext(audioSampleRate)
	}
}

// setzeKlangparameter wird für Kompatibilität akzeptiert.
// Ebitengine verwaltet Audioparameter intern — die Werte werden ignoriert.
func setzeKlangparameter(_ uint32, _, _, _ uint8, _ float64) {
	initAudio()
}

// spieleSound spielt eine WAV-Datei im Hintergrund ab.
func spieleSound(pfad string) {
	initAudio()
	pcm, ok := soundCache[pfad]
	if !ok {
		f, err := os.Open(pfad)
		if err != nil {
			return
		}
		pcm = dekodiereWAV(f)
		f.Close()
		if pcm == nil {
			return
		}
		soundCache[pfad] = pcm
	}
	audioCtx.NewPlayerFromBytes(pcm).Play()
}

// spieleSoundDaten spielt WAV-Daten aus dem Speicher ab.
func spieleSoundDaten(daten []byte, name string) {
	initAudio()
	pcm, ok := soundCache[name]
	if !ok {
		pcm = dekodiereWAV(bytes.NewReader(daten))
		if pcm == nil {
			return
		}
		soundCache[name] = pcm
	}
	audioCtx.NewPlayerFromBytes(pcm).Play()
}

// spieleSoundStream streamt WAV direkt aus einem Reader.
// Kein Heap-Alloc für PCM — ideal für große Musikdateien.
func spieleSoundStream(r io.ReadSeeker) {
	initAudio()
	stream, err := wav.DecodeWithSampleRate(audioSampleRate, r)
	if err != nil {
		return
	}
	player, err := audioCtx.NewPlayer(stream)
	if err != nil {
		return
	}
	player.Play()
}

func dekodiereWAV(r io.ReadSeeker) []byte {
	stream, err := wav.DecodeWithSampleRate(audioSampleRate, r)
	if err != nil {
		return nil
	}
	daten, err := io.ReadAll(stream)
	if err != nil {
		return nil
	}
	return daten
}
