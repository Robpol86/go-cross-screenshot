package main

import (
	"image/png"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/cretz/go-scrap"
	"github.com/sevlyar/go-daemon"
)

func ss() error {
	start := time.Now()

	// Take screenshot.
	if err := scrap.MakeDPIAware(); err != nil {
		panic(err)
	}
	d, err := scrap.PrimaryDisplay()
	if err != nil {
		panic(err)
	}
	c, err := scrap.NewCapturer(d)
	if err != nil {
		panic(err)
	}
	for {
		if img, _, err := c.FrameImage(); img != nil || err != nil {
			// Detach the image so it's safe to use after this method
			if img != nil {
				img.Detach()
			}
			if err != nil {
				panic(err)
			}
			break
		}
		// Sleep 17ms (~1/60th of a second)
		time.Sleep(17 * time.Millisecond)
	}

	// Determine file name.
	prefix := os.Args[1]
	now := time.Now()
	name := prefix + now.Format("_2006-01-02-15-04-05")
	fileName := name + ".png"
	
	// Write screen shot to file. Return here if file already exists.
	file, err := os.OpenFile(fileName, os.O_RDWR | os.O_CREATE | os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	png.Encode(file, img)

	// Done.
	log.Println(fileName)
	log.Println(time.Since(start))
	return file.Sync()
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
		// Run.
		for {
			err := ss()
			if err == nil {
				break
			}
			log.Println("File name with same timestamp exists, sleeping and retrying")
			time.Sleep(1 * time.Second)
		}
		// Done with iteration. Sleep if not the last one.
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
