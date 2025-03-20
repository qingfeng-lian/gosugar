package sugar

import (
	"os"
	"path/filepath"
)

// CreateDirs 递归创建目录
func CreateDirs(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

// EnsureDirForFile 检查目录是否存在则创建，但是传入的是文件路径
func EnsureDirForFile(filename string) error {
	dir := filepath.Dir(filename) // 获取文件路径中的目录部分
	return CreateDirs(dir)        // 使用已定义的 CreateDirs 来确保目录存在
}

// CreateFileWithDirs 创建文件，如果目录不存在递归创建目录
func CreateFileWithDirs(filename string) (*os.File, error) {
	dir := filepath.Dir(filename)
	if err := CreateDirs(dir); err != nil {
		return nil, err
	}
	return os.Create(filename)
}

// AppendToFile 写入文件，文件不存在则创建，存在则追加
func AppendToFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

// WriteToFile 写入文件，文件不存在则创建，存在则覆盖
func WriteToFile(filename string, data []byte) (int, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(data)
}

// FileExists 检查文件是否存在
func FileExists(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return !info.IsDir(), nil
}

// DirExists 检查目录是否存在
func DirExists(dirname string) (bool, error) {
	info, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}
