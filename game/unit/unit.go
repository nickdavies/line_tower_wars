package unit

import (
    "time"
    "sync"
)

import (
    "github.com/nickdavies/line_tower_wars/game/pathing"
)

type UnitType struct {
    Id int
    Name string

    Speed float64
    IncomeDelta int

    Health uint

    Cost uint
}

type Unit struct {
    sync.Mutex

    Type *UnitType

    Loc pathing.Locf

    path *pathing.Path
    return_path *pathing.Path
}

func NewUnit(t *UnitType, path *pathing.Path) *Unit {
    u := &Unit{
        Type: t,
        Loc: path.Startf(),
    }
    u.SetPath(path)

    return u
}

func (u *Unit) Update(timeDelta int64) {
    u.Lock()
    defer u.Unlock()

    var end bool

    // If there is a path to return to the main
    // path take that first
    if u.return_path != nil {
        // if were already on the path then cancel the extra routing
        if u.path.On(u.Loc.ToInt()) {
            u.return_path = nil
        } else {
            u.Loc, end = u.return_path.Move(u.Loc, float64(timeDelta) * u.Type.Speed / float64(time.Second))
            if end {
                u.return_path = nil
            }
            return
        }
    }

    u.Loc, _ = u.path.Move(u.Loc, float64(timeDelta) * u.Type.Speed / float64(time.Second))
}

func (u *Unit) AtEnd() bool {
    return u.Loc.ToInt() == u.path.End()
}

func (u *Unit) SetPath(path *pathing.Path) {
    u.Lock()
    defer u.Unlock()

    u.path = path

    if !u.path.On(u.Loc.ToInt()) {
        u.return_path = u.path.RouteTo(u.Loc.ToInt())
    }
}
