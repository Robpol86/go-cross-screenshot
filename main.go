package main

import (
	"fmt"
	"image/png"
	"os"
	"time"

	"github.com/kbinani/screenshot"
)

func run() {
	start := time.Now()
	bounds := screenshot.GetDisplayBounds(0)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}

	prefix := os.Args[1]
	now := time.Now()
	name := prefix + now.Format("2006-01-02-15-04-05")

	fileName := name + ".png"
	file, _ := os.Create(fileName)
	defer file.Close()
	png.Encode(file, img)

	fmt.Println(fileName)

	fmt.Println(time.Since(start))
}

func main() {
	loopCount := 1
	
	if len(os.Args) > 2 {
		if i, err := strconv.Atioi(os.Args[2]); err == nil {
			fmt.Printf("Looping %d time(s)\n", i)
			loopCount = i
		}
	}
	
	for i := 0; i < loopCount; i++ {
		run()
	}
}
