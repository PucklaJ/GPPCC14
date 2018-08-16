package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"golang.org/x/image/colornames"
)

const INVENTORY_TEXTURE_SIZE float32 = 48.0
const INVENTORY_PADDING float32 = INVENTORY_TEXTURE_SIZE / 8.0

type InventoryBar struct {
	gohome.Sprite2D

	weaponTextures []gohome.Texture
	prevNumWeapons int
}

func (this *InventoryBar) Init() {
	tex := gohome.Render.CreateRenderTexture("InventoryBarTexture", uint32(INVENTORY_TEXTURE_SIZE+INVENTORY_PADDING*2.0), uint32(INVENTORY_TEXTURE_SIZE+INVENTORY_PADDING*2.0), 1, false, false, false, false)
	this.Sprite2D.InitTexture(tex)

	gohome.RenderMgr.AddObject(this)
	gohome.UpdateMgr.AddObject(this)

	this.Depth = 1
	this.NotRelativeToCamera = 0

	this.Transform.Position = gohome.Framew.WindowGetSize().Mul(0.5)
	this.Transform.Position[1] = gohome.Framew.WindowGetSize()[1] - (INVENTORY_PADDING*2.0+INVENTORY_TEXTURE_SIZE)/2.0 - INVENTORY_PADDING
	this.Transform.Origin = [2]float32{0.5, 0.5}
	this.prevNumWeapons = -1
}

func (this *InventoryBar) AddWeapon(w Weapon) {
	this.weaponTextures = append(this.weaponTextures, w.GetInventoryTexture())
}

func (this *InventoryBar) Update(delta_time float32) {
	if this.prevNumWeapons != len(this.weaponTextures) {
		this.renderInventory()
		this.prevNumWeapons = len(this.weaponTextures)
	}
}

func (this *InventoryBar) renderInventory() {
	rt := this.Texture.(gohome.RenderTexture)
	width := INVENTORY_PADDING*2 + INVENTORY_TEXTURE_SIZE

	imgWidth := uint32(len(this.weaponTextures)) * uint32(width)
	imgHeight := uint32(rt.GetHeight())

	rt.ChangeSize(imgWidth, imgHeight)
	this.Transform.Size = [2]float32{float32(imgWidth), float32(imgHeight)}
	this.TextureRegion.Max = this.Transform.Size

	rt.SetAsTarget()
	gohome.Render.ClearScreen(gohome.Color{0, 0, 0, 0})

	gohome.DrawColor = colornames.Gray
	gohome.Filled = true

	gohome.RenderMgr.SetCamera2D(nil, 0)
	proj := gohome.Ortho2DProjection{
		Left:   0.0,
		Right:  float32(rt.GetWidth()),
		Top:    0.0,
		Bottom: float32(rt.GetHeight()),
	}
	prevProj := gohome.RenderMgr.Projection2D
	gohome.RenderMgr.Projection2D = &proj

	for i := 0; i < len(this.weaponTextures); i++ {
		wt := this.weaponTextures[i]
		// fmt.Println("Texture:", i)

		x := width * float32(i)

		gohome.DrawRectangle2D(
			[2]float32{x, INVENTORY_PADDING*2.0 + INVENTORY_TEXTURE_SIZE},
			[2]float32{x + INVENTORY_PADDING, INVENTORY_PADDING*2.0 + INVENTORY_TEXTURE_SIZE},
			[2]float32{x + INVENTORY_PADDING, 0.0},
			[2]float32{x, 0.0})

		gohome.DrawRectangle2D(
			[2]float32{x + INVENTORY_PADDING + INVENTORY_TEXTURE_SIZE, INVENTORY_PADDING*2.0 + INVENTORY_TEXTURE_SIZE},
			[2]float32{x + INVENTORY_PADDING*2 + INVENTORY_TEXTURE_SIZE, INVENTORY_PADDING*2.0 + INVENTORY_TEXTURE_SIZE},
			[2]float32{x + INVENTORY_PADDING*2 + INVENTORY_TEXTURE_SIZE, 0.0},
			[2]float32{x + INVENTORY_PADDING + INVENTORY_TEXTURE_SIZE, 0.0})

		gohome.DrawRectangle2D(
			[2]float32{x + INVENTORY_PADDING, INVENTORY_PADDING},
			[2]float32{x + INVENTORY_PADDING + INVENTORY_TEXTURE_SIZE, INVENTORY_PADDING},
			[2]float32{x + INVENTORY_PADDING + INVENTORY_TEXTURE_SIZE, 0.0},
			[2]float32{x + INVENTORY_PADDING, 0.0})

		gohome.DrawRectangle2D(
			[2]float32{x + INVENTORY_PADDING, INVENTORY_PADDING*2 + INVENTORY_TEXTURE_SIZE},
			[2]float32{x + INVENTORY_PADDING + INVENTORY_TEXTURE_SIZE, INVENTORY_PADDING*2 + INVENTORY_TEXTURE_SIZE},
			[2]float32{x + INVENTORY_PADDING + INVENTORY_TEXTURE_SIZE, INVENTORY_PADDING + INVENTORY_TEXTURE_SIZE},
			[2]float32{x + INVENTORY_PADDING, INVENTORY_PADDING + INVENTORY_TEXTURE_SIZE})
		// fmt.Println("Texture:", wt)
		if wt != nil {
			var spr gohome.Sprite2D
			spr.InitTexture(wt)
			spr.Transform.Position[0] = x + INVENTORY_PADDING
			spr.Transform.Position[1] = INVENTORY_PADDING
			gohome.RenderMgr.RenderRenderObject(&spr)
			// fmt.Println("Rendering To", spr.Transform.Position)
		}
	}

	gohome.RenderMgr.SetCamera2D(&Camera, 0)
	gohome.RenderMgr.Projection2D = prevProj

	rt.UnsetAsTarget()
}

func (this *InventoryBar) Terminate() {
	gohome.RenderMgr.RemoveObject(this)
	gohome.UpdateMgr.RemoveObject(this)
	this.Sprite2D.Terminate()
}
