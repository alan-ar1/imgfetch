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

func GetTmuxRgbSeq(img image.Image, size ImageTermSize) (string, error) {
	config := ImageProtocolConfig{
		Format:     24,
		ColorDepth: 3,
		UseTmux:    true,
	}
	return generateRemoteKittySequence(img, size, config)
}

func GetTmuxRgbaSeq(img image.Image, size ImageTermSize) (string, error) {
	config := ImageProtocolConfig{
		Format:     32,
		ColorDepth: 4,
		UseTmux:    true,
	}
	return generateRemoteKittySequence(img, size, config)
}

func GetUnicodeRgbSeq(img image.Image, size ImageTermSize) (string, error) {
	config := ImageProtocolConfig{
		Format:     24,
		ColorDepth: 3,
		UseTmux:    false,
	}
	return generateRemoteKittySequence(img, size, config)
}

func GetUnicodeRgbaSeq(img image.Image, size ImageTermSize) (string, error) {
	config := ImageProtocolConfig{
		Format:     32,
		ColorDepth: 4,
		UseTmux:    false,
	}
	return generateRemoteKittySequence(img, size, config)
}

func GetTmuxPngSeq(imagePath string, size ImageTermSize) (string, error) {
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

func GetTmuxSeq(imagePath string, size ImageTermSize) (string, error) {
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
		seq, err := GetTmuxPngSeq(imagePath, size)
		if err != nil {
			return "", err
		}
		return seq, err
	}

	file.Seek(0, 0)

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	seq, err := GetTmuxRgbSeq(img, size)
	if err != nil {
		return "", err
	}
	return seq, err
}

func GetUnicdoeSeq(imagePath string, size ImageTermSize) (string, error) {
	absPath, err := filepath.Abs(imagePath)
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
		enc := base64.StdEncoding.EncodeToString([]byte(absPath))
		id, rgb, maskIndex := generateKittyID()
		return fmt.Sprintf(
			"%s_Gf=100,t=f,a=T,c=%d,r=%d,U=1,i=%d,q=2;%s%s",
			esc,
			size.Columns,
			size.Rows,
			id,
			enc,
			st,
		) + encodeImageID(size, rgb, maskIndex), nil
	}
	file.Seek(0, 0)
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	seq, err := GetUnicodeRgbSeq(img, size)
	if err != nil {
		return "", err
	}
	return seq, nil

}

func GetSeq(imagePath string, size ImageTermSize) (string, error) {

	if os.Getenv("TMUX") != "" {
		seq, err := GetTmuxSeq(imagePath, size)
		if err != nil {
			return "", err
		}
		return seq, nil
	}
	seq, err := GetUnicdoeSeq(imagePath, size)
	if err != nil {
		return "", err
	}

	return seq, nil
}

func GetRemoteSeq(img image.Image, format string, size ImageTermSize) (string, error) {
	if os.Getenv("TMUX") != "" {
		if format == "png" {
			seq, err := GetTmuxRgbaSeq(img, size)
			if err != nil {
				return "", err
			}
			return seq, nil
		} else {
			seq, err := GetTmuxRgbSeq(img, size)
			if err != nil {
				return "", err
			}
			return seq, nil
		}
	} else {
		if format == "png" {
			seq, err := GetUnicodeRgbaSeq(img, size)
			if err != nil {
				return "", err
			}

			return seq, nil
		} else {
			seq, err := GetUnicodeRgbSeq(img, size)
			if err != nil {
				return "", err
			}

			return seq, nil
		}
	}
}
