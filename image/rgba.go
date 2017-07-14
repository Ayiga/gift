package image

import (
	"fmt"
	"image"
	"image/color"

	giftcolor "github.com/disintegration/gift/image/color"
)

type F32RGBA struct {
	// Pix holds the image's pixles, in R, G, B, A order. The pixel at
	// (x, y) stars at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []float32

	// Stride is the Pix stride (as index offset) between vertically adjacent pixels.
	Stride int

	// Rect is the image's bounds.
	Rect image.Rectangle
}

// ColorModel implements image.Image
func (*F32RGBA) ColorModel() color.Model {
	return giftcolor.F32RGBAModel
}

// Bounds implements image.Image
func (i *F32RGBA) Bounds() image.Rectangle {
	return i.Rect
}

// At implements image.Image
func (i *F32RGBA) At(x, y int) color.Color {
	return i.F32RGBAAt(x, y)
}

func (i *F32RGBA) F32RGBAAt(x, y int) giftcolor.F32RGBA {
	if !(image.Point{x, y}.In(i.Rect)) {
		return giftcolor.F32RGBA{}
	}

	idx := i.PixOffset(x, y)

	return giftcolor.F32RGBA{
		i.Pix[idx+0],
		i.Pix[idx+1],
		i.Pix[idx+2],
		i.Pix[idx+3],
	}
}

func (i *F32RGBA) PixOffset(x, y int) int {
	return (y-i.Rect.Min.Y)*i.Stride + (x-i.Rect.Min.X)*4
}

func (i *F32RGBA) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(i.Rect)) {
		return
	}

	idx := i.PixOffset(x, y)
	c1 := giftcolor.F32RGBAModel.Convert(c).(giftcolor.F32RGBA)
	i.Pix[idx+0] = c1.R
	i.Pix[idx+1] = c1.G
	i.Pix[idx+2] = c1.B
	i.Pix[idx+3] = c1.A
}

func (i *F32RGBA) SetF32RGBA(x, y int, c giftcolor.F32RGBA) {
	if !(image.Point{x, y}.In(i.Rect)) {
		return
	}

	idx := i.PixOffset(x, y)
	i.Pix[idx+0] = c.R
	i.Pix[idx+1] = c.G
	i.Pix[idx+2] = c.B
	i.Pix[idx+3] = c.A
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (i *F32RGBA) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(i.Rect)
	// If r1 and r2 are Rectangles, r1.Insect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &F32RGBA{}
	}

	idx := i.PixOffset(r.Min.X, r.Min.Y)
	return &F32RGBA{
		Pix:    i.Pix[idx:],
		Stride: i.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (i *F32RGBA) Opaque() bool {
	if i.Rect.Empty() {
		return true
	}
	i0, i1 := 3, i.Rect.Dx()*4
	for y := i.Rect.Min.Y; y < i.Rect.Max.Y; y++ {
		for idx := i0; idx < i1; idx += 4 {
			if i.Pix[idx] != 1.0 {
				return false
			}
		}
		i0 += i.Stride
		i1 += i.Stride
	}
	return true
}

// NewF32RGBA returns a new F32RGBA image with the given bounds.
func NewF32RGBA(r image.Rectangle) *F32RGBA {
	w, h := r.Dx(), r.Dy()
	pix := make([]float32, 4*w*h)
	return &F32RGBA{pix, 4 * w, r}
}

type F64RGBA struct {
	// Pix holds the image's pixles, in R, G, B, A order. The pixel at
	// (x, y) stars at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []float64

	// Stride is the Pix stride (as index offset) between vertically adjacent pixels.
	Stride int

	// Rect is the image's bounds.
	Rect image.Rectangle
}

// ColorModel implements image.Image
func (*F64RGBA) ColorModel() color.Model {
	return giftcolor.F64RGBAModel
}

// Bounds implements image.Image
func (i *F64RGBA) Bounds() image.Rectangle {
	return i.Rect
}

// At implements image.Image
func (i *F64RGBA) At(x, y int) color.Color {
	return i.F64RGBAAt(x, y)
}

func (i *F64RGBA) F64RGBAAt(x, y int) giftcolor.F64RGBA {
	if !(image.Point{x, y}.In(i.Rect)) {
		return giftcolor.F64RGBA{}
	}

	idx := i.PixOffset(x, y)

	return giftcolor.F64RGBA{
		i.Pix[idx+0],
		i.Pix[idx+1],
		i.Pix[idx+2],
		i.Pix[idx+3],
	}
}

func (i *F64RGBA) PixOffset(x, y int) int {
	return (y-i.Rect.Min.Y)*i.Stride + (x-i.Rect.Min.X)*4
}

func (i *F64RGBA) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(i.Rect)) {
		return
	}

	idx := i.PixOffset(x, y)
	c1 := giftcolor.F64RGBAModel.Convert(c).(giftcolor.F64RGBA)
	i.Pix[idx+0] = c1.R
	i.Pix[idx+1] = c1.G
	i.Pix[idx+2] = c1.B
	i.Pix[idx+3] = c1.A
}

