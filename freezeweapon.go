package main

import (
	"github.com/ByteArena/box2d"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/physics2d"
	"github.com/go-gl/mathgl/mgl32"
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
	this.tex.SetAsTarget()
	gohome.Render.ClearScreen(gohome.Color{0, 255, 50, 255})
	this.tex.UnsetAsTarget()
	gohome.UpdateMgr.AddObject(this)
	this.Ammo = FREEZE_AMMO
}

func (this *FreezeWeapon) Update(delta_time float32) {
	for i := 0; i < len(this.times); i++ {
		if this.times[i] > 0.0 {
			this.times[i] -= delta_time
		}
		if this.times[i] <= 0.0 {
			this.bodies[i].SetType(box2d.B2BodyType.B2_staticBody)
		}
	}
	off := [2]float32{FREEZE_OFFSET_X, FREEZE_OFFSET_Y}
	this.Flip = this.Player.Flip
	if this.Flip == gohome.FLIP_HORIZONTAL {
		off[0] = -off[0]
	}
	this.Transform.Position = this.Player.Transform.Position.Add(this.Player.GetWeaponOffset()).Add(off)
}

func (this *FreezeWeapon) Use(target mgl32.Vec2) {
	dir := target.Sub(this.Player.Transform.Position).Normalize()
	body := this.createBox(dir)
	this.bodies = append(this.bodies, body)
	this.times = append(this.times, FREEZE_TIME)
	this.Ammo--
}

func (this *FreezeWeapon) createBox(dir mgl32.Vec2) *box2d.B2Body {
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

	body.SetLinearVelocity(physics2d.ToBox2DDirection(dir.Mul(FREEZE_VELOCITY)))
	return body
}

func (this *FreezeWeapon) Terminate() {
	this.NilWeapon.Terminate()
}
