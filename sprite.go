package grid

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Sprite struct {
	texture       Texture
	vao, ebo, vbo uint32
	shader        Shader
}

func (s *Sprite) Draw() {
	s.shader.Use()
	s.texture.Bind()

	gl.BindVertexArray(s.vao)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, s.ebo)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)
	//gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

}

func CreateSprite(t Texture) (s Sprite, err error) {

	s.texture = t
	s.shader, err = CreateShader(vertexShader, fragmentShader)
	if err != nil {
		return s, err
	}

	gl.GenBuffers(1, &s.ebo)
	gl.GenBuffers(1, &s.vbo)
	gl.GenVertexArrays(1, &s.vao)

	gl.BindVertexArray(s.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, s.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	err = s.shader.SetAttrib("vert", 2, gl.FLOAT, 4*4, 0)
	err = s.shader.SetAttrib("vertTexCoord", 2, gl.FLOAT, 4*4, 2*4)

	gl.BindVertexArray(0)

	return s, err
}

var vertices = []float32{
	// X, Y, U, V
	0, 0, 0, 0, // Bottom Left
	0, 1, 0, 1, // Top Left
	1, 1, 1, 1, // Top Right
	1, 0, 1, 0, // Bottom Right
}

var indices = []uint32{
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
    gl_Position =  vec4(vert, 0, 1);
}
` + "\x00"

var fragmentShader = `
#version 330 core

uniform sampler2D tex;
uniform vec4 tint;

in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
    //outputColor = vec4(fragTexCoord, 1, 1) ;
	outputColor = texture(tex, fragTexCoord);
}
` + "\x00"
