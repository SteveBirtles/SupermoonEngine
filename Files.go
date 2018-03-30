package main

import (
	"encoding/gob"
	"os"
)

func save() {
	f1, err := os.Create(levelFile)
	check(err)
	encoder1 := gob.NewEncoder(f1)
	check(encoder1.Encode(grid))

	f2, err := os.Create(clipboardFile)
	check(err)
	encoder2 := gob.NewEncoder(f2)
	check(encoder2.Encode(clipboard))
	check(encoder2.Encode(clipboardWidth))
	check(encoder2.Encode(clipboardHeight))
}

func load() {
	f1, err := os.Open(levelFile)
	if err != nil {
		f1, err = os.Open(defaultLevelFile)
	}
	if err == nil {
		decoder1 := gob.NewDecoder(f1)
		check(decoder1.Decode(&grid))
	}

	f2, err := os.Open(clipboardFile)
	if err != nil {
		f2, err = os.Open(defaultClipboardFile)
	}
	if err == nil {
		decoder2 := gob.NewDecoder(f2)
		check(decoder2.Decode(&clipboard))
		check(decoder2.Decode(&clipboardWidth))
		check(decoder2.Decode(&clipboardHeight))
	}
}

func backup() {
	f, err := os.Create(levelFile + ".bak")
	if err == nil {
		encoder := gob.NewEncoder(f)
		check(encoder.Encode(grid))
	}
}
