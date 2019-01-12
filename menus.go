package main

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"strconv"
)

const (
	VOLUME_SLIDER_CIRCLE_SIZE float32 = 50.0
	VOLUME_SLIDER_LONG_WIDTH  float32 = 300.0
	VOLUME_SLIDER_LONG_HEIGHT float32 = 25.0
	VOLUME_SLIDER_STEP_SIZE   float32 = 0.1
)

type WinMenu struct {
	backBtn     gohome.Button
	continueBtn gohome.Button
	winText     gohome.Text2D

	direction bool
}

func (this *WinMenu) Init() {
	start := gohome.Render.GetNativeResolution().Mul(0.5)
	this.backBtn.Init([2]float32{
		start.X() - (DEATH_BUTTON_SIZE*2+DEATH_BUTTON_PADDING)/2.0 + DEATH_BUTTON_SIZE/2.0,
		-DEATH_BUTTON_SIZE - DEATH_BUTTON_SIZE/2.0,
	}, "Back")
	this.backBtn.Transform.Size = [2]float32{DEATH_BUTTON_SIZE, DEATH_BUTTON_SIZE}
	this.backBtn.Transform.Origin = [2]float32{0.5, 0.5}
	this.backBtn.PressCallback = func(btn *gohome.Button) {
		gohome.ResourceMgr.GetSound("ButtonPressed").Play(false)
		gohome.SceneMgr.SwitchScene(&LevelSelectScene{})
	}
	this.backBtn.EnterCallback = func(btn *gohome.Button) {
		gohome.ResourceMgr.GetSound("Button").Play(false)
	}
	this.backBtn.Depth = MENU_DEPTH

	this.continueBtn.Init([2]float32{
		start.X() + (DEATH_BUTTON_SIZE*2+DEATH_BUTTON_PADDING)/2.0 - DEATH_BUTTON_SIZE/2.0,
		-DEATH_BUTTON_SIZE - DEATH_BUTTON_SIZE/2.0,
	}, "Continue")
	this.continueBtn.Transform.Size = [2]float32{DEATH_BUTTON_SIZE, DEATH_BUTTON_SIZE}
	this.continueBtn.Transform.Origin = [2]float32{0.5, 0.5}
	this.continueBtn.PressCallback = func(btn *gohome.Button) {
		gohome.ResourceMgr.GetSound("ButtonPressed").Play(false)
		gohome.SceneMgr.SwitchScene(&LevelScene{LevelID: gohome.SceneMgr.GetCurrentScene().(*LevelScene).LevelID + 1})
	}
	this.continueBtn.EnterCallback = func(btn *gohome.Button) {
		gohome.ResourceMgr.GetSound("Button").Play(false)
	}
	this.continueBtn.Depth = MENU_DEPTH

	this.winText.Init(gohome.ButtonFont, gohome.ButtonFontSize*2, "Level Abgeschlossen")
	this.winText.Transform.Origin = [2]float32{0.5, 0.5}
	this.winText.Transform.Position = [2]float32{
		start.X(),
		-DEATH_BUTTON_SIZE - DEATH_BUTTON_SIZE/2.0 - DEATH_BUTTON_SIZE*1.5,
	}
	this.winText.NotRelativeToCamera = 0
	this.winText.Depth = MENU_DEPTH

	gohome.RenderMgr.AddObject(&this.winText)
	gohome.UpdateMgr.AddObject(this)

	this.direction = UP
}

func (this *WinMenu) Update(delta_time float32) {
	var target mgl32.Vec2
	if this.direction == DOWN {
		target = gohome.Render.GetNativeResolution().Mul(0.5)
	} else {
		target = gohome.Render.GetNativeResolution().Mul(0.5)
		target[1] = -DEATH_BUTTON_SIZE - DEATH_BUTTON_SIZE/2.0
	}

	backTarget := target
	backTarget[0] = backTarget[0] - (DEATH_BUTTON_SIZE*2+DEATH_BUTTON_PADDING)/2.0 + DEATH_BUTTON_SIZE/2.0
	continueTarget := target
	continueTarget[0] = continueTarget[0] + (DEATH_BUTTON_SIZE*2+DEATH_BUTTON_PADDING)/2.0 - DEATH_BUTTON_SIZE/2.0
	winTextTarget := target
	winTextTarget[1] = winTextTarget[1] - DEATH_BUTTON_SIZE*1.5

	this.backBtn.Transform.Position = this.backBtn.Transform.Position.Add(backTarget.Sub(this.backBtn.Transform.Position).Mul(0.2))
	this.continueBtn.Transform.Position = this.continueBtn.Transform.Position.Add(continueTarget.Sub(this.continueBtn.Transform.Position).Mul(0.2))
	this.winText.Transform.Position = this.winText.Transform.Position.Add(winTextTarget.Sub(this.winText.Transform.Position).Mul(0.15))
}

func (this *WinMenu) Terminate() {
	this.backBtn.Terminate()
	this.continueBtn.Terminate()
	gohome.RenderMgr.RemoveObject(&this.winText)
	gohome.UpdateMgr.RemoveObject(this)
}

type OptionsMenu struct {
	text         gohome.Text2D
	volumeSlider gohome.Slider
	direction    bool
}

