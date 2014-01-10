package entity

import (
    "fmt"
)

import (
    "code.google.com/p/gcfg"
)

import (
    "github.com/nickdavies/line_tower_wars/game/unit"
    "github.com/nickdavies/line_tower_wars/game/tower"
)

type NoSuchEntityErr struct {
    entity string
}

func (e NoSuchEntityErr) Error() string {
    return fmt.Sprintf("No entity named %s exists", e.entity)
}

type EntityFactory interface {
    GetTower(name string) (tower.TowerType, error)
    GetUnit(name string) (unit.UnitType, error)

    ListTowers() ([]string)
    ListUnits() ([]string)
}
type entityFactoryStruct struct {
    Unit map[string]*unit.UnitType
    Tower map[string]*tower.TowerType
}

func NewEntityFactory (cfgFile string) (EntityFactory, error) {
    ef := &entityFactoryStruct{
        Unit: make(map[string]*unit.UnitType),
        Tower: make(map[string]*tower.TowerType),
    }

    err := gcfg.ReadFileInto(ef, cfgFile)
    if err != nil {
        return nil, err
    }

    for name, u := range ef.Unit {
        if name != u.Name {
            return nil, fmt.Errorf("Unit is listed as '%s' but name = '%s'", name, u.Name)
        }
    }

    for name, t := range ef.Tower {
        if name != t.Name {
            return nil, fmt.Errorf("Tower is listed as '%s' but name = '%s'", name, t.Name)
        }
    }

    return ef, nil
}

func (ef *entityFactoryStruct) GetTower(name string) (tower.TowerType, error) {
    t, ok := ef.Tower[name]
    if !ok {
        return tower.TowerType{}, NoSuchEntityErr{name}
    }
    return *t, nil
}

func (ef *entityFactoryStruct) ListTowers() ([]string) {
    towers := make([]string, len(ef.Tower))

    i := 0
    for name, _ := range ef.Tower {
        towers[i] = name
        i++
    }

    return towers
}

func (ef *entityFactoryStruct) GetUnit(name string) (unit.UnitType, error) {
    u, ok := ef.Unit[name]
    if !ok {
        return unit.UnitType{}, NoSuchEntityErr{name}
    }
    return *u, nil
}
func (ef *entityFactoryStruct) ListUnits() ([]string) {
    units := make([]string, len(ef.Unit))

    i := 0
    for name, _ := range ef.Unit {
        units[i] = name
        i++
    }

    return units
}

