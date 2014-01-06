package player

import (
    "github.com/nickdavies/go-astar/astar"
)

import (
    "github.com/nickdavies/line_tower_wars/game/stage"
    "github.com/nickdavies/line_tower_wars/game/towers"
    "github.com/nickdavies/line_tower_wars/game/unit"
    "github.com/nickdavies/line_tower_wars/game/pathing"
)

var towerWeight = 10000

type Player struct {
    NextPlayer *Player
    myStage *stage.PlayerStage

    astar_source []astar.Point
    astar_target []astar.Point

    Path *pathing.Path

    AStar astar.AStar
    Towers map[pathing.Loc]*towers.Tower

    // the units attacking you
    Units map[int]*unit.Unit
    unitNum int

    lives int
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
        Towers: make(map[pathing.Loc]*towers.Tower),
        Units: make(map[int]*unit.Unit),
    }
}

func (p *Player) GainLife() {
    p.lives++
}

func (p *Player) GetLives() int {
    return p.lives
}

func (p *Player) Update(deltaTime int64) {
    for u_id, u := range p.Units {
        u.Update(deltaTime)

        if u.AtEnd() {
            delete(p.Units, u_id)
            p.lives--
            p.NextPlayer.GainLife()
        }
    }

    for loc, t := range p.Towers {
        t.Update(deltaTime)

        // Kill tower on path
        if p.Path.On(loc) {
            delete(p.Towers, loc)
        }
    }
}

func (p *Player) SpawnUnit(u *unit.Unit) {

}

func (p *Player) Buildable(row, col uint16) (bool) {
    _, ok := p.Towers[pathing.Loc{row, col}]
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

    p.Towers[pathing.Loc{row, col}] = &towers.Tower{}
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

    for _, u := range p.Units {
        u.SetPath(p.Path)
    }
}

