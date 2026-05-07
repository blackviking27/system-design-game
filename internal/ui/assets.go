package ui

import (
	"bytes"
	"image/color"
	"os"

	"github.com/blackviking27/system-design-game/internal/sim"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	NodeImages    = make(map[sim.NodeType]*ebiten.Image)
	NodeImagesRed = make(map[sim.NodeType]*ebiten.Image)
	FontSource    *text.GoTextFaceSource
)

func LoadAssets() error {
	// Loading asset image from assets

	// Load balancer
	imgLB, _, err := ebitenutil.NewImageFromFile("assets/images/network.png")
	if err != nil {
		return err
	}
	imgLBRed, _, err := ebitenutil.NewImageFromFile("assets/images/network_red.png")
	if err != nil {
		return err
	}

	// Server
	imgServer, _, err := ebitenutil.NewImageFromFile("assets/images/server.png")
	if err != nil {
		return err
	}
	imgServerRed, _, err := ebitenutil.NewImageFromFile("assets/images/server_red.png")
	if err != nil {
		return err
	}

	// DB
	imgDB, _, err := ebitenutil.NewImageFromFile("assets/images/db.png")
	if err != nil {
		return err
	}
	imgDBRed, _, err := ebitenutil.NewImageFromFile("assets/images/db_red.png")
	if err != nil {
		return err
	}

	// Queue
	imgQueue, _, err := ebitenutil.NewImageFromFile("assets/images/queue.png")
	if err != nil {
		return err
	}
	imgQueueRed, _, err := ebitenutil.NewImageFromFile("assets/images/queue_red.png")
	if err != nil {
		return err
	}

	// Cache
	imgCache, _, err := ebitenutil.NewImageFromFile("assets/images/cache.png")
	if err != nil {
		return err
	}

	NodeImages[sim.TypeLoadBalancer] = imgLB
	NodeImagesRed[sim.TypeLoadBalancer] = imgLBRed

	NodeImages[sim.TypeServer] = imgServer
	NodeImagesRed[sim.TypeServer] = imgServerRed

	NodeImages[sim.TypeDatabase] = imgDB
	NodeImagesRed[sim.TypeDatabase] = imgDBRed

	NodeImages[sim.TypeMessageQueue] = imgQueue
	NodeImagesRed[sim.TypeMessageQueue] = imgQueueRed

	NodeImages[sim.TypeCache] = imgCache
	NodeImagesRed[sim.TypeCache] = imgCache // No red cache image exists

	// Load Font
	fontBytes, err := os.ReadFile("assets/fonts/JetBrainsMono-Regular.ttf")
	if err != nil {
		return err
	}

	FontSource, err = text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		return err
	}

	return nil
}

// Helper to draw text
func DrawText(screen *ebiten.Image, str string, x, y float64, size float64, clr color.Color) {
	if FontSource == nil {
		return
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)
	op.LineSpacing = size * 1.2
	text.Draw(screen, str, &text.GoTextFace{Source: FontSource, Size: size}, op)
}
