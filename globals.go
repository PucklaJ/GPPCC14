package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

const GAME_WIDTH uint32 = 1280.0
const GAME_HEIGHT uint32 = 720.0

const GRAVITY float32 = 200.0
const NUM_LEVELS uint32 = 1
const ZOOM float32 = 3.0

const PLAYER_CATEGORY uint16 = 1 << 0
const PLAYER_FEET_CATEGORY uint16 = (1 << 1) | PLAYER_CATEGORY
const PLAYER_FEET_SENSOR_CATEGORY uint16 = (1 << 10) | PLAYER_CATEGORY

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

const WEAPON_CATEGORY uint16 = 1 << 3

const ENEMY_CATEGORY uint16 = 1 << 4
const ENEMY_SENSOR_CATEGORY uint16 = 1 << 5
const ENEMY_SMALL_LEFT_SENSOR_CATEGORY uint16 = 1<<6 | ENEMY_SENSOR_CATEGORY
const ENEMY_SMALL_RIGHT_SENSOR_CATEGORY uint16 = 1<<7 | ENEMY_SENSOR_CATEGORY
const ENEMY_BIG_LEFT_SENSOR_CATEGORY uint16 = 1<<8 | ENEMY_SENSOR_CATEGORY
const ENEMY_BIG_RIGHT_SENSOR_CATEGORY uint16 = 1<<9 | ENEMY_SENSOR_CATEGORY

func LoadResources() {
	gohome.ResourceMgr.LoadFont("Ammo", "UbuntuMono-R.ttf")
	// gohome.ResourceMgr.LoadPreloadedResources()
}
