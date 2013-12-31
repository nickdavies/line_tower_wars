package game

import (
    "github.com/banthar/Go-SDL/sdl"
)

type Game interface {

    Setup() error
    Cleanup()

    HandleEvent(event sdl.Event)

    Update(deltaTime int64)
    Render(target *sdl.Surface)

    GetSize() (x uint16, y uint16)
    GetXYOffsets() (x uint16, y uint16)

    setParent(parent Game)

    Run() error
    End()
}

type gameBase struct {
    child Game
    parent Game
}

func (g *gameBase) setParent(parent Game) {
    g.parent = parent
}

// Struct for games with no setup/cleanup
type voidSetup struct {
    *gameBase
}

func (g *voidSetup) Setup() error {
    if g.child != nil {
        return g.child.Setup()
    }

    return nil
}

func (g *voidSetup) Cleanup() {
    if g.child != nil {
        g.child.Cleanup()
    }
}

// Struct for doing nothing on run
type voidRun struct {
    *gameBase
}

func (g *voidRun) Run() error {
    return nil
}

// Struct for passing events though
type voidEvents struct {
    *gameBase
}

func (g *voidEvents) HandleEvent(event sdl.Event) {
    if g.child != nil {
        g.child.HandleEvent(event)
    }
}

// Struct for simply passing end command up
type bubbleEnd struct {
    *gameBase
}

func (g *bubbleEnd) End() {
    if g.parent != nil {
        g.parent.End()
    }
}

type voidOffsets struct {
    *gameBase
}

func (g *voidOffsets) GetXYOffsets() (uint16, uint16) {
    return g.parent.GetXYOffsets()
}

