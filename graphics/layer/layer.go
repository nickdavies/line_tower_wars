package layer

import (
    "github.com/neagix/Go-SDL/sdl"
)

type Layer interface {

    Setup() error
    Cleanup()

    HandleEvent(event interface{})

    Update(deltaTime int64)
    Render(target *sdl.Surface)

    GetSize() (x uint16, y uint16)
    GetXYOffsets() (x uint16, y uint16)

    setParent(parent Layer)

    Run() error
    End()
}

type layerBase struct {
    child Layer
    parent Layer
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

// Struct for doing nothing on run
type voidRun struct {
    *layerBase
}

func (g *voidRun) Run() error {
    return nil
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

// Struct for simply passing end command up
type bubbleEnd struct {
    *layerBase
}

func (g *bubbleEnd) End() {
    if g.parent != nil {
        g.parent.End()
    }
}

type voidOffsets struct {
    *layerBase
}

func (g *voidOffsets) GetXYOffsets() (uint16, uint16) {
    return g.parent.GetXYOffsets()
}

