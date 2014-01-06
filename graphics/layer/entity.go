package layer

import (
    "github.com/neagix/Go-SDL/sdl"
)

import (
    "github.com/nickdavies/line_tower_wars/game"

    "github.com/nickdavies/line_tower_wars/graphics/texture"
)

type entityLayer struct {
    layerBase

    // no offsets
    voidOffsets

    // no setup
    voidSetup

    // no events
    voidEvents

    square_size uint16

    texture_map texture.TextureMap
}

func NewEntityLayer(texture_map texture.TextureMap, square_size uint16, child Layer, g *game.Game) Layer {
    eg := &entityLayer{
        layerBase: layerBase{child: child, game: g},

        texture_map: texture_map,
        square_size: square_size,
    }

    eg.voidOffsets = voidOffsets{&eg.layerBase}
    eg.voidSetup = voidSetup{&eg.layerBase}
    eg.voidEvents = voidEvents{&eg.layerBase}

    if child != nil {
        child.setParent(eg)
    }

    return eg
}

func (g *entityLayer) Update(deltaTime int64) {
    if g.child != nil {
        g.child.Update(deltaTime)
    }
}

func (g *entityLayer) Render(target *sdl.Surface) {

    turret_texture := g.texture_map.GetName("turret_basic").Surface
    unit_texture := g.texture_map.GetName("unit_basic").Surface

    for i := 0; i < g.game.NumPlayers; i++ {
        player := g.game.GetPlayer(i)

        for loc, _ := range player.Towers {
            target.Blit(
                &sdl.Rect{
                    X: int16(loc.Col * g.square_size),
                    Y: int16(loc.Row * g.square_size),
                },
                turret_texture,
                nil,
            )
        }

        for _, u := range player.Units {
            loc := u.Loc
            target.Blit(
                &sdl.Rect{
                    X: int16(loc.Col * float64(g.square_size)) - 32,
                    Y: int16(loc.Row * float64(g.square_size)) - 32,
                },
                unit_texture,
                nil,
            )

        }
    }

    if g.child != nil {
        g.child.Render(target)
    }

}

func (g *entityLayer) GetSize() (uint16, uint16) {
    return g.parent.GetSize()
}


