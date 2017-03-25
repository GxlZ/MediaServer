package routers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"MediaServer/controllers"
	"net/http"
	"MediaServer/config"
)

func Register(router *gin.Engine) {

	//图片相关
	imageRouter := router.Group("/image")
	{
		//二进制图片上传
		imageRouter.POST("/upload/binary", new(controllers.ImageUploadController).Binary())
		//图片获取
		imageRouter.StaticFS("/static", http.Dir(config.ImagePath))
	}

	videoRouter := router.Group("/video")
	{
		//视频上传
		videoRouter.POST("/upload/binary", new(controllers.VideoUploadController).Binary())
		//视频获取
		videoRouter.StaticFS("/static", http.Dir(config.VideoPath))
	}

	//404处理
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code":404,
			"msg":"api not found",
		})
	})

	//性能分析
	//ginpprof.Wrapper(router)
}