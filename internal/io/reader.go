package io

import (
	"errors"
	"io"
	"os"

	"github.com/fedemengo/f2bist/internal/engine"
	"github.com/fedemengo/f2bist/internal/types"
)

func BitsFromStdin() ([]types.Bit, error) {
	return BitsFromReader(os.Stdin)
}

func BitsFromFile(filename string) func() ([]types.Bit, error) {
	return func() ([]types.Bit, error) {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		return BitsFromReader(f)
	}
}

func BitsFromReaderWithCap(r io.Reader, maxBits int) ([]types.Bit, error) {
	bits := []types.Bit{}

	bytes := make([]byte, 8)
	lastCount := -1
	for {
		n, err := r.Read(bytes)
		if errors.Is(err, io.EOF) {
			lastCount = n
			break
		}
		if err != nil {
			return nil, err
		}

		for _, b := range bytes[:n] {
			bitsArray := engine.ByteToBits(b)
			//fmt.Printf("%v `%c`\n", bitsArray, b)
			bits = append(bits, bitsArray[0:8]...)
		}

		if maxBits > 0 && len(bits) >= maxBits {
			lastCount = 0
			break
		}
	}

	if lastCount > 0 {
		for _, b := range bytes[:lastCount] {
			bitsArray := engine.ByteToBits(b)
			bits = append(bits, bitsArray[0:8]...)
		}
	}

	if maxBits > 0 && maxBits < len(bits) {
		bits = bits[:maxBits]
	}

	return bits, nil

}

func BitsFromReader(r io.Reader) ([]types.Bit, error) {
	return BitsFromReaderWithCap(r, -1)
}
