package main

import "C"

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/frameworks/SDL2"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/GoHomeEngine/src/renderers/OpenGLES2"
)

func main() {
	gohome.MainLop.Run(&framework.SDL2Framework{}, &renderer.OpenGLES2Renderer{}, GAME_WIDTH, GAME_HEIGHT, "GPPCC14", &StartupScene{})
}

//export SDL_main
func SDL_main() {
	main()
}
