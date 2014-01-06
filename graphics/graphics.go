package graphics

import (
    "os"
    "runtime"
    "time"
    "fmt"
)

import (
    "github.com/neagix/Go-SDL/sdl"
)

import (
    "github.com/nickdavies/line_tower_wars/game"

    "github.com/nickdavies/line_tower_wars/graphics/layer"
    "github.com/nickdavies/line_tower_wars/graphics/texture"
)

type Graphics struct {
    cfg GraphicsConfig

    g *game.Game

    display *sdl.Surface

    layers map[string]layer.Layer
    topLayer layer.Layer
}

func NewGraphics(cfg GraphicsConfig, g *game.Game, perspective int) (*Graphics, error) {

    gfx := &Graphics{
        cfg: cfg,
        g: g,

        layers: make(map[string]layer.Layer),
    }

    err := gfx.setup()
    if err != nil {
        return nil, err
    }

    tm, err := texture.NewTextureMap(cfg.SquareSize, cfg.TextureDir)
    if err != nil {
        return nil, err
    }

    gfx.layers["entity"] = layer.NewEntityLayer(
        tm,
        cfg.SquareSize,
        nil,
        g,
    )
    next_layer := gfx.layers["entity"]

    if perspective != -1 {
        gfx.layers["build"] = layer.NewBuildLayer(
            g.GetPlayer(perspective),
            tm,
            cfg.SquareSize,
            gfx.layers["entity"],
        )
        next_layer = gfx.layers["build"]
    }

    gfx.layers["stage"] = layer.NewStageLayer(
        g.GetStage(),
        tm,
        cfg.SquareSize,
        next_layer,
    )

    gfx.layers["pan"] = layer.NewPanLayer(
        cfg.PanningOptions.PanXSize,
        cfg.PanningOptions.PanYSize,

        cfg.PanningOptions.StartingX,
        cfg.PanningOptions.StartingY,

        cfg.PanningOptions.PanSpeed,

        gfx.layers["stage"],
    )
    gfx.layers["main"] = layer.NewSdlLayer(
        gfx.display,
        cfg.ScreenX,
        cfg.ScreenY,
        gfx.layers["pan"],
    )

    gfx.topLayer = gfx.layers["main"]

    return gfx, nil
}

func (gfx *Graphics) setup() error {
    var errno = sdl.Init(sdl.INIT_EVERYTHING)
    if errno != 0 {
        return fmt.Errorf("Init failed: %s", sdl.GetError())
    }

    gfx.display = sdl.SetVideoMode(int(gfx.cfg.ScreenX), int(gfx.cfg.ScreenY), 32, sdl.HWSURFACE | sdl.DOUBLEBUF | sdl.FULLSCREEN)
    if gfx.display == nil {
        return fmt.Errorf("No surface created: %s", sdl.GetError())
    }

    return nil
}

func (gfx *Graphics) Cleanup() {
    sdl.Quit()
}

func (gfx *Graphics) Run() error {
    var err error
    var end_ch = make(chan interface{})

    err = gfx.topLayer.Setup()
    if err != nil {
        return err
    }
    defer gfx.topLayer.Cleanup()

    go func() {
        for {
            select {
            case <-end_ch:
                return
            case event, ok := <-sdl.Events:
                if !ok {
                    gfx.g.End()
                    //TODO: fix this
                    os.Exit(0)
                    continue
                }
                kbe, ok := event.(sdl.KeyboardEvent)
                if ok && kbe.Keysym.Sym == sdl.K_F1 {
                    gfx.g.End()
                    //TODO: fix this
                    os.Exit(0)
                    continue
                }
                gfx.topLayer.HandleEvent(event)
            }
        }
    }()
    defer func() {
        end_ch <-nil
    }()

    var last_time = time.Now().UnixNano()

    gfx.g.Unlock()
    for gfx.g.Running() {
        // Update State
        gfx.topLayer.Update(time.Now().UnixNano() - last_time)
        last_time = time.Now().UnixNano()

        // Update Screen
        gfx.topLayer.Render(nil)
        runtime.Gosched()
        <-time.After(10 * time.Millisecond)
    }

    return nil
}
