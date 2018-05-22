package fitsutils

import (
	"errors"
	"image"
	"image/png"
	"io"

	"github.com/astrogo/fitsio"
)

var (
	ErrNoImageInFile = errors.New("no image in fits file")
)

func SavePNG(source io.Reader, destination io.Writer) error {
	fitsFile, err := fitsio.Open(source)
	if err != nil {
		return err
	}

	var img image.Image

	hdus := fitsFile.HDUs()

	for _, hdu := range hdus {
		switch v := hdu.(type) {
		case fitsio.Image:
			img = v.Image()
		}
	}

	if img == nil {
		return ErrNoImageInFile
	}

	return png.Encode(destination, img)
}
