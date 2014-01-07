package layer

import (
    "github.com/neagix/Go-SDL/sdl"
)

import (
    "github.com/nickdavies/line_tower_wars/game/player"
    "github.com/nickdavies/line_tower_wars/game/tower"
    "github.com/nickdavies/line_tower_wars/game/unit"

    "github.com/nickdavies/line_tower_wars/graphics/texture"
    "github.com/nickdavies/line_tower_wars/graphics/util"
)

type buildLayer struct {
    layerBase

    // no offsets
    voidOffsets

    square_size uint16

    buildable bool

    texture_map texture.TextureMap
    buildSurface *sdl.Surface

    player *player.Player
}

func NewBuildLayer(player *player.Player, texture_map texture.TextureMap, square_size uint16, child Layer) Layer {
    eg := &buildLayer{
        layerBase: layerBase{child: child},

        texture_map: texture_map,
        square_size: square_size,

        player: player,
    }

    eg.voidOffsets = voidOffsets{&eg.layerBase}

    if child != nil {
        child.setParent(eg)
    }

    return eg
}

func (g *buildLayer) Setup() (err error) {
    // TODO: make tower size come in properly
    g.buildSurface, err = util.CreateSurface(true, true, 1 * g.square_size, 1 * g.square_size)
    if err != nil {
        return err
    }

    g.buildSurface.SetAlpha(sdl.SRCALPHA, 128)
    if g.child != nil {
        return g.child.Setup()
    }

    return nil
}

func (g *buildLayer) Cleanup() {
    g.buildSurface.Free()

    if g.child != nil {
        g.child.Cleanup()
    }
}

func (g *buildLayer) HandleEvent(event interface{}) {
    switch event.(type) {
    case sdl.MouseButtonEvent:
        e := event.(sdl.MouseButtonEvent)
        if e.Type == sdl.MOUSEBUTTONDOWN {
            x_off, y_off := g.GetXYOffsets()

            base_tower := &tower.TowerType{
                Id: 1,
                Name: "basic_tower",

                Range: 10,
                FireRate: 2,
                Damage: 10,

                Health: 100,

                Cost: 5,
            }

            g.player.BuildTower(base_tower, (e.Y + y_off) / g.square_size, (e.X + x_off) / g.square_size, false)
        }
    case sdl.KeyboardEvent:
        e := event.(sdl.KeyboardEvent)
        if e.Type == sdl.KEYDOWN && e.Keysym.Sym == sdl.K_RETURN {
            base_unit := &unit.UnitType{
                Id: 1,
                Name: "basic_mob",

                Speed: 2,
                IncomeDelta: 1,

                Health: 5,
                Cost: 5,
            }

            g.player.BuyUnit(base_unit, true)
        }
    default:
    }

    if g.child != nil {
        g.child.HandleEvent(event)
    }
}

func (g *buildLayer) Update(deltaTime int64) {
    if g.child != nil {
        g.child.Update(deltaTime)
    }
}

func (g *buildLayer) Render(target *sdl.Surface) {
    x_off, y_off := g.GetXYOffsets()
    square_x, square_y := util.GetMouse(g.square_size, x_off, y_off, false)

    var build_colour uint32
    if g.player.Buildable(square_y / g.square_size, square_x / g.square_size) {
        build_colour = 0x00ff00
    } else {
        build_colour = 0xff0000
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

func (g *buildLayer) GetSize() (uint16, uint16) {
    return g.parent.GetSize()
}


