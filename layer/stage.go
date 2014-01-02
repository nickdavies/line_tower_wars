package layer

import (
    "fmt"
)

import (
    "github.com/neagix/Go-SDL/sdl"
)

import (
    "github.com/nickdavies/line_tower_wars/stage"
    "github.com/nickdavies/line_tower_wars/terrain"
    "github.com/nickdavies/line_tower_wars/texture"
    "github.com/nickdavies/line_tower_wars/util"
)

type stageLayer struct {
    layerBase

    // no run loop
    voidRun

    // proxy events
    voidEvents

    // no offsets
    voidOffsets

    // proxy ends
    bubbleEnd

    size_x uint16
    size_y uint16

    s *stage.Stage
    square_size int

    texture_map texture.TextureMap
    terrain_textures map[terrain.Terrain]*texture.Texture

    surface *sdl.Surface
}

func NewStageLayer(s *stage.Stage, texture_map texture.TextureMap, square_size int, child Layer) Layer {

    sg := &stageLayer{
        layerBase: layerBase{child: child},

        size_x: uint16(square_size * len(s.Tiles)),
        size_y: uint16(square_size * len(s.Tiles[0])),

        s: s,
        square_size: square_size,

        texture_map: texture_map,
        terrain_textures: map[terrain.Terrain]*texture.Texture{
            terrain.T_Grass: texture_map.GetName("grass_center"),
            terrain.T_Wall:  texture_map.GetName("wall_center"),
            terrain.T_Spawn: texture_map.GetName("spawn_center"),
            terrain.T_Goal:  texture_map.GetName("goal_center"),
        },
    }

    sg.voidRun = voidRun{&sg.layerBase}
    sg.voidEvents = voidEvents{&sg.layerBase}
    sg.voidOffsets = voidOffsets{&sg.layerBase}
    sg.bubbleEnd = bubbleEnd{&sg.layerBase}

    if child != nil {
        child.setParent(sg)
    }

    return sg
}

func (g *stageLayer) Setup() (err error) {
    g.surface, err = util.CreateSurface(true, int(g.size_x), int(g.size_y))
    if err != nil {
        return err
    }

    for row := 0; row < len(g.s.Tiles[0]); row++ {
        for col := 0; col < len(g.s.Tiles); col++ {
            texture := g.terrain_textures[g.s.Tiles[col][row]]

            x_tiles := g.square_size / texture.Width
            y_tiles := g.square_size / texture.Height

            for y_tile := 0; y_tile < y_tiles; y_tile++ {
                for x_tile := 0; x_tile < x_tiles; x_tile++ {
                    errno := g.surface.Blit(
                        &sdl.Rect{
                            X: int16(col * g.square_size + x_tile * texture.Width),
                            Y: int16(row * g.square_size + y_tile * texture.Height),
                            W: uint16(texture.Width),
                            H: uint16(texture.Height),
                        },
                        texture.Surface,
                        &sdl.Rect{
                            X: 0,
                            Y: 0,
                            W: uint16(texture.Width),
                            H: uint16(texture.Height),
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

