package main

import (
	"encoding/gob"
	"os"
	"io"
)

func save() {
	f1, err := os.Create(levelFile)
	check(err)
	encoder1 := gob.NewEncoder(f1)
	check(encoder1.Encode(grid))
	check(encoder1.Encode(entities[0]))

	f2, err := os.Create(clipboardFile)
	check(err)
	encoder2 := gob.NewEncoder(f2)
	check(encoder2.Encode(clipboard))
	check(encoder2.Encode(clipboardWidth))
	check(encoder2.Encode(clipboardHeight))
}

func load() {

	f1, err1 := os.Open(levelFile)
	if err1 == nil {
		decoder1 := gob.NewDecoder(f1)
		check(decoder1.Decode(&grid))

		err2 := decoder1.Decode(&entities[0])
		if err2 != io.EOF {
			check(err2)
			for i := range entities[0] {
				if entities[0][i].Id > entityUID {
					entityUID = entities[0][i].Id
				}
				entities[0][i].active = true
				entities[0][i].sprite = [4]int{-1, -1, -1, -1}
				entities[0][i].velocity = 0
				entities[0][i].direction = 'S'
				entities[0][i].distance = 0
				entities[0][i].lastX = entities[0][i].X
				entities[0][i].lastY = entities[0][i].Y
				entities[0][i].lastZ = entities[0][i].Z
				entities[0][i].targetX = entities[0][i].X
				entities[0][i].targetY = entities[0][i].Y
				entities[0][i].targetZ = entities[0][i].Z
				entities[0][i].progress = 0
			}
		}

	}

	f2, err1 := os.Open(clipboardFile)
	if err1 == nil {
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
