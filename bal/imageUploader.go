package bal

import (
	"io"
	"os"
	"MediaServer/imager"
	"MediaServer/utils"
	"MediaServer/config"
	"errors"
	"fmt"
)

type ImageUploader struct {
	Base
}

func (this *ImageUploader) Upload(r io.Reader) (info imager.Info, err error) {
	imger, err := imager.New(r)
	if err != nil {
		return info, err
	}

	var imgName, imgPath, imgDir string
	var dirHashNumber int

	switch imger.ImageType() {
	case imager.IMAGE_TYPE_GIF:
		gif, ok := imger.(*imager.GIF)
		if !ok {
			return info, errors.New("load gif faild")
		}
		img, err := gif.DecodeAll(r)
		if err != nil {
			return info, err
		}
		imgName = utils.RandString("." + imager.IMAGE_TYPE_GIF)

		dirHashNumber, err = utils.GenHashNumber(imgName, 3)
		if err != nil {
			return info, err
		}

		imgDir = fmt.Sprintf("%s/%d", config.ImagePath, dirHashNumber)
		if err := utils.MakeDir(imgDir, 0755); err != nil {
			return info, err
		}
		imgPath = fmt.Sprintf("%s/%s", imgDir, imgName)

		targetFile, err := os.Create(imgPath)
		if err != nil {
			return info, err
		}
		gif.EncodeAll(targetFile, img)

		defer targetFile.Close()

		if err != nil {
			return info, err
		}
		info.ImageType = imager.IMAGE_TYPE_GIF

		go GifCoverImg(r, imgDir, imgName)

	case imager.IMAGE_TYPE_JPG: //非GIF暂时一律转换成jpg
		fallthrough
	case imager.IMAGE_TYPE_BMP:
		fallthrough
	case imager.IMAGE_TYPE_PNG:
		img, err := imger.Decode(r)
		if err != nil {
			return info, err
		}
		imgName = utils.RandString("." + imager.IMAGE_TYPE_JPG)
		dirHashNumber, err = utils.GenHashNumber(imgName, 3)
		if err != nil {
			return info, err
		}

		imgDir = fmt.Sprintf("%s/%d", config.ImagePath, dirHashNumber)
		if err := utils.MakeDir(imgDir, 0755); err != nil {
			return info, err
		}
		imgPath = fmt.Sprintf("%s/%s", imgDir, imgName)

		targetFile, err := os.Create(imgPath)
		if err != nil {
			return info, err
		}
		imger.Encode(targetFile, img, 75)
		defer targetFile.Close()

		if err != nil {
			return info, err
		}

		info.ImageType = imager.IMAGE_TYPE_JPG

	default:
		return info, errors.New("unsupport image type.")
	}

	info.ImageName = imgName
	info.ImagePath = imgPath
	info.ImageDir = imgDir
	info.DirHashNumber = dirHashNumber

	return info, err
}

func GifCoverImg(r io.Reader, imgDir string, imgName string) error {
	r.(io.Seeker).Seek(0, os.SEEK_SET)
	gif := imager.GIF{}
	imgPath := fmt.Sprintf("%s/%s_cover", imgDir, imgName)
	targetFile, err := os.Create(imgPath)
	if err != nil {
		return err
	}
	if err := gif.CoverImage(r, targetFile, 100); err != nil {
		return err
	}
	return nil
}