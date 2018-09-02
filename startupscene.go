package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

type StartupScene struct {
}

func (this *StartupScene) Init() {
	gohome.Init2DShaders()
	gohome.Render.SetBackgroundColor(gohome.Color{52, 101, 255, 255})
	LoadResources()
	gohome.RenderMgr.SetCamera2D(&Camera, 0)
	Camera.Zoom = ZOOM
	gohome.SceneMgr.SwitchScene(&LevelSelectScene{})
}

func (this *StartupScene) Update(delta_time float32) {

}

func (this *StartupScene) Terminate() {
}
