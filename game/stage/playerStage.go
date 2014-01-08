package stage

import (
    "github.com/nickdavies/line_tower_wars/game/terrain"
)

type PlayerStage struct {
    start uint16
    end uint16

    PlayerNum int
    Tiles [][]terrain.Terrain
    Stage *Stage
}

func NewPlayerStage(player_num int, tiles [][]terrain.Terrain, s *Stage) *PlayerStage {

    start := player_num * PlayerStageWidth
    end := start + PlayerStageWidth

    player := &PlayerStage{
        start: uint16(start),
        end: uint16(end),

        PlayerNum: player_num,
        Tiles: tiles[start:end],
        Stage: s,
    }

    player.initMap()

    return player
}

func (ps *PlayerStage) TopLeft() (row, col uint16) {
    return Wall_Top, ps.start + Wall_Side + Shadow_Side
}

func (ps *PlayerStage) BottomRight() (row, col uint16) {
    return PlayerStageHeight - Wall_Bottom, ps.end
}

func (ps *PlayerStage) Buildable(row, col uint16) bool {
    if row > PlayerStageHeight {
        return false
    }

    if col > ps.end || col < ps.start {
        return false
    }

    return ps.Stage.Tiles[col][row].Buildable()
}

func (ps *PlayerStage) initMap() {

    for col := 0; col < PlayerStageWidth; col++ {
        // Setup the top and bottom walls
        for row := 0; row < Wall_Top; row++ {
            ps.Tiles[col][row] = terrain.T_Wall
        }

        for row := 0; row < Wall_Bottom; row++ {
            ps.Tiles[col][PlayerStageHeight - row - 1] = terrain.T_Wall
        }
    }

    // Set both sides to wall
    for col := 0; col < Shadow_Side; col++ {
        for row := 0; row < PlayerStageHeight; row++ {
            ps.Tiles[col][row] = terrain.T_Shadow
            ps.Tiles[PlayerStageWidth - col - 1][row] = terrain.T_Shadow
        }
    }

    // Set both sides to wall
    for col := 0; col < Wall_Side; col++ {
        for row := 0; row < PlayerStageHeight; row++ {
            ps.Tiles[col + Shadow_Side][row] = terrain.T_Wall
            ps.Tiles[PlayerStageWidth - Shadow_Side - col - 1][row] = terrain.T_Wall
        }
    }

    border_size := Wall_Side + Shadow_Side

    // Setup goal and spawn area
    for col := border_size; col < (PlayerStageWidth - border_size); col++ {
        for row := 0; row < Spawn_Size; row++ {
            ps.Tiles[col][Wall_Top + row] = terrain.T_Spawn
        }

        for row := 0; row < Goal_Size; row++ {
            ps.Tiles[col][PlayerStageHeight - row - Wall_Bottom - 1] = terrain.T_Goal
        }
    }

}
