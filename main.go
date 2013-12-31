package main

import "./game"

func main() {

    sg := game.NewStageGame(2, 100, nil)
    pg := game.NewPanGame(90, 90, 0, 0, 3000, sg)
    g := game.NewSdlGame(2560, 1600, pg)

    g.Run()
}

