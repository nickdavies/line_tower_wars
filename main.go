package main

import (
    "fmt"
    "time"
    "math/rand"
)

import (
    "github.com/nickdavies/line_tower_wars/ai"
    "github.com/nickdavies/line_tower_wars/game"
    "github.com/nickdavies/line_tower_wars/graphics"
)

func main() {
    var use_graphics = true
    var seed int64 = time.Now().UnixNano()
    seed = 0
    rand.Seed(seed)

    var players int = 2

    game_cfg := game.GameConfig{
        MoneyConfig: game.MoneyConfig{
            Balance: 120,
            Income: 5,
            MinIncome: 1,
            IncomeInterval: 2,
        },
        EntityDir: "./gfx/entity.cfg",
    }

    gfx_config := graphics.GraphicsConfig{
        ScreenX: 1920,
        ScreenY: 1080,

        SquareSize: 128,

        TextureDir: "./gfx/textures",

        PanningOptions: graphics.PanningOptions{
            PanXSize: 90,
            PanYSize: 90,

            StartingX: 300,
            StartingY: 0,

            PanSpeed: 3000,
        },
    }

    g, controls, err := game.NewGame(game_cfg, players)
    if err != nil {
        panic(err)
    }

    /*
    go func (){
        for {
            controls[0].Tick()
        }
    }()
    */
    go func () {
        err := ai.RunAttackDefenceRatioAI(controls[0], 0.5, int64(time.Second * 2))
        panic(err)
    }()
    go func () {
        err := ai.RunAttackDefenceRatioAI(controls[1], 0.5, int64(time.Second * 2))
        panic(err)
    }()

    if use_graphics {
        g.Lock()
    }

    var gfx_tick = make(chan int64)
    var game_tick = make(chan int64)

    go func() {
        winner := g.Run(game_tick)
        fmt.Println("Winner = ", winner)
    }()

    if use_graphics {
        gfx, err := graphics.NewGraphics(gfx_config, g, 0)
        if err != nil {
            panic(err)
        }
        go func() {
            err := gfx.Run(gfx_tick)
            if err != nil {
                panic(err)
            }
        }()
    }

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

        if use_graphics {
            now = time.Now().UnixNano()
            gfx_tick<-(now - gfx_last)
            gfx_tick<-0
            gfx_last = now
        }

        //fmt.Println("tick", tick_count, (time.Now().UnixNano() - tick_start) / int64(time.Millisecond))
        tick_count++
    }
}

