package game

import (
    "fmt"
)

import (
    "github.com/banthar/Go-SDL/sdl"
)

import (
    "../stage"
)

type stageGame struct {
    gameBase

    // no run loop
    voidRun

    // proxy events
    voidEvents

    // no offsets
    voidOffsets

    // proxy ends
    bubbleEnd

    size_x uint16
    size_y uint16

    s *stage.Stage
    square_size int
    terrain_colours map[stage.Terrain]uint32

    surface *sdl.Surface
}

func NewStageGame(s *stage.Stage, square_size int, child Game) Game {

    sg := &stageGame{
        gameBase: gameBase{child: child},

        size_x: uint16(square_size * len(s.Tiles)),
        size_y: uint16(square_size * len(s.Tiles[0])),

        s: s,
        square_size: square_size,

        terrain_colours: map[stage.Terrain]uint32{
            stage.T_Grass: 0x00ff00,
            stage.T_Wall: 0xd1d1d1,
            stage.T_Spawn: 0x3333cc,
            stage.T_Goal: 0xff0000,
        },
    }

    sg.voidRun = voidRun{&sg.gameBase}
    sg.voidEvents = voidEvents{&sg.gameBase}
    sg.voidOffsets = voidOffsets{&sg.gameBase}
    sg.bubbleEnd = bubbleEnd{&sg.gameBase}

    if child != nil {
        child.setParent(sg)
    }

    return sg
}

func (g *stageGame) Setup() error {
    g.surface = sdl.CreateRGBSurface(sdl.HWSURFACE, int(g.size_x), int(g.size_y), 32, 0x000000ff, 0x0000ff00, 0x00ff0000, 0xff000000)
    if g.surface == nil {
        return fmt.Errorf("No surface created: %s", sdl.GetError())
    }

    g.surface = sdl.DisplayFormat(g.surface)
    if g.surface == nil {
        return fmt.Errorf("No surface created: %s", sdl.GetError())
    }

    for row := 0; row < len(g.s.Tiles[0]); row++ {
        for col := 0; col < len(g.s.Tiles); col++ {
            g.surface.FillRect(&sdl.Rect{
                X: int16(col * g.square_size),
                Y: int16(row * g.square_size),
                W: uint16(g.square_size),
                H: uint16(g.square_size),
            }, g.terrain_colours[g.s.Tiles[col][row]])
        }
    }

    if g.child != nil {
        return g.child.Setup()
    }

    return nil
}

func (g *stageGame) Cleanup() {
    g.surface.Free()
    if g.child != nil {
        g.child.Cleanup()
    }
}


func (g *stageGame) Update(deltaTime int64) {
    if g.child != nil {
        g.child.Update(deltaTime)
    }
}

func (g *stageGame) Render(target *sdl.Surface) {

    target.Blit(
        &sdl.Rect{
            X: 0,
            Y: 0,
            W: uint16(g.size_x),
            H: uint16(g.size_y),
        },
        g.surface,
        nil,
    )

    if g.child != nil {
        g.child.Render(target)
    }

}

func (g *stageGame) GetSize() (uint16, uint16) {
    return g.size_x, g.size_y
}

