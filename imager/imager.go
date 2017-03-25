package imager

import (
	"image"
	"io"
	"fmt"
	"os"
	"github.com/GxlZ/imaging"
)

const (
	IMAGE_TYPE_JPG = "jpg"
	IMAGE_TYPE_PNG = "png"
	IMAGE_TYPE_GIF = "gif"
	IMAGE_TYPE_BMP = "bmp"
)

type Info struct {
	ImageName     string
	ImagePath     string
	ImageDir      string
	ImageType     string
	DirHashNumber int
	Width         int
	Height        int
	Quality       int
	Size          int64
}

type Imager interface {
	Decode(r io.Reader) (img image.Image, err error)
	Encode(w io.Writer, img image.Image, quality int) (err error)
	Resize(w io.Writer, img image.Image, width, height int, filter imaging.ResampleFilter) (err error)
	ImageType() (imageType string)
}

//判断类型 创建对应实例
func New(r io.Reader) (imager Imager, err error) {
	imageType, err := GetImageType(r)
	switch imageType {
	case IMAGE_TYPE_JPG:
		jpg := &JPG{}
		return jpg, nil
	case IMAGE_TYPE_PNG:
		png := &PNG{}
		return png, nil
	case IMAGE_TYPE_GIF:
		gif := &GIF{}
		return gif, nil
	case IMAGE_TYPE_BMP:
	default:
		return nil, fmt.Errorf("unsupport image type.")
	}
	return imager, err
}

//获取文件类型
func GetImageType(f io.Reader) (imageType string, err error) {
	bytes := make([]byte, 4)
	n, err := f.Read(bytes)
	if (err != nil) {
		return "", err
	}
	defer f.(io.ReadSeeker).Seek(0, os.SEEK_SET)

	if n < 4 {
		return "", fmt.Errorf("image read faild.")
	}
	if bytes[0] == 0x89 && bytes[1] == 0x50 && bytes[2] == 0x4E && bytes[3] == 0x47 {
		return IMAGE_TYPE_PNG, nil
	}
	if bytes[0] == 0xFF && bytes[1] == 0xD8 {
		return IMAGE_TYPE_JPG, nil
	}
	if bytes[0] == 0x47 && bytes[1] == 0x49 && bytes[2] == 0x46 && bytes[3] == 0x38 {
		return IMAGE_TYPE_GIF, nil
	}
	if bytes[0] == 0x42 && bytes[1] == 0x4D {
		return IMAGE_TYPE_BMP, nil
	}
	return "", fmt.Errorf("unknow image type.")
}