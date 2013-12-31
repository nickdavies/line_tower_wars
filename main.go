package main

import (
    "./game"
    "./stage"
)

func main() {

    s := stage.NewStage(2)

    eg := game.NewEntityGame(s, 100, nil)
    sg := game.NewStageGame(s, 100, eg)
    pg := game.NewPanGame(90, 90, 0, 0, 3000, sg)
    g := game.NewSdlGame(2560, 1600, pg)

    g.Run()
}

