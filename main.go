package main

import (
//    "os"
    "fmt"
    "time"
    "math/rand"
)

import (
    "github.com/nickdavies/line_tower_wars/game/unit"

    "github.com/nickdavies/line_tower_wars/game"
    "github.com/nickdavies/line_tower_wars/graphics"
)

func main() {
    var seed int64 = time.Now().UnixNano()
    seed = 0
    rand.Seed(seed)

    var players int = 2

    /*
    go func() {
        <-time.After(1 * time.Minute)
        os.Exit(0)
    }()
    */

    game_cfg := game.GameConfig{
        MoneyConfig: game.MoneyConfig{
            Balance: 120,
            Income: 5,
            MinIncome: 1,
            IncomeInterval: 5,
        },
    }

    gfx_config := graphics.GraphicsConfig{
        ScreenX: 1920,
        ScreenY: 1080,

        SquareSize: 128,

        TextureDir: "./gfx/textures",

        PanningOptions: graphics.PanningOptions{
            PanXSize: 90,
            PanYSize: 90,

            StartingX: 0,
            StartingY: 0,

            PanSpeed: 3000,
        },
    }

    g, controls := game.NewGame(game_cfg, players)

    // TODO: make graphics use this
    go func () {
        for {
            controls[0].Tick()
        }
    }()

    go func () {
        gap := 0
        min_gap := 10

        player := controls[1].GetPlayer()
        base_unit := &unit.UnitType{
            Id: 1,
            Name: "basic_mob",

            Speed: 1,
            IncomeDelta: 1,

            Health: 25,
            Cost: 5,
        }

        var start int64 = -1
        for {
            controls[1].Tick()

            if start == -1 {
                start = time.Now().UnixNano()
            }

            if time.Now().UnixNano() - start < int64(time.Second * 5) {
                continue
            }
            gap += 1

            if gap >= min_gap && player.Money.Get() > base_unit.Cost {
                player.BuyUnit(base_unit, false)
                gap = 0
            }
        }
    }()

    gfx, err := graphics.NewGraphics(gfx_config, g, 0)
    if err != nil {
        panic(err)
    }

    var gfx_tick = make(chan int64)
    var game_tick = make(chan int64)

    g.Lock()
    go func() {
        err := gfx.Run(gfx_tick)
        if err != nil {
            panic(err)
        }
    }()

    go func() {
        winner := g.Run(game_tick)
        fmt.Println("Winner = ", winner)
    }()

    var gfx_last int64 = time.Now().UnixNano()
    var game_last int64 = time.Now().UnixNano()

    var tick_count int64
    var tick_start int64
    _ = tick_start
    for {
        now := time.Now().UnixNano()
        tick_start = now

        game_tick<-(now - game_last)
        game_tick<-0
        game_last = now

        now = time.Now().UnixNano()
        gfx_tick<-(now - gfx_last)
        gfx_tick<-0
        gfx_last = now

        //fmt.Println("tick", tick_count, (time.Now().UnixNano() - tick_start) / int64(time.Millisecond))
        tick_count++
    }
}

