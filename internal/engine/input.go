package engine

import (
	"fmt"

	"github.com/blackviking27/system-design-game/internal/sim"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Process all mouse and keyboard interactions
func (this *GameplayScene) HandleInput() {
	x, y := ebiten.CursorPosition()
	this.mouseX, this.mouseY = x, y

	sh := this.screenHeight
	if sh == 0 {
		sh = 600
	}
	trayY := sh - 100

	// 1. DRAG and DROP a node
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

		if y >= trayY {
			// Tray click
			// Check if catalog item is clicked
			itemIndex := (x - 150) / 180
			if x >= 150 && itemIndex >= 0 && itemIndex < len(sim.Catalog) {
				template := sim.Catalog[itemIndex]
				if this.CurrentBudget >= template.Cost {
					// Buying the item
					this.CurrentBudget -= template.Cost

					// Generate unique id based on map size
					id := fmt.Sprintf("%s-%d", template.Type, len(this.Network.Nodes)+1)
					newNode := sim.NewNode(id, template.Type, template.MaxRam, template.ProcessPower, template.Cost)
					newNode.X, newNode.Y = float64(x), float64(y)

					this.draggingNode = newNode
					this.dragOffsetX, this.dragOffsetY = 0, 0
				}
			}

		} else {
			// Canvas click: Drag a node on the canvas
			if node := this.getNodeAt(x, y); node != nil {
				this.draggingNode = node

				//Calculate tlhe offset
				this.dragOffsetX = node.X - float64(x)
				this.dragOffsetY = node.Y - float64(y)
			}
		}

	}

	// Updating the node position
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if this.draggingNode != nil {
			this.draggingNode.X = float64(x) + this.dragOffsetX
			this.draggingNode.Y = float64(y) + this.dragOffsetY
		}
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if this.draggingNode != nil {
			// Did the user drop the node
			if y >= trayY {
				// Dropped back into the tray
				// New node does not exist in the network, refund the money into the tray
				if _, exists := this.Network.Nodes[this.draggingNode.ID]; !exists {
					this.CurrentBudget += this.draggingNode.Cost
				} else {
					// Existing node, remove the node from network
					if this.draggingNode.ID != "LB_MAIN" {
						this.Network.RemoveNode(this.draggingNode.ID)
						this.CurrentBudget += this.draggingNode.Cost
					}
				}
			} else {
				// Dropped the component on the canvas
				if _, exists := this.Network.Nodes[this.draggingNode.ID]; !exists {
					this.Network.Nodes[this.draggingNode.ID] = this.draggingNode
				}
			}
			this.draggingNode = nil
		}
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

	// 3. Node deletion via keyboard (Hover + Backspace)
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) || inpututil.IsKeyJustPressed(ebiten.KeyDelete) {
		if node := this.getNodeAt(x, y); node != nil {
			this.Network.RemoveNode(node.ID)
			this.CurrentBudget += node.Cost
		}
	}

}

// Check if the given X, Y node intersect with a node
func (this *GameplayScene) getNodeAt(x, y int) *sim.Node {
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
