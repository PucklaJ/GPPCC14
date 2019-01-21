package main

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"strconv"
)

const (
	AMMO_TEXT_PADDING = 1
	AMMO_TEXT_SIZE    = 16
	AMMO_TEXT_SCALE   = 1
)

var ()

type AmmoText struct {
	gohome.Sprite2D

	Number    uint32
	OldNumber uint32
}

func (this *AmmoText) Init(num uint32) {
	this.Sprite2D.Init("")
	this.Number = num
	this.OldNumber = num
	this.updateTexture()
	this.Flip = gohome.FLIP_VERTICAL
}

func (this *AmmoText) updateTexture() {
	if this.Texture != nil {
		this.Texture.Terminate()
	}
	str := strconv.FormatUint(uint64(this.Number), 10)
	rt := gohome.Render.CreateRenderTexture("AmmoTextTexture", int((float32(len(str))*(AMMO_TEXT_SIZE+AMMO_TEXT_PADDING)-AMMO_TEXT_PADDING)*AMMO_TEXT_SCALE), AMMO_TEXT_SIZE*AMMO_TEXT_SCALE, 1, false, false, false, false)
	rt.SetFiltering(gohome.FILTERING_NEAREST)
	prevProj := gohome.RenderMgr.Projection2D
	rt.SetAsTarget()
	gohome.RenderMgr.SetProjection2DToTexture(rt)
	var spr gohome.Sprite2D
	spr.Init("AmmoFont")
	spr.Transform.Size = [2]float32{AMMO_TEXT_SIZE * AMMO_TEXT_SCALE, AMMO_TEXT_SIZE * AMMO_TEXT_SCALE}
	spr.NotRelativeToCamera = 0
	for i := 0; i < len(str); i++ {
		c := str[i] - 48
		reg := getRegionForNumber(c)
		spr.TextureRegion = reg
		spr.Transform.Position = [2]float32{float32(i) * (AMMO_TEXT_SIZE*AMMO_TEXT_SCALE + AMMO_TEXT_PADDING*AMMO_TEXT_SCALE), 0.0}
		gohome.RenderMgr.RenderRenderObject(&spr)
	}

	rt.UnsetAsTarget()
	gohome.RenderMgr.Projection2D = prevProj

	this.Texture = rt
	this.Transform.Size = [2]float32{float32(this.Texture.GetWidth()), float32(this.Texture.GetHeight())}
	this.TextureRegion.Max = [2]float32{float32(rt.GetWidth()), float32(rt.GetHeight())}
}

func getRegionForNumber(c byte) (reg gohome.TextureRegion) {
	reg.Max[1] = AMMO_TEXT_SIZE
	reg.Min[1] = 0.0

	reg.Min[0] = float32(c) * AMMO_TEXT_SIZE
	reg.Max[0] = reg.Min[0] + AMMO_TEXT_SIZE
	return
}

func (this *AmmoText) Render() {
	if this.OldNumber != this.Number {
		this.updateTexture()
		this.OldNumber = this.Number
	}
	this.Sprite2D.Render()
}
