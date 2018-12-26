package main

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/frameworks/GLFW"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/GoHomeEngine/src/renderers/OpenGL"
)

func main() {
	gohome.MainLop.Run(&framework.GLFWFramework{},&renderer.OpenGLRenderer{},1280,720,"Schieße den Weg",&StartupScene{})
}
