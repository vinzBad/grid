package grid

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Shader struct {
	Program uint32
}

func (s Shader) Use() {
	gl.UseProgram(s.Program)
}

func (s Shader) UniformLocation(name string) (int32, error) {
	loc := gl.GetUniformLocation(s.Program, gl.Str(name+"\x00"))
	if loc < 0 {
		return 0, fmt.Errorf("Unable to get index of uniform '%s'", name)
	}
	return loc, nil
}

func (s Shader) AttribLocation(name string) (uint32, error) {
	loc := gl.GetAttribLocation(s.Program, gl.Str(name+"\x00"))
	if loc < 0 {
		return 0, fmt.Errorf("Unable to get index of attribute '%s'", name)
	}
	return uint32(loc), nil
}

func (s Shader) SetAttrib(name string, count int32, xtype uint32, strideInByte int32, offsetInByte int) error {
	loc, err := s.AttribLocation(name)
	if err != nil {
		return err
	}
	gl.EnableVertexAttribArray(loc)
	gl.VertexAttribPointer(loc,
		count,
		xtype,
		false,
		strideInByte,
		gl.PtrOffset(offsetInByte))

	return nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csource := gl.Str(source + "\x00")
	gl.ShaderSource(shader, 1, &csource, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func CreateShader(vertexShaderSource, fragmentShaderSource string) (Shader, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return Shader{}, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return Shader{}, err
	}

	shader := Shader{Program: gl.CreateProgram()}

	gl.AttachShader(shader.Program, vertexShader)
	gl.AttachShader(shader.Program, fragmentShader)
	gl.BindFragDataLocation(shader.Program, 0, gl.Str("outColor\x00"))
	gl.LinkProgram(shader.Program)

	var status int32
	gl.GetProgramiv(shader.Program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shader.Program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shader.Program, logLength, nil, gl.Str(log))

		return Shader{}, errors.New(fmt.Sprintf("failed to link program: %v", log))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return shader, nil
}
