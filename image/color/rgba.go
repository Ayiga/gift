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

// F32RGBA is an RGBA color where every channel is
// represented by a float32.
//
// The accepted color range is between -1.0 and 1.0
type F32RGBA struct {
	R, G, B, A float32
}

// RGBA implements image/color.Color
func (c F32RGBA) RGBA() (r, g, b, a uint32) {
	return uint32(0.5 * f32conv * (c.R + 1)),
		uint32(0.5 * f32conv * (c.G + 1)),
		uint32(0.5 * f32conv * (c.B + 1)),
		uint32(0.5 * f32conv * (c.A + 1))
}

// F64RGBA is an RGBA color where every channel is
// represented by a float64.
//
// The accepted color range is between -1.0 and 1.0
type F64RGBA struct {
	R, G, B, A float64
}

// RGBA implements image/color.Color
func (c F64RGBA) RGBA() (r, g, b, a uint32) {
	return uint32(0.5 * f64conv * (c.R + 1)),
		uint32(0.5 * f64conv * (c.G + 1)),
		uint32(0.5 * f64conv * (c.B + 1)),
		uint32(0.5 * f64conv * (c.A + 1))
}

// C64RGBA is an RGBA color where every channel is
// represented by a complex64.
type C64RGBA struct {
	R, G, B, A complex64
}

// RGBA implements image/color.Color
func (c C64RGBA) RGBA() (r, g, b, a uint32) {
	return uint32(f32conv * 0.5 * (real(c.R) + 1)),
		uint32(f32conv * 0.5 * (real(c.G) + 1)),
		uint32(f32conv * 0.5 * (real(c.B) + 1)),
		uint32(f32conv * 0.5 * (real(c.A) + 1))
}

// C128RGBA is an RGBA color where every channel is
// represented by a complex128.
type C128RGBA struct {
	R, G, B, A complex128
}

// RGBA implements image/color.Color
func (c C128RGBA) RGBA() (r, g, b, a uint32) {
	return uint32(f64conv * 0.5 * real(c.R)),
		uint32(f64conv * 0.5 * real(c.G)),
		uint32(f64conv * 0.5 * real(c.B)),
		uint32(f64conv * 0.5 * real(c.A))
}

var _ color.Color = F32RGBA{}
var _ color.Color = F64RGBA{}
var _ color.Color = C64RGBA{}
var _ color.Color = C128RGBA{}

var (
	// F32RGBAModel is a color.Model that will translate any given color
	// into the range of [0, 1]
	//
	// Any floating point number color model defined within this package will
	// be maintained for values below 0
	F32RGBAModel = color.ModelFunc(f32rgbaModel)

	// F64RGBAModel is a color.Model that will translate any given color
	// into the range of [0, 1]
	//
	// Any floating point number color model defined within this package will
	// be maintained for values below 0
	F64RGBAModel = color.ModelFunc(f64rgbaModel)

	// C64RGBAModel is a color.Model that will translate any given color
	// into the range of [0, 1] in only the real components
	//
	// Any floating point number color model defined within this package will
	// be maintained for values below 0
	C64RGBAModel = color.ModelFunc(c64rgbaModel)

	// C128RGBAModel is a color.Model that will translate any given color
	// into the range of [0, 1] in only the real components
	//
	// Any floating point number color model defined within this package will
	// be maintained for values below 0
	C128RGBAModel = color.ModelFunc(c128rgbaModel)

	// RealRGBAModel is a color.Model that will convert floating point colors
	// in the real number domain to an NRGBA color.
	//
	// This translates the range of [0, 1] to [0, 255]
	RealRGBAModel = color.ModelFunc(clampRealColor)

	// ComplexRGBAModel is a color.Model that will convert floating point colors
	// in the imaginary number domain to an NRGBA color.
	//
	// This translates the range of [0, 1] to [0, 255]
	ComplexRGBAModel = color.ModelFunc(clampComplexColor)
)

func f32ToUint8(v float32) uint8 {
	if v < 0 {
		return 0
	}

	return uint8(v * f32conv)
}

func f64ToUint8(v float64) uint8 {
	if v < 0 {
		return 0
	}

	return uint8(v * f64conv)
}

// clampRealColor converts channels from the range [0, 1] to [0, 255]
func clampRealColor(c color.Color) color.Color {
	switch t := c.(type) {
	case F32RGBA:
		return color.NRGBA{
			R: f32ToUint8(t.R),
			G: f32ToUint8(t.G),
			B: f32ToUint8(t.B),
			A: f32ToUint8(t.A),
		}
	case F64RGBA:
		return color.NRGBA{
			R: f64ToUint8(t.R),
			G: f64ToUint8(t.G),
			B: f64ToUint8(t.B),
			A: f64ToUint8(t.A),
		}
	case C64RGBA:
		return color.NRGBA{
			R: f32ToUint8(real(t.R)),
			G: f32ToUint8(real(t.G)),
			B: f32ToUint8(real(t.B)),
			A: f32ToUint8(real(t.A)),
		}
	case C128RGBA:
		return color.NRGBA{
			R: f64ToUint8(real(t.R)),
			G: f64ToUint8(real(t.G)),
			B: f64ToUint8(real(t.B)),
			A: f64ToUint8(real(t.A)),
		}

	default:
		return color.NRGBAModel.Convert(c)
	}
}

func clampComplexColor(c color.Color) color.Color {
	switch t := c.(type) {
	case F32RGBA, F64RGBA:
		return color.NRGBA{}

	case C64RGBA:
		return color.NRGBA{
			R: f32ToUint8(imag(t.R)),
			G: f32ToUint8(imag(t.G)),
			B: f32ToUint8(imag(t.B)),
			A: f32ToUint8(imag(t.A)),
		}
	case C128RGBA:
		return color.NRGBA{
			R: f64ToUint8(imag(t.R)),
			G: f64ToUint8(imag(t.G)),
			B: f64ToUint8(imag(t.B)),
			A: f64ToUint8(imag(t.A)),
		}

	default:
		return color.NRGBAModel.Convert(c)
	}
}

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
