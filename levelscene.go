package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/physics2d"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	DEATH_BUTTON_SIZE    float32 = 100.0
	DEATH_BUTTON_PADDING float32 = 100.0
	DEATH_TEXT_PADDING   float32 = 25.0
)

type LevelScene struct {
	PhysicsMgr physics2d.PhysicsManager2D
	LevelID    uint32
	Map        gohome.TiledMap
	Player     Player
	Enemies    []*Enemy

	debugDraw physics2d.PhysicsDebugDraw2D

	deathBtns     [2]*gohome.Button
	deathText     *gohome.Text2D
	menuInited    bool
	menuDirection bool
	paused        bool
}

func (this *LevelScene) Init() {
	physics2d.PIXEL_PER_METER = 10.0
	gohome.ResourceMgr.LoadTMXMap("Level", LEVELS_TMX_MAPS[this.LevelID])

	this.Map.Init("Level")
	gohome.RenderMgr.AddObject(&this.Map)

	this.PhysicsMgr.Init([2]float32{0.0, GRAVITY})
	gohome.UpdateMgr.AddObject(&this.PhysicsMgr)
	this.debugDraw = this.PhysicsMgr.GetDebugDraw()
	this.debugDraw.Visible = false
	gohome.RenderMgr.AddObject(&this.debugDraw)

	groundBodies := this.PhysicsMgr.LayerToCollision(&this.Map, "Collision")
	for i := 0; i < len(groundBodies); i++ {
		b := groundBodies[i]
		if b == nil {
			continue
		}
		for f := b.GetFixtureList(); f != nil; f = f.GetNext() {
			filter := f.GetFilterData()
			filter.CategoryBits = GROUND_CATEGORY
			filter.MaskBits = 0xffff
			f.SetFilterData(filter)
			f.SetFriction(GROUND_FRICTION)
		}
	}

	var playerStart [2]float32

	ls := this.Map.Layers
	for i := 0; i < len(ls); i++ {
		l := ls[i]
		if l.Name == "Settings" {
			objs := l.Objects
			for j := 0; j < len(objs); j++ {
				o := objs[j]
				if o.Name == "start" {
					playerStart[0] = float32(o.X)
					playerStart[1] = float32(o.Y)
				} else if o.Name == "enemy" {
					enemy := &Enemy{}
					enemy.Sprite2D.Init("")
					enemy.Transform.Position = [2]float32{float32(o.X), float32(o.Y)}
					this.Enemies = append(this.Enemies, enemy)
				}
			}
		}
	}

	this.Player.Init(playerStart, &this.PhysicsMgr)
	for i := 0; i < len(this.Enemies); i++ {
		this.Enemies[i].Init(this.Enemies[i].Transform.Position, &this.Player)
	}

	Camera.Position = [2]float32{-CAMERA_BOX_WIDTH, -CAMERA_BOX_HEIGHT}
}

func (this *LevelScene) terminateMenu() {
	for _, btn := range this.deathBtns {
		if btn != nil {
			btn.Terminate()
			btn.Sprite2D.Terminate()
		}
	}
	if this.deathText != nil {
		gohome.RenderMgr.RemoveObject(this.deathText)
		this.deathText.Terminate()
	}
	this.menuInited = false
}

func (this *LevelScene) initMenu(death bool, inMid bool) {
	if this.menuInited {
		return
	}

	if this.deathBtns[0] != nil {
		this.terminateMenu()
	}

	var restartBtn, backBtn gohome.Button

	width := 2.0*DEATH_BUTTON_SIZE + DEATH_BUTTON_PADDING
	mid := gohome.Framew.WindowGetSize().Mul(0.5)
	var restartPos, backPos, deathTextPos mgl32.Vec2

	if !inMid {
		restartPos = mid.Add([2]float32{
			-width/2.0 + DEATH_BUTTON_SIZE/2.0,
			-mid.Y() - DEATH_BUTTON_SIZE,
		})
		backPos = mid.Add([2]float32{
			width/2.0 - DEATH_BUTTON_SIZE/2.0,
			-mid.Y() - DEATH_BUTTON_SIZE,
		})
		deathTextPos = mid.Add([2]float32{
			10.0,
			-mid.Y() - DEATH_BUTTON_SIZE*1.5 - DEATH_TEXT_PADDING,
		})
		this.menuDirection = DOWN

	} else {
		restartPos = mid.Add([2]float32{
			-width/2.0 + DEATH_BUTTON_SIZE/2.0,
			0.0,
		})
		backPos = mid.Add([2]float32{
			width/2.0 - DEATH_BUTTON_SIZE/2.0,
			0.0,
		})
		deathTextPos = mid.Add([2]float32{
			10.0,
			-DEATH_BUTTON_SIZE - DEATH_TEXT_PADDING,
		})
		this.menuDirection = UP
	}

	restartBtn.Init(restartPos, "")
	restartBtn.Transform.Origin = [2]float32{0.5, 0.5}
	restartBtn.Transform.Size = [2]float32{DEATH_BUTTON_SIZE, DEATH_BUTTON_SIZE}
	restartBtn.PressCallback = func(btn *gohome.Button) {
		this.Restart()
	}

	backBtn.Init(backPos, "")
	backBtn.Transform.Origin = [2]float32{0.5, 0.5}
	backBtn.Transform.Size = [2]float32{DEATH_BUTTON_SIZE, DEATH_BUTTON_SIZE}
	backBtn.PressCallback = func(btn *gohome.Button) {
		gohome.SceneMgr.SwitchScene(&LevelSelectScene{})
	}

	this.deathBtns[0] = &restartBtn
	this.deathBtns[1] = &backBtn

	if death {
		this.deathText = &gohome.Text2D{}
		this.deathText.Init(gohome.ButtonFont, uint32(float32(gohome.ButtonFontSize)*1.5), "Sie sind gestorben")
		this.deathText.NotRelativeToCamera = 0
		this.deathText.Transform.Origin = [2]float32{0.5, 0.5}
		this.deathText.Transform.Position = deathTextPos
		gohome.RenderMgr.AddObject(this.deathText)
	}

	this.menuInited = true
}

