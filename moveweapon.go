package main

import (
	"github.com/ByteArena/box2d"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/physics2d"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

const (
	MOVE_WEAPON_AMMO           uint32  = 2
	MOVE_WEAPON_SPEED          float32 = 50.0
	MOVE_WEAPON_DISTANCE       float32 = 100.0
	MOVE_WEAPON_WIDTH          float32 = 48.0
	MOVE_WEAPON_HEIGHT         float32 = 6.0
	MOVE_WEAPON_WEIGHT         float64 = 1.0
	MOVE_WEAPON_VELOCITY       float32 = 200.0
	MOVE_WEAPON_TIME           float32 = 0.5
	MOVE_WEAPON_FRICTION       float64 = GROUND_FRICTION * 2.0
	MOVE_WEAPON_RESTITUTION    float64 = DEFAULT_WEAPON_RESTITUTION
	MOVE_WEAPON_ROTATE_SPEED   float64 = 2.0
	MOVE_WEAPON_VELOCITY_SPEED float64 = 0.5
	MOVE_WEAPON_MIN_DISTANCE   float32 = 10.0
	MOVE_WEAPON_SLOW_DISTANCE  float32 = 32.0

	MOVE_WEAPON_OFFSET_X float32 = 1.0
	MOVE_WEAPON_OFFSET_Y float32 = -2.0
)

type MoveWeapon struct {
	NilWeapon

	platforms []MovePlatform
}

func (this *MoveWeapon) OnAdd(p *Player) {
	this.Sprite2D.Init("MoveWeapon")
	this.Transform.Origin = [2]float32{0.5, 0.5}

	this.NilWeapon.OnAdd(p)
	this.Ammo = MOVE_WEAPON_AMMO

	gohome.UpdateMgr.AddObject(this)
}

func (this *MoveWeapon) GetInventoryTexture() gohome.Texture {
	return gohome.ResourceMgr.GetTexture("MoveWeaponInv")
}

func (this *MoveWeapon) Use(target mgl32.Vec2) {
	dir := target.Sub(this.Player.Transform.Position).Normalize()
	body := this.createBox(dir)

	this.platforms = append(this.platforms, MovePlatform{
		body,
		0.0,
		[2]float32{0.0, 0.0},
		false,
		box2d.B2Vec2{0.0, 0.0},
		box2d.B2Vec2{0.0, 0.0},
		dir.X() >= 0.0,
		this.Player,
	})

	// this.Ammo--
}

func (this *MoveWeapon) Update(delta_time float32) {
	for i := 0; i < len(this.platforms); i++ {
		this.platforms[i].Update(delta_time)
	}

	off := [2]float32{MOVE_WEAPON_OFFSET_X, MOVE_WEAPON_OFFSET_Y}
	this.Flip = this.Player.Flip
	if this.Flip == gohome.FLIP_HORIZONTAL {
		off[0] = -off[0]
	}
	this.Transform.Position = this.Player.Transform.Position.Add(this.Player.GetWeaponOffset()).Add(off)
}

func (this *MoveWeapon) createBox(dir mgl32.Vec2) *box2d.B2Body {
	pos := this.Player.Transform.Position.Add(dir.Mul(PLAYER_WIDTH * 2.0))
	size := [2]float32{MOVE_WEAPON_WIDTH, MOVE_WEAPON_HEIGHT}

	bodyDef := box2d.MakeB2BodyDef()
	bodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	bodyDef.Position = physics2d.ToBox2DCoordinates(pos)
	bodyDef.Angle = -float64(dir.Angle())
	fdef := box2d.MakeB2FixtureDef()
	fdef.Friction = MOVE_WEAPON_FRICTION
	fdef.Density = 1.0 / (physics2d.ScalarToBox2D(MOVE_WEAPON_WIDTH) * physics2d.ScalarToBox2D(MOVE_WEAPON_HEIGHT)) * MOVE_WEAPON_WEIGHT
	fdef.Restitution = MOVE_WEAPON_RESTITUTION
	fdef.Filter.CategoryBits = WEAPON_CATEGORY
	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox(physics2d.ScalarToBox2D(size[0])/2.0, physics2d.ScalarToBox2D(size[1])/2.0)
	fdef.Shape = &shape
	body := this.Player.PhysicsMgr.World.CreateBody(&bodyDef)
	body.CreateFixtureFromDef(&fdef)

	body.SetLinearVelocity(physics2d.ToBox2DDirection(dir.Mul(MOVE_WEAPON_VELOCITY)))
	return body
}

const (
	RIGHT = true
	LEFT  = false
)

type MovePlatform struct {
	Body               *box2d.B2Body
	Time               float32
	PrevPosition       mgl32.Vec2
	IsMoving           bool
	TargetPosition     box2d.B2Vec2
	PrevTargetPosition box2d.B2Vec2
	Direction          bool
	Player             *Player
}

func (this *MovePlatform) HoldRotation() {
	b := this.Body
	const NUM_TARGET_ANGLES uint8 = 5

	var targetAngles = [NUM_TARGET_ANGLES]float64{
		0.0,
		-math.Pi * 2.0,
		-math.Pi,
		math.Pi,
		math.Pi * 2.0,
	}
	curAngle := b.GetAngle()
	var smallestError float64 = targetAngles[0] - curAngle
	for j := uint8(1); j < NUM_TARGET_ANGLES; j++ {
		if math.Abs(targetAngles[j]-curAngle) < math.Abs(smallestError) {
			smallestError = targetAngles[j] - curAngle
		}
	}
	b.SetAngularVelocity(smallestError * MOVE_WEAPON_ROTATE_SPEED)
}

func (this *MovePlatform) HoldPosition() {
	b := this.Body
	targetPos := this.TargetPosition
	curPos := b.GetPosition()
	errorPos := box2d.B2Vec2{targetPos.X - curPos.X, targetPos.Y - curPos.Y}
	b.SetLinearVelocity(box2d.B2Vec2{
		b.GetLinearVelocity().X,
		errorPos.Y * MOVE_WEAPON_VELOCITY_SPEED,
	})
}

func (this *MovePlatform) Move() {

	minDist := physics2d.ScalarToBox2D(MOVE_WEAPON_MIN_DISTANCE)
	pos := this.Body.GetPosition()
	target := this.TargetPosition
	rel := box2d.B2Vec2Sub(target, pos)
	dist := math.Abs(rel.X)

	dist1 := math.Abs(box2d.B2Vec2Sub(this.PrevTargetPosition, pos).X)
	mdist := physics2d.ScalarToBox2D(MOVE_WEAPON_SLOW_DISTANCE)

	speed := physics2d.ScalarToBox2D(MOVE_WEAPON_SPEED)
	if this.Direction == LEFT {
		speed = -speed
	}

	var mul float64
	if dist < dist1 {
		mul = dist / mdist
	} else {
		mul = dist1 / mdist
	}
	if mul > 1.0 {
		mul = 1.0
	}
	speed *= mul

	this.Body.SetLinearVelocity(box2d.B2Vec2{
		speed,
		this.Body.GetLinearVelocity().Y,
	})

	if dist <= minDist {
		this.PrevTargetPosition = this.TargetPosition
		if this.Direction == RIGHT {
			this.TargetPosition = box2d.B2Vec2Sub(this.TargetPosition, box2d.B2Vec2{physics2d.ScalarToBox2D(MOVE_WEAPON_DISTANCE), 0.0})
		} else {
			this.TargetPosition = box2d.B2Vec2Add(this.TargetPosition, box2d.B2Vec2{physics2d.ScalarToBox2D(MOVE_WEAPON_DISTANCE), 0.0})
		}
		this.Direction = !this.Direction
	}
}

func (this *MovePlatform) Update(delta_time float32) {
	if !this.IsMoving {
		this.Time += delta_time
		if this.Time > MOVE_WEAPON_TIME {
			this.IsMoving = true
			this.TargetPosition = this.Body.GetPosition()
			dist := physics2d.ScalarToBox2D(MOVE_WEAPON_DISTANCE)
			if this.Direction == RIGHT {
				this.TargetPosition.X += dist
			} else {
				this.TargetPosition.X -= dist
			}
			this.PrevPosition = physics2d.ToPixelCoordinates(this.Body.GetPosition())
		}
	} else {
		this.Move()
		this.HoldRotation()
		this.HoldPosition()
	}
}
