package layer

import (
    "fmt"
)

import (
    "github.com/neagix/Go-SDL/sdl"
    "github.com/neagix/Go-SDL/ttf"
)

import (
    "github.com/nickdavies/line_tower_wars/game"
)

type guiLayer struct {
    layerBase

    // no offsets
    voidOffsets

    // no events
    voidEvents

    font *ttf.Font
}

func NewGuiLayer(game *game.Game, child Layer) Layer {
    gl := &guiLayer{
        layerBase: layerBase{child: child, game: game},
    }

    gl.voidOffsets = voidOffsets{&gl.layerBase}
    gl.voidEvents = voidEvents{&gl.layerBase}

    if child != nil {
        child.setParent(gl)
    }

    return gl
}

func (g *guiLayer) Setup() (err error) {
    errno := ttf.Init()
    if errno != 0 {
        return fmt.Errorf("ttf.Init failed: %s", sdl.GetError())
    }

    // TODO: fix this
    g.font = ttf.OpenFont("./gfx/DejaVuSansMono.ttf", 12)
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

