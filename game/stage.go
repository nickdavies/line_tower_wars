package game

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

    // proxy setup/cleanup
    voidSetup

    // proxy ends
    bubbleEnd

    s *stage.Stage
    square_size int
    terrain_colours map[stage.Terrain]uint32
}

func NewStageGame(players int, square_size int, child Game) Game {
    s := stage.NewStage(players)

    sg := &stageGame{
        gameBase: gameBase{child: child},

        s: s,
        square_size: square_size,

        terrain_colours: map[stage.Terrain]uint32{
            stage.T_Grass: 0x00ff00,
            stage.T_Wall: 0xd1d1d1,
            stage.T_Spawn: 0x3333cc,
            stage.T_Goal: 0xff0000,
        },
    }
    if child != nil {
        child.setParent(sg)
    }

    return sg
}


func (g *stageGame) Update(deltaTime int64) {
    if g.child != nil {
        g.child.Update(deltaTime)
    }
}

func (g *stageGame) Render(target *sdl.Surface) {

    for row := 0; row < len(g.s.Tiles[0]); row++ {
        for col := 0; col < len(g.s.Tiles); col++ {
            target.FillRect(&sdl.Rect{
                X: int16(col * g.square_size),
                Y: int16(row * g.square_size),
                W: uint16(g.square_size),
                H: uint16(g.square_size),
            }, g.terrain_colours[g.s.Tiles[col][row]])
        }
    }

    if g.child != nil {
        g.child.Render(target)
    }

}
