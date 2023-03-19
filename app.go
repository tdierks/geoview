package main

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"os"
	"time"
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		480, 480, sdl.WINDOW_SHOWN | sdl.WINDOW_FULLSCREEN | sdl.WINDOW_BORDERLESS)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)
  
  images := make([]*sdl.Surface,0)
  for i := 1 ; i < len(os.Args) ; i++ {
    rwJpeg := sdl.RWFromFile(os.Args[i], "r")
    if rwJpeg == nil {
      panic(sdl.GetError())
    }
    defer rwJpeg.Close()
    
    jpeg, err := img.LoadJPGRW(rwJpeg)
    if err != nil {
      panic(err)
    }
    defer jpeg.Free()
    
    scaled, err := sdl.CreateRGBSurfaceWithFormat(0, surface.W, surface.H, 
      int32(surface.Format.BitsPerPixel), surface.Format.Format)
    if err != nil {
      panic(err)
    }

    err = jpeg.BlitScaled(nil, scaled, nil)
    if err != nil {
      panic(err)
    }
    
    images = append(images, scaled)
  }
  
  frames := 0
  t0 := time.Now()
    
//   var i int32 = 0
  j := 0
  
	running := true
	for running {
	  images[j].Blit(nil, surface, nil)
	  j = j+1
	  if j >= len(images) {
	    j = 0
	  }

    window.UpdateSurface()
    
    frames++
    t1 := time.Now()
    elapsed := t1.Sub(t0)
    if elapsed.Milliseconds() > 5000 {
      fmt.Printf("%d frames in %d ms, %.1f/sec\n", frames, elapsed.Milliseconds(), float64(frames)/elapsed.Seconds())
      t0 = t1
      frames = 0
    }

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
			fmt.Println("Got event", event)
		}
	}
}