func (i *F64RGBA) SetF64RGBA(x, y int, c giftcolor.F64RGBA) {
	if !(image.Point{x, y}.In(i.Rect)) {
		return
	}

	idx := i.PixOffset(x, y)
	i.Pix[idx+0] = c.R
	i.Pix[idx+1] = c.G
	i.Pix[idx+2] = c.B
	i.Pix[idx+3] = c.A
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (i *F64RGBA) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(i.Rect)
	// If r1 and r2 are Rectangles, r1.Insect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &F64RGBA{}
	}

	idx := i.PixOffset(r.Min.X, r.Min.Y)
	return &F64RGBA{
		Pix:    i.Pix[idx:],
		Stride: i.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (i *F64RGBA) Opaque() bool {
	if i.Rect.Empty() {
		return true
	}
	i0, i1 := 3, i.Rect.Dx()*4
	for y := i.Rect.Min.Y; y < i.Rect.Max.Y; y++ {
		for idx := i0; idx < i1; idx += 4 {
			if i.Pix[idx] != 1.0 {
				return false
			}
		}
		i0 += i.Stride
		i1 += i.Stride
	}
	return true
}

// NewF64RGBA returns a new F64RGBA image with the given bounds.
func NewF64RGBA(r image.Rectangle) *F64RGBA {
	w, h := r.Dx(), r.Dy()
	pix := make([]float64, 4*w*h)
	return &F64RGBA{pix, 4 * w, r}
}

type C64RGBA struct {
	// Pix holds the image's pixles, in R, G, B, A order. The pixel at
	// (x, y) stars at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []complex64

	// Stride is the Pix stride (as index offset) between vertically adjacent pixels.
	Stride int

	// Rect is the image's bounds.
	Rect image.Rectangle
}

// ColorModel implements image.Image
func (*C64RGBA) ColorModel() color.Model {
	return giftcolor.C64RGBAModel
}

// Bounds implements image.Image
func (i *C64RGBA) Bounds() image.Rectangle {
	return i.Rect
}

// At implements image.Image
func (i *C64RGBA) At(x, y int) color.Color {
	return i.C64RGBAAt(x, y)
}

func (i *C64RGBA) C64RGBAAt(x, y int) giftcolor.C64RGBA {
	if !(image.Point{x, y}.In(i.Rect)) {
		return giftcolor.C64RGBA{}
	}

	idx := i.PixOffset(x, y)

	return giftcolor.C64RGBA{
		R: i.Pix[idx+0],
		G: i.Pix[idx+1],
		B: i.Pix[idx+2],
		A: i.Pix[idx+3],
	}
}

func (i *C64RGBA) PixOffset(x, y int) int {
	return (y-i.Rect.Min.Y)*i.Stride + (x-i.Rect.Min.X)*4
}

func (i *C64RGBA) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(i.Rect)) {
		return
	}

	idx := i.PixOffset(x, y)
	c1 := giftcolor.C64RGBAModel.Convert(c).(giftcolor.C64RGBA)
	i.Pix[idx+0] = c1.R
	i.Pix[idx+1] = c1.G
	i.Pix[idx+2] = c1.B
	i.Pix[idx+3] = c1.A
}

