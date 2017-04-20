package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"MediaServer/errorCode"
	"log"
	"fmt"
	"MediaServer/bal"
	"MediaServer/imager"
	"MediaServer/utils"
)

type ImageUploadController struct {
	BaseController
}

type ImageSize struct {
	width  int
	height int
}

const maxImageSize = 20 * 1 << 20  //文件最大20M

//二进制文件上传
func (this *ImageUploadController) Binary() gin.HandlerFunc {
	return func(c *gin.Context) {
		f, _, err := c.Request.FormFile("img")

		//客户端用于判断图片上传顺序的字段
		index := c.DefaultPostForm("index", "0")

		if err != nil {
			log.Println(err)
			c.JSON(200, gin.H{
				"code":errorCode.IMAGE_UPLOAD_FAILD.No,
				"msg":errorCode.IMAGE_UPLOAD_FAILD.Msg,
			})
			return
		}
		defer f.Close()

		size, err := utils.FileSize(f)
		if size <= 0 || err != nil {
			log.Println(err)
			c.JSON(200, gin.H{
				"code":errorCode.VIDEO_SIZE_ERROR.No,
				"msg":errorCode.VIDEO_SIZE_ERROR.Msg,
			})
			return
		}

		//图片大小验证
		if imageSize := size; imageSize > maxImageSize {
			errMsg := fmt.Sprintf("%s size: %d , maxSize : %d", errorCode.IMAGE_OUT_OF_SIZE.Msg, imageSize, maxImageSize)
			log.Println(errMsg)
			c.JSON(200, gin.H{
				"code":errorCode.IMAGE_OUT_OF_SIZE.No,
				"msg":errMsg,
			})
			return
		}

		//保存原尺寸图片
		uploadBal := bal.ImageUploader{}
		uploadInfo, err := uploadBal.Upload(f)

		if err != nil {
			log.Println(err)
			c.JSON(200, gin.H{
				"code":errorCode.IMAGE_SOURCE_SAVE_FAILD.No,
				"msg":errorCode.IMAGE_SOURCE_SAVE_FAILD.Msg,
			})
			return
		}

		//生成缩略图
		resizer := bal.ImageResizer{}
		resizeInfos := getResizeInfos(uploadInfo)
		resizer.Resize(uploadInfo.ImageName, uploadInfo.ImagePath, resizeInfos...)

		host := c.Request.Host

		imgInfos := map[string]string{
			"imageUrl":getImgUrlInfo(host, uploadInfo.DirHashNumber, uploadInfo.ImageName, 0, 0),
		}

		if uploadInfo.ImageType == imager.IMAGE_TYPE_GIF {
			coverImgName := fmt.Sprintf("%s_cover", uploadInfo.ImageName)
			imgInfos["imageUrl_cover"] = getImgUrlInfo(host, uploadInfo.DirHashNumber, coverImgName, 0, 0)
		}

		for _, resizeInfo := range resizeInfos {
			k := fmt.Sprintf("imageUrl_%dx%d", resizeInfo.Width, resizeInfo.Height)
			imgInfos[k] = getImgUrlInfo(host, uploadInfo.DirHashNumber, uploadInfo.ImageName, resizeInfo.Width, resizeInfo.Height)
		}

		c.JSON(200, gin.H{
			"code":errorCode.SUCCESS.No,
			"msg":errorCode.SUCCESS.Msg,
			"data":imgInfos,
			"index":index,
		})

		c.Next()
	}
}

func getImgUrlInfo(host string, dirHashNumber int, imageName string, width int, height int) (url string) {
	if width == 0 && height == 0 {
		url = fmt.Sprintf("http://%s/image/static/%d/%s", host, dirHashNumber, imageName)
	} else {
		url = fmt.Sprintf("http://%s/image/static/%d/%s_%dx%d", host, dirHashNumber, imageName, width, height)
	}
	return url
}

//压缩列表
func getResizeInfos(imgInfo imager.Info) (infos []imager.Info) {
	infos = []imager.Info{
		imager.Info{Width:360, Height:0, Quality:100},
		imager.Info{Width:720, Height:0, Quality:100},
	}

	for k, _ := range infos {
		infos[k].ImageDir = imgInfo.ImageDir
		infos[k].ImagePath = imgInfo.ImagePath
		infos[k].ImageName = imgInfo.ImageName
	}

	return infos
}