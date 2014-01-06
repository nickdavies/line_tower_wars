package game

type MoneyConfig struct {
    Balance uint
    Income uint
    MinIncome uint
}

type StageConfig struct {
    PlayerStageWidth uint
    PlayerStageHeight uint

    Wall_Side uint
    Wall_Top uint
    Wall_Bottom uint

    Shadow_Side uint

    Spawn_Size uint
    Goal_Size uint
}

type TowerConfig struct {
    Size_Row uint
    Size_Col uint
}

type GameConfig struct {
    MoneyConfig
    StageConfig
    TowerConfig
}

