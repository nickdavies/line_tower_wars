package layer

import (
    "github.com/neagix/Go-SDL/sdl"
)

type entityLayer struct {
    layerBase

    // no offsets
    voidOffsets

    // no setup
    voidSetup

    // no events
    voidEvents
}

func init() {
    registerLayer("entity", func(base layerBase, cfg interface{}) Layer {
        l := &entityLayer{
            layerBase: base,
        }

        l.voidOffsets = voidOffsets{&l.layerBase}
        l.voidSetup = voidSetup{&l.layerBase}
        l.voidEvents = voidEvents{&l.layerBase}

        return l
    })
}

func (g *entityLayer) Update(deltaTime int64) {
    if g.child != nil {
        g.child.Update(deltaTime)
    }
}

func (g *entityLayer) Render(target *sdl.Surface) {

    var units = make(map[string]*sdl.Surface)

    turret_texture := g.texture_map.GetName("turret_basic").Surface

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
            unit_name := u.Type.Name
            _, ok := units[unit_name]
            if !ok {
                units[unit_name] = g.texture_map.GetName(unit_name).Surface
            }

            loc := u.Loc
            target.Blit(
                &sdl.Rect{
                    X: int16(loc.Col * float64(g.square_size)) - 32,
                    Y: int16(loc.Row * float64(g.square_size)) - 32,
                },
                units[unit_name],
                nil,
            )

            health := 64 * (float64(u.Health) / float64(u.Type.Health))

            target.FillRect(
                &sdl.Rect{
                    X: int16(loc.Col * float64(g.square_size)) - 32,
                    Y: int16(loc.Row * float64(g.square_size)) - 32,
                    W: 5,
                    H: 64,
                },
                0x000000,
            )

            if health > 0 {
                target.FillRect(
                    &sdl.Rect{
                        X: int16(loc.Col * float64(g.square_size)) - 32,
                        Y: int16(loc.Row * float64(g.square_size)) - 32,
                        W: 5,
                        H: uint16(health),
                    },
                    0x00ff00,
                )
            }

        }
    }

    if g.child != nil {
        g.child.Render(target)
    }

}

func (g *entityLayer) GetSize() (uint16, uint16) {
    return g.parent.GetSize()
}


