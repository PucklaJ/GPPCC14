package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"strconv"
)

const (
	LEVEL_BUTTON_PADDING float32 = 50.0
	LEVEL_BUTTON_SIZE    float32 = 100.0
	LEVEL_BUTTON_PER_ROW uint32  = 3
)

type LevelSelectScene struct {
	levelBtns []*gohome.Button
}

func selectLevel(btn *gohome.Button) {
	id, _ := strconv.ParseInt(btn.Text, 10, 32)
	id -= 1
	gohome.SceneMgr.SwitchScene(&LevelScene{LevelID: uint32(id)})
}

func (this *LevelSelectScene) Init() {
	lbr := float32(LEVEL_BUTTON_PER_ROW)
	lbc := float32(NUM_LEVELS / LEVEL_BUTTON_PER_ROW)
	start := gohome.Framew.WindowGetSize().Mul(0.5)
	start = start.Sub([2]float32{
		(lbr*LEVEL_BUTTON_SIZE+(lbr-1.0)*LEVEL_BUTTON_PADDING)/2.0 - LEVEL_BUTTON_SIZE/2.0,
		(lbc*LEVEL_BUTTON_SIZE+(lbc-1.0)*LEVEL_BUTTON_PADDING)/2.0 - LEVEL_BUTTON_SIZE/2.0,
	})
	for i := uint32(0); i < NUM_LEVELS; i++ {
		this.levelBtns = append(this.levelBtns, &gohome.Button{})
		btn := this.levelBtns[len(this.levelBtns)-1]
		x := float32(i%LEVEL_BUTTON_PER_ROW) * (LEVEL_BUTTON_SIZE + LEVEL_BUTTON_PADDING)
		y := float32(i/LEVEL_BUTTON_PER_ROW) * (LEVEL_BUTTON_SIZE + LEVEL_BUTTON_PADDING)

		btn.Text = strconv.FormatInt(int64(i+1), 10)
		btn.Init(start.Add([2]float32{x, y}), "")
		btn.Transform.Size = [2]float32{LEVEL_BUTTON_SIZE, LEVEL_BUTTON_SIZE}
		btn.Transform.Origin = [2]float32{0.5, 0.5}
		btn.PressCallback = selectLevel
	}
}

func (this *LevelSelectScene) Update(delta_time float32) {

}

func (this *LevelSelectScene) Terminate() {
	for i := 0; i < len(this.levelBtns); i++ {
		this.levelBtns[i].Terminate()
		this.levelBtns[i].Sprite2D.Terminate()
	}
}
