package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/mathgl/mgl32"
)

type Weapon interface {
	Init(p *Player)
	Use(target mgl32.Vec2)
	GetInventoryTexture() gohome.Texture
	Terminate()
}

type NilWeapon struct {
	gohome.NilRenderObject

	Player *Player
}

func (this *NilWeapon) Init(p *Player) {
	gohome.RenderMgr.AddObject(this)
	this.Player = p
}

func (this *NilWeapon) Use(target mgl32.Vec2) {
	var shape2d gohome.Shape2D
	shape2d.Init()
	var line gohome.Line2D
	line[0].Make(this.Player.Transform.Position, gohome.Color{255, 0, 0, 255})
	line[1].Make(target, gohome.Color{255, 0, 0, 255})
	shape2d.AddLines([]gohome.Line2D{line})
	shape2d.Load()
	shape2d.SetDrawMode(gohome.DRAW_MODE_LINES)
	gohome.RenderMgr.AddObject(&shape2d)
}

func (this *NilWeapon) GetInventoryTexture() gohome.Texture {
	return nil
}

func (this *NilWeapon) Terminate() {
	gohome.RenderMgr.RemoveObject(this)
}

func (this *NilWeapon) GetType() gohome.RenderType {
	return gohome.TYPE_2D_NORMAL
}
