package imager

import (
	"io"
	"image"
	"image/jpeg"
	"image/png"
	"github.com/GxlZ/imaging"
)

type PNG struct {
	image.Image
}

func (this *PNG) Decode(r io.Reader) (img image.Image, err error) {
	return png.Decode(r)
}

func (this *PNG) Encode(w io.Writer, img image.Image, quality int) (err error) {
	return jpeg.Encode(w, img, &jpeg.Options{quality})
}

func (this *PNG) Resize(w io.Writer, img image.Image, width, height int, filter imaging.ResampleFilter) (err error) {
	return nil
}

func (this *PNG) ImageType() (imageType string) {
	return IMAGE_TYPE_PNG;
}

func (this *PNG) Config() (info Info, err error) {
	return info, err
}
