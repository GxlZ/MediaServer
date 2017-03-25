package bal

import (
	"io"
	"MediaServer/videoer"
	"log"
	"errors"
	"MediaServer/utils"
)

type VideoUploader struct {
	Base
}

func (this *VideoUploader) Upload(r io.Reader) (video *videoer.Videoer, err error) {
	//保存视频到临时目录
	video, err = videoer.SaveToTmp(r)
	if err != nil {
		log.Println(err)
		return video, errors.New("save video to tmp faild.")
	}

	//生成正式视频保存目录
	if err := utils.MakeDir(video.VideoDir, 0755); err != nil {
		log.Println(err)
		return video, err
	}

	if err := video.GenCoverImageFromTmp(); err != nil {
		log.Println(err)
		return video, errors.New("generate cover image faild.")
	}

	//保存视频到临时文件夹
	video.GenCoverImageFromTmp()

	return video, err
}