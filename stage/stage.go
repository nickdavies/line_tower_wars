package stage

const PlayerStageWidth = 25
const PlayerStageHeight = 60

type Stage struct {
    Rows int
    Cols int

    NumPlayers int
    players []*PlayerStage

    Tiles [][]Terrain
}

func NewStage(num_players int) *Stage {

    // Workout field dimentions
    cols := num_players * PlayerStageWidth
    rows := PlayerStageHeight

    // make all the columns
    tiles := make([][]Terrain, cols)

    // Fill each column with a column
    for i := 0; i < cols; i++ {
        tiles[i] = make([]Terrain, rows)
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

