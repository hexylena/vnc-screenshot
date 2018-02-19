package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net"
	"os"

	"github.com/mitchellh/go-vnc"
	"github.com/urfave/cli"
)

var (
	version   string
	builddate string
)

func screenshot(address, filename string){
	nc, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("Error connecting to VNC: %s", err)
	}
	defer nc.Close()

	ch := make(chan vnc.ServerMessage)

	c, err := vnc.Client(nc, &vnc.ClientConfig{
		Exclusive:       false,
		ServerMessageCh: ch,
		ServerMessages:  []vnc.ServerMessage{new(vnc.FramebufferUpdateMessage)},
	})
	if err != nil {
		fmt.Printf("Error handshaking with VNC: %s", err)
	}
	defer c.Close()
	fmt.Printf("Connected to VNC desktop: %s [res:%dx%d]\n", c.DesktopName, c.FrameBufferWidth, c.FrameBufferHeight)

	// Move the mouse to wake it up?
	//c.PointerEvent(0, 0, 0)
	//c.PointerEvent(0, 1, 1)

	// Then send a buffer updat request!
	err = c.FramebufferUpdateRequest(false, 0, 0, c.FrameBufferWidth, c.FrameBufferHeight)

	if err != nil {
		fmt.Printf("Error handshaking with VNC: %s", err)
	}

	msg := <-ch

	rects := msg.(*vnc.FramebufferUpdateMessage).Rectangles
	fmt.Println()

	w := int(rects[0].Width)
	h := int(rects[0].Height)
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	enc := rects[0].Enc.(*vnc.RawEncoding)
	i := 0
	x := 0
	y := 0
	for _, v := range enc.Colors {
		x = i % w
		y = i / w
		r := uint8(v.R)
		g := uint8(v.G)
		b := uint8(v.B)

		img.Set(x, y, color.RGBA{r, g, b, 255})
		i++
	}

	// Save to out.png
	f, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, img)
}

func main() {
	app := cli.NewApp()
	app.Name = "vnc-screenshot"
	app.Usage = "Take screenshots of VNC servers from the command line."
	app.Version = fmt.Sprintf("%s (%s)", version, builddate)

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "serverAddr",
			Value:  "127.0.0.1:5000",
			Usage:  "Server Address",
			EnvVar: "VNC_SERVER_ADDR",
		},
		cli.StringFlag{
			Name:   "out",
			Value:  "out.png",
			Usage:  "Output File Name",
		},
	}
	app.Action = func(c *cli.Context) {
		screenshot(c.String("serverAddr"), c.String("out"))
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}

}
