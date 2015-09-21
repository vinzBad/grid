package grid

type Sprite struct {
	Texture Texture
}

func CreateSprite(t Texture) {

}

var vertices = [...]float32{
	// X, Y, U, V
	0, 0, 0, 0, // Bottom Left
	0, 1, 0, 1, // Top Left
	1, 1, 1, 1, // Top Right
	1, 0, 1, 0, // Bottom Right
}

var indices = [...]float32{
	0, 1, 2,
	0, 3, 2,
}

var vertexShader string = `
#version 330 core

uniform mat4 scene;
uniform mat4 model;

in vec2 vert;
in vec2 vertTexCoord;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = scene * model * vec4(vert, 0, 1);
}
` + "\x00"

var fragmentShader = `
#version 330 core

uniform sampler2D tex;
uniform vec4 tint

in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
    outputColor = texture(tex, fragTexCoord) * tint;
}
` + "\x00"
