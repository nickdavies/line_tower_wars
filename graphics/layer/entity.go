package layer

import (
    "github.com/neagix/Go-SDL/sdl"
)

import (
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

func NewEntityLayer(texture_map texture.TextureMap, square_size uint16, child Layer) Layer {
    eg := &entityLayer{
        layerBase: layerBase{child: child},

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

    /*
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
    */

    if g.child != nil {
        g.child.Render(target)
    }

}

func (g *entityLayer) GetSize() (uint16, uint16) {
    return g.parent.GetSize()
}


