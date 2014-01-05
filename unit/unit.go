package unit

import (
    "time"
    "sync"
)

import (
    "github.com/nickdavies/line_tower_wars/pathing"
)

type Unit struct {
    sync.Mutex

    Loc pathing.Locf

    Speed float64

    path *pathing.Path

    return_path *pathing.Path
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
            u.Loc, end = u.return_path.Move(u.Loc, float64(timeDelta) * u.Speed / float64(time.Second))
            if end {
                u.return_path = nil
            }
            return
        }
    }

    u.Loc, end = u.path.Move(u.Loc, float64(timeDelta) * u.Speed / float64(time.Second))

    if end {
        u.Loc = u.path.Start()
    }
}

func (u *Unit) SetPath(path *pathing.Path) {
    u.Lock()
    defer u.Unlock()

    u.path = path

    if !u.path.On(u.Loc.ToInt()) {
        u.return_path = u.path.RouteTo(u.Loc.ToInt())
    }
}
