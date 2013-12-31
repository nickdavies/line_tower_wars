package game

import "fmt"

import (
    "github.com/banthar/Go-SDL/sdl"
)

import (
    "../stage"
)

type entityGame struct {
    gameBase

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

func NewEntityGame(s *stage.Stage, square_size int, child Game) Game {
    eg := &entityGame{
        gameBase: gameBase{child: child},

        size_x: uint16(square_size * len(s.Tiles)),
        size_y: uint16(square_size * len(s.Tiles[0])),

        s: s,
        square_size: square_size,
    }

    eg.voidRun = voidRun{&eg.gameBase}
    eg.voidOffsets = voidOffsets{&eg.gameBase}
    eg.voidSetup = voidSetup{&eg.gameBase}
    eg.bubbleEnd = bubbleEnd{&eg.gameBase}

    if child != nil {
        child.setParent(eg)
    }

    return eg
}

func (g *entityGame) HandleEvent(event sdl.Event) {
    switch event.(type) {
    case *sdl.MouseButtonEvent:
        e := event.(*sdl.MouseButtonEvent)
        if e.Type == sdl.MOUSEBUTTONDOWN {
            fmt.Println("click", e.X, e.Y)
        }
    default:
    }

    if g.child != nil {
        g.child.HandleEvent(event)
    }
}

func (g *entityGame) Update(deltaTime int64) {
    if g.child != nil {
        g.child.Update(deltaTime)
    }
}

func (g *entityGame) Render(target *sdl.Surface) {

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

func (g *entityGame) GetSize() (uint16, uint16) {
    return g.size_x, g.size_y
}


