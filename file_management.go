package main

import (
	"fmt"
	"io"
	"os"
)

func save() {
	fo, err := os.Create(os.Args[2])

	if err != nil {
		panic(err)
	}

	defer fo.Close()

	for _, v := range Result {
		fo.WriteString(v)
	}
}

func load() (instructions []uint32) {
	fileIn, _ := os.Open(os.Args[1])

	defer fileIn.Close()

	var this uint32

	for {
		_, err := fmt.Fscanf(fileIn, "%v\n", &this)
		if err != nil && err == io.EOF {
			break
		}
		instructions = append(instructions, this)
	}
	return
}
