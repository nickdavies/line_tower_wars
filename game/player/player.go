package player

import (
    "math/rand"
    "time"
    "fmt"
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
    AStar := astar.NewAStar(stage.Spawn_Size + stage.Grass_Rows + 1, stage.Grass_Cols)

    astar_source := []astar.Point{astar.Point{Row: stage.Spawn_Size - 1, Col: stage.Grass_Cols / 2}}
    astar_target := make([]astar.Point, stage.Grass_Cols)
    for i := 0; i < stage.Grass_Cols; i++ {
        astar_target[i].Row = stage.Spawn_Size + stage.Grass_Rows + 1
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
            p.Money.Add(u.Type.Cost / 2)
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
            row_off, col_off := p.myStage.TopLeft()
            p.AStar.ClearTile(astar.Point{int(loc.Row - row_off), int(loc.Col - col_off)})
            //delete(p.Towers, loc)
        }
    }
}

func (p *Player) BuyUnit(t *unit.UnitType, no_spawn bool) error {
    err := p.Money.Spend(t.Cost)
    if err != nil {
        return err
    }

    if !no_spawn {
        p.NextPlayer.SpawnUnit(t)
    }
    if t.IncomeDelta > 0 {
        p.Money.IncreaseIncome(uint(t.IncomeDelta))
    } else if t.IncomeDelta < 0 {
        p.Money.DecreaseIncome(uint(-1 * t.IncomeDelta))
    }

    return nil
}

func (p *Player) SpawnUnit(t *unit.UnitType) {
    row_off, col_off := p.myStage.TopLeft()
    start := pathing.Locf{
        Row: float64(row_off) + float64(rand.Intn(stage.Spawn_Size)) + 0.5,
        Col: float64(col_off) + float64(rand.Intn(stage.Grass_Cols)) + 0.5,
    }
    u := unit.NewUnit(t, p.Path, &start)

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

func (p *Player) BuildTower(t *tower.TowerType, row, col uint16, no_repath bool) error {
    row_off, col_off := p.myStage.TopLeft()

    abs_row := row + row_off
    abs_col := col + col_off

    if !p.Buildable(abs_row, abs_col) {
        return fmt.Errorf("(%d, %d) is not buildable \n", row, col)
    }

    err := p.Money.Spend(t.Cost)
    if err != nil {
        return err
    }

    p.Towers[pathing.Loc{abs_row, abs_col}] = tower.NewTower(t, pathing.Loc{abs_row, abs_col})

    p.AStar.FillTile(astar.Point{int(row), int(col)}, towerWeight)

    if !no_repath {
        p.Repath()
    }

    return nil
}

func (p *Player) Repath() {
    row_off, col_off := p.myStage.TopLeft()

    path := p.AStar.FindPath(astar.NewRowToRow(), p.astar_source, p.astar_target)
    p.Path = pathing.NewPath(p.AStar, path, row_off, col_off)

    for _, u := range p.Units {
        u.SetPath(p.Path)
    }
}

