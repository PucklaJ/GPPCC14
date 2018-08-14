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
}

func (this *LevelScene) Init() {
	gohome.Init2DShaders()

	gohome.ResourceMgr.LoadTMXMap("Level", LEVELS_TMX_MAPS[this.LevelID])

	this.Map.Init("Level")
	gohome.RenderMgr.AddObject(&this.Map)

	this.PhysicsMgr.Init([2]float32{0.0, GRAVITY})
	gohome.UpdateMgr.AddObject(&this.PhysicsMgr)
	debug := this.PhysicsMgr.GetDebugDraw()
	gohome.RenderMgr.AddObject(&debug)

	this.PhysicsMgr.LayerToCollision(&this.Map, "Collision")

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
				}
			}
		}
	}

	this.Player.Init(playerStart, &this.PhysicsMgr)
}

func (this *LevelScene) Update(delta_time float32) {
	Camera.Zoom += float32(gohome.InputMgr.Mouse.Wheel[1]) * 0.1
}

func (this *LevelScene) Terminate() {
	this.PhysicsMgr.Terminate()
	gohome.ResourceMgr.DeleteTMXMap("Level")
	this.Map.Terminate()
	this.Player.Terminate()
}
