package main

import (
	"github.com/ByteArena/box2d"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/physics2d"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

const BALL_WEAPON_RADIUS float32 = 16.0
const BALL_WEAPON_FRICTION float64 = 1.0
const BALL_WEAPON_WEIGHT float64 = 0.5
const BALL_WEAPON_RESTITUTION float64 = 0.0
const BALL_WEAPON_VELOCITY float32 = 50.0
const BALL_WEAPON_AMMO uint32 = 10
const BALL_WEAPON_ANGLE_VELOCITY float32 = 120.0

type BallWeapon struct {
	NilWeapon

	bodies []*box2d.B2Body
	vels   []float64
}

func (this *BallWeapon) OnAdd(p *Player) {
	this.NilWeapon.OnAdd(p)
	this.tex.SetAsTarget()
	gohome.Render.ClearScreen(gohome.Color{100, 20, 255, 255})
	this.tex.UnsetAsTarget()
	this.Ammo = BALL_WEAPON_AMMO

	gohome.UpdateMgr.AddObject(this)
}

func (this *BallWeapon) Use(target mgl32.Vec2) {
	dir := target.Sub(this.Player.Transform.Position).Normalize()
	this.bodies = append(this.bodies, this.createBall(dir))
	var vel float64
	if dir.X() > 0.0 {
		vel = -float64(mgl32.DegToRad(BALL_WEAPON_ANGLE_VELOCITY))
	} else {
		vel = float64(mgl32.DegToRad(BALL_WEAPON_ANGLE_VELOCITY))
	}
	this.vels = append(this.vels, vel)
	this.Ammo--
}

func (this *BallWeapon) Update(delta_time float32) {
	for i := 0; i < len(this.bodies); i++ {
		b := this.bodies[i]
		v := this.vels[i]
		av := b.GetAngularVelocity()
		if (v > 0.0 && av < v) || (v < 0.0 && av > v) {
			b.SetAngularVelocity(v)
		}
	}
}

func (this *BallWeapon) createBall(dir mgl32.Vec2) *box2d.B2Body {
	pos := this.Player.Transform.Position.Add(dir.Mul(PLAYER_WIDTH * 2.0))

	bodyDef := box2d.MakeB2BodyDef()
	bodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	bodyDef.Position = physics2d.ToBox2DCoordinates(pos)
	bodyDef.Angle = -float64(dir.Angle())
	fdef := box2d.MakeB2FixtureDef()
	fdef.Friction = BALL_WEAPON_FRICTION
	fdef.Density = 1.0 / (2.0 * math.Pi * physics2d.ScalarToBox2D(BALL_WEAPON_RADIUS) * physics2d.ScalarToBox2D(BALL_WEAPON_RADIUS)) * BALL_WEAPON_WEIGHT
	fdef.Restitution = BALL_WEAPON_RESTITUTION
	fdef.Filter.CategoryBits = WEAPON_CATEGORY
	shape := box2d.MakeB2CircleShape()
	shape.SetRadius(physics2d.ScalarToBox2D(BALL_WEAPON_RADIUS))
	fdef.Shape = &shape
	body := this.Player.PhysicsMgr.World.CreateBody(&bodyDef)
	body.CreateFixtureFromDef(&fdef)

	body.SetLinearVelocity(physics2d.ToBox2DDirection(dir.Mul(BALL_WEAPON_VELOCITY)))
	return body
}
