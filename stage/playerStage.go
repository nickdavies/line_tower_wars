package stage

type PlayerStage struct {
    Player int
    Tiles [][]Terrain
}

func (ps *PlayerStage) InitMap() {

    for col := 0; col < PlayerStageWidth; col++ {
        // Setup the top and bottom walls
        for row := 0; row < Wall_Top; row++ {
            ps.Tiles[col][row] = T_Wall
        }

        for row := 0; row < Wall_Bottom; row++ {
            ps.Tiles[col][PlayerStageHeight - row] = T_Wall
        }
    }

    // Set both sides to wall
    for col := 0; col < Wall_Side; col++ {
        for row := 0; row < PlayerStageHeight; row++ {
            ps.Tiles[col][row] = T_Wall
            ps.Tiles[PlayerStageWidth - col][row] = T_Wall
        }
    }

    // Setup goal and spawn area
    for col := Wall_Side; col < (PlayerStageWidth - Wall_Side); col++ {
        for row := 0; row < Spawn_Size; row++ {
            ps.Tiles[col][row] = T_Spawn
        }

        for row := 0; row < Wall_Bottom; row++ {
            ps.Tiles[col][PlayerStageHeight - row] = T_Goal
        }
    }

}
