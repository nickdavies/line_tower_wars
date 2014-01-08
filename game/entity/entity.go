package entity

import (
    "errors"
)

import (
    "code.google.com/p/gcfg"
)

import (
    "github.com/nickdavies/line_tower_wars/game/unit"
    "github.com/nickdavies/line_tower_wars/game/tower"
)

var NoSuchEntityErr = error

type UnitFactory struct {
    items map[string]unit.UnitType
}

func NewUnitFactory (cfgFile string) (*UnitFactory, error) {
    uf := &UnitFactory{
        items: make(map[string]unit.UnitType)
    }

    err := gcfg.ReadFileInto(uf, cfgFile)
    if err != nil {
        return nil, err
    }

    return uf, nil
}

func (uf *UnitFactory) Get(name string) (unit.UnitType, error) {

    u, ok := uf.items[name]
    if !ok {
        return unit.UnitType{}, NoSuchEntityErr
    }
    return u, nil
}

type TowerFactory struct {
    items map[string]tower.TowerType
}

func NewTowerFactory (cfgFile string) (*TowerFactory, error) {
    uf := &TowerFactory{
        items: make(map[string]tower.TowerType)
    }

    err := gcfg.ReadFileInto(uf, cfgFile)
    if err != nil {
        return nil, err
    }

    return uf, nil
}

func (uf *TowerFactory) Get(name string) (tower.TowerType, error) {

    u, ok := uf.items[name]
    if !ok {
        return tower.TowerType{}, NoSuchEntityErr
    }
    return u, nil
}
