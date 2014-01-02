package stage

import (
    "github.com/nickdavies/line_tower_wars/terrain"
)

const PlayerStageWidth = 21
const PlayerStageHeight = 60

const Wall_Side = 1
const Wall_Top = 1
const Wall_Bottom = 1

const Shadow_Side = 4

const Spawn_Size = 4
const Goal_size = 2

type Stage struct {
    Rows uint16
    Cols uint16

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

    s := &Stage{
        Rows: uint16(rows),
        Cols: uint16(cols),

        NumPlayers: num_players,

        Tiles: tiles,
    }

    // Build Players
    players := make([]*PlayerStage, num_players)
    for i := 0; i < num_players; i++ {
        players[i] = NewPlayerStage(i, tiles, s)
    }

    s.players = players

    return s
}

func (s *Stage) GetPlayer(player int) *PlayerStage {

    if player < 0 || player >= s.NumPlayers {
        panic("Player out of bounds")
    }

    return s.players[player]
}

func (s *Stage) GetTerrain(row, col uint16) terrain.Terrain {
    if row > s.Rows {
        return terrain.T_NIL
    }

    if col > s.Cols {
        return terrain.T_NIL
    }

    return s.Tiles[col][row]
}

func (s *Stage) GetOwner(row, col uint16) int {
    return int(col / PlayerStageWidth)
}


