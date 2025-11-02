package main

import (
	"fmt"
	imgfetch "github.com/alan-ar1/imgfetch/pkg/imgfetch"
)

func main() {
	seq, err := imgfetch.GetNativeImageSeq("../imgs/dice.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(seq)
}
