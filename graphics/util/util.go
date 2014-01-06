package util

import (
    "fmt"
)

import (
    "github.com/neagix/Go-SDL/sdl"
)

func CreateSurface(hw bool, display_format bool, x, y uint16) (*sdl.Surface, error) {
    var flags uint32 = 0
    if hw {
        flags = sdl.HWSURFACE
    } else {
        flags = sdl.SWSURFACE
    }

    surface := sdl.CreateRGBSurface(flags, int(x), int(y), 32, 0x000000ff, 0x0000ff00, 0x00ff0000, 0xff000000)
    if surface == nil {
        return nil, fmt.Errorf("Failed to load texture: %s", sdl.GetError())
    }

    if display_format {
        surface = surface.DisplayFormat()
        if surface == nil {
            return nil, fmt.Errorf("Failed to convert texture: %s", sdl.GetError())
        }
    }

    return surface, nil
}

func GetMouse(square_size, x_off, y_off uint16, tile_id bool) (col uint16, row uint16) {
    var raw_mouse_x int
    var raw_mouse_y int

    sdl.GetMouseState(&raw_mouse_x, &raw_mouse_y)

    mouse_x := uint16(raw_mouse_x) + x_off
    mouse_y := uint16(raw_mouse_y) + y_off

    if tile_id {
        return uint16(mouse_x / square_size), uint16(mouse_y / square_size)

    } else {
        square_x := mouse_x - (mouse_x % square_size)
        square_y := mouse_y - (mouse_y % square_size)

        return uint16(square_x), uint16(square_y)
    }
}
