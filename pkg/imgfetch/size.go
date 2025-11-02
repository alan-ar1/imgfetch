package imgfetch

import (
	"image"
	"os"

	"golang.org/x/sys/unix"
)

func CalculateImageTermSize(imgWidth, imgHeight, widthFactor int) (ImageTermSize, error) {
	f, err := os.OpenFile("/dev/tty", unix.O_NOCTTY|unix.O_CLOEXEC|unix.O_NDELAY|unix.O_RDWR, 0666)
	if err != nil {
		return ImageTermSize{}, err
	}
	defer f.Close()

	sz, err := unix.IoctlGetWinsize(int(f.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return ImageTermSize{}, err
	}

	pixelsPerCol := float32(sz.Xpixel) / float32(sz.Col)
	pixelsPerRow := float32(sz.Ypixel) / float32(sz.Row)

	imgCols := float32(imgWidth) / pixelsPerCol
	imgRows := float32(imgHeight) / pixelsPerRow

	newImgCols := int(sz.Col) / widthFactor
	newImgRows := int(float32(newImgCols) * imgRows / imgCols)

	return ImageTermSize{Columns: newImgCols, Rows: newImgRows}, nil
}

func GetImageTermSize(imagePath string, widthFactor int) (ImageTermSize, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return ImageTermSize{}, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return ImageTermSize{}, err
	}

	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	return CalculateImageTermSize(w, h, widthFactor)
}