func (i *C64RGBA) SetC64RGBA(x, y int, c giftcolor.C64RGBA) {
	if !(image.Point{x, y}.In(i.Rect)) {
		return
	}

	idx := i.PixOffset(x, y)
	i.Pix[idx+0] = c.R
	i.Pix[idx+1] = c.G
	i.Pix[idx+2] = c.B
	i.Pix[idx+3] = c.A
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (i *C64RGBA) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(i.Rect)
	// If r1 and r2 are Rectangles, r1.Insect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &C64RGBA{}
	}

	idx := i.PixOffset(r.Min.X, r.Min.Y)
	return &C64RGBA{
		Pix:    i.Pix[idx:],
		Stride: i.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (i *C64RGBA) Opaque() bool {
	if i.Rect.Empty() {
		return true
	}
	i0, i1 := 3, i.Rect.Dx()*4
	for y := i.Rect.Min.Y; y < i.Rect.Max.Y; y++ {
		for idx := i0; idx < i1; idx += 4 {
			if real(i.Pix[idx]) != 1.0 {
				return false
			}
		}
		i0 += i.Stride
		i1 += i.Stride
	}
	return true
}

// NewC64RGBA returns a new C64RGBA image with the given bounds.
func NewC64RGBA(r image.Rectangle) *C64RGBA {
	w, h := r.Dx(), r.Dy()
	pix := make([]complex64, 4*w*h)
	return &C64RGBA{pix, 4 * w, r}
}

type C128RGBA struct {
	// Pix holds the image's pixles, in R, G, B, A order. The pixel at
	// (x, y) stars at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []complex128

	// Stride is the Pix stride (as index offset) between vertically adjacent pixels.
	Stride int

	// Rect is the image's bounds.
	Rect image.Rectangle
}

// ColorModel implements image.Image
func (*C128RGBA) ColorModel() color.Model {
	return giftcolor.C128RGBAModel
}

// Bounds implements image.Image
func (i *C128RGBA) Bounds() image.Rectangle {
	return i.Rect
}

// At implements image.Image
func (i *C128RGBA) At(x, y int) color.Color {
	return i.C128RGBAAt(x, y)
}

func (i *C128RGBA) C128RGBAAt(x, y int) giftcolor.C128RGBA {
	if !(image.Point{x, y}.In(i.Rect)) {
		return giftcolor.C128RGBA{}
	}

	idx := i.PixOffset(x, y)

	return giftcolor.C128RGBA{
		R: i.Pix[idx+0],
		G: i.Pix[idx+1],
		B: i.Pix[idx+2],
		A: i.Pix[idx+3],
	}
}

func (i *C128RGBA) PixOffset(x, y int) int {
	return (y-i.Rect.Min.Y)*i.Stride + (x-i.Rect.Min.X)*4
}

func (i *C128RGBA) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(i.Rect)) {
		return
	}

	idx := i.PixOffset(x, y)
	c1 := giftcolor.C128RGBAModel.Convert(c).(giftcolor.C128RGBA)
	i.Pix[idx+0] = c1.R
	i.Pix[idx+1] = c1.G
	i.Pix[idx+2] = c1.B
	i.Pix[idx+3] = c1.A
}

func (i *C128RGBA) SetC128RGBA(x, y int, c giftcolor.C128RGBA) {
	if !(image.Point{x, y}.In(i.Rect)) {
		return
	}

	idx := i.PixOffset(x, y)
	i.Pix[idx+0] = c.R
	i.Pix[idx+1] = c.G
	i.Pix[idx+2] = c.B
	i.Pix[idx+3] = c.A
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (i *C128RGBA) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(i.Rect)
	// If r1 and r2 are Rectangles, r1.Insect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &C128RGBA{}
	}

	idx := i.PixOffset(r.Min.X, r.Min.Y)
	return &C128RGBA{
		Pix:    i.Pix[idx:],
		Stride: i.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (i *C128RGBA) Opaque() bool {
	if i.Rect.Empty() {
		return true
	}
	i0, i1 := 3, i.Rect.Dx()*4
	for y := i.Rect.Min.Y; y < i.Rect.Max.Y; y++ {
		for idx := i0; idx < i1; idx += 4 {
			if real(i.Pix[idx]) != 1.0 {
				return false
			}
		}
		i0 += i.Stride
		i1 += i.Stride
	}
	return true
}

func (i *C128RGBA) getChannel(c int) [][]complex128 {
	w, h := i.Rect.Dx(), i.Rect.Dy()
	r := make([][]complex128, h)
	i0, i1 := c, i.Rect.Dx()*4
	for y, yi := i.Rect.Min.Y, 0; y < i.Rect.Max.Y; y, yi = y+1, yi+1 {
		r[yi] = make([]complex128, w)
		for idx, ii := i0, 0; idx < i1; idx, ii = idx+4, ii+1 {
			r[yi][ii] = i.Pix[idx]
		}
		i0 += i.Stride
		i1 += i.Stride
	}

	return r
}

func (i *C128RGBA) RedChannel() [][]complex128 {
	return i.getChannel(0)
}

func (i *C128RGBA) GreenChannel() [][]complex128 {
	return i.getChannel(1)
}

func (i *C128RGBA) BlueChannel() [][]complex128 {
	return i.getChannel(2)
}

func (i *C128RGBA) AlphaChannel() [][]complex128 {
	return i.getChannel(3)
}

// NewC128RGBA returns a new C128RGBA image with the given bounds.
func NewC128RGBA(r image.Rectangle) *C128RGBA {
	w, h := r.Dx(), r.Dy()
	pix := make([]complex128, 4*w*h)
	return &C128RGBA{pix, 4 * w, r}
}

func NewC128RGBAFromChannels(r, g, b, a [][]complex128) *C128RGBA {
	h := len(r)
	w := len(r[0])

	fmt.Printf("image dimensions: %d x %d\n", w, h)

	i := NewC128RGBA(image.Rect(0, 0, w, h))

	for y := i.Rect.Min.Y; y < i.Rect.Max.Y; y++ {
		for x := i.Rect.Min.X; x < i.Rect.Max.X; x++ {
			idx := i.PixOffset(x, y)
			i.Pix[idx+0] = r[y][x]
			i.Pix[idx+1] = g[y][x]
			i.Pix[idx+2] = b[y][x]
			i.Pix[idx+3] = a[y][x]
		}
	}

	return i
}
