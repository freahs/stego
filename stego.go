package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"

	stego "github.com/freahs/stego/pkg"
)

func encode(argc int, argv []string) {
	inFile, err := os.Open(argv[2])
	if err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()

	inImg, imgType, err := image.Decode(inFile)
	if err != nil {
		log.Fatal(err)
	}

	var bytes []byte
	if argc < 5 {
		if bytes, err = ioutil.ReadAll(os.Stdin); err != nil {
			log.Fatal(err)
		}
	} else {
		bytes = []byte(argv[4])
	}

	outImg, err := stego.Encode(bytes, inImg, &stego.DefaultScrambler{})
	if err != nil {
		log.Fatal(err)
	}

	var outFile *os.File
	if argc < 4 {
		outFile = os.Stdout
	} else {
		if outFile, err = os.Create(argv[3]); err != nil {
			log.Fatal(err)
		}
	}
	defer outFile.Close()

	switch imgType {
	case "png":
		png.Encode(outFile, outImg)
	}

}

func decode(argc int, argv []string) {

	var inFile *os.File
	var err error

	if argc < 3 {
		inFile = os.Stdin
	} else {
		if inFile, err = os.Open(argv[2]); err != nil {
			log.Fatal(err)
		}
	}
	defer inFile.Close()

	img, _, err := image.Decode(inFile)
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := stego.Decode(img, &stego.DefaultScrambler{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(bytes))
}

func usage() {
	fmt.Println(`
usage: stego <encode | decode> [args]

ENCODE:
    stego encode <INPUT IMAGE PATH> [OUTPUT IMAGE PATH] [MESSAGE]
	
Without MESSAGE the message to encode will be read from stdin. Without OUTPUT IMAGE PATH the resulting image will be forwarded to stdout
		

DECODE:
    stego decode [INPUT IMAGE PATH]

Without INPUT IMAGE PATH the image will be read from stdin. The resulting messge will always be printed on stdout.`)

}

func main() {
	num_args := len(os.Args)
	if num_args >= 3 && strings.ToLower(os.Args[1]) == "encode" && num_args <= 5 {
		encode(num_args, os.Args)
	} else if num_args >= 2 && strings.ToLower(os.Args[1]) == "decode" && num_args <= 3 {
		decode(num_args, os.Args)
	} else {
		usage()
	}

}
