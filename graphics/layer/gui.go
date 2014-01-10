package layer

import (
    "fmt"
)

import (
    "github.com/neagix/Go-SDL/sdl"
    "github.com/neagix/Go-SDL/ttf"
)

type GuiLayerCfg struct {
    FontFile string
    FontSize int
}

type guiLayer struct {
    layerBase

    // no offsets
    voidOffsets

    // no events
    voidEvents

    fontFile string
    fontSize int
    font *ttf.Font
}

func init() {
    registerLayer("gui", func(base layerBase, cfg interface{}) Layer {
        gui_cfg := cfg.(GuiLayerCfg)

        l := &guiLayer{
            layerBase: base,
            fontFile: gui_cfg.FontFile,
            fontSize: gui_cfg.FontSize,
        }

        l.voidOffsets = voidOffsets{&l.layerBase}
        l.voidEvents = voidEvents{&l.layerBase}

        return l
    })
}

func (g *guiLayer) Setup() (err error) {
    errno := ttf.Init()
    if errno != 0 {
        return fmt.Errorf("ttf.Init failed: %s", sdl.GetError())
    }

    // TODO: fix this
    g.font = ttf.OpenFont(g.fontFile, g.fontSize)
    if g.font == nil {
        return fmt.Errorf("OpenFont failed: %s", sdl.GetError())
    }

    if g.child != nil {
        return g.child.Setup()
    }

    return nil
}

func (g *guiLayer) Cleanup() {
    ttf.Quit()
    g.font.Close()

    if g.child != nil {
        g.child.Cleanup()
    }
}

func (g *guiLayer) Update(deltaTime int64) {
    if g.child != nil {
        g.child.Update(deltaTime)
    }
}

func (g *guiLayer) Render(target *sdl.Surface) {
    if g.child != nil {
        g.child.Render(target)
    }

    target.FillRect(
        &sdl.Rect{
            X: 20,
            Y: 300,
            H: uint16(30 + (30 * g.game.NumPlayers)),
            W: 200,
        },
        0xffffff,
    )

    for i := 0; i < g.game.NumPlayers; i++ {
        p := g.game.GetPlayer(i)
        str := fmt.Sprintf("Player %d: %d - %d - %d", i, p.Money.Get(), p.Money.GetIncome(), p.GetLives())

        txt_surface := ttf.RenderText_Solid(g.font, str, sdl.Color{})
        if txt_surface == nil {
            fmt.Printf("RenderText Solid failed: %s", sdl.GetError())
        }

        target.Blit(
            &sdl.Rect{
                X: 50,
                Y: int16(320 + (i * 30)),
            },
            txt_surface,
            nil,
        )
    }
}

func (g *guiLayer) GetSize() (uint16, uint16) {
    return g.parent.GetSize()
}

