package main

import (
    "fmt"
    "time"
    "math/rand"
)

import (
    "github.com/neagix/Go-SDL/sdl"
)

import (
    "github.com/nickdavies/line_tower_wars/layer"
    "github.com/nickdavies/line_tower_wars/stage"
    "github.com/nickdavies/line_tower_wars/player"
    "github.com/nickdavies/line_tower_wars/unit"
    "github.com/nickdavies/line_tower_wars/texture"
)

func setup(screen_x, screen_y uint16) (*sdl.Surface, error){
    var errno = sdl.Init(sdl.INIT_EVERYTHING)
    if errno != 0 {
        return nil, fmt.Errorf("Init failed: %s", sdl.GetError())
    }

    display := sdl.SetVideoMode(int(screen_x), int(screen_y), 32, sdl.HWSURFACE | sdl.DOUBLEBUF | sdl.FULLSCREEN)
    if display == nil {
        return nil, fmt.Errorf("No surface created: %s", sdl.GetError())
    }

    return display, nil
}

func cleanup(display *sdl.Surface) {
    sdl.Quit()
}

func main() {
    var seed int64 = time.Now().UnixNano()
    seed = 0
    rand.Seed(seed)

    var players int = 2
    var square_size uint16 = 128
    var screen_x uint16 = 2560
    var screen_y uint16 = 1600
    var texture_dir string = "./gfx/textures"

    s := stage.NewStage(players)
    me := player.NewPlayer(s.GetPlayer(0))

    var left uint16 = 0
    var row uint16
    var col uint16

    for row = 1; row < stage.PlayerStageHeight; row += 2 {
        for col = 0; col < 10; col++ {
            if rand.Intn(100) < 90 {
                me.BuildTower(row, col + stage.Shadow_Side + stage.Wall_Side + left, true)
            }
        }
        if left == 0 {
            left = 1
        } else {
            left = 0
        }
    }

    me.Repath()
    me.AntiCheat = &unit.Unit{Loc: me.Path.Start(), Speed: 1}
    me.AntiCheat.SetPath(me.Path)

    display, err := setup(screen_x, screen_y)
    if err != nil {
        panic(err)
    }

    tm, err := texture.NewTextureMap(square_size, texture_dir)
    if err != nil {
        panic(err)
    }

    eg := layer.NewEntityLayer(me, tm, square_size, nil)
    sg := layer.NewStageLayer(s, tm, square_size, eg)
    pg := layer.NewPanLayer(90, 90, 0, 0, 3000, sg)
    g := layer.NewSdlLayer(display, screen_x, screen_y, pg)

    err = g.Run()
    if err != nil {
        panic(err)
    }
}

