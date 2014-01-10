package graphics

import (
    "os"
    "fmt"
)

import (
    "github.com/neagix/Go-SDL/sdl"
)

import (
    "github.com/nickdavies/line_tower_wars/game"

    "github.com/nickdavies/line_tower_wars/graphics/config"
    "github.com/nickdavies/line_tower_wars/graphics/layer"
    "github.com/nickdavies/line_tower_wars/graphics/texture"
)

type Graphics struct {
    cfg config.GraphicsConfig

    g *game.Game

    display *sdl.Surface

    layers map[string]layer.Layer
    topLayer layer.Layer
}

func NewGraphics(cfg config.GraphicsConfig, g *game.Game, controls game.PlayerControls) (*Graphics, error) {

    gfx := &Graphics{
        cfg: cfg,
        g: g,

        layers: make(map[string]layer.Layer),
    }

    err := gfx.setup()
    if err != nil {
        return nil, err
    }

    texture_map, err := texture.NewTextureMap(cfg.SquareSize, cfg.TextureDir)
    if err != nil {
        return nil, err
    }

    layers := []string{"entity", "build", "stage", "pan", "gui", "sdl"}

    layer_cfg := map[string]interface{}{
        "entity": nil,
        "build": nil,
        "stage": nil,
        "pan": cfg.PanningOptions,
        "gui": layer.GuiLayerCfg{cfg.FontFile, cfg.FontSize},
        "sdl": layer.SdlLayerCfg{gfx.display, cfg.ScreenX, cfg.ScreenY},
    }

    layer_disabled := make(map[string]bool)

    if controls == nil {
        layer_disabled["build"] = true
    }


    builder := layer.NewLayerBuilder(g, controls, texture_map, cfg.SquareSize)

    var prev layer.Layer
    for _, name := range layers {
        if !layer_disabled[name] {
            prev = builder.Build(name, layer_cfg[name], prev)
        }
    }
    gfx.topLayer = prev

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

func (gfx *Graphics) Run(tick_ch chan int64) error {
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
                if ok && kbe.Keysym.Sym == sdl.K_ESCAPE {
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

    gfx.g.Unlock()
    for delta := range tick_ch {
        // Update State
        gfx.topLayer.Update(delta)

        // Update Screen
        gfx.topLayer.Render(nil)

        <-tick_ch
    }

    return nil
}
