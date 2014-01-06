package main

import (
    "time"
    "math/rand"
)

import (
    "github.com/nickdavies/line_tower_wars/game"
    "github.com/nickdavies/line_tower_wars/graphics"
)

func main() {
    var seed int64 = time.Now().UnixNano()
    seed = 0
    rand.Seed(seed)

    var players int = 2
    var square_size uint16 = 128
    var screen_x uint16 = 2560
    var screen_y uint16 = 1600
    var texture_dir string = "./gfx/textures"

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

    g := game.NewGame(game.GameConfig{}, players)

    gfx, err := graphics.NewGraphics(gfx_config, g, -1)
    if err != nil {
        panic(err)
    }

    err = g.Run()
    if err != nil {
        panic(err)
    }
}

