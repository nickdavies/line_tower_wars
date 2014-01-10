package layer

import (
    "github.com/nickdavies/line_tower_wars/game"

    "github.com/nickdavies/line_tower_wars/graphics/texture"
)

type layerBase struct {
    child Layer
    parent Layer

    game *game.Game
    controls game.PlayerControls
    texture_map texture.TextureMap

    square_size uint16
}

func (g *layerBase) setParent(parent Layer) {
    g.parent = parent
}

// Struct for Layers with no setup/cleanup
type voidSetup struct {
    *layerBase
}

func (g *voidSetup) Setup() error {
    if g.child != nil {
        return g.child.Setup()
    }

    return nil
}

func (g *voidSetup) Cleanup() {
    if g.child != nil {
        g.child.Cleanup()
    }
}

// Struct for passing events though
type voidEvents struct {
    *layerBase
}

func (g *voidEvents) HandleEvent(event interface{}) {
    if g.child != nil {
        g.child.HandleEvent(event)
    }
}

type voidOffsets struct {
    *layerBase
}

func (g *voidOffsets) GetXYOffsets() (uint16, uint16) {
    return g.parent.GetXYOffsets()
}
