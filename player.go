package main

import (
	"github.com/ByteArena/box2d"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/physics2d"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

const (
	PLAYER_RESTITUITION float64 = 0.0
	PLAYER_FRICTION     float64 = 1.0
	PLAYER_HEIGHT       float32 = 28.0
	PLAYER_WIDTH        float32 = 13.0
	PLAYER_VELOCITY     float32 = 500.0
	PLAYER_JUMP_FORCE   float32 = 25.0
	PLAYER_DAMPING      float64 = 0.0
	PLAYER_MAX_VELOCITY float32 = 75.0
	PLAYER_WEIGHT       float64 = 0.08

	PLAYER_FEET_SENSOR_WIDTH    float32 = PLAYER_WIDTH / 2.0
	PLAYER_FEET_SENSOR_HEIGHT   float32 = 5.0
	PLAYER_FEET_SENSOR_OFFSET_X float32 = 0.0
	PLAYER_FEET_SENSOR_OFFSET_Y float32 = 16.0 - PLAYER_FEET_SENSOR_HEIGHT/2.0

	PLAYER_ENEMY_BOUNCE float32 = -100.0

	NO_ANIM    uint8 = 0
	ANIM_WALK  uint8 = 1
	ANIM_FALL  uint8 = 2
	ANIM_SHOOT uint8 = 3

	PLAYER_FRAME_WIDTH  float32 = 14.0
	PLAYER_FRAME_HEIGHT float32 = 29.0
	PLAYER_FRAME_TIME   float32 = 1.0 / 7.0

	PLAYER_STAND_THRESHOLD float32 = 10.0
	PLAYER_PREVX_THRESHOLD float32 = 2.0
	PLAYER_JUMP_THRESHOLD  float32 = 10.0
)

type Player struct {
	gohome.Sprite2D
	connector       physics2d.PhysicsConnector2D
	body            *box2d.B2Body
	targetCameraPos mgl32.Vec2
	PhysicsMgr      *physics2d.PhysicsManager2D
	Inventory       InventoryBar

	weapons       []Weapon
	currentWeapon uint8
	terminated    bool

	currentAnimation *gohome.Tweenset
	currentAnim      uint8

	walkAnimation  gohome.Tweenset
	fallAnimation  gohome.Tweenset
	shootAnimation gohome.Tweenset
}

func (this *Player) Init(pos mgl32.Vec2, pmgr *physics2d.PhysicsManager2D) {
	this.Sprite2D.Init("Player")
	this.Transform.Position = pos
	this.Transform.Origin = [2]float32{0.5, 0.5}

	this.createBody(pmgr)
	this.connector.Init(this.Transform, this.body)

	gohome.UpdateMgr.AddObject(this)
	gohome.UpdateMgr.AddObject(&this.connector)
	gohome.RenderMgr.AddObject(this)

	this.PhysicsMgr = pmgr

	this.Inventory.Init()
	this.addWeapons()
	this.setupAnimations()

	this.terminated = false
}

func (this *Player) setupAnimations() {
	this.walkAnimation = gohome.SpriteAnimation2DRegions([]gohome.TextureRegion{
		gohome.TextureRegion{
			[2]float32{0, 0},
			[2]float32{PLAYER_FRAME_WIDTH * 1, PLAYER_FRAME_HEIGHT * 1},
		},
		gohome.TextureRegion{
			[2]float32{PLAYER_FRAME_WIDTH, PLAYER_FRAME_HEIGHT},
			[2]float32{PLAYER_FRAME_WIDTH * 2, PLAYER_FRAME_HEIGHT * 2},
		},
		gohome.TextureRegion{
			[2]float32{PLAYER_FRAME_WIDTH, 0},
			[2]float32{PLAYER_FRAME_WIDTH * 2, PLAYER_FRAME_HEIGHT * 1},
		},
		gohome.TextureRegion{
			[2]float32{0, PLAYER_FRAME_HEIGHT},
			[2]float32{PLAYER_FRAME_WIDTH * 1, PLAYER_FRAME_HEIGHT * 2},
		},
	}, PLAYER_FRAME_TIME)
	this.fallAnimation = gohome.SpriteAnimation2DRegions([]gohome.TextureRegion{
		gohome.TextureRegion{
			[2]float32{0, PLAYER_FRAME_HEIGHT * 2},
			[2]float32{PLAYER_FRAME_WIDTH * 1, PLAYER_FRAME_HEIGHT * 3},
		},
		gohome.TextureRegion{
			[2]float32{PLAYER_FRAME_WIDTH * 1, PLAYER_FRAME_HEIGHT * 2},
			[2]float32{PLAYER_FRAME_WIDTH * 2, PLAYER_FRAME_HEIGHT * 3},
		},
	}, PLAYER_FRAME_TIME)
	this.shootAnimation = gohome.SpriteAnimation2DRegions([]gohome.TextureRegion{
		gohome.TextureRegion{
			[2]float32{0, PLAYER_FRAME_HEIGHT * 3},
			[2]float32{PLAYER_FRAME_WIDTH * 1, PLAYER_FRAME_HEIGHT * 4},
		},
		gohome.TextureRegion{
			[2]float32{PLAYER_FRAME_WIDTH * 1, PLAYER_FRAME_HEIGHT * 3},
			[2]float32{PLAYER_FRAME_WIDTH * 2, PLAYER_FRAME_HEIGHT * 4},
		},
	}, PLAYER_FRAME_TIME)

	this.walkAnimation.Loop = true
	this.fallAnimation.Loop = true

	this.walkAnimation.SetParent(&this.Sprite2D)
	this.fallAnimation.SetParent(&this.Sprite2D)
	this.shootAnimation.SetParent(&this.Sprite2D)

	gohome.UpdateMgr.AddObject(&this.walkAnimation)
	gohome.UpdateMgr.AddObject(&this.fallAnimation)
	gohome.UpdateMgr.AddObject(&this.shootAnimation)

	this.StopAnimation()
}

func (this *Player) SetAnimation(anim uint8) {
	if anim == this.currentAnim {
		return
	}

	this.StopAnimation()

	switch anim {
	case ANIM_WALK:
		this.currentAnimation = &this.walkAnimation
	case ANIM_FALL:
		this.currentAnimation = &this.fallAnimation
	case ANIM_SHOOT:
		this.currentAnimation = &this.shootAnimation
	}

	this.currentAnimation.Start()
	this.currentAnim = anim
}

func (this *Player) StopAnimation() {
	this.walkAnimation.Stop()
	this.fallAnimation.Stop()
	this.shootAnimation.Stop()

	this.TextureRegion = gohome.TextureRegion{
		[2]float32{0, 0},
		[2]float32{PLAYER_FRAME_WIDTH * 1, PLAYER_FRAME_HEIGHT * 1},
	}
	this.Transform.Size = [2]float32{PLAYER_FRAME_WIDTH, PLAYER_FRAME_HEIGHT}
	this.currentAnim = NO_ANIM
}

func (this *Player) updateAnimation() {
	x := this.body.GetLinearVelocity().X
	px := physics2d.ScalarToPixel(math.Abs(x))

	this.walkAnimation.LoopBackwards = (px > 0.0 && this.Flip == gohome.FLIP_HORIZONTAL) || (px < 0.0 && this.Flip == gohome.FLIP_NONE)

	if this.shootAnimation.Done() {
		if this.IsGrounded() {
			if px < PLAYER_STAND_THRESHOLD {
				this.StopAnimation()
			} else {
				this.SetAnimation(ANIM_WALK)
			}
		} else {
			this.SetAnimation(ANIM_FALL)
		}
	}
}

func (this *Player) addWeapons() {
	this.currentWeapon = 0
	this.addWeapon(&DefaultWeapon{})
	this.addWeapon(&FreezeWeapon{})
	this.addWeapon(&BallWeapon{})
	this.addWeapon(&MoveWeapon{})
	this.addWeapon(&DeleteWeapon{})
}

func (this *Player) createBody(pmgr *physics2d.PhysicsManager2D) {
	bdef := box2d.MakeB2BodyDef()
	bdef.FixedRotation = true
	bdef.Type = box2d.B2BodyType.B2_dynamicBody
	bdef.Position = physics2d.ToBox2DCoordinates(this.Transform.Position)

	radius := physics2d.ScalarToBox2D(PLAYER_WIDTH / 2.0)

	fdef := box2d.MakeB2FixtureDef()
	fdef.Density = 1.0 / (2.0 * math.Pi * radius * radius) * PLAYER_WEIGHT
	fdef.Friction = PLAYER_FRICTION
	fdef.Restitution = PLAYER_RESTITUITION
	fdef.Filter.CategoryBits = PLAYER_FEET_CATEGORY

	circleShape := box2d.MakeB2CircleShape()
	circleShape.SetRadius(radius)
	circleShape.M_p = physics2d.ToBox2DDirection([2]float32{0.0, PLAYER_HEIGHT / 4.0})

	fdef.Shape = &circleShape

	this.body = pmgr.World.CreateBody(&bdef)
	this.body.CreateFixtureFromDef(&fdef)

	fdef.Friction = 0.0
	fdef.Filter.CategoryBits = PLAYER_CATEGORY
	circleShape.M_p = physics2d.ToBox2DDirection([2]float32{0.0, -PLAYER_HEIGHT / 4.0})

	this.body.CreateFixtureFromDef(&fdef)

	boxShape := box2d.MakeB2PolygonShape()
	w, h := physics2d.ScalarToBox2D(PLAYER_WIDTH/2.0), physics2d.ScalarToBox2D(PLAYER_HEIGHT/4.0)
	boxShape.SetAsBox(w, h)

	fdef.Shape = &boxShape
	fdef.Density = 1.0 / (w * 2.0 * h * 2.0) * PLAYER_WEIGHT

	this.body.CreateFixtureFromDef(&fdef)

	boxShape.SetAsBox(physics2d.ScalarToBox2D(PLAYER_FEET_SENSOR_WIDTH)/2.0, physics2d.ScalarToBox2D(PLAYER_FEET_SENSOR_HEIGHT)/2.0)
	offset := physics2d.ToBox2DDirection([2]float32{PLAYER_FEET_SENSOR_OFFSET_X, PLAYER_FEET_SENSOR_OFFSET_Y})
	for i := 0; i < 4; i++ {
		v := &boxShape.M_vertices[i]
		*v = box2d.B2Vec2Add(*v, offset)
	}
	fdef.IsSensor = true
	fdef.Filter.CategoryBits = PLAYER_FEET_SENSOR_CATEGORY
	this.body.CreateFixtureFromDef(&fdef)

	this.body.SetLinearDamping(PLAYER_DAMPING)
}

func (this *Player) updateVelocity(delta_time float32) {
	vel := this.body.GetLinearVelocity()
	pvel := physics2d.ToPixelDirection(vel).X()
	if gohome.InputMgr.IsPressed(KEY_RIGHT) {
		if pvel < PLAYER_MAX_VELOCITY {
			force := physics2d.ToBox2DDirection([2]float32{PLAYER_VELOCITY * delta_time, 0.0})
			vel.X += force.X
			this.body.SetLinearVelocity(vel)
		}
	} else if gohome.InputMgr.IsPressed(KEY_LEFT) {
		if pvel > -PLAYER_MAX_VELOCITY {
			force := physics2d.ToBox2DDirection([2]float32{-PLAYER_VELOCITY * delta_time, 0.0})
			vel.X += force.X
			this.body.SetLinearVelocity(vel)
		}
	}
}

func (this *Player) handleJump() {
	if (gohome.InputMgr.JustPressed(KEY_JUMP) || gohome.InputMgr.JustPressed(KEY_JUMP1)) && this.IsGrounded() {
		this.body.ApplyLinearImpulseToCenter(physics2d.ToBox2DDirection([2]float32{0.0, -PLAYER_JUMP_FORCE}), true)
	}
}

const UP = true
const DOWN = false

func (this *Player) handleAngle(mpos mgl32.Vec2) {
	angle := 360.0 - mgl32.RadToDeg(mpos.Sub(this.Transform.Position).Angle())
	if angle > 90.0 && angle < 270.0 {
		this.Flip = gohome.FLIP_HORIZONTAL
	} else {
		this.Flip = gohome.FLIP_NONE
	}
}

func (this *Player) handleWeapon() {
	mpos := gohome.InputMgr.Mouse.ToWorldPosition2D()
	this.handleAngle(mpos)
	w := this.weapons[this.currentWeapon]
	if gohome.InputMgr.JustPressed(KEY_SHOOT) && w.GetAmmo() > 0 {
		w.Use(mpos)
		if this.currentAnim == NO_ANIM {
			this.SetAnimation(ANIM_SHOOT)
		}
	}

	if gohome.InputMgr.Mouse.Wheel[1] > 0 {
		for i := 0; i < int(mgl32.Abs(float32(gohome.InputMgr.Mouse.Wheel[1]))); i++ {
			this.changeWeapon(UP)
		}
	} else if gohome.InputMgr.Mouse.Wheel[1] < 0 {
		for i := 0; i < int(mgl32.Abs(float32(gohome.InputMgr.Mouse.Wheel[1]))); i++ {
			this.changeWeapon(DOWN)
		}
	}
}

func (this *Player) addWeapon(w Weapon) {
	w.OnAdd(this)
	if len(this.weapons) == 0 {
		w.OnChange()
	}
	this.weapons = append(this.weapons, w)
	this.Inventory.AddWeapon(w)
}

func (this *Player) changeWeapon(dir bool) {
	w := this.weapons[this.currentWeapon]
	w.Terminate()
	if dir == UP {
		this.currentWeapon++
	} else {
		if this.currentWeapon == 0 {
			this.currentWeapon = uint8(len(this.weapons) - 1)
		} else {
			this.currentWeapon--
		}
	}

	if this.currentWeapon > uint8(len(this.weapons)-1) {
		this.currentWeapon = 0
	}

	w = this.weapons[this.currentWeapon]
	w.OnChange()
	this.Inventory.SetCurrent(dir)
}

func (this *Player) Update(delta_time float32) {
	this.updateVelocity(delta_time)
	this.handleJump()
	this.handleWeapon()
	this.updateCamera(delta_time)
	this.checkEnemy()
	this.updateAnimation()
}

func (this *Player) updateCamera(delta_time float32) {
	boxXNC := float64(this.Transform.Position[0] / CAMERA_BOX_WIDTH)
	boxYNC := float64(this.Transform.Position[1] / CAMERA_BOX_HEIGHT)
	boxX := float32(math.Floor(boxXNC))
	boxY := float32(math.Floor(boxYNC))
	this.targetCameraPos[0] = boxX*CAMERA_BOX_WIDTH + CAMERA_OFFSET[0]
	this.targetCameraPos[1] = boxY*CAMERA_BOX_HEIGHT + CAMERA_OFFSET[1]
	var zero float32 = 0.0
	mgl32.SetMax(&this.targetCameraPos[0], &zero)
	mgl32.SetMax(&this.targetCameraPos[1], &zero)

	Camera.Position = Camera.Position.Add(this.targetCameraPos.Sub(Camera.Position).Mul((1.0 / CAMERA_SPEED) * delta_time))
}

func (this *Player) IsGrounded() (grounded bool) {
	for ce := this.body.GetContactList(); ce != nil; ce = ce.Next {
		c := ce.Contact
		if !c.IsTouching() {
			continue
		}
		fa := c.GetFixtureA()
		fb := c.GetFixtureB()
		if fb.GetFilterData().CategoryBits&PLAYER_CATEGORY == PLAYER_CATEGORY {
			fa, fb = fb, fa
		}

		if fa.GetFilterData().CategoryBits&PLAYER_FEET_CATEGORY == PLAYER_FEET_CATEGORY {
			grounded = true
			return
		}
	}
	grounded = false
	return
}

func (this *Player) checkEnemy() {
	var fc, bc bool = false, false
	var enemy *Enemy
	for ce := this.body.GetContactList(); ce != nil; ce = ce.Next {
		c := ce.Contact
		if !c.IsTouching() {
			continue
		}
		fa := c.GetFixtureA()
		fb := c.GetFixtureB()

		if fb.GetFilterData().CategoryBits&PLAYER_CATEGORY != 0 {
			fa, fb = fb, fa
		}
		if fb.GetFilterData().CategoryBits != ENEMY_CATEGORY {
			continue
		}

		switch fa.GetFilterData().CategoryBits {
		case PLAYER_CATEGORY:
			bc = true
		case PLAYER_FEET_SENSOR_CATEGORY, PLAYER_FEET_CATEGORY:
			enemy = fb.GetBody().GetUserData().(*Enemy)
			fc = true
		default:
			break
		}
	}

	if bc {
		this.Terminate()
	}

	if fc {
		vel := this.body.GetLinearVelocity()
		vel.Y = -physics2d.ScalarToBox2D(PLAYER_ENEMY_BOUNCE)
		this.body.SetLinearVelocity(vel)
		enemy.Terminate()
	}
}

func (this *Player) IsMoving() bool {
	return gohome.InputMgr.IsPressed(KEY_RIGHT) || gohome.InputMgr.IsPressed(KEY_LEFT) ||
		gohome.InputMgr.JustPressed(KEY_JUMP) || gohome.InputMgr.JustPressed(KEY_JUMP1)
}

func (this *Player) Terminate() {
	if this.terminated {
		return
	}

	gohome.UpdateMgr.RemoveObject(this)
	gohome.UpdateMgr.RemoveObject(&this.connector)
	gohome.UpdateMgr.RemoveObject(&this.walkAnimation)
	gohome.UpdateMgr.RemoveObject(&this.fallAnimation)
	gohome.RenderMgr.RemoveObject(this)

	this.weapons[this.currentWeapon].Terminate()
	this.Inventory.Terminate()
	this.PhysicsMgr.World.DestroyBody(this.body)
	this.terminated = true
}
