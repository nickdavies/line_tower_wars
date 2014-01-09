package ai

import (
    "errors"
)

import (
    "github.com/nickdavies/line_tower_wars/game"
    "github.com/nickdavies/line_tower_wars/game/stage"
)

type towerBuilder struct {
    row int16
    col int16
    assending bool
}

var AtEndErr = errors.New("Cant make more towers already at the end")

func (tb *towerBuilder) Build(control game.PlayerControls, name string) error {
    row := tb.row
    col := tb.col
    assending := tb.assending

    if int(row) >= stage.Grass_Rows + stage.Spawn_Size {
        return AtEndErr
    }

    if assending {
        col++
    } else {
        col--
    }

    // if were at the start and want to goto the next
    // you must move down
    if col == 0 && !assending {
        row += 2
        col -= 1

        tb.row = row
        tb.col = col
        tb.assending = !assending

        err := tb.Build(control, name)
        if err != nil {
            return err
        }

        return nil
    }

    // if were at the start and want to goto the next
    // you must move down
    if int(col) == stage.Grass_Cols - 1 && assending {
        row += 2
        col += 1

        tb.row = row
        tb.col = col
        tb.assending = !assending

        err := tb.Build(control, name)
        if err != nil {
            return err
        }

        return nil
    }

    tb.row = row
    tb.col = col

    err := control.BuyTower(name, uint16(row), uint16(col))
    if err != nil {
        return err
    }

    return nil
}

func RunAttackDefenceRatioAI(control game.PlayerControls, attackPercent float32, graceTime int64) error {
    money := control.GetBalance()
    entities := control.GetEntityFactory()

    tower, err := entities.GetTower("basic_tower")
    if err != nil {
        return err
    }

    unit, err := entities.GetUnit("basic_unit")
    if err != nil {
        return err
    }

    tb := towerBuilder{
        row: stage.Spawn_Size,
        col: -1,
        assending: true,
    }

    var total_spend int64
    var unit_spend int64

    var elapsed int64 = -1
    var at_end = false
    for {
        delta, running := control.Tick()
        if !running {
            return nil
        }
        elapsed += delta

        if elapsed < graceTime {
            continue
        }

        for {
            if float32(unit_spend) / float32(total_spend) > attackPercent && !at_end {
                if tower.Cost > money.Get() {
                    break
                }

                err := tb.Build(control, "basic_tower")
                if err == AtEndErr {
                    at_end = true
                } else if err != nil {
                    panic(err)
                }

                total_spend += int64(tower.Cost)
            } else {
                if unit.Cost > money.Get() {
                    break
                }

                err := control.BuyUnit("basic_unit")
                if err != nil {
                    panic(err)
                }

                total_spend += int64(unit.Cost)
                unit_spend += int64(unit.Cost)
            }
        }
        if !at_end && tower.Cost < money.Get() {
            err := tb.Build(control, "basic_tower")
            if err == AtEndErr {
                at_end = true
            } else if err != nil {
                panic(err)
            }
        }

    }
}
