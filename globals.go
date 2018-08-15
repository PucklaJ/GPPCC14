package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

const GAME_WIDTH uint32 = 1280.0
const GAME_HEIGHT uint32 = 720.0

const GRAVITY float32 = 150.0
const NUM_LEVELS uint32 = 1
const ZOOM float32 = 3.0

const PLAYER_RESTITUITION float64 = 0.0
const PLAYER_FRICTION float64 = 1.0
const PLAYER_HEIGHT float32 = 32.0
const PLAYER_WIDTH float32 = PLAYER_HEIGHT / 2.0
const PLAYER_VELOCITY float32 = 500.0
const PLAYER_JUMP_FORCE float32 = 25.0
const PLAYER_DAMPING float64 = 0.0
const PLAYER_MAX_VELOCITY float32 = 50.0
const PLAYER_CATEGORY uint16 = 1 << 0
const PLAYER_FEET_CATEGORY uint16 = (1 << 1) | PLAYER_CATEGORY
const PLAYER_DENSITY float64 = 3.0

var (
	LEVELS_TMX_MAPS = [NUM_LEVELS]string{
		"test_map.tmx",
	}
)

var Camera gohome.Camera2D

const KEY_RIGHT = gohome.KeyD
const KEY_LEFT = gohome.KeyA
const KEY_DOWN = gohome.KeyS
const KEY_JUMP = gohome.KeyW
const KEY_JUMP1 = gohome.KeySpace
const KEY_SHOOT = gohome.MouseButtonLeft

const CAMERA_BOX_WIDTH float32 = float32(GAME_WIDTH) / ZOOM
const CAMERA_BOX_HEIGHT float32 = float32(GAME_HEIGHT) / ZOOM
const CAMERA_SPEED float32 = 0.1

var CAMERA_OFFSET = [2]float32{0.0, 0.0}

const GROUND_CATEGORY uint16 = 1 << 2
const GROUND_FRICTION float64 = 1.8

func LoadResources() {
}
