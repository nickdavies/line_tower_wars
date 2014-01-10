package layer

import (
    "github.com/neagix/Go-SDL/sdl"
)

import (
    "github.com/nickdavies/line_tower_wars/game"

    "github.com/nickdavies/line_tower_wars/graphics/texture"
)

type Layer interface {
    Setup() error
    Cleanup()

    HandleEvent(event interface{})

    Update(deltaTime int64)
    Render(target *sdl.Surface)

    GetSize() (x uint16, y uint16)
    GetXYOffsets() (x uint16, y uint16)

    setParent(parent Layer)
}

type LayerBuilder interface {
    Build(name string, cfg interface{}, child Layer) Layer
}

func NewLayerBuilder(game *game.Game, controls game.PlayerControls, texture_map texture.TextureMap, square_size uint16) LayerBuilder {
    return &layerBuilderStruct{
        game: game,
        controls: controls,
        texture_map: texture_map,
        square_size: square_size,
    }
}

type layerBuilderStruct struct {
    game *game.Game
    controls game.PlayerControls
    texture_map texture.TextureMap
    square_size uint16
}

func (lb *layerBuilderStruct) Build(name string, cfg interface{}, child Layer) Layer {
    l, ok := layers[name]
    if !ok {
        panic("No such layer " + name)
    }

    base := layerBase{
        game: lb.game,
        controls: lb.controls,
        texture_map: lb.texture_map,
        square_size: lb.square_size,
        child: child,
    }

    new_layer := l(base, cfg)
    if child != nil {
        child.setParent(new_layer)
    }

    return new_layer
}

// store of all registered layers
var layers = make(map[string]newLayerFunc)

// All layers must have this construction format
type newLayerFunc func(base layerBase, cfg interface{}) Layer

// all layers must call register so they can be built later
func registerLayer(name string , new_func newLayerFunc) {
    _, ok := layers[name]
    if ok {
        panic("Cant register layer " + name + " twice")
    }
    layers[name] = new_func
}
