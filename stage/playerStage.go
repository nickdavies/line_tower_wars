package stage

type PlayerStage struct {
    PlayerNum int
    Tiles [][]Terrain
}

func NewPlayerStage(player_num int, tiles [][]Terrain) *PlayerStage {

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
            ps.Tiles[col][row] = T_Wall
        }

        for row := 0; row < Wall_Bottom; row++ {
            ps.Tiles[col][PlayerStageHeight - row - 1] = T_Wall
        }
    }

    // Set both sides to wall
    for col := 0; col < Wall_Side; col++ {
        for row := 0; row < PlayerStageHeight; row++ {
            ps.Tiles[col][row] = T_Wall
            ps.Tiles[PlayerStageWidth - col - 1][row] = T_Wall
        }
    }

    // Setup goal and spawn area
    for col := Wall_Side; col < (PlayerStageWidth - Wall_Side); col++ {
        for row := 0; row < Spawn_Size; row++ {
            ps.Tiles[col][Wall_Top + row] = T_Spawn
        }

        for row := 0; row < Wall_Bottom; row++ {
            ps.Tiles[col][PlayerStageHeight - row - Wall_Bottom - 1] = T_Goal
        }
    }

}
