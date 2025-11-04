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

type fileInfo struct {
	name        string
	fileType    string
	absFilePath string
	size        int64
	modTime     time.Time
}
type imageInfo struct {
	fileInfo
	width  int
	height int
}

func getImageInfo(filePath string) (imageInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return imageInfo{}, err
	}
	defer file.Close()

	fileInfo, err := getFileInfo(file, filePath)
	if err != nil {
		return imageInfo{}, err
	}

	if !strings.HasPrefix(fileInfo.fileType, "image") {
		return imageInfo{}, fmt.Errorf("file is not an image")
	}

	config, _, err := image.DecodeConfig(file)

	return imageInfo{fileInfo, config.Width, config.Height}, err
}

func getFileInfo(file *os.File, filePath string) (fileInfo, error) {
	fileType, err := detectFileType(file)
	if err != nil {
		return fileInfo{}, err
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return fileInfo{}, err
	}

	fileStat, err := os.Stat(filePath)
	if err != nil {
		return fileInfo{}, err
	}

	name := fileStat.Name()
	size := fileStat.Size()
	modTime := fileStat.ModTime()
	file.Seek(0, 0)

	return fileInfo{name, fileType, absFilePath, size, modTime}, nil
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
