package layer

import (
    "github.com/neagix/Go-SDL/sdl"
)

type sdlLayer struct {
    // Get parent/child fields
    layerBase

    // no offsets
    voidOffsets

    // no setup
    voidSetup

    // no events
    voidEvents

    running bool

    // Display size
    x uint16
    y uint16

    // main display
    display *sdl.Surface
}


func NewSdlLayer(display *sdl.Surface, x, y uint16, child Layer) Layer {
    g := &sdlLayer{
        layerBase: layerBase{child: child},

        x: x,
        y: y,

        display: display,
    }

    g.voidOffsets = voidOffsets{&g.layerBase}
    g.voidSetup = voidSetup{&g.layerBase}
    g.voidEvents = voidEvents{&g.layerBase}

    if child != nil {
        child.setParent(g)
    }

    return g
}

func (g *sdlLayer) HandleEvent(event interface{}) {
    if g.child != nil {
        g.child.HandleEvent(event)
    }
}

func (g *sdlLayer) Update(deltaTime int64) {
    if g.child != nil {
        g.child.Update(deltaTime)
    }
}

func (g *sdlLayer) Render(target *sdl.Surface) {
    if g.child != nil {
        g.child.Render(g.display)
    }

    g.display.Flip()
}

func (g *sdlLayer) GetSize() (uint16, uint16) {
    return g.x, g.y
}
