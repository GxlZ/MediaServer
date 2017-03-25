package utils

import "os"

//创建目录
func MakeDir(path string, mode os.FileMode) error {

	//目录已存在
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	//目录不存在
	if err := os.Mkdir(path, mode); err != nil {
		return err
	}

	return nil
}
