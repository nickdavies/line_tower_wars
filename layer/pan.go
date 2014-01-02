package layer

import (
    "fmt"
    "time"
)

import (
    "github.com/neagix/Go-SDL/sdl"
)

type panLayer struct {
    layerBase

    // no run loop
    voidRun

    // proxy events
    voidEvents

    // proxy ends
    bubbleEnd

    pan_region_x uint16
    pan_region_y uint16
    pan_speed uint16

    view_x uint16
    view_y uint16

    child_x uint16
    child_y uint16

    surface *sdl.Surface
}

func NewPanLayer(pan_region_x, pan_region_y, starting_x, starting_y, pan_speed uint16, child Layer) Layer {

    if child == nil {
        panic(fmt.Errorf("You must give a child to pan Layer"))
    }

    pg := &panLayer{
        layerBase: layerBase{child: child},

        pan_region_x: pan_region_x,
        pan_region_y: pan_region_y,

        pan_speed: pan_speed,

        view_x: starting_x,
        view_y: starting_y,
    }

    pg.voidRun = voidRun{&pg.layerBase}
    pg.voidEvents = voidEvents{&pg.layerBase}
    pg.bubbleEnd = bubbleEnd{&pg.layerBase}

    child.setParent(pg)

    return pg
}

func (g *panLayer) Setup() error {
    g.child_x, g.child_y = g.child.GetSize()

    g.surface = sdl.CreateRGBSurface(sdl.HWSURFACE, int(g.child_x), int(g.child_y), 32, 0x000000ff, 0x0000ff00, 0x00ff0000, 0xff000000)
    if g.surface == nil {
        return fmt.Errorf("No surface created: %s", sdl.GetError())
    }

    g.surface = g.surface.DisplayFormat()
    if g.surface == nil {
        return fmt.Errorf("No surface created: %s", sdl.GetError())
    }

    return g.child.Setup()
}

func (g *panLayer) Cleanup() {
    g.surface.Free()
    g.child.Cleanup()
}

func (g *panLayer) calculatePan(mouse, pan_region_size, current_pan, pannable_size, window_size uint16, deltaTime int64) int16 {
    max_pan := float32((int64(g.pan_speed) * deltaTime) / int64(time.Second))

    if mouse < pan_region_size {
        pan_size := uint16(max_pan * float32(pan_region_size - mouse) / float32(pan_region_size))

        if current_pan < pan_size {
            return -1 * int16(current_pan)
        } else {
            return -1 * int16(pan_size)
        }
    }

    if mouse > (window_size - pan_region_size) {
        pan_size := uint16(max_pan * float32(pan_region_size - (window_size - mouse)) / float32(pan_region_size))

        if current_pan + window_size + pan_size > pannable_size {
            return int16(pannable_size - current_pan - window_size)
        } else {
            return int16(pan_size)
        }
    }

    return 0
}

func (g *panLayer) Update(deltaTime int64) {
    var mouse_x int
    var mouse_y int

    sdl.GetMouseState(&mouse_x, &mouse_y)

    parent_x, parent_y := g.parent.GetSize()

    g.view_x = uint16(int16(g.view_x) + g.calculatePan(uint16(mouse_x), g.pan_region_x, g.view_x, g.child_x, parent_x, deltaTime))
    g.view_y = uint16(int16(g.view_y) + g.calculatePan(uint16(mouse_y), g.pan_region_y, g.view_y, g.child_y, parent_y, deltaTime))
}

func (g *panLayer) Render(target *sdl.Surface) {
    g.child.Render(g.surface)
    g.surface.Flip()

    parent_x, parent_y := g.parent.GetSize()

    target.Blit(
        &sdl.Rect{
            X: 0,
            Y: 0,
            W: parent_x,
            H: parent_y,
        },
        g.surface,
        &sdl.Rect{
            X: int16(g.view_x),
            Y: int16(g.view_y),
            W: parent_x,
            H: parent_y,
        },
    )
}

func (g *panLayer) GetSize() (uint16, uint16) {
    return g.child.GetSize()
}

func (g *panLayer) GetXYOffsets() (uint16, uint16) {
    return g.view_x, g.view_y
}

