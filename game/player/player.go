package player

import (
    "time"
)

import (
    "github.com/nickdavies/go-astar/astar"
)

import (
    "github.com/nickdavies/line_tower_wars/game/stage"
    "github.com/nickdavies/line_tower_wars/game/money"
    "github.com/nickdavies/line_tower_wars/game/tower"
    "github.com/nickdavies/line_tower_wars/game/unit"
    "github.com/nickdavies/line_tower_wars/game/pathing"
)

var towerWeight = 10000

type Player struct {
    NextPlayer *Player
    PrevPlayer *Player

    myStage *stage.PlayerStage

    astar_source []astar.Point
    astar_target []astar.Point

    Path *pathing.Path
    Money money.PlayerBalance

    AStar astar.AStar
    Towers map[pathing.Loc]*tower.Tower

    // the units attacking you
    Units map[int]*unit.Unit
    unitNum int

    lives int
    died bool

    moneyTimer int64
}

func NewPlayer(myStage *stage.PlayerStage, balance money.PlayerBalance) *Player {

    AStar := astar.NewAStar(stage.Grass_Rows + 2, stage.Grass_Cols)

    astar_source := []astar.Point{astar.Point{Row: 0, Col: stage.Grass_Cols / 2}}
    astar_target := make([]astar.Point, stage.Grass_Cols)
    for i := 0; i < stage.Grass_Cols; i++ {
        astar_target[i].Row = stage.Grass_Rows + 1
        astar_target[i].Col = i
    }

    p := &Player{
        myStage: myStage,

        astar_source: astar_source,
        astar_target: astar_target,

        AStar: AStar,
        Towers: make(map[pathing.Loc]*tower.Tower),
        Units: make(map[int]*unit.Unit),

        // TODO: make this a proper variable
        lives: 50,

        Money: balance,
    }
    p.Repath()

    return p
}

func (p *Player) GainLife() {
    p.lives++
}

func (p *Player) GetLives() int {
    return p.lives
}

func (p *Player) Die() {
    if p.died {
        return
    }

    for loc, _ := range p.Towers {
        delete(p.Towers, loc)
    }
    p.Repath()

    p.PrevPlayer.NextPlayer = p.NextPlayer
    p.NextPlayer.PrevPlayer = p.PrevPlayer
}

func (p *Player) Update(deltaTime int64) {
    //fmt.Println(p.moneyTimer / int64(time.Millisecond), deltaTime, p.moneyTimer / int64(time.Second))
    p.moneyTimer += deltaTime

    if p.moneyTimer / int64(time.Second) > p.Money.IncomeInterval() {
        p.Money.PayIncome()
        p.moneyTimer = 0
    }

    for u_id, u := range p.Units {
        if u.Health <= 0 {
            delete(p.Units, u_id)
            continue
        }

        u.Update(deltaTime)

        if u.AtEnd() {
            delete(p.Units, u_id)
            if p.lives > 0 {
                p.lives--
                p.PrevPlayer.GainLife()
            }
        }
    }

    for loc, t := range p.Towers {
        t.Update(deltaTime, p.Units)

        // Kill tower on path
        if p.Path.On(loc) {
            delete(p.Towers, loc)
        }
    }
}

func (p *Player) BuyUnit(t *unit.UnitType, no_spawn bool) bool {
    ok := p.Money.Spend(t.Cost)
    if ok {
        if !no_spawn {
            p.NextPlayer.SpawnUnit(unit.NewUnit(t, p.NextPlayer.Path))
        }
        if t.IncomeDelta > 0 {
            p.Money.IncreaseIncome(uint(t.IncomeDelta))
        } else if t.IncomeDelta < 0 {
            p.Money.DecreaseIncome(uint(-1 * t.IncomeDelta))
        }

        return true
    }
    return false
}

func (p *Player) SpawnUnit(u *unit.Unit) {
    p.unitNum++
    p.Units[p.unitNum] = u
}

func (p *Player) Buildable(row, col uint16) (bool) {
    _, ok := p.Towers[pathing.Loc{row, col}]
    if ok {
        return false
    }

    return p.myStage.Buildable(row, col)
}

func (p *Player) BuildTower(t *tower.TowerType, row, col uint16, no_repath bool) (user_message string) {
    if !p.Buildable(row, col) {
        return "not buildable"
    }

    ok := p.Money.Spend(t.Cost)
    if !ok {
        return "no money"
    }

    row_off, col_off := p.myStage.FirstGrass()

    p.Towers[pathing.Loc{row, col}] = &tower.Tower{
        Type: t,
    }
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

