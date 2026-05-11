package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"os"
	"sort"

	"github.com/blackviking27/system-design-game/internal/sim"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Rect struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type FrameData struct {
	Frame Rect `json:"frame"`
}

type SpriteSheetData struct {
	Frames map[string]FrameData `json:"frames"`
}

var (
	NodeFrames    = make(map[sim.NodeType][]*ebiten.Image)
	NodeRedFrames = make(map[sim.NodeType][]*ebiten.Image)
	FontSource    *text.GoTextFaceSource
)

func loadComponentAtlas(baseName string) ([]*ebiten.Image, []*ebiten.Image, error) {
	pngPath := fmt.Sprintf("assets/images/%s.png", baseName)
	redPngPath := fmt.Sprintf("assets/images/%s_red.png", baseName)
	jsonPath := fmt.Sprintf("assets/images/%s.json", baseName)
	redJsonPath := fmt.Sprintf("assets/images/%s_red.json", baseName)

	// 1. Load the sprite sheet PNG
	spriteSheet, _, err := ebitenutil.NewImageFromFile(pngPath)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to load png %s: %v", pngPath, err)
	}

	redSpriteSheet, _, err := ebitenutil.NewImageFromFile(redPngPath)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to load png %s: %v", redPngPath, err)
	}

	// 2. Load and parse the JSON file
	jsonFile, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load json %s: %v", jsonPath, err)
	}

	redJsonFile, err := os.ReadFile(redJsonPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load json %s: %v", redJsonPath, err)
	}

	var atlasData SpriteSheetData
	if err := json.Unmarshal(jsonFile, &atlasData); err != nil {
		return nil, nil, fmt.Errorf("failed to parse json %s: %v", jsonPath, err)
	}

	var redAtlasData SpriteSheetData
	if err := json.Unmarshal(redJsonFile, &redAtlasData); err != nil {
		return nil, nil, fmt.Errorf("failed to parse json %s: %v", redJsonPath, err)
	}

	// 3. Extract all frames from the map
	var frameKeys []string
	for keys := range atlasData.Frames {
		frameKeys = append(frameKeys, keys)
	}

	var redFrameKeys []string
	for keys := range redAtlasData.Frames {
		redFrameKeys = append(redFrameKeys, keys)
	}

	// 4. sort the keys
	sort.Strings(frameKeys)
	sort.Strings(redFrameKeys)

	// 5. Extract frames in the sorted order
	var frames []*ebiten.Image
	for _, key := range frameKeys {
		data := atlasData.Frames[key]
		rect := image.Rect(data.Frame.X, data.Frame.Y, data.Frame.X+data.Frame.W, data.Frame.Y+data.Frame.H)
		frames = append(frames, spriteSheet.SubImage(rect).(*ebiten.Image))
	}

	var redFrames []*ebiten.Image
	for _, key := range redFrameKeys {
		data := redAtlasData.Frames[key]
		rect := image.Rect(data.Frame.X, data.Frame.Y, data.Frame.X+data.Frame.W, data.Frame.Y+data.Frame.H)
		redFrames = append(redFrames, redSpriteSheet.SubImage(rect).(*ebiten.Image))
	}

	// Sanity check
	if len(frames) == 0 {
		return nil, nil, fmt.Errorf("Found 0 frames for %s in %s", baseName, jsonPath)
	}

	if len(redFrames) == 0 {
		return nil, nil, fmt.Errorf("Found 0 frames for %s in %s", baseName, redJsonPath)
	}

	return frames, redFrames, nil
}

func LoadAssets() error {
	// Map internal node types of their base filenames
	components := map[sim.NodeType]string{
		sim.TypeLoadBalancer: "load_balancer",
		sim.TypeServer:       "server",
		sim.TypeDatabase:     "db",
		sim.TypeMessageQueue: "queue",
		sim.TypeCache:        "cache",
	}

	// Load each component's atlas
	for nodeType, baseName := range components {
		frames, redFrames, err := loadComponentAtlas(baseName)
		if err != nil {
			return fmt.Errorf("Error while loading assets: %v", err)
		}
		NodeFrames[nodeType] = frames
		NodeRedFrames[nodeType] = redFrames
	}

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
