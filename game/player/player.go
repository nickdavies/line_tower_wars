package player

import (
    "github.com/nickdavies/go-astar/astar"
)

import (
    "github.com/nickdavies/line_tower_wars/stage"
    "github.com/nickdavies/line_tower_wars/towers"
    "github.com/nickdavies/line_tower_wars/unit"
    "github.com/nickdavies/line_tower_wars/pathing"
)

var towerWeight = 10000

type loc struct {
    Row uint16
    Col uint16
}

type Player struct {
    myStage *stage.PlayerStage

    astar_source []astar.Point
    astar_target []astar.Point

    Path *pathing.Path

    AStar astar.AStar
    Towers map[loc]*towers.Tower

    AntiCheat *unit.Unit
}

func NewPlayer(myStage *stage.PlayerStage) *Player {

    AStar := astar.NewAStar(stage.Grass_Rows + 2, stage.Grass_Cols)

    astar_source := []astar.Point{astar.Point{Row: 0, Col: stage.Grass_Cols / 2}}
    astar_target := make([]astar.Point, stage.Grass_Cols)
    for i := 0; i < stage.Grass_Cols; i++ {
        astar_target[i].Row = stage.Grass_Rows + 1
        astar_target[i].Col = i
    }

    return &Player{
        myStage: myStage,

        astar_source: astar_source,
        astar_target: astar_target,

        AStar: AStar,
        Towers: make(map[loc]*towers.Tower),
    }
}

func (p *Player) Buildable(row, col uint16) (bool) {
    _, ok := p.Towers[loc{row, col}]
    if ok {
        return false
    }

    return p.myStage.Buildable(row, col)
}

func (p *Player) BuildTower(row, col uint16, no_repath bool) (user_message string) {
    if !p.Buildable(row, col) {
        return "not buildable"
    }
    row_off, col_off := p.myStage.FirstGrass()

    p.Towers[loc{row, col}] = &towers.Tower{}
    p.AStar.FillTile(astar.Point{int(row - row_off + 1), int(col - col_off)}, towerWeight)

    if !no_repath {
        p.Repath()
    }

    return ""
}

func (p *Player) Repath() {
    row_off, col_off := p.myStage.FirstGrass()

    path := p.AStar.FindPath(astar.NewRowToRow(), p.astar_source, p.astar_target)
    p.Path = pathing.NewPath(p.AStar, path, row_off - 1, col_off)

    if p.AntiCheat != nil {
        p.AntiCheat.SetPath(p.Path)
    }
}


