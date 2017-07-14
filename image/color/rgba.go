package color

import (
	"image/color"
)

const (
	f32conv  = float32(0xffff)
	f32iconv = 1 / f32conv

	f64conv  = float64(0xffff)
	f64iconv = 1 / f64conv
)

// Colors

type F32RGBA struct {
	R, G, B, A float32
}

// RGBA implements image/color.Color
func (c F32RGBA) RGBA() (r, g, b, a uint32) {
	return uint32(((c.R + 1) * 0.5) * f32conv),
		uint32(((c.G + 1) * 0.5) * f32conv),
		uint32(((c.B + 1) * 0.5) * f32conv),
		uint32(((c.A + 1) * 0.5) * f32conv)
}

type F64RGBA struct {
	R, G, B, A float64
}

// RGBA implements image/color.Color
func (c F64RGBA) RGBA() (r, g, b, a uint32) {
	return uint32(((c.R + 1) * 0.5) * f64conv),
		uint32(((c.G + 1) * 0.5) * f64conv),
		uint32(((c.B + 1) * 0.5) * f64conv),
		uint32(((c.A + 1) * 0.5) * f64conv)
}

type C64RGBA struct {
	R, G, B, A complex64
	Imag       bool
}

// RGBA implements image/color.Color
func (c C64RGBA) RGBA() (r, g, b, a uint32) {
	if c.Imag {
		return uint32(f32conv * imag(c.R)),
			uint32(f32conv * imag(c.G)),
			uint32(f32conv * imag(c.B)),
			uint32(f32conv * imag(c.A))
	}
	return uint32(f32conv * real(c.R)),
		uint32(f32conv * real(c.G)),
		uint32(f32conv * real(c.B)),
		uint32(f32conv * real(c.A))
}

type C128RGBA struct {
	R, G, B, A complex128
	Imag       bool
}

// RGBA implements image/color.Color
func (c C128RGBA) RGBA() (r, g, b, a uint32) {
	if c.Imag {
		return uint32(f64conv * imag(c.R)),
			uint32(f64conv * imag(c.G)),
			uint32(f64conv * imag(c.B)),
			uint32(f64conv * imag(c.A))
	}
	return uint32(f64conv * real(c.R)),
		uint32(f64conv * real(c.G)),
		uint32(f64conv * real(c.B)),
		uint32(f64conv * real(c.A))
}

var _ color.Color = F32RGBA{}
var _ color.Color = F64RGBA{}
var _ color.Color = C64RGBA{}
var _ color.Color = C128RGBA{}

var (
	F32RGBAModel  color.Model = color.ModelFunc(f32rgbaModel)
	F64RGBAModel  color.Model = color.ModelFunc(f64rgbaModel)
	C64RGBAModel  color.Model = color.ModelFunc(c64rgbaModel)
	C128RGBAModel color.Model = color.ModelFunc(c128rgbaModel)
)

func f32rgbaModel(c color.Color) color.Color {
	r, g, b, a := c.RGBA()

	return F32RGBA{
		R: float32(r) * f32iconv,
		G: float32(g) * f32iconv,
		B: float32(b) * f32iconv,
		A: float32(a) * f32iconv,
	}
}

func f64rgbaModel(c color.Color) color.Color {
	r, g, b, a := c.RGBA()

	return F64RGBA{
		R: float64(r) * f64iconv,
		G: float64(g) * f64iconv,
		B: float64(b) * f64iconv,
		A: float64(a) * f64iconv,
	}
}

func c64rgbaModel(c color.Color) color.Color {
	switch c := c.(type) {
	case C64RGBA:
		return c

	case C128RGBA:
		return C64RGBA{
			R: complex(float32(real(c.R)), float32(imag(c.R))),
			G: complex(float32(real(c.G)), float32(imag(c.G))),
			B: complex(float32(real(c.B)), float32(imag(c.B))),
			A: complex(float32(real(c.A)), float32(imag(c.A))),
		}
	}

	r, g, b, a := c.RGBA()

	return C64RGBA{
		R: complex(float32(r)*f32iconv, 0),
		G: complex(float32(g)*f32iconv, 0),
		B: complex(float32(b)*f32iconv, 0),
		A: complex(float32(a)*f32iconv, 0),
	}
}

func c128rgbaModel(c color.Color) color.Color {
	switch c := c.(type) {
	case C64RGBA:
		return C128RGBA{
			R: complex(float64(real(c.R)), float64(imag(c.R))),
			G: complex(float64(real(c.G)), float64(imag(c.G))),
			B: complex(float64(real(c.B)), float64(imag(c.B))),
			A: complex(float64(real(c.A)), float64(imag(c.A))),
		}

	case C128RGBA:
		return c
	}

	r, g, b, a := c.RGBA()

	return C128RGBA{
		R: complex(float64(r)*f64iconv, 0),
		G: complex(float64(g)*f64iconv, 0),
		B: complex(float64(b)*f64iconv, 0),
		A: complex(float64(a)*f64iconv, 0),
	}
}
