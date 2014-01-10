package layer

import (
    "fmt"
)

import (
    "github.com/neagix/Go-SDL/sdl"
)

import (
    "github.com/nickdavies/line_tower_wars/game/terrain"

    "github.com/nickdavies/line_tower_wars/graphics/texture"
    "github.com/nickdavies/line_tower_wars/graphics/util"
)

type stageLayer struct {
    layerBase

    // proxy events
    voidEvents

    // no offsets
    voidOffsets

    size_x uint16
    size_y uint16

    terrain_textures map[terrain.Terrain]*texture.Texture

    surface *sdl.Surface
}

func init() {
    registerLayer("stage", func(base layerBase, cfg interface{}) Layer {
        stage := base.game.GetStage()
        l := &stageLayer{
            layerBase: base,

            size_x: base.square_size * uint16(len(stage.Tiles)),
            size_y: base.square_size * uint16(len(stage.Tiles[0])),

            terrain_textures: map[terrain.Terrain]*texture.Texture{
                terrain.T_Grass: base.texture_map.GetName("grass_center"),
                terrain.T_Wall:  base.texture_map.GetName("wall_center"),
                terrain.T_Shadow:  base.texture_map.GetName("shadow_center"),
                terrain.T_Spawn: base.texture_map.GetName("spawn_center"),
                terrain.T_Goal:  base.texture_map.GetName("goal_center"),
            },
        }

        l.voidEvents = voidEvents{&l.layerBase}
        l.voidOffsets = voidOffsets{&l.layerBase}

        return l
    })
}

func (g *stageLayer) Setup() (err error) {
    g.surface, err = util.CreateSurface(true, true, g.size_x, g.size_y)
    if err != nil {
        return err
    }

    var row uint16
    var col uint16
    var y_tile uint16
    var x_tile uint16
    var y_tiles uint16
    var x_tiles uint16

    stage := g.game.GetStage()
    for row = 0; row < uint16(len(stage.Tiles[0])); row++ {
        for col = 0; col < uint16(len(stage.Tiles)); col++ {
            texture := g.terrain_textures[stage.Tiles[col][row]]

            x_tiles = g.square_size / texture.Width
            y_tiles = g.square_size / texture.Height

            if x_tiles == 0 || y_tiles == 0 {
                panic(fmt.Errorf("terrain (%s) texture is larger than square_size (%d, %d) > %d", texture.Name, texture.Width, texture.Height, g.square_size))
            }

            for y_tile = 0; y_tile < y_tiles; y_tile++ {
                for x_tile = 0; x_tile < x_tiles; x_tile++ {
                    errno := g.surface.Blit(
                        &sdl.Rect{
                            X: int16(col * g.square_size + x_tile * texture.Width),
                            Y: int16(row * g.square_size + y_tile * texture.Height),
                            W: texture.Width,
                            H: texture.Height,
                        },
                        texture.Surface,
                        &sdl.Rect{
                            X: 0,
                            Y: 0,
                            W: texture.Width,
                            H: texture.Height,
                        },
                    )

                    if errno != 0 {
                        return fmt.Errorf("Blit failed: %s", sdl.GetError())
                    }
                }
            }
        }
    }

    if g.child != nil {
        return g.child.Setup()
    }

    return nil
}

func (g *stageLayer) Cleanup() {
    g.surface.Free()
    if g.child != nil {
        g.child.Cleanup()
    }
}


func (g *stageLayer) Update(deltaTime int64) {
    if g.child != nil {
        g.child.Update(deltaTime)
    }
}

func (g *stageLayer) Render(target *sdl.Surface) {

    target.Blit(
        &sdl.Rect{
            X: 0,
            Y: 0,
            W: uint16(g.size_x),
            H: uint16(g.size_y),
        },
        g.surface,
        nil,
    )

    if g.child != nil {
        g.child.Render(target)
    }

}

func (g *stageLayer) GetSize() (uint16, uint16) {
    return g.size_x, g.size_y
}

