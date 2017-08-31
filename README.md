# stego
A simple package for image steganography

The repo contains a package for image steganography and main package which compiles to a command line tool for encoding and decoding messages in images

## package usage

The package is well commented and most of it should be self explanatory.

It provides an interface for the distribution of data in an image called ´Scrambler´.

There are two exported function: `Encode([]byte, image.Image, *Scrambler) (image.Image, error)` and `Decode(image.Image, *Scrambler) ([]byte, error)`.

`Encode` will encode the bytes in the image using the scrambler and return a new image.

`Decode` will decode the image using the scrambler and return a byte array with the decoded data

Both functions returns a `StegoError` if something goes wrong

## Command line utility usage

Encodeing is done with `stego encode <INPUT IMAGE PATH> [OUTPUT IMAGE PATH] [MESSAGE]` where MESSAGE and OUTPUT IMAGE PATH are optional. If ommited, the message will be read from stdin and the resulting image will be forwarded to stdout.
		
Decoding are done with `stego decode [INPUT IMAGE PATH]` where INPUT IMAGE PATH are optional. If ommited the image to decode will be read from stdin.

