package main

import (
	"github.com/ByteArena/box2d"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/physics2d"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

const (
	ENEMY_RADIUS      float32 = 10.0
	ENEMY_FRICTION    float64 = 1.0
	ENEMY_RESTITUTION float64 = 0.0
	ENEMY_WEIGHT      float64 = 0.1
	ENEMY_VELOCITY    float32 = 30.0

	ENEMY_SMALL_SENSOR_WIDTH    float32 = 3.0
	ENEMY_SMALL_SENSOR_HEIGHT   float32 = 5.0
	ENEMY_SMALL_SENSOR_OFFSET_X float32 = -10.0 - ENEMY_SMALL_SENSOR_WIDTH/2.0
	ENEMY_SMALL_SENSOR_OFFSET_Y float32 = 15.0 - ENEMY_SMALL_SENSOR_HEIGHT/2.0

	ENEMY_BIG_SENSOR_WIDTH    float32 = 1.5
	ENEMY_BIG_SENSOR_HEIGHT   float32 = 10.0
	ENEMY_BIG_SENSOR_OFFSET_X float32 = ENEMY_SMALL_SENSOR_OFFSET_X
	ENEMY_BIG_SENSOR_OFFSET_Y float32 = 0.0
)

type Enemy struct {
	gohome.Sprite2D
	Body      *box2d.B2Body
	connector physics2d.PhysicsConnector2D
	Player    *Player
	direction bool
}

func (this *Enemy) Init(pos mgl32.Vec2, player *Player) {
	this.Sprite2D.Init("")
	this.Transform.Position = pos
	this.Player = player
	this.direction = RIGHT

	this.createBody()

	gohome.UpdateMgr.AddObject(this)
	gohome.RenderMgr.AddObject(this)
	this.connector.Init(this.Transform, this.Body)
	gohome.UpdateMgr.AddObject(&this.connector)
}

func (this *Enemy) createBody() {
	bdef := box2d.MakeB2BodyDef()
	bdef.Type = box2d.B2BodyType.B2_dynamicBody
	bdef.Position = physics2d.ToBox2DCoordinates(this.Transform.Position)
	bdef.FixedRotation = true

	radius := physics2d.ScalarToBox2D(ENEMY_RADIUS)

	fdef := box2d.MakeB2FixtureDef()
	fdef.Filter.CategoryBits = ENEMY_CATEGORY
	fdef.Filter.MaskBits = 0xffff
	fdef.Friction = ENEMY_FRICTION
	fdef.Density = 1.0 / (2.0 * math.Pi * radius * radius) * ENEMY_WEIGHT
	fdef.Restitution = ENEMY_RESTITUTION

	shape := box2d.MakeB2CircleShape()
	shape.SetRadius(radius)

	fdef.Shape = &shape

	this.Body = this.Player.PhysicsMgr.World.CreateBody(&bdef)
	this.Body.CreateFixtureFromDef(&fdef)

	fdef.IsSensor = true
	fdef.Filter.CategoryBits = ENEMY_SMALL_LEFT_SENSOR_CATEGORY
	sshape := box2d.MakeB2PolygonShape()
	offset := physics2d.ToBox2DDirection([2]float32{ENEMY_SMALL_SENSOR_OFFSET_X, ENEMY_SMALL_SENSOR_OFFSET_Y})
	sshape.SetAsBox(physics2d.ScalarToBox2D(ENEMY_SMALL_SENSOR_WIDTH)/2.0, physics2d.ScalarToBox2D(ENEMY_SMALL_SENSOR_HEIGHT)/2.0)
	for i := 0; i < 4; i++ {
		v := &sshape.M_vertices[i]
		*v = box2d.B2Vec2Add(*v, offset)
	}
	fdef.Shape = &sshape

	this.Body.CreateFixtureFromDef(&fdef)

	for i := 0; i < 4; i++ {
		v := &sshape.M_vertices[i]
		v.X -= offset.X * 2.0
	}
	fdef.Filter.CategoryBits = ENEMY_SMALL_RIGHT_SENSOR_CATEGORY

	this.Body.CreateFixtureFromDef(&fdef)

	sshape.SetAsBox(physics2d.ScalarToBox2D(ENEMY_BIG_SENSOR_WIDTH)/2.0, physics2d.ScalarToBox2D(ENEMY_BIG_SENSOR_HEIGHT)/2.0)
	offset = physics2d.ToBox2DDirection([2]float32{ENEMY_BIG_SENSOR_OFFSET_X, ENEMY_BIG_SENSOR_OFFSET_Y})
	for i := 0; i < 4; i++ {
		v := &sshape.M_vertices[i]
		*v = box2d.B2Vec2Add(*v, offset)
	}
	fdef.Filter.CategoryBits = ENEMY_BIG_LEFT_SENSOR_CATEGORY

	this.Body.CreateFixtureFromDef(&fdef)

	for i := 0; i < 4; i++ {
		v := &sshape.M_vertices[i]
		v.X -= offset.X * 2.0
	}
	fdef.Filter.CategoryBits = ENEMY_BIG_RIGHT_SENSOR_CATEGORY

	this.Body.CreateFixtureFromDef(&fdef)
}

func (this *Enemy) updateDirection() {
	var sl, sr, bl, br bool = false, false, false, false
	for ce := this.Body.GetContactList(); ce != nil; ce = ce.Next {
		c := ce.Contact
		if !c.IsTouching() {
			continue
		}
		fa := c.GetFixtureA()
		fb := c.GetFixtureB()
		if fa.GetFilterData().CategoryBits&ENEMY_CATEGORY != 0 ||
			fb.GetFilterData().CategoryBits&ENEMY_CATEGORY != 0 {
			continue
		}
		if fb.GetFilterData().CategoryBits&ENEMY_SENSOR_CATEGORY != 0 {
			fa, fb = fb, fa
		}
		if fb.GetFilterData().CategoryBits != GROUND_CATEGORY &&
			fb.GetFilterData().CategoryBits != WEAPON_CATEGORY {
			continue
		}

		switch fa.GetFilterData().CategoryBits {
		case ENEMY_SMALL_LEFT_SENSOR_CATEGORY:
			sl = true
		case ENEMY_SMALL_RIGHT_SENSOR_CATEGORY:
			sr = true
		case ENEMY_BIG_LEFT_SENSOR_CATEGORY:
			bl = true
		case ENEMY_BIG_RIGHT_SENSOR_CATEGORY:
			br = true
		}
	}
	if !sl {
		this.direction = RIGHT
	}
	if !sr {
		this.direction = LEFT
	}
	if bl {
		this.direction = RIGHT
	}
	if br {
		this.direction = LEFT
	}
}

func (this *Enemy) updateVelocity() {
	vel := this.Body.GetLinearVelocity()
	evel := physics2d.ScalarToBox2D(ENEMY_VELOCITY)
	if this.direction == RIGHT {
		this.Body.SetLinearVelocity(box2d.B2Vec2{evel, vel.Y})
	} else {
		this.Body.SetLinearVelocity(box2d.B2Vec2{-evel, vel.Y})
	}
}

func (this *Enemy) Update(delta_time float32) {
	this.updateDirection()
	this.updateVelocity()
}

func (this *Enemy) Terminate() {
	gohome.UpdateMgr.RemoveObject(this)
	gohome.RenderMgr.RemoveObject(this)
	gohome.UpdateMgr.RemoveObject(&this.connector)
}
