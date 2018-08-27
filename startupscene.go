package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

type StartupScene struct {
}

func (this *StartupScene) Init() {
	LoadResources()
	gohome.RenderMgr.SetCamera2D(&Camera, 0)
	Camera.Zoom = ZOOM
}

func (this *StartupScene) Update(delta_time float32) {
	gohome.SceneMgr.SwitchScene(&LevelScene{LevelID: 0})
}

func (this *StartupScene) Terminate() {
}
