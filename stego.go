package stego

import (
	"encoding/binary"
	"image"
)

type StegoError string

func (e StegoError) Error() string {
	return string(e)
}

const (
	SOH        byte = 1 // Start of header
	STX        byte = 2 // Start of text (end of header)
	headerSize int  = 6 // Bytes in header
)

// Scrambler is used to scramble data in an image, but an implementation should not alter the image
// itself.
//
// The Init method should put the scrambler in a state from which Next should always generate the
// same sequence of results.
//
// Next should return whichever pixel and color channel that should be modified and/or read next.
//
// Cap should return the maximum nunmber of bits for which the scrambler can produce color
// coordinates given the image used when last calling Init. If Init hasn't been called before Cap,
// Cap should return 0.
type Scrambler interface {
	Init(image image.Image)
	Next() (x, y, c int)
	Cap() int
}

// Encode encodes a slice of Byte in the Image using the Scrambler. If it fails for some reason,
// an error is returned.
func Encode(bytes []byte, img image.Image, scrambler Scrambler) (image.Image, error) {

	source := img.(*image.RGBA)
	scrambler.Init(source)

	// Make a slice of header and argument bytes, to encode in the image.
	data := make([]byte, headerSize)
	data[0] = SOH
	binary.LittleEndian.PutUint32(data[1:5], uint32(len(bytes)))
	data[5] = STX
	data = append(data, bytes...)

	// Make sure the scrambler can fit all data in the image.
	if len(data)*8 > scrambler.Cap() {
		return nil, StegoError("image to small.")
	}

	for _, b := range data {
		writeByte(b, source, scrambler)
	}

	output := image.NewRGBA(source.Bounds())
	output.Pix = source.Pix
	return output, nil
}

// Decode reads bytes from a Image using a Scrambler. If it fails for some reason, an error is
// returned
func Decode(img image.Image, scrambler Scrambler) ([]byte, error) {
	source := img.(*image.RGBA)
	scrambler.Init(source)

	// Read the header...
	header := make([]byte, headerSize)
	for i := 0; i < len(header); i++ {
		header[i] = readByte(source, scrambler)

	}

	// ...abort if it's not formated correct.
	if header[0] != SOH || header[5] != STX {
		return nil, StegoError("Unknown header.")
	}

	// If the size is too large for the image, the header is probably not correct. Either way
	size := int(binary.LittleEndian.Uint32(header[1:5]))
	if size*8 > scrambler.Cap() {
		return nil, StegoError("Unknown header.")
	}

	// Read as many bytes as specified in the header and return the result.
	ret := make([]byte, size)
	for i := 0; i < size; i++ {
		ret[i] = readByte(source, scrambler)
	}
	return ret, nil
}

// readByte reads a byte from a image using a scrambler.
func readByte(img *image.RGBA, scrambler Scrambler) (b byte) {
	var i uint8
	for i = 0; i < 8; i++ {
		x, y, c := scrambler.Next() // Next pixel coords and channel.
		p := img.PixOffset(x, y) + c
		if img.Pix[p]&0x01 == 1 { // Read current bit.
			b += (1 << i) // Icrease byte if necessary.
		}
	}
	return
}

// writeByte writes a byte to a image using a scrambler.
func writeByte(b byte, img *image.RGBA, scrambler Scrambler) {
	var i uint8
	for i = 0; i < 8; i++ { // Step throug each bit
		bit := (b>>i)&0x01 == 1
		x, y, c := scrambler.Next()
		p := img.PixOffset(x, y) + c
		if bit { // Set the bit.
			img.Pix[p] = img.Pix[p] | 0x01 // LSB = 1
		} else {
			img.Pix[p] = img.Pix[p] & 0xFE // LSB = 0
		}
	}
}
