package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"strconv"
)

type DebugInfo struct {
	gohome.Text2D
}

func (this *DebugInfo) Init() {
	this.Text2D.Init(gohome.ButtonFont, 24, "FPS: 60\nUOBJs: 0\nROBJs: 0")
	this.Text2D.NotRelativeToCamera = 0
	gohome.UpdateMgr.AddObject(this)
	gohome.RenderMgr.AddObject(this)

	this.Depth = 255
}

func (this *DebugInfo) Update(delta_time float32) {
	if this.Visible {
		this.Text = "FPS: " + strconv.FormatFloat(float64(1.0/delta_time), 'f', 1, 32) + "\n" +
			"UOBJs: " + strconv.FormatUint(uint64(gohome.UpdateMgr.NumUpdateObjects()), 10) + "\n" +
			"ROBJs: " + strconv.FormatUint(uint64(gohome.RenderMgr.NumRenderObjects()), 10)
	}
}

func (this *DebugInfo) Terminate() {
	this.Text2D.Terminate()
	gohome.UpdateMgr.RemoveObject(this)
	gohome.RenderMgr.RemoveObject(this)
}
