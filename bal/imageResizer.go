package bal

import (
	"github.com/GxlZ/imaging"
	"fmt"
	"os"
	"MediaServer/imager"
	"errors"
	"io"
	"log"
)

type ImageResizer struct {
	Base
}

func (this *ImageResizer) Resize(sourceFileName string, sourceFilePath string, resizeInfos ...imager.Info) (infos map[string]imager.Info, err error) {
	for _, resizeInfo := range resizeInfos {
		go this.resize(resizeInfo)
	}

	return infos, err
}

func (this *ImageResizer)  resize(resizeInfo imager.Info) (err error) {
	f, err := os.Open(resizeInfo.ImagePath)
	if err != nil {
		log.Println(err)
		return err
	}

	defer f.Close()

	imger, err := imager.New(f)
	if err != nil {
		log.Println(err)
		return err
	}

	targetImgName := fmt.Sprintf("%s_%dx%d", resizeInfo.ImageName, resizeInfo.Width, resizeInfo.Height)
	targetImgPath := fmt.Sprintf("%s/%s", resizeInfo.ImageDir, targetImgName)

	switch imger.ImageType() {
	case imager.IMAGE_TYPE_PNG:
		fallthrough
	case imager.IMAGE_TYPE_BMP:
		fallthrough
	case imager.IMAGE_TYPE_JPG:
		img, err := imger.Decode(f)
		if err != nil {
			log.Println(err)
			return err
		}

		targetImg := imaging.Resize(img, resizeInfo.Width, resizeInfo.Height, imaging.Lanczos)
		if err := imaging.Save(targetImg, "jpg", targetImgPath); err != nil {
			return err
		}
	case imager.IMAGE_TYPE_GIF:
		targetImg, err := os.Create(targetImgPath)
		if err != nil {
			return err
		}
		defer targetImg.Close()

		if _, err := io.Copy(targetImg, f); err != nil {
			return err
		}
	default:
		return errors.New("unsupport image type")
	}

	return err
}