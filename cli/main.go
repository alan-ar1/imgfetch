package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alan-ar1/imgfetch/pkg/imgfetch"
	"github.com/charmbracelet/lipgloss"
)

func main() {

	if len(os.Args) <= 1 {
		fmt.Println("Provide an image path")
		return
	}

	filePath := os.Args[1]

	imageSeq := ""

	info, file, err := GetFileInfo(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("4"))

	infoStr := fmt.Sprintf("%s %s\n%s %s\n%s %s\n%s %d Bytes\n%s %s\n",
		labelStyle.Render("Name:"), info.Name,
		labelStyle.Render("Path:"), info.AbsFilePath,
		labelStyle.Render("Type:"), info.FileType,
		labelStyle.Render("Size:"), info.Size,
		labelStyle.Render("Modified:"), info.ModTime.Format(time.DateTime))

	if strings.HasPrefix(info.FileType, "image") {
		imageInfo, err := GetImageInfo(file, info)
		if err != nil {
			fmt.Println(err)
			return
		}

		imageSeq, err = imgfetch.GetImageSeq(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}

		infoStr += fmt.Sprintf("%s %dx%d", labelStyle.Render("Dimensions:"), imageInfo.Width, imageInfo.Height)
	}
	image := lipgloss.NewStyle().
		Padding(1).
		Render(imageSeq)

	text := lipgloss.NewStyle().
		Padding(1).
		Render(infoStr)

	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, image, text))
}
