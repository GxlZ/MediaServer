package bal

import (
	"MediaServer/videoer"
	"log"
)

type VideoResizer struct {
	Base
}

func (this *VideoResizer) Resize(video *videoer.Videoer) error {
	if err := video.CompressToFormal(640, 480); err != nil {
		log.Print("video compress faild.")
		return err
	}
	return nil
}