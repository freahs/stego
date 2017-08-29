package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"flag"

	"github.com/freahs/stego"
)

func main() {
	var inPath, outPath string
	var message string
	flag.StringVar(&inPath, "in", "", "input image")
	flag.StringVar(&outPath, "out", "", "output image")
	flag.StringVar(&message, "msg", "", "message to encode")

	flag.Parse()

	if inPath == "" {
		fmt.Println("path to input image must be provided with -in=")
		return
	}

	var inFile, outFile *os.File
	var err error

	if inFile, err = os.Open(inPath); err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()

	inImg, imgType, err := image.Decode(inFile)
	if err != nil {
		log.Fatal(err)
	}

	var bytes []byte
	if message == "" {
		if bytes, err = ioutil.ReadAll(os.Stdin); err != nil {
			log.Fatal(err)
		}
	} else {
		bytes = []byte(message)
	}

	outImg, err := stego.Encode(bytes, inImg, &stego.DefaultScrambler{})
	if err != nil {
		log.Fatal(err)
	}

	if outPath == "" {
		outFile = os.Stdout
	} else {
		if outFile, err = os.Create(outPath); err != nil {
			log.Fatal(err)
		}
		defer outFile.Close()
	}

	switch imgType {
	case "png":
		png.Encode(outFile, outImg)
	}

}
