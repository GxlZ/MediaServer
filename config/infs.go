package config

//获取文件大小的接口
type Sizer interface {
	Size() int64
}