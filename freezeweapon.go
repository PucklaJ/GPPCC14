package main

import (
	"github.com/ByteArena/box2d"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/physics2d"
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

const (
	FREEZE_TIME float32 = 0.5
	FREEZE_AMMO uint32  = 50

	FREEZE_WIDTH    float32 = 32.0
	FREEZE_HEIGHT   float32 = 10.0
	FREEZE_FRICTION float64 = 0.3

	FREEZE_VELOCITY float32 = 250.0

	FREEZE_OFFSET_X float32 = 2.0
	FREEZE_OFFSET_Y float32 = -1.0

	FREEZE_FRAME_WIDTH  float32 = 36.0
	FREEZE_FRAME_HEIGHT float32 = 16.0
)

type FreezeWeapon struct {
	NilWeapon
	bodies []*box2d.B2Body
	times  []float32
}

func (this *FreezeWeapon) OnAdd(p *Player) {
	this.Sprite2D.Init("FreezeWeapon")
	this.Transform.Origin = [2]float32{0.5, 0.5}

	this.NilWeapon.OnAdd(p)
	gohome.UpdateMgr.AddObject(this)
	this.Ammo = FREEZE_AMMO
}

func (this *FreezeWeapon) GetInventoryTexture() gohome.Texture {
	return gohome.ResourceMgr.GetTexture("FreezeWeaponInv")
}

func (this *FreezeWeapon) Update(delta_time float32) {
	off := [2]float32{FREEZE_OFFSET_X, FREEZE_OFFSET_Y}
	this.Flip = this.Player.Flip
	if this.Flip == gohome.FLIP_HORIZONTAL {
		off[0] = -off[0]
	}
	this.Transform.Position = this.Player.Transform.Position.Add(this.Player.GetWeaponOffset()).Add(off)

	if this.paused {
		return
	}

	for i := 0; i < len(this.times); i++ {
		if this.times[i] > 0.0 {
			this.times[i] -= delta_time
		}
		if this.times[i] <= 0.0 {
			this.bodies[i].SetType(box2d.B2BodyType.B2_staticBody)
			block := this.bodies[i].GetUserData().(*WeaponBlock)
			block.Sprite.TextureRegion.Min[0], block.Sprite.TextureRegion.Max[0] = FREEZE_FRAME_WIDTH, FREEZE_FRAME_WIDTH*2
		}
	}

}

func (this *FreezeWeapon) Use(target mgl32.Vec2, energy float32) {
	dir := target.Sub(this.Player.Transform.Position).Normalize()
	body := this.createBox(dir, energy)
	this.bodies = append(this.bodies, body)
	this.times = append(this.times, FREEZE_TIME)
	this.Ammo--
}

func (this *FreezeWeapon) createBox(dir mgl32.Vec2, energy float32) *box2d.B2Body {
	pos := this.Player.Transform.Position.Add(dir.Mul(PLAYER_WIDTH * 2.0))
	size := [2]float32{FREEZE_WIDTH, FREEZE_HEIGHT}

	bodyDef := box2d.MakeB2BodyDef()
	bodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	bodyDef.Position = physics2d.ToBox2DCoordinates(pos)
	bodyDef.Angle = -float64(dir.Angle())
	fdef := box2d.MakeB2FixtureDef()
	fdef.Friction = FREEZE_FRICTION
	fdef.Density = 1.0 / (physics2d.ScalarToBox2D(FREEZE_WIDTH) * physics2d.ScalarToBox2D(FREEZE_HEIGHT)) * DEFAULT_WEAPON_WEIGHT
	fdef.Restitution = DEFAULT_WEAPON_RESTITUTION
	fdef.Filter.CategoryBits = WEAPON_CATEGORY
	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox(physics2d.ScalarToBox2D(size[0])/2.0, physics2d.ScalarToBox2D(size[1])/2.0)
	fdef.Shape = &shape
	body := this.Player.PhysicsMgr.World.CreateBody(&bodyDef)
	body.CreateFixtureFromDef(&fdef)

	body.SetLinearVelocity(physics2d.ToBox2DDirection(dir.Mul(FREEZE_VELOCITY * energy)))
	body.SetLinearVelocity(box2d.B2Vec2Add(this.Player.body.GetLinearVelocity(), body.GetLinearVelocity()))

	var spr gohome.Sprite2D
	var con physics2d.PhysicsConnector2D

	spr.Init("FreezeWeaponBlock")
	spr.TextureRegion.Max[0] = FREEZE_FRAME_WIDTH
	spr.Transform.Size[0], spr.Transform.Size[1] = FREEZE_FRAME_WIDTH, FREEZE_FRAME_HEIGHT
	spr.Transform.Origin = [2]float32{0.5, 0.5}
	con.Init(spr.Transform, body)

	gohome.RenderMgr.AddObject(&spr)
	gohome.UpdateMgr.AddObject(&con)

	var block WeaponBlock
	block.Sprite = &spr
	block.Connector = &con
	this.blocks = append(this.blocks, block)

	body.SetUserData(&this.blocks[len(this.blocks)-1])

	con.Update(0.0)

	return body
}

func (this *FreezeWeapon) OnDie() {
	gohome.UpdateMgr.RemoveObject(this)
	gohome.RenderMgr.RemoveObject(&this.NilWeapon)
}

func (this *FreezeWeapon) Terminate() {
	this.NilWeapon.Terminate()
	gohome.UpdateMgr.RemoveObject(this)
}
