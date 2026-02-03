package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type FileInfo struct {
	Name        string
	AbsFilePath string
	Size        int64
	ModTime     time.Time
}

type RemoteFileInfo struct {
	Size         int64
	LastModified string
}

type ImageSpecInfo struct {
	Width  int
	Height int
}

type ImageInfo struct {
	FileInfo
	ImageSpecInfo
}

type RemoteImageInfo struct {
	RemoteFileInfo
	ImageSpecInfo
}

func GetImageSpecInfo(imagePath string) (ImageSpecInfo, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return ImageSpecInfo{}, err
	}

	defer file.Close()

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return ImageSpecInfo{}, err
	}

	return ImageSpecInfo{config.Width, config.Height}, nil
}

func GetFileInfo(filePath string) (FileInfo, error) {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return FileInfo{}, err
	}

	fileStat, err := os.Stat(filePath)
	if err != nil {
		return FileInfo{}, err
	}

	name := fileStat.Name()
	size := fileStat.Size()
	modTime := fileStat.ModTime()

	return FileInfo{name, absFilePath, size, modTime}, nil
}

func DetectFileType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}
	fileType := http.DetectContentType(buffer[:n])

	return fileType, nil
}

func DetectUrlFileType(rawUrl string) (string, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	format := strings.ToLower(path.Ext(u.Path))[1:]

	switch format {
	case "jpg", "jpeg", "png":
		return "image/" + format, nil
	case "mp4", "avi", "mov", "mkv", "webm", "wmv", "flv":
		return "video/" + format, nil
	}

	return "", fmt.Errorf("Not supported. Only image and video http urls are supported")
}

func GetImageFromURL(url string, includeInfo bool) (image.Image, string, RemoteImageInfo, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, "", RemoteImageInfo{}, fmt.Errorf("http.Get failed: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, "", RemoteImageInfo{}, fmt.Errorf("bad HTTP status: %s", response.Status)
	}

	img, format, err := image.Decode(response.Body)
	if err != nil {
		return nil, "", RemoteImageInfo{}, fmt.Errorf("image.Decode failed: %w", err)
	}

	if !includeInfo {
		return img, format, RemoteImageInfo{}, nil
	}

	size := response.ContentLength
	lastModified := response.Header.Get("Last-Modified")
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	return img, format, RemoteImageInfo{RemoteFileInfo{size, lastModified}, ImageSpecInfo{width, height}}, nil
}

func GetVideoInfoFromUrl(url string) (RemoteFileInfo, error) {
	resp, err := http.Head(url)
	if err != nil {
		return RemoteFileInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Bad status:", resp.Status)
		return RemoteFileInfo{}, err
	}

	return RemoteFileInfo{resp.ContentLength, resp.Header.Get("Last-Modified")}, nil
}
