package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

type StartupScene struct {
}

func (this *StartupScene) Init() {
	gohome.Init2DShaders()

	audio := gohome.Framew.GetAudioManager()
	audio.Init()
	audio.SetVolume(0.5)

	LoadResources()

	gohome.Render.SetBackgroundColor(gohome.Color{52, 101, 255, 255})
	gohome.RenderMgr.SetCamera2D(&Camera, 0)
	Camera.Zoom = ZOOM
	gohome.SceneMgr.SwitchScene(&LevelSelectScene{})
}

func (this *StartupScene) Update(delta_time float32) {

}

func (this *StartupScene) Terminate() {
}
