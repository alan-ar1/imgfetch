package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"os/exec"
)

func getVideoThumbnail(videoPath string) (image.Image, error) {
	cmd := exec.Command("ffmpeg", "-ss", "00:00:00.001", "-i", videoPath, "-vframes", "1", "-f", "image2pipe", "-vcodec", "mjpeg", "pipe:1")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return nil, err

	}

	img, err := jpeg.Decode(&stdout)
	if err != nil {
		return nil, err
	}

	return img, nil
}
