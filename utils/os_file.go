package utils

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
)

// AppendFilesToZipWriter 添加文件到ZipWriter
func AppendFilesToZipWriter(filename string, zipWriter *zip.Writer) error {
	file, err := os.Open(filename)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to open %s: %s", filename, err))
	}
	defer file.Close()

	info, _ := file.Stat()                  // 获取文件信息
	header, err := zip.FileInfoHeader(info) // 创建ZIP文件中的文件头
	header.Name = filename                  // 设置ZIP文件中的文件名
	header.Method = zip.Deflate

	wr, err := zipWriter.CreateHeader(header) // 将文件头写入ZIP文件

	//wr, err := zipWriter.Create(filename)
	//if err != nil {
	//	return errors.New(fmt.Sprintf("Failed to create entry for %s in zip file: %s", filename, err))
	//}

	// 将文件内容复制到ZIP文件中
	if _, err = io.Copy(wr, file); err != nil {
		return errors.New(fmt.Sprintf("Failed to write %s to zip: %s", filename, err))
	}

	return nil
}

// MakeDir 创建文件夹
func MakeDir(path string) error {
	return os.Mkdir(path, 0755)
}

// RemoveDir 删除文件夹
func RemoveDir(path string) error {
	return os.RemoveAll(path)
}

// FileOrDirExists 文件/文件夹是否存在
func FileOrDirExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
