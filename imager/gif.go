package imager

import (
	"io"
	"image"
	"image/gif"
	"github.com/GxlZ/imaging"
)

type GIF struct {
	image.Image
}

func (this *GIF) Decode(r io.Reader) (img image.Image, err error) {
	return gif.Decode(r)
}

func (this *GIF) DecodeAll(r io.Reader) (img *gif.GIF, err error) {
	return gif.DecodeAll(r)
}

func (this *GIF) Encode(w io.Writer, img image.Image, quality int) (err error) {
	return gif.Encode(w, img, &gif.Options{})
}

func (this *GIF) EncodeAll(w io.Writer, img *gif.GIF) (err error) {
	return gif.EncodeAll(w, img)
}

func (this *GIF) Resize(w io.Writer, img image.Image, width, height int, filter imaging.ResampleFilter) (err error) {
	return nil
}

func (this *GIF) ImageType() (imageType string) {
	return IMAGE_TYPE_GIF;
}

func (this *GIF) CoverImage(r io.Reader, w io.Writer, quality int) (err error) {
	img, err := this.Decode(r)
	if err != nil {
		return err
	}
	if err := this.Encode(w, img, quality); err != nil {
		return err
	}
	return nil
}

func (this *GIF) Config() (info Info, err error) {
	return info, err
}
