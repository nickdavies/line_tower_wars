package player

import (
    "github.com/nickdavies/line_tower_wars/stage"
    "github.com/nickdavies/line_tower_wars/towers"
)

type loc struct {
    Row uint16
    Col uint16
}

type Player struct {
    myStage *stage.PlayerStage
    Towers map[loc]*towers.Tower
}

func NewPlayer(myStage *stage.PlayerStage) *Player {
    return &Player{
        myStage: myStage,
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

func (t *Player) BuildTower(row, col uint16) (user_message string) {
    t.Towers[loc{row, col}] = &towers.Tower{}
    return ""
}
