package tower

import (
    "time"
)

import (
    "github.com/nickdavies/line_tower_wars/game/pathing"
    "github.com/nickdavies/line_tower_wars/game/unit"
)

type TowerType struct {
    Id int
    Name string

    Range float64
    FireRate int
    Damage int

    Health uint

    Cost uint
}

type Tower struct {
    Type *TowerType
    Loc pathing.Loc

    attackTimer int64
}

func NewTower(t *TowerType, loc pathing.Loc) *Tower {
    return &Tower{
        Type: t,
        Loc: loc,
    }
}

func (t *Tower) Update(deltaTime int64, units map[int]*unit.Unit) {
    t.attackTimer += deltaTime

    if t.attackTimer / int64(time.Second) > int64(t.Type.FireRate) {
        var best_u *unit.Unit
        locf := t.Loc.ToFloat(0.5)

        for _, u := range units {
            dist := locf.Dist(u.Loc)
            if dist < t.Type.Range {
                if best_u == nil || u.Loc.Row > best_u.Loc.Row {
                    best_u = u
                }
            }
        }

        if best_u != nil {
            best_u.Health -= int(t.Type.Damage)
        }

        t.attackTimer = 0
    }
}
