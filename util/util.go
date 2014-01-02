package util

import (
    "fmt"
)

import (
    "github.com/neagix/Go-SDL/sdl"
)

func CreateSurface(hw bool, x, y int) (*sdl.Surface, error) {
    var flags uint32 = 0
    if hw {
        flags = sdl.HWSURFACE
    } else {
        flags = sdl.SWSURFACE
    }

    surface := sdl.CreateRGBSurface(flags, x, y, 32, 0x000000ff, 0x0000ff00, 0x00ff0000, 0xff000000)
    if surface == nil {
        return nil, fmt.Errorf("Failed to load texture: %s", sdl.GetError())
    }

    surface = surface.DisplayFormat()
    if surface == nil {
        return nil, fmt.Errorf("Failed to convert texture: %s", sdl.GetError())
    }

    return surface, nil
}