func (this *LevelScene) Pause() {
	if this.Player.Died() {
		return
	}

	this.initMenu(false, false)
	this.menuDirection = DOWN
	this.paused = true

	this.Player.Pause()
	for _, e := range this.Enemies {
		e.paused = true
	}
	this.PhysicsMgr.Paused = true
}

func (this *LevelScene) Resume() {
	this.menuDirection = UP
	this.paused = false

	this.Player.Resume()
	for _, e := range this.Enemies {
		e.paused = false
	}
	this.PhysicsMgr.Paused = false
}

func (this *LevelScene) Restart() {
	prevCamPos := Camera.Position
	died := this.Player.Died()
	scn := &LevelScene{LevelID: this.LevelID}
	gohome.SceneMgr.SwitchScene(scn)
	if died {
		scn.initMenu(true, true)
		scn.menuInited = false
	}
	Camera.Position = prevCamPos
}

func (this *LevelScene) updateMenu() {
	restartBtn := this.deathBtns[0]
	backBtn := this.deathBtns[1]

	if restartBtn == nil || backBtn == nil {
		return
	}

	width := 2.0*DEATH_BUTTON_SIZE + DEATH_BUTTON_PADDING
	mid := gohome.Framew.WindowGetSize().Mul(0.5)

	var restartTarget, backTarget, deathTextTarget mgl32.Vec2

	if this.menuDirection == DOWN {
		restartTarget = mid.Add([2]float32{
			-width/2.0 + DEATH_BUTTON_SIZE/2.0,
			0.0,
		})
		backTarget = mid.Add([2]float32{
			width/2.0 - DEATH_BUTTON_SIZE/2.0,
			0.0,
		})
		deathTextTarget = mid.Add([2]float32{
			10.0,
			-DEATH_BUTTON_SIZE - DEATH_TEXT_PADDING,
		})
	} else {
		restartTarget = mid.Add([2]float32{
			-width/2.0 + DEATH_BUTTON_SIZE/2.0,
			-mid.Y() - DEATH_BUTTON_SIZE,
		})
		backTarget = mid.Add([2]float32{
			width/2.0 - DEATH_BUTTON_SIZE/2.0,
			-mid.Y() - DEATH_BUTTON_SIZE,
		})
		deathTextTarget = mid.Add([2]float32{
			10.0,
			-mid.Y() - DEATH_BUTTON_SIZE*1.5 - DEATH_TEXT_PADDING,
		})
	}

	var btnSpeed, textSpeed float32
	if this.menuDirection == DOWN {
		btnSpeed = 0.05
		textSpeed = 0.04
	} else {
		btnSpeed = 0.1
		textSpeed = 0.08
	}

	restartBtn.Transform.Position = restartBtn.Transform.Position.Add(restartTarget.Sub(restartBtn.Transform.Position).Mul(btnSpeed))
	backBtn.Transform.Position = backBtn.Transform.Position.Add(backTarget.Sub(backBtn.Transform.Position).Mul(btnSpeed))
	if this.deathText != nil {

		this.deathText.Transform.Position = this.deathText.Transform.Position.Add(deathTextTarget.Sub(this.deathText.Transform.Position).Mul(textSpeed))
	}
	if this.menuDirection == UP && restartBtn.Transform.Position[1]+DEATH_BUTTON_SIZE/2.0 < 0.0 {
		this.terminateMenu()
	}
}

func (this *LevelScene) Update(delta_time float32) {
	if gohome.InputMgr.JustPressed(gohome.KeyF3) {
		this.debugDraw.Visible = !this.debugDraw.Visible
	} else if gohome.InputMgr.JustPressed(gohome.KeyR) {
		this.Restart()
	} else if gohome.InputMgr.JustPressed(gohome.KeyU) {
		gohome.SceneMgr.SwitchScene(&gohome.NilScene{})
	}
	if gohome.InputMgr.JustPressed(gohome.KeyK) {
		this.menuDirection = !this.menuDirection
	}
	if gohome.InputMgr.JustPressed(gohome.KeyP) {
		if this.paused {
			this.Resume()
		} else {
			this.Pause()
		}
	}
	this.updateMenu()
	if this.Player.Died() {
		this.initMenu(true, false)
	}
}

func (this *LevelScene) Terminate() {
	gohome.UpdateMgr.RemoveObject(&this.PhysicsMgr)
	gohome.RenderMgr.RemoveObject(&this.Map)
	gohome.RenderMgr.RemoveObject(&this.debugDraw)

	gohome.ResourceMgr.DeleteTMXMap("Level")

	this.terminateMenu()
	for i := 0; i < len(this.Enemies); i++ {
		this.Enemies[i].Terminate()
	}
	this.Player.Terminate()
	this.Map.Terminate()
	this.PhysicsMgr.Terminate()
}
