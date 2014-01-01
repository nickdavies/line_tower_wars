package main

import (
    "./layer"
    "./stage"
)

func main() {

    s := stage.NewStage(2)

    eg := layer.NewEntityLayer(s, 100, nil)
    sg := layer.NewStageLayer(s, 100, eg)
    pg := layer.NewPanLayer(90, 90, 0, 0, 3000, sg)
    g := layer.NewSdlLayer(2560, 1600, pg)

    g.Run()
}

