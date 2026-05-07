# Assets and Font Integration Plan & Code Changes

## Setup

First, download the font:
```bash
wget "https://github.com/JetBrains/JetBrainsMono/raw/master/fonts/ttf/JetBrainsMono-Regular.ttf" -O assets/fonts/JetBrainsMono-Regular.ttf
```

## 1. Create `internal/assets/assets.go`

To avoid import cycles between `internal/engine` and `internal/ui`, we create a new package for assets.

```go
package assets

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

func Load() error {
	// Load Images
	imgLB, _, err := ebitenutil.NewImageFromFile("assets/fonts/network.png")
	if err != nil { return err }
	imgLBRed, _, err := ebitenutil.NewImageFromFile("assets/fonts/lb_red.png")
	if err != nil { return err }

	imgServer, _, err := ebitenutil.NewImageFromFile("assets/fonts/server.png")
	if err != nil { return err }
	imgServerRed, _, err := ebitenutil.NewImageFromFile("assets/fonts/server_red.png")
	if err != nil { return err }

	imgDB, _, err := ebitenutil.NewImageFromFile("assets/fonts/db.png")
	if err != nil { return err }
	imgDBRed, _, err := ebitenutil.NewImageFromFile("assets/fonts/db_red.png")
	if err != nil { return err }

	imgQueue, _, err := ebitenutil.NewImageFromFile("assets/fonts/queue.png")
	if err != nil { return err }
	imgQueueRed, _, err := ebitenutil.NewImageFromFile("assets/fonts/queue_red.png")
	if err != nil { return err }

	imgCache, _, err := ebitenutil.NewImageFromFile("assets/fonts/cache.png")
	if err != nil { return err }

	NodeImages[sim.TypeLoadBalancer] = imgLB
	NodeImagesRed[sim.TypeLoadBalancer] = imgLBRed

	NodeImages[sim.TypeServer] = imgServer
	NodeImagesRed[sim.TypeServer] = imgServerRed

	NodeImages[sim.TypeDatabase] = imgDB
	NodeImagesRed[sim.TypeDatabase] = imgDBRed

	NodeImages[sim.TypeMessageQueue] = imgQueue
	NodeImagesRed[sim.TypeMessageQueue] = imgQueueRed

	NodeImages[sim.TypeCache] = imgCache
	NodeImagesRed[sim.TypeCache] = imgCache // No red cache

	// Load Font
	fontBytes, err := os.ReadFile("assets/fonts/JetBrainsMono-Regular.ttf")
	if err != nil { return err }
	
	FontSource, err = text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil { return err }

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
	text.Draw(screen, str, &text.GoTextFace{Source: FontSource, Size: size}, op)
}
```

## 2. Initialize in `cmd/game/main.go`

```go
package main

import (
	"log"

	"github.com/blackviking27/system-design-game/internal/assets"
	"github.com/blackviking27/system-design-game/internal/engine"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	if err := assets.Load(); err != nil {
		log.Fatalf("Failed to load assets: %v", err)
	}

	game := &engine.Game{
		CurrentScene: &engine.MainMenuScene{},
	}
	// ... rest of main
}
```

## 3. Render Nodes in `internal/engine/render.go`

Replace the `vector.FillRect` block inside `DrawNetwork` with image drawing:

```go
import (
	"fmt"
	"image/color"

	"github.com/blackviking27/system-design-game/internal/assets"
	"github.com/blackviking27/system-design-game/internal/sim"
	"github.com/blackviking27/system-design-game/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ...

	// Draw nodes(servers)
	for _, node := range game.Network.Nodes {
		isFailing := len(node.Queue) >= node.MaxRam

		var img *ebiten.Image
		if isFailing {
			img = assets.NodeImagesRed[node.Type]
		} else {
			img = assets.NodeImages[node.Type]
		}

		w, h := float64(80), float64(50)
		startX, startY := node.X-w/2, node.Y-h/2

		if img != nil {
			op := &ebiten.DrawImageOptions{}
			
			// Optional: Scale image to fit 80x50 exactly if it's too big/small
			bounds := img.Bounds()
			scaleX := w / float64(bounds.Dx())
			scaleY := h / float64(bounds.Dy())
			op.GeoM.Scale(scaleX, scaleY)
			
			op.GeoM.Translate(startX, startY)
			screen.DrawImage(img, op)
		} else {
			// Fallback
			vector.FillRect(screen, float32(startX), float32(startY), float32(w), float32(h), color.RGBA{100, 255, 250, 255}, true)
		}

		// Stats for node
		stats := fmt.Sprintf("%s\nRAM: %d/%d\nDrop: %d", node.ID, len(node.Queue), node.MaxRam, node.DroppedCount)
		
		assets.DrawText(screen, stats, startX, startY-50, 12, color.White)
	}
// ...
```

## 4. UI Elements using Custom Font

In `internal/ui/hud.go`, replace `ebitenutil.DebugPrintAt` with `assets.DrawText`:
```go
import (
	"fmt"
	"image/color"

	"github.com/blackviking27/system-design-game/internal/assets"
	"github.com/blackviking27/system-design-game/internal/sim"
	"github.com/hajimehoshi/ebiten/v2"
)

func DrawHUD(screen *ebiten.Image, network *sim.Network, levelName string, targetTick, maxDropped int, isGameOver, isVictory bool) {
	// ...
	assets.DrawText(screen, stats, 10, 10, 14, colorText)

	if isGameOver || isVictory {
		// ...
		assets.DrawText(screen, msg, float64(w/2-100), float64(h/2-20), 20, colorText)
	}
}
```

In `internal/ui/tray.go`:
```go
func DrawTray(screen *ebiten.Image, budget int) {
	// ...
	assets.DrawText(screen, fmt.Sprintf("BUDGET: $%d", budget), 20, float64(trayY+40), 16, color.White)

	// ...
	for i, template := range sim.Catalog {
		// ... draw icon ...
		
		// Draw label and Cost
		label := fmt.Sprintf("%s\n%d", template.Name, template.Cost)
		assets.DrawText(screen, label, float64(x+50), float64(y+50), 12, color.White)
	}
}
```

In `internal/engine/mainmenu.go`:
```go
func (this *MainMenuScene) Draw(screen *ebiten.Image) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	// ...
	assets.DrawText(screen, "=== SYSTEM DESIGN SANDBOX ===", float64(w/2-150), float64(h/2-150), 20, color.White)
	// ...
	assets.DrawText(screen, "START LEVEL 1", float64(w/2-55), float64(h/2-8), 16, color.White)
}
```