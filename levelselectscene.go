package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"strconv"
)

const (
	LEVEL_BUTTON_PADDING float32 = 50.0
	LEVEL_BUTTON_SIZE    float32 = 100.0
	LEVEL_BUTTON_PER_ROW uint32  = 3
)

type LevelSelectScene struct {
	levelBtns    []*gohome.Button
	targetBtnPos []mgl32.Vec2
	title        *gohome.Text2D
}

func selectLevel(btn *gohome.Button) {
	id, _ := strconv.ParseInt(btn.Text, 10, 32)
	id -= 1
	gohome.ResourceMgr.GetSound("ButtonPressed").Play(false)
	gohome.SceneMgr.SwitchScene(&LevelScene{LevelID: uint32(id)})
}

func (this *LevelSelectScene) initButtons() {
	lbr := float32(LEVEL_BUTTON_PER_ROW)
	lbc := float32(NUM_LEVELS / LEVEL_BUTTON_PER_ROW)
	start := gohome.Render.GetNativeResolution().Mul(0.5)
	start = start.Sub([2]float32{
		(lbr*LEVEL_BUTTON_SIZE+(lbr-1.0)*LEVEL_BUTTON_PADDING)/2.0 - LEVEL_BUTTON_SIZE/2.0,
		(lbc*LEVEL_BUTTON_SIZE+(lbc-1.0)*LEVEL_BUTTON_PADDING)/2.0 - LEVEL_BUTTON_SIZE/2.0,
	})
	for i := uint32(0); i < NUM_LEVELS; i++ {
		this.levelBtns = append(this.levelBtns, &gohome.Button{})
		btn := this.levelBtns[len(this.levelBtns)-1]
		x := float32(i%LEVEL_BUTTON_PER_ROW) * (LEVEL_BUTTON_SIZE + LEVEL_BUTTON_PADDING)
		y := float32(i/LEVEL_BUTTON_PER_ROW) * (LEVEL_BUTTON_SIZE + LEVEL_BUTTON_PADDING)
		maxy := float32((NUM_LEVELS-1)/LEVEL_BUTTON_PER_ROW) * (LEVEL_BUTTON_SIZE + LEVEL_BUTTON_PADDING)

		btn.Text = strconv.FormatInt(int64(i+1), 10)
		this.targetBtnPos = append(this.targetBtnPos, start.Add([2]float32{x, y}))
		btn.Init(start.Add([2]float32{x, y}), "LevelButton1")
		btn.Transform.Position[1] = -LEVEL_BUTTON_SIZE/2.0 - (maxy - y)
		btn.Transform.Size = [2]float32{LEVEL_BUTTON_SIZE, LEVEL_BUTTON_SIZE}
		btn.Transform.Origin = [2]float32{0.5, 0.5}
		btn.PressCallback = selectLevel
		btn.EnterCallback = func(button *gohome.Button) {
			button.Texture = gohome.ResourceMgr.GetTexture("LevelButtonPressed")
			gohome.ResourceMgr.GetSound("Button").Play(false)
		}
		btn.LeaveCallback = func(button *gohome.Button) {
			button.Texture = gohome.ResourceMgr.GetTexture("LevelButton1")
		}
		btn.EnterModColor = nil
		btn.PressModColor = nil
	}
}

func (this *LevelSelectScene) initTitle() {
	lbc := float32(NUM_LEVELS / LEVEL_BUTTON_PER_ROW)
	start := gohome.Render.GetNativeResolution().Mul(0.5)
	start = start.Sub([2]float32{
		0.0,
		(lbc*LEVEL_BUTTON_SIZE+(lbc-1.0)*LEVEL_BUTTON_PADDING)/2.0 - LEVEL_BUTTON_SIZE/2.0,
	})
	maxy := float32((NUM_LEVELS-1)/LEVEL_BUTTON_PER_ROW) * (LEVEL_BUTTON_SIZE + LEVEL_BUTTON_PADDING)

	this.title = &gohome.Text2D{}
	this.title.Init(gohome.ButtonFont, gohome.ButtonFontSize*2, "Wähle einen Level")
	this.title.Transform.Origin = [2]float32{0.5, 0.5}
	this.title.Transform.Position = [2]float32{gohome.Render.GetNativeResolution().X()/2.0 + 10.0, -LEVEL_BUTTON_SIZE/2.0 - maxy - (start[1] - 100.0)}
	this.title.NotRelativeToCamera = 0
	gohome.RenderMgr.AddObject(this.title)
}

func (this *LevelSelectScene) updateButtons() {
	for i := 0; i < len(this.levelBtns); i++ {
		btn := this.levelBtns[i]
		tpos := this.targetBtnPos[i]
		btn.Transform.Position = btn.Transform.Position.Add(tpos.Sub(btn.Transform.Position).Mul(0.1))
	}
}

func (this *LevelSelectScene) updateTitle() {
	titleTarget := mgl32.Vec2{gohome.Render.GetNativeResolution().X()/2.0 + 10.0, 100.0}
	this.title.Transform.Position = this.title.Transform.Position.Add(titleTarget.Sub(this.title.Transform.Position).Mul(0.08))
}

func (this *LevelSelectScene) Init() {
	this.initButtons()
	this.initTitle()
}

func (this *LevelSelectScene) Update(delta_time float32) {
	this.updateButtons()
	this.updateTitle()
}

func (this *LevelSelectScene) Terminate() {
	for i := 0; i < len(this.levelBtns); i++ {
		this.levelBtns[i].Terminate()
	}
	gohome.RenderMgr.RemoveObject(this.title)
	this.title.Terminate()
}
