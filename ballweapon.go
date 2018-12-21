package main

import (
	"github.com/ByteArena/box2d"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/GoHomeEngine/src/physics2d"
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"math"
)

const (
	BALL_WEAPON_RADIUS         float32 = 16.0
	BALL_WEAPON_FRICTION       float64 = 1.0
	BALL_WEAPON_WEIGHT         float64 = 0.5
	BALL_WEAPON_RESTITUTION    float64 = 0.0
	BALL_WEAPON_VELOCITY       float32 = 100.0
	BALL_WEAPON_AMMO           uint32  = 3
	BALL_WEAPON_ANGLE_VELOCITY float32 = 120.0

	BALL_WEAPON_OFFSET_X float32 = 2.0
	BALL_WEAPON_OFFSET_Y float32 = -2.0

	BALL_WEAPON_FRAME_TIME float32 = 1.0 / 4.0
	BALL_WEAPON_ANIM_WAIT  float32 = 1.0
)

type BallWeaponBlock struct {
	WeaponBlock
	anim gohome.Tweenset
}

func (this *BallWeaponBlock) Terminate() {
	this.WeaponBlock.Terminate()
	gohome.UpdateMgr.RemoveObject(&this.anim)
	this.Connector.Terminate()
}

type BallWeapon struct {
	NilWeapon

	bodies     []*box2d.B2Body
	vels       []float64
	ballBlocks []*BallWeaponBlock
}

func (this *BallWeapon) OnAdd(p *Player) {
	this.Sprite2D.Init("BallWeapon")
	this.Transform.Origin = [2]float32{0.5, 0.5}

	this.NilWeapon.OnAdd(p)
	this.Ammo = BALL_WEAPON_AMMO

	gohome.UpdateMgr.AddObject(this)
}

func (this *BallWeapon) GetInventoryTexture() gohome.Texture {
	return gohome.ResourceMgr.GetTexture("BallWeaponInv")
}

func (this *BallWeapon) Use(target mgl32.Vec2, energy float32) {
	dir := target.Sub(this.Player.Transform.Position).Normalize()
	this.bodies = append(this.bodies, this.createBall(dir, energy))
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

	off := [2]float32{BALL_WEAPON_OFFSET_X, BALL_WEAPON_OFFSET_Y}
	this.Flip = this.Player.Flip
	if this.Flip == gohome.FLIP_HORIZONTAL {
		off[0] = -off[0]
	}
	this.Transform.Position = this.Player.Transform.Position.Add(this.Player.GetWeaponOffset()).Add(off)
}

func (this *BallWeapon) createBall(dir mgl32.Vec2, energy float32) *box2d.B2Body {
	pos := this.Player.Transform.Position.Add(dir.Mul(PLAYER_WIDTH * 2.0))

	bodyDef := box2d.MakeB2BodyDef()
	bodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	bodyDef.Position = physics2d.ToBox2DCoordinates(pos)
	bodyDef.Angle = -float64(dir.Angle())
	fdef := box2d.MakeB2FixtureDef()
	fdef.Friction = BALL_WEAPON_FRICTION
	fdef.Density = 1.0 / (2.0 * math.Pi * physics2d.ScalarToBox2D(BALL_WEAPON_RADIUS) * physics2d.ScalarToBox2D(BALL_WEAPON_RADIUS)) * BALL_WEAPON_WEIGHT
	fdef.Restitution = BALL_WEAPON_RESTITUTION
	fdef.Filter.CategoryBits = WEAPON_CATEGORY | BALL_CATEGORY
	shape := box2d.MakeB2CircleShape()
	shape.SetRadius(physics2d.ScalarToBox2D(BALL_WEAPON_RADIUS))
	fdef.Shape = &shape
	body := this.Player.PhysicsMgr.World.CreateBody(&bodyDef)
	body.CreateFixtureFromDef(&fdef)

	body.SetLinearVelocity(physics2d.ToBox2DDirection(dir.Mul(BALL_WEAPON_VELOCITY * energy)))
	body.SetLinearVelocity(box2d.B2Vec2Add(this.Player.body.GetLinearVelocity(), body.GetLinearVelocity()))

	var spr gohome.Sprite2D
	var con physics2d.PhysicsConnector2D

	spr.Init("BallWeaponBlock")
	spr.TextureRegion.Max[0] = float32(spr.Texture.GetWidth()) / 7.0
	spr.Transform.Size[0] = spr.TextureRegion.Max[0]
	con.Init(spr.Transform, body, &PhysicsMgr)

	gohome.RenderMgr.AddObject(&spr)

	var block BallWeaponBlock
	block.Sprite = &spr
	block.Connector = &con
	block.anim = gohome.SpriteAnimation2D(spr.Texture, 7, 1, BALL_WEAPON_FRAME_TIME)
	block.anim.Tweens = append(block.anim.Tweens, &gohome.TweenRegion2D{
		TweenType: gohome.TWEEN_TYPE_AFTER_PREVIOUS,
		Destination: gohome.TextureRegion{
			[2]float32{0.0, 0.0},
			[2]float32{float32(spr.Texture.GetWidth()) / 7.0, float32(spr.Texture.GetHeight())},
		},
		Time: 0.0,
	})
	block.anim.Tweens = append(block.anim.Tweens, &gohome.TweenWait{
		Time:      BALL_WEAPON_ANIM_WAIT,
		TweenType: gohome.TWEEN_TYPE_AFTER_PREVIOUS,
	})
	block.anim.Loop = true
	block.anim.SetParent(&spr)
	block.anim.Start()
	gohome.UpdateMgr.AddObject(&block.anim)
	this.ballBlocks = append(this.ballBlocks, &block)

	body.SetUserData(this.ballBlocks[len(this.ballBlocks)-1])

	con.Update()

	return body
}

func (this *BallWeapon) OnDie() {
	gohome.UpdateMgr.RemoveObject(this)
	gohome.RenderMgr.RemoveObject(&this.NilWeapon)
}

func (this *BallWeapon) Terminate() {
	this.NilWeapon.Terminate()
	gohome.UpdateMgr.RemoveObject(this)
	for _, block := range this.ballBlocks {
		block.Terminate()
	}
}
