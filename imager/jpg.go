package imager

import (
	"io"
	"image"
	"image/jpeg"
	"github.com/GxlZ/imaging"
)

type JPG struct {
	image.Image
}

func (this *JPG) Decode(r io.Reader) (img image.Image, err error) {
	return jpeg.Decode(r)
}

func (this *JPG) Encode(w io.Writer, img image.Image, quality int) (err error) {
	return jpeg.Encode(w, img, &jpeg.Options{quality})
}

func (this *JPG) Resize(w io.Writer, img image.Image, width, height int, filter imaging.ResampleFilter) (err error) {
	return nil
}

func (this *JPG) ImageType() (imageType string) {
	return IMAGE_TYPE_JPG;
}

func (this *JPG) Config() (info Info, err error) {
	return info, err
}
