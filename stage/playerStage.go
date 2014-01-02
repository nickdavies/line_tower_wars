package stage

import (
    "github.com/nickdavies/line_tower_wars/terrain"
)

type PlayerStage struct {
    PlayerNum int
    Tiles [][]terrain.Terrain
}

func NewPlayerStage(player_num int, tiles [][]terrain.Terrain) *PlayerStage {

    start := player_num * PlayerStageWidth
    end := start + PlayerStageWidth

    player := &PlayerStage{
        PlayerNum: player_num,
        Tiles: tiles[start:end],
    }

    player.initMap()

    return player
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
    for col := 0; col < Wall_Side; col++ {
        for row := 0; row < PlayerStageHeight; row++ {
            ps.Tiles[col][row] = terrain.T_Wall
            ps.Tiles[PlayerStageWidth - col - 1][row] = terrain.T_Wall
        }
    }

    // Setup goal and spawn area
    for col := Wall_Side; col < (PlayerStageWidth - Wall_Side); col++ {
        for row := 0; row < Spawn_Size; row++ {
            ps.Tiles[col][Wall_Top + row] = terrain.T_Spawn
        }

        for row := 0; row < Wall_Bottom; row++ {
            ps.Tiles[col][PlayerStageHeight - row - Wall_Bottom - 1] = terrain.T_Goal
        }
    }

}
