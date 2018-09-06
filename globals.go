package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

const GAME_WIDTH uint32 = 1280.0
const GAME_HEIGHT uint32 = 720.0

const GRAVITY float32 = 200.0
const NUM_LEVELS uint32 = 3
const ZOOM float32 = 3.0

const DELETE_RAY_DEPTH uint8 = 2
const PLAYER_DEPTH uint8 = 1
const WEAPON_DEPTH uint8 = 3
const INVENTORY_DEPTH uint8 = 4
const MENU_DEPTH uint8 = 5

const PLAYER_CATEGORY uint16 = 1 << 0
const PLAYER_FEET_CATEGORY uint16 = (1 << 1) | PLAYER_CATEGORY
const PLAYER_FEET_SENSOR_CATEGORY uint16 = (1 << 10) | PLAYER_CATEGORY

var (
	LEVELS_TMX_MAPS = [NUM_LEVELS]string{
		"test_map.tmx",
		"test_map1.tmx",
		"test_map2.tmx",
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

const WIN_CONDITION_TARGET bool = true
const WIN_CONDITION_ENEMY bool = false

var CURRENT_WIN_CONDITION bool

func LoadResources() {
	gohome.ResourceMgr.PreloadFont("Button", "/usr/share/fonts/truetype/ubuntu/UbuntuMono-R.ttf")
	gohome.ResourceMgr.PreloadTexture("Player", "GPPCC14_Player.png")
	gohome.ResourceMgr.PreloadTexture("DefaultWeapon", "GPPCC14_DefaultWeapon.png")
	gohome.ResourceMgr.PreloadTexture("FreezeWeapon", "GPPCC14_FreezeWeapon.png")
	gohome.ResourceMgr.PreloadTexture("BallWeapon", "GPPCC14_BallWeapon.png")
	gohome.ResourceMgr.PreloadTexture("MoveWeapon", "GPPCC14_MoveWeapon.png")
	gohome.ResourceMgr.PreloadTexture("DeleteWeapon", "GPPCC14_DeleteWeapon.png")
	gohome.ResourceMgr.PreloadTexture("DefaultWeaponInv", "GPPCC14_DefaultWeaponInv.png")
	gohome.ResourceMgr.PreloadTexture("FreezeWeaponInv", "GPPCC14_FreezeWeaponInv.png")
	gohome.ResourceMgr.PreloadTexture("BallWeaponInv", "GPPCC14_BallWeaponInv.png")
	gohome.ResourceMgr.PreloadTexture("MoveWeaponInv", "GPPCC14_MoveWeaponInv.png")
	gohome.ResourceMgr.PreloadTexture("DeleteWeaponInv", "GPPCC14_DeleteWeaponInv.png")
	gohome.ResourceMgr.PreloadTexture("DefaultWeaponBlock", "GPPCC14_DefaultWeaponBlock.png")
	gohome.ResourceMgr.PreloadTexture("FreezeWeaponBlock", "GPPCC14_FreezeWeaponBlock.png")
	gohome.ResourceMgr.PreloadTexture("BallWeaponBlock", "GPPCC14_BallWeaponBlock.png")
	gohome.ResourceMgr.PreloadTexture("MoveWeaponBlock", "GPPCC14_MoveWeaponBlock.png")
	gohome.ResourceMgr.PreloadTexture("Enemy", "GPPCC14_Enemy.png")
	gohome.ResourceMgr.PreloadTexture("Explosion", "GPPCC14_Explosion.png")
	gohome.ResourceMgr.PreloadTexture("Disappear", "GPPCC14_Disappear.png")
	gohome.ResourceMgr.PreloadTexture("Restart", "GPPCC14_Restart.png")
	gohome.ResourceMgr.PreloadTexture("Back", "GPPCC14_Back.png")
	gohome.ResourceMgr.PreloadTexture("Pause", "GPPCC14_Pause.png")
	gohome.ResourceMgr.PreloadTexture("Resume", "GPPCC14_Resume.png")
	gohome.ResourceMgr.PreloadTexture("LevelButton1", "GPPCC14_LevelButton1.png")
	gohome.ResourceMgr.PreloadTexture("LevelButtonPressed", "GPPCC14_LevelButtonPressed.png")
	gohome.ResourceMgr.PreloadTexture("AmmoFont", "GPPCC14_AmmoFont.png")
	gohome.ResourceMgr.LoadPreloadedResources()

	gohome.ResourceMgr.GetTexture("Player").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("DefaultWeapon").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("FreezeWeapon").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("BallWeapon").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("MoveWeapon").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("DeleteWeapon").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("DefaultWeaponInv").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("FreezeWeaponInv").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("BallWeaponInv").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("MoveWeaponInv").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("DeleteWeaponInv").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("DefaultWeaponBlock").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("FreezeWeaponBlock").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("BallWeaponBlock").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("MoveWeaponBlock").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("Enemy").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("Explosion").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("Disappear").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("Restart").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("Back").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("Pause").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("Resume").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("LevelButton1").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("LevelButtonPressed").SetFiltering(gohome.FILTERING_NEAREST)
	gohome.ResourceMgr.GetTexture("AmmoFont").SetFiltering(gohome.FILTERING_NEAREST)
}
