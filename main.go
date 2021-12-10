package main

import (
	"image/png"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/kbinani/screenshot"
	"github.com/sevlyar/go-daemon"
)

func ss() {
	start := time.Now()
	bounds := screenshot.GetDisplayBounds(0)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}

	prefix := os.Args[1]
	now := time.Now()
	name := prefix + now.Format("_2006-01-02-15-04-05")

	fileName := name + ".png"
	file, _ := os.Create(fileName)
	defer file.Close()
	png.Encode(file, img)

	log.Println(fileName)

	log.Println(time.Since(start))
}

func run() {
	loopCount := 1
	
	if len(os.Args) > 2 {
		if i, err := strconv.Atoi(os.Args[2]); err == nil {
			log.Printf("Looping %d time(s)\n", i)
			loopCount = i
		}
	}
	
	for i := 1; i <= loopCount; i++ {
		ss()
		if i != loopCount {
			time.Sleep(2 * time.Second)
		}
	}
	
	log.Println("Done")
}

func main() {
	if os.Args[len(os.Args)-1] != "-d" {
		run()
		return
	}
	
	// Daemon
	cntxt := &daemon.Context{
		PidFileName: "ss.pid",
		PidFilePerm: 0644,
		LogFileName: "ss.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Print("- - - - - - - - - - - - - - -")
	log.Print("daemon started")
	
	run()
}
