package main

import (
	"github.com/ByteArena/box2d"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/GoHomeEngine/src/physics2d"
)

const (
	TARGET_FRAME_TIME         float32 = 1.0 / 7.0
	TARGET_COLLECT_FRAME_TIME float32 = 1.0 / 8.0
)

type Target struct {
	gohome.Sprite2D
	anim gohome.Tweenset
}

type TargetCollect struct {
	gohome.Sprite2D
	anim gohome.Tweenset
}

func (this *TargetCollect) Init() {
	this.Sprite2D.Init("TargetCollect")
	this.anim = gohome.SpriteAnimation2D(this.Texture, 4, 1, TARGET_COLLECT_FRAME_TIME)
	this.anim.SetParent(&this.Sprite2D)
	this.anim.Start()

	this.TextureRegion.Max[0] = 32.0
	this.Transform.Size = [2]float32{32.0, 32.0}
	this.Transform.Origin = [2]float32{0.5, 0.5}
	this.Depth = SPECIAL_DEPTH

	gohome.UpdateMgr.AddObject(this)
	gohome.UpdateMgr.AddObject(&this.anim)
	gohome.RenderMgr.AddObject(this)

	gohome.ResourceMgr.GetSound("TargetCollect").Play(false)
}

func (this *TargetCollect) Update(delta_time float32) {
	if this.anim.Done() {
		this.Terminate()
	}
}

func (this *TargetCollect) Terminate() {
	gohome.UpdateMgr.RemoveObject(this)
	gohome.UpdateMgr.RemoveObject(&this.anim)
	gohome.RenderMgr.RemoveObject(this)
}

func (this *Target) Init(texName string) {
	this.Sprite2D.Init(texName)
	this.Depth = SPECIAL_DEPTH

	this.anim = gohome.SpriteAnimation2D(this.Texture, 3, 1, TARGET_FRAME_TIME)
	this.anim.Loop = true
	this.anim.SetParent(&this.Sprite2D)
	this.anim.Start()
	gohome.UpdateMgr.AddObject(&this.anim)
}

func (this *Target) Terminate() {
	gohome.RenderMgr.RemoveObject(this)
	gohome.UpdateMgr.RemoveObject(&this.anim)
}

type Sparcles struct {
	gohome.Sprite2D
	anim   gohome.Tweenset
	body   *box2d.B2Body
	world  *box2d.B2World
	weapon *DeleteWeapon
	paused bool
}

func (this *Sparcles) Update(delta_time float32) {
	if this.paused {
		return
	}

	if this.anim.Done() {
		t, ok := this.body.GetUserData().(TerminateObject)
		if ok {
			t.Terminate()
		}
		this.world.DestroyBody(this.body)
		this.Terminate()
	} else {
		this.Transform.Position = physics2d.ToPixelCoordinates(this.body.GetPosition())
	}
}

func (this *Sparcles) Terminate() {
	gohome.RenderMgr.RemoveObject(this)
	gohome.UpdateMgr.RemoveObject(this)
	gohome.UpdateMgr.RemoveObject(&this.anim)
	for i := 0; i < len(this.weapon.sparcles); i++ {
		if this.weapon.sparcles[i] == this {
			this.weapon.sparcles = append(this.weapon.sparcles[:i], this.weapon.sparcles[i+1:]...)
			return
		}
	}
}

func disappear(body *box2d.B2Body, world *box2d.B2World, wp *DeleteWeapon) *Sparcles {
	for f := body.GetFixtureList(); f != nil; f = f.GetNext() {
		var i int = 1
		f.SetUserData(&i)
	}

	var sp Sparcles
	sp.Init("Disappear")
	sp.Depth = SPECIAL_DEPTH
	sp.Transform.Position = physics2d.ToPixelCoordinates(body.GetPosition())
	sp.Transform.Origin = [2]float32{0.5, 0.5}
	sp.world = world
	sp.body = body
	sp.weapon = wp
	sp.anim = gohome.SpriteAnimation2D(sp.Texture, 3, 2, 1.0/8.0)
	sp.anim.SetParent(&sp.Sprite2D)
	sp.anim.Start()
	sp.anim.Update(0.0)
	gohome.RenderMgr.AddObject(&sp)
	gohome.UpdateMgr.AddObject(&sp.anim)
	gohome.UpdateMgr.AddObject(&sp)

	return &sp
}

type Explosion struct {
	gohome.Sprite2D
	anim gohome.Tweenset
}

func (this *Explosion) Init(texName string) {
	this.Sprite2D.Init(texName)
	this.Depth = SPECIAL_DEPTH
	gohome.ResourceMgr.GetSound("Explosion").Play(false)
}

func (this *Explosion) Update(delta_time float32) {
	if this.anim.Done() {
		gohome.RenderMgr.RemoveObject(this)
		gohome.UpdateMgr.RemoveObject(&this.anim)
		gohome.UpdateMgr.RemoveObject(this)
	}
}
