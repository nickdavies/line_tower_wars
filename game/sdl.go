package game

import (
    "time"
    "fmt"
)

import (
    "github.com/banthar/Go-SDL/sdl"
)

func NewSdlGame(x, y uint16, child Game) Game {
    g := &sdlGame{
        gameBase: gameBase{child: child},

        x: x,
        y: y,
    }

    if child != nil {
        child.setParent(g)
    }

    return g
}

type sdlGame struct {
    // Get parent/child fields
    gameBase

    running bool

    // Display size
    x uint16
    y uint16

    // main display
    display *sdl.Surface
}

func (g *sdlGame) Setup() error {
    var errno = sdl.Init(sdl.INIT_EVERYTHING)
    if errno != 0 {
        return fmt.Errorf("Init failed: %s", sdl.GetError())
    }

    g.display = sdl.SetVideoMode(int(g.x), int(g.y), 32, sdl.HWSURFACE | sdl.DOUBLEBUF | sdl.FULLSCREEN)
    if g.display == nil {
        return fmt.Errorf("No surface created: %s", sdl.GetError())
    }

    if g.child != nil {
        return g.child.Setup()
    }

    return nil
}

func (g *sdlGame) Cleanup() {
    if g.child != nil {
        g.child.Cleanup()
    }
    sdl.Quit()
}

func (g *sdlGame) HandleEvent(event sdl.Event) {
    switch event.(type) {
    case *sdl.QuitEvent:
        g.End()
        return
    case *sdl.KeyboardEvent:
        if event.(*sdl.KeyboardEvent).Keysym.Sym == sdl.K_F1 {
            g.End()
        }
        return
    default:
    }
    if g.child != nil {
        g.child.HandleEvent(event)
    }
}

func (g *sdlGame) Update(deltaTime int64) {
    if g.child != nil {
        g.child.Update(deltaTime)
    }
}

func (g *sdlGame) Render(target *sdl.Surface) {
    if g.child != nil {
        g.child.Render(g.display)
    }

    g.display.Flip()
}

func (g *sdlGame) GetSize() (uint16, uint16) {
    return g.x, g.y
}

func (g *sdlGame) Run() error {
    var err error

    err = g.Setup()
    if err != nil {
        return err
    }
    defer g.Cleanup()

    var last_time = time.Now().UnixNano()

    g.running = true
    for g.running {
        // Process Events
        for {
            event := sdl.PollEvent()
            if event == nil {
                break
            }
            g.HandleEvent(event)
        }

        // Update State
        g.Update(time.Now().UnixNano() - last_time)
        last_time = time.Now().UnixNano()

        // Update Screen
        g.Render(nil)
    }

    return nil
}

func (g *sdlGame) End() {
    g.running = false
}

