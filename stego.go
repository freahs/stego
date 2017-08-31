package stego

// A simple package for image steganography which probably shouldn't be used for actual sensitive
// data.

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

// stegoImage wraps the underlying byte array of an image and provides convenience methods
// for reading and writing bits to it.
type stegoImage struct {
	pix       []uint8
	pixOffset func(x, y int) int
	step      int
}

// pos returns the index in the byte array corresponding to the first byte of color c in the pixel at x, y.
func (s *stegoImage) pos(x, y, c int) int {
	return s.pixOffset(x, y) + c*s.step
}

//Reads the LSB of the color c in the pixel at x, y.
func (s *stegoImage) Read(x, y, c int) bool {
	p := s.pos(x, y, c)
	return s.pix[p]&0x01 == 1
}

// Enables the LSB (sets it to 1) of the color c in the pixel at x, y.
func (s *stegoImage) Enable(x, y, c int) {
	p := s.pos(x, y, c)
	s.pix[p] = s.pix[p] | 0x01 // LSB = 1
}

// Disables the LSB (sets it to 0) of the color c in the pixel at x, y.
func (s *stegoImage) Disable(x, y, c int) {
	p := s.pos(x, y, c)
	s.pix[p] = s.pix[p] & 0xFE // LSB = 0
}

func newStegoImage(img image.Image) (*stegoImage, error) {
	switch i := img.(type) {
	case *image.RGBA:
		return &stegoImage{i.Pix, i.PixOffset, 1}, nil
	case *image.RGBA64:
		return &stegoImage{i.Pix, i.PixOffset, 2}, nil
	case *image.NRGBA:
		return &stegoImage{i.Pix, i.PixOffset, 1}, nil
	case *image.NRGBA64:
		return &stegoImage{i.Pix, i.PixOffset, 2}, nil
	}
	return nil, StegoError("Image format not supported")
}

// Scrambler is used to distribute data in an image, but an implementation should not alter the
// image itself.
//
// The Init method should put the scrambler in a state from which Next should always generate the
// same sequence of results given the same image.
//
// Next should return whichever pixel and colour channel that should be modified or read next.
// Calling next should always result in unique color coordinates (x, y and c), as long as it's not
// called more times than the capacity returned by Cap.
//
// Cap should return the maximum number of times Next can be guaranteed to produce unique colour
// coordinates. If Init hasn't been called before Cap, Cap should return 0.
type Scrambler interface {
	Init(image image.Image)
	Next() (x, y, c int)
	Cap() int
}

// Encode encodes a slice of Bytes in the Image using the Scrambler. If it fails for some reason,
// an error is returned.
func Encode(bytes []byte, img image.Image, scrambler Scrambler) (image.Image, error) {

	w, err := newStegoImage(img)
	if err != nil {
		return nil, err
	}
	scrambler.Init(img)

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
		writeByte(b, w, scrambler)
	}

	return img, nil
}

// Decode reads bytes from a Image using a Scrambler. If it fails for some reason, an error is
// returned
func Decode(img image.Image, scrambler Scrambler) ([]byte, error) {
	w, err := newStegoImage(img)
	if err != nil {
		return nil, err
	}
	scrambler.Init(img)

	// Read the header...
	header := make([]byte, headerSize)
	for i := 0; i < len(header); i++ {
		header[i] = readByte(w, scrambler)

	}

	// ...abort if it's not formatted correct.
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
		ret[i] = readByte(w, scrambler)
	}
	return ret, nil
}

// readByte reads a byte from a image using a scrambler.
func readByte(w *stegoImage, scrambler Scrambler) (b byte) {
	var i uint8
	for i = 0; i < 8; i++ {
		x, y, c := scrambler.Next() // Next pixel coords and channel.
		if w.Read(x, y, c) {        // Read current bit.
			b += (1 << i) // Increase byte if necessary.
		}
	}
	return
}

// writeByte writes a byte to a image using a scrambler.
func writeByte(b byte, w *stegoImage, scrambler Scrambler) {
	var i uint8
	for i = 0; i < 8; i++ { // Step through each bit
		x, y, c := scrambler.Next()
		if (b>>i)&0x01 == 1 { // Set the bit.
			w.Enable(x, y, c)
		} else {
			w.Disable(x, y, c)
		}
	}
}
