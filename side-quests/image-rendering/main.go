package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"strings"
	"time"

	"image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func renderImage(m image.Image) {
	fmt.Print("\033[H\033[2J")
	output := ""
	bounds := m.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			red, green, blue, _ := m.At(x, y).RGBA()
			output += fmt.Sprintf("\033[38;2;%d;%d;%dm█", red/255, green/255, blue/255)
			// fmt.Printf("\033[38;2;%d;%d;%dm█", red/255, green/255, blue/255)
		}
		output += "\n"
		// fmt.Println()
	}
	fmt.Print(output)

}

func main() {
	imageName := os.Args[1]
	reader, err := os.Open(imageName)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	if strings.HasSuffix(imageName, ".gif") {
		gifs, err := gif.DecodeAll(reader)
		if err != nil {
			log.Fatal(err)
		}
		images := gifs.Image
		for _, m := range images {
			renderImage(m)
			time.Sleep(time.Second / 10)
		}

	} else {

		var images []*image.Image
		imageObject, _, err := image.Decode(reader)
		if err != nil {
			log.Fatal(err)
		}
		images = append(images, &imageObject)

		for _, m := range images {
			renderImage(*m)
		}
	}
}
