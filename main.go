package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net"
	"os"
	"strconv"

	"github.com/mitchellh/go-vnc"
)

func main() {
	vncPortStr := os.Args[1]
	vncPort, _ := strconv.Atoi(vncPortStr)

	nc, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", vncPort))
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
	f, _ := os.OpenFile("out.png", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, img)
}
