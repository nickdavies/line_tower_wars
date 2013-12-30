package stage

const PlayerStageWidth = 30
const PlayerStageHeight = 30

type Stage struct {
    rows int
    cols int

    Players int
    Tiles [][]Terrain
}

func NewStage(players int) *Stage {

    // Workout field dimentions
    cols := players * PlayerStageWidth
    rows := PlayerStageHeight

    // make all the columns
    tiles := make([][]Terrain, cols)

    // Fill each column with a column
    for i := 0; i < cols; i++ {
        tiles[i] = make([]Terrain, rows)
    }

    s := &Stage{
        rows: rows,
        cols: cols,

        Players: players,
        Tiles: tiles,
    }

    s.InitMap()

    return s
}

func (s *Stage) InitMap() {

    for i := 0; i < s.Players; i++ {
        p := s.GetPlayer(i)
        p.InitMap()
    }

}

func (s *Stage) GetPlayer(player int) *PlayerStage {

    if player < 0 || player >= s.Players {
        panic("Player out of bounds")
    }

    start := player * PlayerStageWidth
    end := start + PlayerStageWidth

    return &PlayerStage{
        Player: player,
        Tiles: s.Tiles[start:end],
    }
}

