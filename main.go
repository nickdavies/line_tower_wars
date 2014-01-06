package main

import (
    "os"
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

    go func() {
        <-time.After(1 * time.Minute)
        os.Exit(0)
    }()

    gfx_config := graphics.GraphicsConfig{
        ScreenX: 2560,
        ScreenY: 1600,

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

    g, controls := game.NewGame(game.GameConfig{}, players)

    for i := 0; i < players; i++ {
        go func (p_id int) {
            for {
                controls[p_id].Tick()
            }
        }(i)
    }

    gfx, err := graphics.NewGraphics(gfx_config, g, 0)
    if err != nil {
        panic(err)
    }

    g.Lock()
    go func() {
        err := gfx.Run()
        if err != nil {
            panic(err)
        }
    }()


    go func () {
        <-time.After(500 * time.Millisecond)
        p := g.GetPlayer(0)
        u := &unit.Unit{
            Loc: p.Path.Startf(),
            Speed: 5,
        }
        u.SetPath(p.Path)
        p.SpawnUnit(u)
    }()

    winner := g.Run()
    fmt.Println("Winner = ", winner)
}

