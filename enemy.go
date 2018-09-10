package main

import (
	"github.com/ByteArena/box2d"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/physics2d"
	"github.com/PucklaMotzer09/mathgl/mgl32"
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

	ENEMY_FRAME_WIDTH  float32 = 23.0
	ENEMY_FRAME_HEIGHT float32 = 20.0
	ENEMY_FRAME_TIME   float32 = 1.0 / 5.0

	ENEMY_OFFSET_X         float32 = 1.0
	ENEMY_OFFSET_Y         float32 = 0.0
	ENEMY_DESTRUCTION_TIME float32 = 0.5
	ENEMY_FALL_DESTRUCTION float32 = 175.0
)

type Enemy struct {
	gohome.Sprite2D
	Body            *box2d.B2Body
	connector       physics2d.PhysicsConnector2D
	Player          *Player
	direction       bool
	terminated      bool
	destructionTime float32
	destructed      bool
	paused          bool

	anim gohome.Tweenset
}

func (this *Enemy) Init(pos mgl32.Vec2, player *Player) {
	this.Sprite2D.Init("Enemy")
	this.Transform.Position = pos
	this.Transform.Origin = [2]float32{0.5, 0.5}
	this.Player = player
	this.direction = RIGHT
	this.terminated = false

	this.createBody()

	gohome.UpdateMgr.AddObject(this)
	gohome.RenderMgr.AddObject(this)
	this.connector.Init(this.Transform, this.Body)
	this.connector.Offset = [2]float32{ENEMY_OFFSET_X, ENEMY_OFFSET_Y}
	gohome.UpdateMgr.AddObject(&this.connector)

	this.anim = gohome.SpriteAnimation2D(this.Texture, 3, 4, ENEMY_FRAME_TIME)
	this.anim.Loop = true
	this.anim.SetParent(&this.Sprite2D)
	this.anim.Start()
	gohome.UpdateMgr.AddObject(&this.anim)
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
	this.Body.SetUserData(this)
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

type Explosion struct {
	gohome.Sprite2D
	anim gohome.Tweenset
}

func (this *Explosion) Init(texName string) {
	this.Sprite2D.Init(texName)
	gohome.ResourceMgr.GetSound("Explosion").Play(false)
}

func (this *Explosion) Update(delta_time float32) {
	if this.anim.Done() {
		gohome.RenderMgr.RemoveObject(this)
		gohome.UpdateMgr.RemoveObject(&this.anim)
		gohome.UpdateMgr.RemoveObject(this)
	}
}

func (this *Enemy) Die() {
	var exp Explosion
	exp.Init("Explosion")
	exp.Transform.Origin = [2]float32{0.5, 0.5}
	exp.Transform.Position = this.Transform.Position
	exp.anim = gohome.SpriteAnimation2D(exp.Texture, 5, 1, 1.0/8.0)
	exp.anim.Tweens = append(exp.anim.Tweens, &gohome.TweenWait{
		TweenType: gohome.TWEEN_TYPE_AFTER_PREVIOUS,
		Time:      0.5,
	})
	exp.anim.SetParent(&exp.Sprite2D)
	exp.anim.Start()
	gohome.RenderMgr.AddObject(&exp)
	gohome.UpdateMgr.AddObject(&exp.anim)
	gohome.UpdateMgr.AddObject(&exp)
}

func (this *Enemy) checkCollisions() {
	var sl, sr, bl, br bool = false, false, false, false
	for ce := this.Body.GetContactList(); ce != nil; ce = ce.Next {
		c := ce.Contact
		if !c.IsTouching() {
			continue
		}
		fa := c.GetFixtureA()
		fb := c.GetFixtureB()

		if fb.GetBody() == this.Body {
			fa, fb = fb, fa
		}

		if fb.GetBody() == this.Body {
			continue
		}

		if fa.GetFilterData().CategoryBits&ENEMY_SENSOR_CATEGORY == 0 {
			continue
		}

		if fb.GetFilterData().CategoryBits != GROUND_CATEGORY &&
			fb.GetFilterData().CategoryBits != WEAPON_CATEGORY &&
			fb.GetFilterData().CategoryBits != ENEMY_CATEGORY {
			continue
		}

		switch fa.GetFilterData().CategoryBits {
		case ENEMY_SMALL_LEFT_SENSOR_CATEGORY:
			if fb.GetFilterData().CategoryBits == ENEMY_CATEGORY {
				continue
			}
			sl = true
		case ENEMY_SMALL_RIGHT_SENSOR_CATEGORY:
			if fb.GetFilterData().CategoryBits == ENEMY_CATEGORY {
				continue
			}
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

	if bl && br {
		this.destructed = true
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

func (this *Enemy) updateAnimation() {
	x := this.Body.GetLinearVelocity().X
	if x < 0.0 {
		this.Flip = gohome.FLIP_HORIZONTAL
		this.connector.Offset[0] = -ENEMY_OFFSET_X
	} else {
		this.Flip = gohome.FLIP_NONE
		this.connector.Offset[0] = ENEMY_OFFSET_X
	}
}

func (this *Enemy) Update(delta_time float32) {
	if this.paused {
		return
	}

	this.checkCollisions()
	this.updateVelocity()
	this.updateAnimation()

	if this.destructed {
		this.destructionTime += delta_time
	}

	if (this.destructed && this.destructionTime >= ENEMY_DESTRUCTION_TIME) || physics2d.ScalarToPixel(this.Body.GetLinearVelocity().Y) < -ENEMY_FALL_DESTRUCTION {
		this.Die()
		this.Terminate()
	}
}

func (this *Enemy) Terminate() {
	if this.terminated {
		return
	}

	this.Player.PhysicsMgr.World.DestroyBody(this.Body)
	gohome.UpdateMgr.RemoveObject(this)
	gohome.RenderMgr.RemoveObject(this)
	gohome.UpdateMgr.RemoveObject(&this.connector)
	gohome.UpdateMgr.RemoveObject(&this.anim)

	this.terminated = true
}
