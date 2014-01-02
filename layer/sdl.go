package layer

import (
    "time"
)

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

    running bool

    // Display size
    x uint16
    y uint16

    // main display
    display *sdl.Surface
}


func NewSdlLayer(display *sdl.Surface, x, y int, child Layer) Layer {
    g := &sdlLayer{
        layerBase: layerBase{child: child},

        x: uint16(x),
        y: uint16(y),

        display: display,
    }

    g.voidOffsets = voidOffsets{&g.layerBase}
    g.voidSetup = voidSetup{&g.layerBase}

    if child != nil {
        child.setParent(g)
    }

    return g
}

func (g *sdlLayer) HandleEvent(event interface{}) {
    switch event.(type) {
    case sdl.QuitEvent:
        g.End()
        return
    case sdl.KeyboardEvent:
        if event.(sdl.KeyboardEvent).Keysym.Sym == sdl.K_F1 {
            g.End()
        }
        return
    default:
    }
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

func (g *sdlLayer) Run() error {
    var err error
    var end_ch = make(chan interface{})

    err = g.Setup()
    if err != nil {
        return err
    }
    defer g.Cleanup()

    go func() {
        for {
            select {
            case <-end_ch:
                return
            case event, ok := <-sdl.Events:
                if !ok {
                    g.running = false
                    continue
                }
                g.HandleEvent(event)
            }
        }
    }()
    defer func() {
        end_ch <-nil
    }()

    var last_time = time.Now().UnixNano()

    g.running = true
    for g.running {
        // Update State
        g.Update(time.Now().UnixNano() - last_time)
        last_time = time.Now().UnixNano()

        // Update Screen
        g.Render(nil)
    }

    return nil
}

func (g *sdlLayer) End() {
    g.running = false
}

