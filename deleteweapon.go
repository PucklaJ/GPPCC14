package main

import (
	"github.com/ByteArena/box2d"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/physics2d"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/image/colornames"
)

const (
	DELETE_WEAPON_DISTANCE float32 = 400.0
	DELETE_WEAPON_AMMO     uint32  = 10

	DELETE_RAYS_SPEED float32 = 0.3
	DELETE_RAYS_WIDTH float32 = 5.0

	DELETE_WEAPON_OFFSET_X float32 = 5.0
	DELETE_WEAPON_OFFSET_Y float32 = -2.0
)

type DeleteWeapon struct {
	NilWeapon

	rays []*DeleteRay
}

func (this *DeleteWeapon) OnAdd(p *Player) {
	this.Sprite2D.Init("DeleteWeapon")
	this.Transform.Origin = [2]float32{0.5, 0.5}

	this.NilWeapon.OnAdd(p)
	this.Ammo = DELETE_WEAPON_AMMO

	gohome.UpdateMgr.AddObject(this)
}

func (this *DeleteWeapon) GetInventoryTexture() gohome.Texture {
	return gohome.ResourceMgr.GetTexture("DeleteWeaponInv")
}

func (this *DeleteWeapon) castRay(dir mgl32.Vec2) {
	pmgr := this.Player.PhysicsMgr
	w := &pmgr.World
	input := box2d.MakeB2RayCastInput()
	input.P1 = physics2d.ToBox2DCoordinates(this.Player.Transform.Position)
	input.P2 = physics2d.ToBox2DCoordinates(this.Player.Transform.Position.Add(dir.Mul(DELETE_WEAPON_DISTANCE)))
	input.MaxFraction = 1.0
	output := box2d.MakeB2RayCastOutput()
	var bodies []*box2d.B2Body
	for b := w.GetBodyList(); b != nil; b = b.GetNext() {
		for f := b.GetFixtureList(); f != nil; f = f.GetNext() {
			if f.GetFilterData().CategoryBits&WEAPON_CATEGORY == WEAPON_CATEGORY {
				hits := f.RayCast(&output, input, 0)
				if hits {
					bodies = append(bodies, b)
				}
			}
		}
	}

	for i := 0; i < len(bodies); i++ {
		w.DestroyBody(bodies[i])
	}
}

func (this *DeleteWeapon) Update(delta_time float32) {
	off := [2]float32{DELETE_WEAPON_OFFSET_X, DELETE_WEAPON_OFFSET_Y}
	this.Flip = this.Player.Flip
	if this.Flip == gohome.FLIP_HORIZONTAL {
		off[0] = -off[0]
	}
	this.Transform.Position = this.Player.Transform.Position.Add(this.Player.GetWeaponOffset()).Add(off)
}

func (this *DeleteWeapon) Use(target mgl32.Vec2) {
	dir := target.Sub(this.Player.Transform.Position).Normalize()

	var ray DeleteRay
	ray.Init()
	ray.Transform.Position = this.Player.Transform.Position.Add(dir.Mul(DELETE_WEAPON_DISTANCE / 2.0)).Add(this.Player.GetWeaponOffset()).Sub([2]float32{0.0, DELETE_RAYS_WIDTH / 2.0})

	// this.Ammo--

	ray.Transform.Rotation = mgl32.RadToDeg(-dir.Angle())
	this.rays = append(this.rays, &ray)
	this.castRay(dir)
}

func (this *DeleteWeapon) Render() {
	for i := 0; i < len(this.rays); i++ {
		if this.rays[i].time >= DELETE_RAYS_SPEED {
			gohome.UpdateMgr.RemoveObject(this.rays[i])
		}
	}
}

type DeleteRay struct {
	gohome.Shape2D
	time float32
}

func (this *DeleteRay) Init() {
	this.Shape2D.Init()
	var rect gohome.Rectangle2D
	rect[0].Make([2]float32{-1.0, 1.0}, colornames.Red)
	rect[1].Make([2]float32{1.0, 1.0}, colornames.Red)
	rect[2].Make([2]float32{1.0, -1.0}, colornames.Red)
	rect[3].Make([2]float32{-1.0, -1.0}, colornames.Red)
	tris := rect.ToTriangles()
	this.AddTriangles(tris[:])
	this.Load()
	this.SetDrawMode(gohome.DRAW_MODE_TRIANGLES)

	gohome.RenderMgr.AddObject(this)
	gohome.UpdateMgr.AddObject(this)

	this.Transform.Size = [2]float32{DELETE_WEAPON_DISTANCE, DELETE_RAYS_WIDTH}
	this.Depth = DELETE_RAY_DEPTH
}

func (this *DeleteRay) Update(delta_time float32) {
	this.time += delta_time
	if this.time >= DELETE_RAYS_SPEED {
		gohome.RenderMgr.RemoveObject(this)
	}
	width := DELETE_RAYS_WIDTH * (1.0 - this.time/DELETE_RAYS_SPEED)
	this.Transform.Size[1] = width
}
