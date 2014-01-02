package stage

import (
    "github.com/nickdavies/line_tower_wars/terrain"
)

const PlayerStageWidth = 20
const PlayerStageHeight = 60

const Wall_Side = 4
const Wall_Top = 2
const Wall_Bottom = 2

const Spawn_Size = 4
const Goal_size = 4

type Stage struct {
    Rows int
    Cols int

    NumPlayers int
    players []*PlayerStage

    Tiles [][]terrain.Terrain
}

func NewStage(num_players int) *Stage {

    // Workout field dimentions
    cols := num_players * PlayerStageWidth
    rows := PlayerStageHeight

    // make all the columns
    tiles := make([][]terrain.Terrain, cols)

    // Fill each column with a column
    for i := 0; i < cols; i++ {
        tiles[i] = make([]terrain.Terrain, rows)
    }


    // Build Players
    players := make([]*PlayerStage, num_players)
    for i := 0; i < num_players; i++ {
        players[i] = NewPlayerStage(i, tiles)
    }

    return &Stage{
        Rows: rows,
        Cols: cols,

        NumPlayers: num_players,
        players: players,

        Tiles: tiles,
    }
}

func (s *Stage) GetPlayer(player int) *PlayerStage {

    if player < 0 || player >= s.NumPlayers {
        panic("Player out of bounds")
    }

    return s.players[player]
}

