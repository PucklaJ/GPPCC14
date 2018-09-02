package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"golang.org/x/image/colornames"
	"image/color"
	"strconv"
)

const (
	INVENTORY_TEXTURE_SIZE float32 = 48.0
	INVENTORY_PADDING      float32 = INVENTORY_TEXTURE_SIZE / 8.0
	AMMO_TEXT_FONT_SIZE    uint32  = 20
	AMMO_TEXT_OFFSET_X     float32 = INVENTORY_PADDING + 8
	AMMO_TEXT_OFFSET_Y     float32 = 2
	AMMO_TEXT_POS_X        float32 = INVENTORY_TEXTURE_SIZE - AMMO_TEXT_OFFSET_X
	AMMO_TEXT_POS_Y        float32 = INVENTORY_PADDING - AMMO_TEXT_OFFSET_Y
	AMMO_TEXT_ORIGIN_X     float32 = 0.0
	AMMO_TEXT_ORIGIN_Y     float32 = 0.0
)

type InventoryBar struct {
	gohome.Sprite2D

	weapons   []Weapon
	ammoTexts []*gohome.Text2D
	current   uint8

	prevNumWeapons int
	prevCurrent    uint8
	prevAmmos      []uint32
}

func (this *InventoryBar) Init() {
	tex := gohome.Render.CreateRenderTexture("InventoryBarTexture", uint32(INVENTORY_TEXTURE_SIZE+INVENTORY_PADDING*2.0), uint32(INVENTORY_TEXTURE_SIZE+INVENTORY_PADDING*2.0), 1, false, false, false, false)
	this.Sprite2D.InitTexture(tex)

	gohome.RenderMgr.AddObject(this)
	gohome.UpdateMgr.AddObject(this)

	this.Depth = INVENTORY_DEPTH
	this.NotRelativeToCamera = 0

	this.Transform.Position = gohome.Framew.WindowGetSize().Mul(0.5)
	this.Transform.Position[1] = gohome.Framew.WindowGetSize()[1] - (INVENTORY_PADDING*2.0+INVENTORY_TEXTURE_SIZE)/2.0 - INVENTORY_PADDING
	this.Transform.Origin = [2]float32{0.5, 0.5}
	this.prevNumWeapons = -1

	this.current = 0
	this.prevCurrent = 1
}

func (this *InventoryBar) AddWeapon(w Weapon) {
	this.weapons = append(this.weapons, w)
	text := &gohome.Text2D{}
	text.Init("Ammo", AMMO_TEXT_FONT_SIZE, strconv.FormatUint(uint64(w.GetAmmo()), 10))
	text.Transform.Origin[0] = AMMO_TEXT_ORIGIN_X
	text.Transform.Origin[1] = AMMO_TEXT_ORIGIN_Y
	this.ammoTexts = append(this.ammoTexts, text)
	this.prevAmmos = append(this.prevAmmos, w.GetAmmo())
}

func (this *InventoryBar) SetCurrent(dir bool) {
	if dir == UP {
		this.current++
	} else {
		if this.current == 0 {
			this.current = uint8(len(this.weapons) - 1)
		} else {
			this.current--
		}
	}

	if this.current > uint8(len(this.weapons)-1) {
		this.current = 0
	}
}

func (this *InventoryBar) hasChanged() bool {
	if this.current != this.prevCurrent || this.prevNumWeapons != len(this.weapons) {
		return true
	}

	for i := 0; i < len(this.weapons); i++ {
		if this.prevAmmos[i] != this.weapons[i].GetAmmo() {
			return true
		}
	}

	return false
}

func (this *InventoryBar) updateValues() {
	this.prevCurrent = this.current
	this.prevNumWeapons = len(this.weapons)
	for i := 0; i < len(this.weapons); i++ {
		this.prevAmmos[i] = this.weapons[i].GetAmmo()
	}
}

func (this *InventoryBar) Update(delta_time float32) {
	if this.hasChanged() {
		this.renderInventory()
		this.updateValues()
	}
}

func (this *InventoryBar) setRenderTarget() (gohome.Projection, float32) {
	rt := this.Texture.(gohome.RenderTexture)
	width := INVENTORY_PADDING*2 + INVENTORY_TEXTURE_SIZE

	imgWidth := uint32(len(this.weapons)) * uint32(width)
	imgHeight := uint32(rt.GetHeight())

	rt.ChangeSize(imgWidth, imgHeight)
	this.Transform.Size = [2]float32{float32(imgWidth), float32(imgHeight)}
	this.TextureRegion.Max = this.Transform.Size

	rt.SetAsTarget()
	gohome.Render.ClearScreen(gohome.Color{0, 0, 0, 0})

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

	return prevProj, width
}

func (this *InventoryBar) unsetRenderTarget(prevProj gohome.Projection) {
	rt := this.Texture.(gohome.RenderTexture)
	gohome.RenderMgr.SetCamera2D(&Camera, 0)
	gohome.RenderMgr.Projection2D = prevProj

	rt.UnsetAsTarget()
}

func (this *InventoryBar) renderBox(col color.Color, x float32) {
	gohome.DrawColor = col
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
}

func (this *InventoryBar) renderBar() {
	width := INVENTORY_PADDING*2 + INVENTORY_TEXTURE_SIZE
	for i := 0; i < len(this.weapons); i++ {
		x := width * float32(i)
		this.renderBox(colornames.Gray, x)
	}
}

func (this *InventoryBar) renderTextures() {
	width := INVENTORY_PADDING*2 + INVENTORY_TEXTURE_SIZE
	for i := 0; i < len(this.weapons); i++ {
		wt := this.weapons[i].GetInventoryTexture()
		x := width * float32(i)
		if wt != nil {
			var spr gohome.Sprite2D
			spr.InitTexture(wt)
			spr.Flip = gohome.FLIP_VERTICAL
			spr.Transform.Position[0] = x + INVENTORY_PADDING
			spr.Transform.Position[1] = INVENTORY_PADDING
			gohome.RenderMgr.RenderRenderObject(&spr)
		}
	}
}

func (this *InventoryBar) renderCurrent() {
	width := INVENTORY_PADDING*2 + INVENTORY_TEXTURE_SIZE
	x := width * float32(this.current)
	this.renderBox(colornames.Gold, x)
}

func (this *InventoryBar) renderAmmoTexts() {
	width := INVENTORY_PADDING*2 + INVENTORY_TEXTURE_SIZE
	for i := 0; i < len(this.weapons); i++ {
		x := width * float32(i)
		text := this.ammoTexts[i]
		text.Text = strconv.FormatUint(uint64(this.weapons[i].GetAmmo()), 10)
		text.Flip = gohome.FLIP_VERTICAL
		x += AMMO_TEXT_POS_X
		y := AMMO_TEXT_POS_Y

		text.Transform.Position = [2]float32{x, y}
		gohome.RenderMgr.RenderRenderObject(text)
	}
}

func (this *InventoryBar) renderInventory() {
	prevProj, _ := this.setRenderTarget()

	this.renderBar()
	this.renderTextures()
	this.renderCurrent()
	this.renderAmmoTexts()

	this.unsetRenderTarget(prevProj)
}

func (this *InventoryBar) Terminate() {
	gohome.RenderMgr.RemoveObject(this)
	gohome.UpdateMgr.RemoveObject(this)
	this.Sprite2D.Terminate()
}
