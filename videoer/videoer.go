package videoer

import (
	"io"
	"fmt"
	"os"
	"MediaServer/utils"
	"MediaServer/config"
	"os/exec"
	"bytes"
	"errors"
	"log"
	"github.com/pquerna/ffjson/ffjson"
	"strconv"
	"strings"
	"runtime"
)

type Videoer struct {
	Height         int     //宽
	Width          int     //高
	Duration       float64 //视频时长
	CoverImagePath string  //封面图路径
	CoverImageName string  //封面图名称
	Size           int64   //大小
	VideoName      string  //视频文件名
	VideoPath      string  //视频全路径
	VideoDir       string  //视频保存目录
	VideoTmpPath   string  //视频临时保存目录
	DirHashNumber  int     //视频文件名 hash 值
	Ext            string  //扩展格式
}

var videoInfo map[string]interface{}

//加载已存在video
func Load(videoPath string) Videoer {
	return Videoer{
		VideoPath:videoPath,
	}
}

//保存video到临时文件夹
func SaveToTmp(r io.Reader) (videoer *Videoer, err error) {
	videoer = &Videoer{}

	videoName := utils.RandString(".mp4")
	videoer.VideoName = videoName

	//hash值
	hashNumber, _ := utils.GenHashNumber(videoName, 4)
	videoer.DirHashNumber = hashNumber

	//视频最终保存路径(压缩后保存路径)
	videoer.VideoDir = fmt.Sprintf("%s/%d", config.VideoPath, hashNumber)
	videoer.VideoPath = fmt.Sprintf("%s/%s", videoer.VideoDir, videoName)

	//视频临时保存路径
	videoer.VideoTmpPath = fmt.Sprintf("%s/%s", config.TmpPath, videoName)

	tmpFile, err := os.Create(videoer.VideoTmpPath)
	if err != nil {
		return videoer, err
	}
	defer tmpFile.Close()

	if err != nil {
		return videoer, err
	}

	if _, err := io.Copy(tmpFile, r); err != nil {
		return videoer, err
	}

	//获取原始视频信息
	videoer.SourceConfig()

	return videoer, nil
}

//移动&压缩video到正式文件夹
func (this *Videoer) CompressToFormal(width int, height int) error {
	//压缩命令
	tmpPath := fmt.Sprintf("%s/tmp_%s", this.VideoDir, this.VideoName)
	resolvingPower := fmt.Sprintf("%d*%d", width, height)
	cmd := exec.Command("ffmpeg", "-threads", strconv.Itoa(runtime.NumCPU()), "-i", this.VideoTmpPath, "-vcodec", "libx264", "-preset", "fast", "-crf", "32", "-y", "-acodec", "libmp3lame", "-ab", "96k", "-s", resolvingPower, tmpPath)

	var errInfo bytes.Buffer
	cmd.Stderr = &errInfo

	if err := cmd.Run(); err != nil {
		log.Println(errInfo.String())
		return errors.New("exec command faild.")
	}

	//压缩成功 重命名文件为正式名
	if _, err := os.Stat(tmpPath); err != nil {
		return errors.New("video compress faild.")
	}

	if err := os.Rename(tmpPath, this.VideoPath); err != nil {
		return errors.New("video rename faild.")
	}

	return nil
}



//从临时文件生成封面图
func (this *Videoer) GenCoverImageFromTmp() error {
	videoTmpPath := this.VideoTmpPath
	videoPath := this.VideoPath

	//封面图路径
	coverImagePath := fmt.Sprintf("%s_cover", videoPath)

	//制定封面图截取时间点
	beginTime := int(this.Duration / 2)

	//生成封面图
	cmd := exec.Command("ffmpeg", "-ss", strconv.Itoa(beginTime), "-i", videoTmpPath, "-y", "-f", "image2", "-vframes", "1", coverImagePath)

	var errInfo bytes.Buffer
	var jsonInfo bytes.Buffer
	cmd.Stdout = &jsonInfo
	cmd.Stderr = &errInfo

	if err := cmd.Run(); err != nil {
		log.Println(errInfo.String())
		return errors.New("exec command faild.")
	}

	this.CoverImagePath = coverImagePath
	coverImageName := fmt.Sprintf("%s_cover", this.VideoName)
	this.CoverImageName = coverImageName

	return nil
}

//获取视频信息
func (this *Videoer) SourceConfig() error {
	videoTmpPath := this.VideoTmpPath
	cmd := exec.Command("ffprobe", "-print_format", "json", "-show_format", "-show_streams", "-i", videoTmpPath)

	var errInfo bytes.Buffer
	var jsonInfo bytes.Buffer
	cmd.Stdout = &jsonInfo
	cmd.Stderr = &errInfo

	if err := cmd.Run(); err != nil {
		log.Println(errInfo.String())
		return errors.New("exec command faild.")
	}

	if err := ffjson.Unmarshal(jsonInfo.Bytes(), &videoInfo); err != nil {
		log.Println(err)
		return errors.New("parse video from json info faild.")
	}

	streams := videoInfo["streams"].([]interface{})
	height, width, err := this.parseStreams(streams)
	if err != nil {
		log.Println(err)
		return errors.New("parse streams faild")
	}

	this.Height = height
	this.Width = width

	formatInfo := videoInfo["format"].(map[string]interface{})

	//大小
	formatInfoSize, fund := formatInfo["size"]
	if !fund {
		log.Println("parse size info faild.")
		return errors.New("parse size info faild.")
	}
	size, err := strconv.ParseInt(formatInfoSize.(string), 10, 64)
	if err != nil {
		log.Println(err)
		return errors.New("format size to int64 faild.")
	}
	this.Size = size

	//时长
	formatInfoDuration, fund := formatInfo["duration"]
	if !fund {
		log.Println("parse duration info faild.")
		return errors.New("parse duration info faild.")
	}
	duration, err := strconv.ParseFloat(formatInfoDuration.(string), 64)
	if err != nil {
		log.Println("format duration to float faild.")
		return errors.New("format duration to float faild.")
	}
	this.Duration = duration

	//扩展格式
	formatInfoForamtName, fund := formatInfo["format_name"]
	if !fund {
		log.Println("format names string faild.")
		return errors.New("format names string faild.")
	}
	names := strings.Split(formatInfoForamtName.(string), ",")
	if len(names) < 1 {
		log.Println("format name info faild.")
		return errors.New("format name info faild.")
	}

	this.Ext = names[0]

	return nil
}

//解析streams 得到 宽,高 数据
func (this *Videoer)parseStreams(streams []interface{}) (height int, width int, err error) {
	for _, stream := range streams {
		streamHeight, err := stream.(map[string]interface{})["height"]
		if !err {
			continue
		}

		streamWidth, err := stream.(map[string]interface{})["width"]
		if !err {
			continue
		}

		height := int(streamHeight.(float64))
		width := int(streamWidth.(float64))

		return height, width, nil

	}

	return height, width, errors.New("parse streams height or width info faild.")
}