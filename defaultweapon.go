package main

import (
	"github.com/ByteArena/box2d"
	"github.com/PucklaMotzer09/gohomeengine/src/physics2d"
	"github.com/go-gl/mathgl/mgl32"
)

const DEFAULT_WEAPON_WIDTH float32 = 32.0
const DEFAULT_WEAPON_HEIGHT float32 = 16.0
const DEFAULT_WEAPON_FRICTION float64 = 3.0
const DEFAULT_WEAPON_WEIGHT float64 = 0.5
const DEFAULT_WEAPON_RESTITUTION float64 = 0.0
const DEFAULT_WEAPON_VELOCITY float32 = 200.0

type DefaultWeapon struct {
	NilWeapon
}

func (this *DefaultWeapon) Use(target mgl32.Vec2) {
	dir := target.Sub(this.Player.Transform.Position).Normalize()
	this.createBox(dir)
}

func (this *DefaultWeapon) createBox(dir mgl32.Vec2) {
	pos := this.Player.Transform.Position.Add(dir.Mul(PLAYER_WIDTH * 2.0))
	size := [2]float32{DEFAULT_WEAPON_WIDTH, DEFAULT_WEAPON_HEIGHT}

	bodyDef := box2d.MakeB2BodyDef()
	bodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	bodyDef.Position = physics2d.ToBox2DCoordinates(pos)
	bodyDef.Angle = -float64(dir.Angle())
	fdef := box2d.MakeB2FixtureDef()
	fdef.Friction = DEFAULT_WEAPON_FRICTION
	fdef.Density = 1.0 / (physics2d.ScalarToBox2D(DEFAULT_WEAPON_WIDTH) * physics2d.ScalarToBox2D(DEFAULT_WEAPON_HEIGHT)) * DEFAULT_WEAPON_WEIGHT
	fdef.Restitution = DEFAULT_WEAPON_RESTITUTION
	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox(physics2d.ScalarToBox2D(size[0])/2.0, physics2d.ScalarToBox2D(size[1])/2.0)
	fdef.Shape = &shape
	body := this.Player.PhysicsMgr.World.CreateBody(&bodyDef)
	body.CreateFixtureFromDef(&fdef)

	body.SetLinearVelocity(physics2d.ToBox2DDirection(dir.Mul(DEFAULT_WEAPON_VELOCITY)))
}
