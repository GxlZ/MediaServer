package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"log"
	"MediaServer/errorCode"
	"MediaServer/utils"
	"MediaServer/bal"
	"fmt"
)

type VideoUploadController struct {
	BaseController
}

//二进制文件上传
func (this *VideoUploadController) Binary() gin.HandlerFunc {
	return func(c *gin.Context) {
		f, _, err := c.Request.FormFile("video")
		if err != nil {
			log.Println(err)
			c.JSON(200, gin.H{
				"code":errorCode.VIDEO_SOURCE_SAVE_FAILD.No,
				"msg":errorCode.VIDEO_SOURCE_SAVE_FAILD.Msg,
			})
			return
		}
		defer f.Close()

		size, err := utils.FileSize(f);
		if size <= 0 || err != nil {
			log.Println(err)
			c.JSON(200, gin.H{
				"code":errorCode.VIDEO_SIZE_ERROR.No,
				"msg":errorCode.VIDEO_SIZE_ERROR.Msg,
			})
			return
		}

		//保存原视频
		uploadBal := bal.VideoUploader{}
		video, err := uploadBal.Upload(f)

		//压缩视频

		host := c.Request.Host
		resizeBal := bal.VideoResizer{}
		go resizeBal.Resize(video)

		c.JSON(200, gin.H{
			"code":errorCode.SUCCESS.No,
			"msg":errorCode.SUCCESS.Msg,
			"data":map[string]interface{}{
				"videoUrl":fmt.Sprintf("http://%s/video/static/%d/%s", host, video.DirHashNumber, video.VideoName),
				"videoUrl_cover":fmt.Sprintf("http://%s/video/static/%d/%s", host, video.DirHashNumber, video.CoverImageName),
				"videoDuration":video.Duration,
			},
		})

		c.Next()

	}
}