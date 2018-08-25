package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/physics2d"
)

type LevelScene struct {
	PhysicsMgr physics2d.PhysicsManager2D
	LevelID    uint32
	Map        gohome.TiledMap
	Player     Player
	Enemies    []*Enemy

	debugDraw physics2d.PhysicsDebugDraw2D
}

func (this *LevelScene) Init() {
	gohome.Init2DShaders()
	physics2d.PIXEL_PER_METER = 10.0
	gohome.ResourceMgr.LoadTMXMap("Level", LEVELS_TMX_MAPS[this.LevelID])

	this.Map.Init("Level")
	gohome.RenderMgr.AddObject(&this.Map)

	this.PhysicsMgr.Init([2]float32{0.0, GRAVITY})
	gohome.UpdateMgr.AddObject(&this.PhysicsMgr)
	this.debugDraw = this.PhysicsMgr.GetDebugDraw()
	gohome.RenderMgr.AddObject(&this.debugDraw)
	this.debugDraw.OnlyDrawDynamic = true
	this.debugDraw.DrawBodies = false

	groundBodies := this.PhysicsMgr.LayerToCollision(&this.Map, "Collision")
	for i := 0; i < len(groundBodies); i++ {
		b := groundBodies[i]
		if b == nil {
			continue
		}
		for f := b.GetFixtureList(); f != nil; f = f.GetNext() {
			filter := f.GetFilterData()
			filter.CategoryBits = GROUND_CATEGORY
			filter.MaskBits = 0xffff
			f.SetFilterData(filter)
			f.SetFriction(GROUND_FRICTION)
		}
	}

	var playerStart [2]float32

	ls := this.Map.Layers
	for i := 0; i < len(ls); i++ {
		l := ls[i]
		if l.Name == "Settings" {
			objs := l.Objects
			for j := 0; j < len(objs); j++ {
				o := objs[j]
				if o.Name == "start" {
					playerStart[0] = float32(o.X)
					playerStart[1] = float32(o.Y)
				} else if o.Name == "enemy" {
					enemy := &Enemy{}
					enemy.Sprite2D.Init("")
					enemy.Transform.Position = [2]float32{float32(o.X), float32(o.Y)}
					this.Enemies = append(this.Enemies, enemy)
				}
			}
		}
	}

	this.Player.Init(playerStart, &this.PhysicsMgr)
	for i := 0; i < len(this.Enemies); i++ {
		this.Enemies[i].Init(this.Enemies[i].Transform.Position, &this.Player)
	}
}

func (this *LevelScene) Update(delta_time float32) {
}

func (this *LevelScene) Terminate() {
	gohome.UpdateMgr.RemoveObject(&this.PhysicsMgr)
	gohome.RenderMgr.RemoveObject(&this.Map)
	gohome.RenderMgr.RemoveObject(&this.debugDraw)

	gohome.ResourceMgr.DeleteTMXMap("Level")

	for i := 0; i < len(this.Enemies); i++ {
		this.Enemies[i].Terminate()
	}
	this.Player.Terminate()
	this.Map.Terminate()
	this.PhysicsMgr.Terminate()
}
