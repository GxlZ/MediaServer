package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"MediaServer/routers"
	"flag"
	"MediaServer/config"
	"os"
	"fmt"
	"log"
	"MediaServer/utils"
)

func main() {
	initENV()

	initDir()

	log.SetFlags(log.LstdFlags | log.Llongfile)

	//注册路由
	router := gin.Default()
	routers.Register(router)

	//启动web 服务
	router.Run(":12345")
}

func initENV() {
	//环境变量
	env := os.Getenv("IMAGE_SERVER_ENV")
	if env == "" {
		env = "pro"
	}

	config.ENV = env
	fmt.Println("******** current mode:" + env + "********")

	switch env {
	case "dev":
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
}

func initDir() {
	//数据存放根目录
	dataPath := flag.String("d", "", "set image save path")
	flag.Parse()

	if *dataPath == "" {
		log.Fatalln("data save path invalid! use ` -d <data path> ` set data save path or use ` -h ` get help.")
	}
	_, err := os.Stat(*dataPath)
	if err != nil {
		log.Fatalf("data save path invalid! `%s`", err.Error())
	}

	//文件存放目录
	imagePath := fmt.Sprintf("%s/imgs", *dataPath)
	if err := utils.MakeDir(imagePath, 0755); err != nil {
		log.Fatalf("image save path invalid! `%s`", err.Error())
	}

	//视频存放目录
	videoPath := fmt.Sprintf("%s/videos", *dataPath)
	if err := utils.MakeDir(videoPath, 0755); err != nil {
		log.Fatalf("video save path invalid! `%s`", err.Error())
	}

	//临时目录
	tmpPath := fmt.Sprintf("%s/tmp", *dataPath)
	if err := utils.MakeDir(tmpPath, 0755); err != nil {
		log.Fatalf("tmp path invalid! `%s`", err.Error())
	}

	config.DataPath = *dataPath
	config.ImagePath = imagePath
	config.VideoPath = videoPath
	config.TmpPath = tmpPath
}
