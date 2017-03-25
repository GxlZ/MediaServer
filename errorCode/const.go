package errorCode

type E struct {
	No  int
	Msg string
}

var (
	SUCCESS = E{0, "success"}
	//图片校验
	IMAGE_NOT_SUPPORT = E{1001, "image type not support."}
	IMAGE_OUT_OF_SIZE = E{1002, "image out of size."}
	IMAGE_UPLOAD_FAILD = E{1004, "image upload faild."}
	//图片处理
	IMAGE_SOURCE_SAVE_FAILD = E{1003, "source image save faild."}

	//视频校验
	VIDEO_SIZE_ERROR = E{2001, "source video size get faild."}

	//视频处理
	VIDEO_SOURCE_SAVE_FAILD = E{2003, "source video save faild."}

)