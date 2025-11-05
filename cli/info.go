package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileInfo struct {
	Name        string
	FileType    string
	AbsFilePath string
	Size        int64
	ModTime     time.Time
}
type ImageInfo struct {
	FileInfo
	Width  int
	Height int
}

func GetImageInfo(file *os.File, fileInfo FileInfo) (ImageInfo, error) {
	defer file.Close()

	if !strings.HasPrefix(fileInfo.FileType, "image") {
		return ImageInfo{}, fmt.Errorf("file is not an image")
	}

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return ImageInfo{}, err
	}

	return ImageInfo{fileInfo, config.Width, config.Height}, nil
}

func GetFileInfo(filePath string) (FileInfo, *os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return FileInfo{}, nil, err
	}

	fileType, err := detectFileType(file)
	if err != nil {
		return FileInfo{}, nil, err
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return FileInfo{}, nil, err
	}

	fileStat, err := os.Stat(filePath)
	if err != nil {
		return FileInfo{}, nil, err
	}

	name := fileStat.Name()
	size := fileStat.Size()
	modTime := fileStat.ModTime()

	file.Seek(0, 0)

	return FileInfo{name, fileType, absFilePath, size, modTime}, file, nil
}

func detectFileType(file *os.File) (string, error) {

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}
	fileType := http.DetectContentType(buffer[:n])

	return fileType, nil
}
