package layer

import (
    "github.com/neagix/Go-SDL/sdl"
)

type SdlLayerCfg struct {
    Display *sdl.Surface
    ScreenX uint16
    ScreenY uint16
}

type sdlLayer struct {
    // Get parent/child fields
    layerBase

    // no offsets
    voidOffsets

    // no setup
    voidSetup

    // no events
    voidEvents

    // Display size
    x uint16
    y uint16

    // main display
    display *sdl.Surface
}

func init() {
    registerLayer("sdl", func(base layerBase, cfg interface{}) Layer {
        sdl_cfg := cfg.(SdlLayerCfg)

        l := &sdlLayer{
            layerBase: base,

            x: sdl_cfg.ScreenX,
            y: sdl_cfg.ScreenY,

            display: sdl_cfg.Display,
        }

        l.voidOffsets = voidOffsets{&l.layerBase}
        l.voidSetup = voidSetup{&l.layerBase}
        l.voidEvents = voidEvents{&l.layerBase}

        return l
    })
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
