package layer

import (
    "github.com/neagix/Go-SDL/sdl"
)

import (
    "github.com/nickdavies/line_tower_wars/player"
    "github.com/nickdavies/line_tower_wars/towers"
    "github.com/nickdavies/line_tower_wars/texture"
    "github.com/nickdavies/line_tower_wars/util"
)

type entityLayer struct {
    layerBase

    // no run loop
    voidRun

    // no offsets
    voidOffsets

    // proxy ends
    bubbleEnd

    square_size uint16

    buildable bool
    player *player.Player

    texture_map texture.TextureMap
    buildSurface *sdl.Surface
}

func NewEntityLayer(p *player.Player, texture_map texture.TextureMap, square_size uint16, child Layer) Layer {
    eg := &entityLayer{
        layerBase: layerBase{child: child},

        texture_map: texture_map,
        square_size: square_size,

        player: p,
    }

    eg.voidRun = voidRun{&eg.layerBase}
    eg.voidOffsets = voidOffsets{&eg.layerBase}
    eg.bubbleEnd = bubbleEnd{&eg.layerBase}

    if child != nil {
        child.setParent(eg)
    }

    return eg
}

func (g *entityLayer) Setup() (err error) {
    g.buildSurface, err = util.CreateSurface(true, true, g.square_size * towers.SIZE_ROW, g.square_size * towers.SIZE_COL)
    if err != nil {
        return err
    }

    g.buildSurface.SetAlpha(sdl.SRCALPHA, 128)
    if g.child != nil {
        return g.child.Setup()
    }

    return nil
}

func (g *entityLayer) Cleanup() {
    g.buildSurface.Free()

    if g.child != nil {
        g.child.Cleanup()
    }
}

func (g *entityLayer) HandleEvent(event interface{}) {
    switch event.(type) {
    case sdl.MouseButtonEvent:
        e := event.(sdl.MouseButtonEvent)
        if e.Type == sdl.MOUSEBUTTONDOWN {
            x_off, y_off := g.GetXYOffsets()
            g.player.BuildTower( (e.Y + y_off) / g.square_size, (e.X + x_off) / g.square_size)
        }
    default:
    }

    if g.child != nil {
        g.child.HandleEvent(event)
    }
}

func (g *entityLayer) Update(deltaTime int64) {

    if g.child != nil {
        g.child.Update(deltaTime)
    }
}

func (g *entityLayer) Render(target *sdl.Surface) {

    x_off, y_off := g.GetXYOffsets()
    square_x, square_y := util.GetMouse(g.square_size, x_off, y_off, false)

    var build_colour uint32
    if g.player.Buildable(square_y / g.square_size, square_x / g.square_size) {
        build_colour = 0x00ff00
    } else {
        build_colour = 0xff0000
    }

    for loc, _ := range g.player.Towers {
        target.Blit(
            &sdl.Rect{
                X: int16(loc.Col * g.square_size),
                Y: int16(loc.Row * g.square_size),
            },
            g.texture_map.GetName("turret_basic").Surface,
            nil,
        )
    }

    g.buildSurface.FillRect(
        &sdl.Rect{
            X: 0,
            Y: 0,
            H: 128,
            W: 128,
        },
        build_colour,
    )

    g.buildSurface.Blit(
        &sdl.Rect{
            X: 0,
            Y: 0,
        },
        g.texture_map.GetName("turret_basic").Surface,
        nil,
    )

    target.Blit(
        &sdl.Rect{
            X: int16(square_x),
            Y: int16(square_y),
        },
        g.buildSurface,
        nil,
    )

    if g.child != nil {
        g.child.Render(target)
    }

}

func (g *entityLayer) GetSize() (uint16, uint16) {
    return g.parent.GetSize()
}


