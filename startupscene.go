package main

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"time"
)

type GlobalUpdate struct {
}

func (this *GlobalUpdate) Update(delta_time float32) {
	if gohome.InputMgr.JustPressed(gohome.KeyF11) {
		if !gohome.Framew.WindowIsFullscreen() {
			ms := gohome.Framew.MonitorGetSize()
			gohome.Framew.WindowSetSize(ms)
			time.Sleep(time.Millisecond * 100)
			gohome.Framew.WindowSetFullscreen(true)
		} else {
			gohome.Framew.WindowSetFullscreen(false)
			time.Sleep(time.Millisecond * 100)
			gohome.Framew.WindowSetSize([2]float32{
				float32(GAME_WIDTH),
				float32(GAME_HEIGHT),
			})
		}
	}
}

type StartupScene struct {
}

func (this *StartupScene) Init() {
	gohome.Init2DShaders()

	audio := gohome.Framew.GetAudioManager()
	audio.Init()
	audio.SetVolume(0.5)

	LoadResources()

	gohome.RenderMgr.GetBackBuffer().SetFiltering(gohome.FILTERING_NEAREST)
	gohome.UpdateMgr.AddObject(&GlobalUpdate{})

	gohome.Render.SetBackgroundColor(gohome.Color{52, 101, 255, 255})
	gohome.RenderMgr.SetCamera2D(&Camera, 0)
	Camera.Zoom = ZOOM
	gohome.SceneMgr.SwitchScene(&LevelSelectScene{})
}

func (this *StartupScene) Update(delta_time float32) {

}

func (this *StartupScene) Terminate() {
}