func (this *OptionsMenu) Init() {
	mid := gohome.Render.GetNativeResolution().Div(2.0)

	gohome.UpdateMgr.AddObject(this)

	this.volumeSlider.Init(mid.Sub([2]float32{VOLUME_SLIDER_LONG_WIDTH / 2.0, mid.Y() + VOLUME_SLIDER_LONG_HEIGHT + (VOLUME_SLIDER_CIRCLE_SIZE/2.0 - VOLUME_SLIDER_LONG_HEIGHT/2.0) + VOLUME_SLIDER_CIRCLE_SIZE}), "", "")
	this.volumeSlider.Circle.Transform.Size = [2]float32{VOLUME_SLIDER_CIRCLE_SIZE, VOLUME_SLIDER_CIRCLE_SIZE}
	this.volumeSlider.Long.Transform.Size = [2]float32{VOLUME_SLIDER_LONG_WIDTH, VOLUME_SLIDER_LONG_HEIGHT}
	this.volumeSlider.Circle.Depth = MENU_DEPTH
	this.volumeSlider.Long.Depth = MENU_DEPTH
	this.volumeSlider.ValueChangedCallback = func(sld *gohome.Slider) {
		gohome.AudioMgr.SetVolume(sld.Value)
	}
	this.volumeSlider.Value = gohome.AudioMgr.GetVolume()
	this.volumeSlider.StepSize = VOLUME_SLIDER_STEP_SIZE

	this.text.Init(gohome.ButtonFont, gohome.ButtonFontSize*2.0, "Lautstärke")
	this.text.NotRelativeToCamera = 0
	this.text.Transform.Origin = [2]float32{0.5, 0.5}
	this.text.Transform.Position = mid.Sub([2]float32{0.0, mid.Y() + VOLUME_SLIDER_LONG_HEIGHT + this.text.Transform.Size[1]*this.text.Transform.Scale[1]})

	gohome.RenderMgr.AddObject(&this.text)

	this.direction = UP

}

func (this *OptionsMenu) Update(delta_time float32) {
	var target mgl32.Vec2
	mid := gohome.Render.GetNativeResolution().Div(2.0)

	if this.direction == UP {
		target = mid.Sub([2]float32{VOLUME_SLIDER_LONG_WIDTH / 2.0, mid.Y() + VOLUME_SLIDER_LONG_HEIGHT + (VOLUME_SLIDER_CIRCLE_SIZE/2.0 - VOLUME_SLIDER_LONG_HEIGHT/2.0) + VOLUME_SLIDER_CIRCLE_SIZE})
	} else {
		target = mid.Sub([2]float32{VOLUME_SLIDER_LONG_WIDTH / 2.0, -VOLUME_SLIDER_LONG_HEIGHT / 2.0})
	}

	this.volumeSlider.Long.Transform.Position = this.volumeSlider.Long.Transform.Position.Add(target.Sub(this.volumeSlider.Long.Transform.Position).Mul(0.07))

	target1 := target.Sub([2]float32{-165.0, VOLUME_SLIDER_LONG_HEIGHT + this.text.Transform.Size[1]*this.text.Transform.Scale[1]})

	this.text.Transform.Position = this.text.Transform.Position.Add(target1.Sub(this.text.Transform.Position).Mul(0.06))
}

func (this *OptionsMenu) Terminate() {
	this.volumeSlider.Terminate()
	this.volumeSlider.Long.Terminate()
	this.volumeSlider.Circle.Terminate()
	gohome.UpdateMgr.RemoveObject(this)
	this.text.Terminate()
	gohome.RenderMgr.RemoveObject(&this.text)
}

type LevelTitle struct {
	gohome.Text2D

	Level        uint8
	WinCondition bool

	direction bool
}

func (this *LevelTitle) Init() {
	this.Text2D.Init(gohome.ButtonFont, 2*gohome.ButtonFontSize, this.WCS()+"Level "+strconv.FormatUint(uint64(this.Level), 10)+this.WCS()+"\n"+this.winConditionToString())
	this.Transform.Origin = [2]float32{0.5, 0.5}
	this.NotRelativeToCamera = 0
	this.Depth = MENU_DEPTH

	gohome.RenderMgr.AddObject(this)
	gohome.UpdateMgr.AddObject(this)

	this.direction = DOWN
	mid := gohome.Render.GetNativeResolution().Div(2.0)
	this.Transform.Position = mid.Sub([2]float32{0.0, mid.Y() + this.Text2D.Transform.Size[1]/2.0 + 20.0})
}

func (this *LevelTitle) Update(delta_time float32) {
	mid := gohome.Render.GetNativeResolution().Div(2.0)
	target := mid
	if this.direction == DOWN {
		target[1] -= mid.Y() / 2.0
	} else {
		target[1] -= mid.Y() + this.Text2D.Transform.Size[1]/2.0 + 20.0
	}

	this.Transform.Position = this.Transform.Position.Add(target.Sub(this.Transform.Position).Mul(0.03))

	if this.direction == DOWN {
		if this.Transform.Position.Sub(target).Len() < 3.0 {
			this.direction = !this.direction
		}
	}
}

func (this *LevelTitle) Terminate() {
	gohome.RenderMgr.RemoveObject(this)
	gohome.UpdateMgr.RemoveObject(this)
	this.Text2D.Terminate()
}

func (this *LevelTitle) winConditionToString() string {
	switch this.WinCondition {
	case WIN_CONDITION_TARGET:
		return "Sammle alle Flaggen"
	case WIN_CONDITION_ENEMY:
		return "Besiege alle Gegner"
	default:
		return "Schließe den Level ab"
	}
}

func (this *LevelTitle) WCS() string {
	switch this.WinCondition {
	case WIN_CONDITION_TARGET:
		return "      "
	case WIN_CONDITION_ENEMY:
		return "      "
	default:
		return "       "
	}
}
