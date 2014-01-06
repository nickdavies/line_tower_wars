package layer

import (
    "fmt"
)

import (
    "github.com/neagix/Go-SDL/sdl"
)

import (
    "github.com/nickdavies/line_tower_wars/game/stage"
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

    s *stage.Stage
    square_size uint16

    texture_map texture.TextureMap
    terrain_textures map[terrain.Terrain]*texture.Texture

    surface *sdl.Surface
}

func NewStageLayer(s *stage.Stage, texture_map texture.TextureMap, square_size uint16, child Layer) Layer {

    sg := &stageLayer{
        layerBase: layerBase{child: child},

        size_x: square_size * uint16(len(s.Tiles)),
        size_y: square_size * uint16(len(s.Tiles[0])),

        s: s,
        square_size: square_size,

        texture_map: texture_map,
        terrain_textures: map[terrain.Terrain]*texture.Texture{
            terrain.T_Grass: texture_map.GetName("grass_center"),
            terrain.T_Wall:  texture_map.GetName("wall_center"),
            terrain.T_Shadow:  texture_map.GetName("shadow_center"),
            terrain.T_Spawn: texture_map.GetName("spawn_center"),
            terrain.T_Goal:  texture_map.GetName("goal_center"),
        },
    }

    sg.voidEvents = voidEvents{&sg.layerBase}
    sg.voidOffsets = voidOffsets{&sg.layerBase}

    if child != nil {
        child.setParent(sg)
    }

    return sg
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

    for row = 0; row < uint16(len(g.s.Tiles[0])); row++ {
        for col = 0; col < uint16(len(g.s.Tiles)); col++ {
            texture := g.terrain_textures[g.s.Tiles[col][row]]

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

