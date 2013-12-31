package game

import (
    "fmt"
    "time"
)

import (
    "github.com/banthar/Go-SDL/sdl"
)

type panGame struct {
    gameBase

    // no run loop
    voidRun

    // proxy events
    voidEvents

    // proxy ends
    bubbleEnd

    pan_region_x uint16
    pan_region_y uint16
    pan_speed uint16

    view_x uint16
    view_y uint16

    child_x uint16
    child_y uint16

    surface *sdl.Surface
}

func NewPanGame(pan_region_x, pan_region_y, starting_x, starting_y, pan_speed uint16, child Game) Game {

    if child == nil {
        panic(fmt.Errorf("You must give a child to pan game"))
    }

    pg := &panGame{
        gameBase: gameBase{child: child},

        pan_region_x: pan_region_x,
        pan_region_y: pan_region_y,

        pan_speed: pan_speed,

        view_x: starting_x,
        view_y: starting_y,
    }

    child.setParent(pg)

    return pg
}

func (g *panGame) Setup() error {
    g.child_x, g.child_y = g.child.GetSize()

    g.surface = sdl.CreateRGBSurface(sdl.SWSURFACE, int(g.child_x), int(g.child_y), 32, 0x000000ff, 0x0000ff00, 0x00ff0000, 0xff000000)
    if g.surface == nil {
        return fmt.Errorf("No surface created: %s", sdl.GetError())
    }

    g.surface = sdl.DisplayFormat(g.surface)
    if g.surface == nil {
        return fmt.Errorf("No surface created: %s", sdl.GetError())
    }

    return g.child.Setup()
}

func (g *panGame) Cleanup() {
    g.surface.Free()
    g.child.Cleanup()
}

func (g *panGame) Update(deltaTime int64) {
    var mouse_x int
    var mouse_y int

    errno := sdl.GetMouseState(&mouse_x, &mouse_y)
    if errno != 0 {
        panic(fmt.Errorf("Error getting mouse position: %s", sdl.GetError()))
    }

    parent_x, parent_y := g.parent.GetSize()

    pan_amount := uint16((int64(g.pan_speed) * deltaTime) / int64(time.Second))

    if uint16(mouse_x) < g.pan_region_x {
        if g.view_x < pan_amount {
            g.view_x = 0
        } else {
            g.view_x -= pan_amount
        }
    }

    if uint16(mouse_x) > (parent_x - g.pan_region_x) {
        max_pan := g.child_x - parent_x
        if g.view_x + pan_amount > max_pan {
            g.view_x = max_pan
        } else {
            g.view_x += pan_amount
        }
    }

    if uint16(mouse_y) < g.pan_region_y {
        if g.view_y < pan_amount {
            g.view_y = 0
        } else {
            g.view_y -= pan_amount
        }
    }

    if uint16(mouse_y) > (parent_y - g.pan_region_y) {
        max_pan := g.child_y - parent_y
        if g.view_y + pan_amount > max_pan {
            g.view_y = max_pan
        } else {
            g.view_y += pan_amount
        }
    }
}

func (g *panGame) Render(target *sdl.Surface) {
    g.child.Render(g.surface)
    g.surface.Flip()

    parent_x, parent_y := g.parent.GetSize()

    target.Blit(
        &sdl.Rect{
            X: 0,
            Y: 0,
            W: parent_x,
            H: parent_y,
        },
        g.surface,
        &sdl.Rect{
            X: int16(g.view_x),
            Y: int16(g.view_y),
            W: parent_x,
            H: parent_y,
        },
    )
}

func (g *panGame) GetSize() (uint16, uint16) {
    return g.child.GetSize()
}

