package kitty

import (
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
)

func KittyUnicodeRGB(imagePath string, size ImageTermSize) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	imgPixels := make([]byte, w*h*3)

	pix := 0
	for i := range h {
		for j := range w {
			r, g, b, _ := img.At(j, i).RGBA()
			imgPixels[pix] = byte(r >> 8)
			imgPixels[pix+1] = byte(g >> 8)
			imgPixels[pix+2] = byte(b >> 8)
			pix += 3
		}
	}

	encPixelData := base64.StdEncoding.EncodeToString(imgPixels)
	id, rgb, maskIndex := generateKittyID()

	encPixelDataLength := len(encPixelData)
	chunkSize := 4096

	m := 1
	if encPixelDataLength <= chunkSize {
		m = 0
		chunkSize = encPixelDataLength
	}

	escapeCode, err := passthrough(fmt.Sprintf("%s_Gf=24,m=%d,a=T,U=1,q=2,i=%d,c=%d,r=%d,s=%d,v=%d;%s%s", esc, 1, id, size.Columns, size.Rows, w, h, encPixelData[0:4096], st))
	if err != nil {
		return "", err
	}

	chunkEnd := chunkSize * 2
	for i := chunkSize; i < encPixelDataLength; {
		if chunkEnd >= encPixelDataLength {
			chunkEnd = encPixelDataLength
			m = 0
		}
		pass, err := passthrough(fmt.Sprintf("%s_Gm=%d;%s%s", esc, m, encPixelData[i:chunkEnd], st))
		if err != nil {
			return "", nil
		}
		escapeCode += pass
		i += chunkSize
		chunkEnd += chunkSize
	}

	seq := escapeCode + encodeImageID(size, rgb, maskIndex)
	return seq, nil

}

func KittyUnicodePNG(imagePath string, size ImageTermSize) (string, error) {
	absPath, err := filepath.Abs(imagePath)
	if err != nil {
		return "", err
	}
	encPath := base64.StdEncoding.EncodeToString([]byte(absPath))
	id, rgb, maskIndex := generateKittyID()

	passedthroughSeq, err := passthrough(fmt.Sprintf("%s_Gf=100,t=f,a=T,U=1,q=2,i=%d,c=%d,r=%d;%s%s", esc, id, size.Columns, size.Rows, encPath, st))
	if err != nil {
		return "", err
	}
	seq := passedthroughSeq + encodeImageID(size, rgb, maskIndex)
	return seq, nil
}

func KittyUnicode(imagePath string, size ImageTermSize) (string, error) {
	cmd := exec.Command("tmux", "set", "-p", "allow-passthrough", "on")
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, format, err := image.DecodeConfig(file)
	if err != nil {
		return "", err
	}

	if format == "png" {
		seq, err := KittyUnicodePNG(imagePath, size)
		if err != nil {
			return "", err
		}
		return seq, err
	}
	seq, err := KittyUnicodeRGB(imagePath, size)
	if err != nil {
		return "", err
	}
	return seq, err
}

func KittyStandardRGB(imagePath string, size ImageTermSize) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	imgPixels := make([]byte, w*h*3)

	pix := 0
	for i := range h {
		for j := range w {
			r, g, b, _ := img.At(j, i).RGBA()
			imgPixels[pix] = byte(r >> 8)
			imgPixels[pix+1] = byte(g >> 8)
			imgPixels[pix+2] = byte(b >> 8)
			pix += 3
		}
	}

	encPixelData := base64.StdEncoding.EncodeToString(imgPixels)

	seq := ""
	encPixelDataLength := len(encPixelData)
	chunkSize := 4096

	m := 1
	if encPixelDataLength <= chunkSize {
		m = 0
		chunkSize = encPixelDataLength
	}

	seq += fmt.Sprintf("%s_Gf=24,a=T,c=%d,r=%d,s=%d,v=%d,m=%d;%s%s", esc, size.Columns, size.Rows, w, h, m, encPixelData[0:chunkSize], st)

	chunkEnd := chunkSize * 2
	for i := chunkSize; i < encPixelDataLength; {
		if chunkEnd >= encPixelDataLength {
			chunkEnd = encPixelDataLength
			m = 0
		}

		seq += fmt.Sprintf("%s_Gm=%d;%s%s", esc, m, encPixelData[i:chunkEnd], st)
		i += chunkSize
		chunkEnd += chunkSize
	}

	return seq, nil
}

func KittyStandard(imagePath string, size ImageTermSize) (string, error) {
	absPath, err := filepath.Abs(imagePath)

	if err != nil {
		return "", err
	}
	seq := ""

	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, format, err := image.DecodeConfig(file)
	if err != nil {
		return "", err
	}

	if format == "png" {
		enc := base64.StdEncoding.EncodeToString([]byte(absPath))
		seq = fmt.Sprintf(
			"%s_Gf=100,t=f,a=T,c=%d,r=%d;%s%s",
			esc,
			size.Columns,
			size.Rows,
			enc,
			st,
		)
	} else {
		str, err := KittyStandardRGB(imagePath, size)
		if err != nil {
			return "", err
		}
		seq = str
	}
	return seq, nil
}

func Kitty(imagePath string, size ImageTermSize) (string, error) {
	if os.Getenv("TMUX") != "" {
		seq, err := KittyUnicode(imagePath, size)
		if err != nil {
			return "", err
		}
		return seq, nil
	}
	seq, err := KittyStandard(imagePath, size)
	if err != nil {
		return "", err
	}
	return seq, nil
}
