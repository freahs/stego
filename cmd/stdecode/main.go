package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"flag"

	"github.com/freahs/stego"
)

func main() {
	var inPath string
	flag.StringVar(&inPath, "in", "", "input image")
	flag.Parse()

	var inFile *os.File
	var err error

	if inPath == "" {
		inFile = os.Stdin
	} else {
		if inFile, err = os.Open(inPath); err != nil {
			log.Fatal(err)
		}
		defer inFile.Close()
	}

	img, _, err := image.Decode(inFile)
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := stego.Decode(img, &stego.DefaultScrambler{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bytes))
}
