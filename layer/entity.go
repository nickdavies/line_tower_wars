package layer

import "fmt"

import (
    "github.com/neagix/Go-SDL/sdl"
)

import (
    "github.com/nickdavies/line_tower_wars/stage"
)

type entityLayer struct {
    layerBase

    // no run loop
    voidRun

    // no offsets
    voidOffsets

    // no setup
    voidSetup

    // proxy ends
    bubbleEnd

    size_x uint16
    size_y uint16

    s *stage.Stage
    square_size int
}

func NewEntityLayer(s *stage.Stage, square_size int, child Layer) Layer {
    eg := &entityLayer{
        layerBase: layerBase{child: child},

        size_x: uint16(square_size * len(s.Tiles)),
        size_y: uint16(square_size * len(s.Tiles[0])),

        s: s,
        square_size: square_size,
    }

    eg.voidRun = voidRun{&eg.layerBase}
    eg.voidOffsets = voidOffsets{&eg.layerBase}
    eg.voidSetup = voidSetup{&eg.layerBase}
    eg.bubbleEnd = bubbleEnd{&eg.layerBase}

    if child != nil {
        child.setParent(eg)
    }

    return eg
}

func (g *entityLayer) HandleEvent(event interface{}) {
    switch event.(type) {
    case sdl.MouseButtonEvent:
        e := event.(sdl.MouseButtonEvent)
        if e.Type == sdl.MOUSEBUTTONDOWN {
            fmt.Println("click", e.X, e.Y)
        }
    default:
    }

    if g.child != nil {
        g.child.HandleEvent(event)
    }
}

func (g *entityLayer) Update(deltaTime int64) {
    if g.child != nil {
        g.child.Update(deltaTime)
    }
}

func (g *entityLayer) Render(target *sdl.Surface) {

    var mouse_x int
    var mouse_y int

    sdl.GetMouseState(&mouse_x, &mouse_y)

    x_off, y_off := g.GetXYOffsets()

    mouse_x += int(x_off)
    mouse_y += int(y_off)

    square_x := mouse_x - (mouse_x % g.square_size)
    square_y := mouse_y - (mouse_y % g.square_size)

    target.FillRect(
        &sdl.Rect{
            X: int16(square_x),
            Y: int16(square_y),
            W: uint16(g.square_size),
            H: uint16(g.square_size),
        },
        0x000000,
    )

    if g.child != nil {
        g.child.Render(target)
    }

}

func (g *entityLayer) GetSize() (uint16, uint16) {
    return g.size_x, g.size_y
}


