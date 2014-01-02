package main

import (
    "github.com/nickdavies/line_tower_wars/layer"
    "github.com/nickdavies/line_tower_wars/stage"
)

func main() {

    players := 2
    square_size := 100

    s := stage.NewStage(players)

    eg := layer.NewEntityLayer(s, square_size, nil)
    sg := layer.NewStageLayer(s, square_size, eg)
    pg := layer.NewPanLayer(90, 90, 0, 0, 3000, sg)
    g := layer.NewSdlLayer(2560, 1600, pg)

    g.Run()
}

