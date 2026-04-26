package engine

import (
	"github.com/blackviking27/system-design-game/internal/sim"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Process all mouse and keyboard interactions
func (this *Game) HandleInput() {
	x, y := ebiten.CursorPosition()
	this.mouseX, this.mouseY = x, y

	// 1. DRAG and DROP a node
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if node := this.getNodeAt(x, y); node != nil {
			this.draggingNode = node

			//Calculate tlhe offset
			this.dragOffsetX = node.X - float64(x)
			this.dragOffsetY = node.Y - float64(y)
		}
	}

	// Updating the node position
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if this.draggingNode != nil {
			this.draggingNode.X = float64(x) + this.dragOffsetX
			this.draggingNode.Y = float64(y) + this.dragOffsetY
		}
	} else {
		this.draggingNode = nil
	}

	// 2. Dynamic linking
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		if node := this.getNodeAt(x, y); node != nil {
			this.linkingNode = node
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		if this.linkingNode != nil {
			targetNode := this.getNodeAt(x, y)
			if targetNode != nil && targetNode != this.linkingNode {
				this.linkingNode.LinkTo(targetNode)
			}
			this.linkingNode = nil // Cancel the draw state
		}
	}

}

// Check if the given X, Y node intersect with a node
func (this *Game) getNodeAt(x, y int) *sim.Node {
	fx, fy := float64(x), float64(y)
	w, h := float64(80), float64(50)

	for _, node := range this.Network.Nodes {
		left := node.X - w/2
		right := node.X + w/2
		top := node.Y - h/2
		bottom := node.Y + h/2

		if fx >= left && fx <= right && fy >= top && fy <= bottom {
			return node
		}
	}
	return nil
}
