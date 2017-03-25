package utils

import (
	"fmt"
	"github.com/dchest/uniuri"
	"hash/crc32"
	"strconv"
	"errors"
	"os"
	"mime/multipart"
	"MediaServer/config"
)

//产生随机字符串
func RandString(suffix string) string {
	return fmt.Sprintf("%s%s", uniuri.New(), suffix)
}

//生成hash目录名
//@s hash源, @l 随机值长度
func GenHashNumber(s string, l int) (number int, err error) {
	if l > 9 {
		l = 9
	}
	if s == "" {
		return number, errors.New("imageName can not be empty.")
	}
	n := crc32.ChecksumIEEE([]byte(s));
	ns := strconv.FormatUint(uint64(n), 10)
	sLen := len(ns)
	partS := ns[sLen - l:sLen]
	res, err := strconv.Atoi(partS)
	if err != nil {
		return number, err
	}
	return res, err
}

//获取文件大小
func FileSize(f interface{}) (size int64, err error) {
	switch f.(type){
	case *os.File:
		info, err := f.(*os.File).Stat()
		if err != nil {
			return size, err
		}
		return info.Size(), err
	case multipart.File:
		size := f.(config.Sizer).Size()
		return size, err
	default:
		return size, errors.New("get file size faild.")
	}
}