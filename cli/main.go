package main

import (
	"fmt"
	"time"

	imgfetch "github.com/alan-ar1/imgfetch/pkg/imgfetch"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	imagePath := "../imgs/dice.png"

	seq, err := imgfetch.GetNativeImageSeq(imagePath)

	if err != nil {
		fmt.Println(err)
		return
	}

	info, err := getImageInfo(imagePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	infoStr := fmt.Sprintf("Name: %s\nPath: %s\nType: %s\nSize: %d Bytes\nModified: %s\nDimension: %dx%d", info.name, info.absFilePath, info.fileType, info.size, info.modTime.Format(time.DateTime), info.width, info.height)

	image := lipgloss.NewStyle().
		Padding(1).
		Render(seq)

	text := lipgloss.NewStyle().
		Padding(1).
		Render(infoStr)

	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, image, text))

}
