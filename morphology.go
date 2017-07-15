package gift

import (
	"image"
	"image/draw"
	"math"
)

func prepareMorphologyMasks(kernel []float32) (int, []uvweight) {
	size := int(math.Sqrt(float64(len(kernel))))
	if size%2 == 0 {
		size--
	}
	if size < 1 {
		return 0, []uvweight{}
	}
	center := size / 2

	weights := []uvweight{}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			k := j*size + i
			var w float32
			if k < len(kernel) {
				w = kernel[k]
			}
			if w != 0 {
				weights = append(weights, uvweight{u: i - center, v: j - center, weight: w})
			}
		}
	}

	return size, weights
}

type morphologyErosion struct {
	kernel []float32
}

func (*morphologyErosion) Bounds(srcBounds image.Rectangle) (dstBounds image.Rectangle) {
	dstBounds = image.Rect(0, 0, srcBounds.Dx(), srcBounds.Dy())
	return
}

func (p *morphologyErosion) Draw(dst draw.Image, src image.Image, options *Options) {
	if options == nil {
		options = &defaultOptions
	}

	srcb := src.Bounds()
	dstb := dst.Bounds()

	if srcb.Dx() <= 0 || srcb.Dy() <= 0 {
		return
	}

	ksize, masks := prepareMorphologyMasks(p.kernel)
	kcenter := ksize / 2

	if ksize < 1 {
		copyimage(dst, src, options)
		return
	}

	_ = masks != nil

	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	parallelize(options.Parallelization, srcb.Min.Y, srcb.Max.Y, func(pmin, pmax int) {
		// init temp rows
		starty := pmin
		rows := make([][]pixel, ksize)
		for i := 0; i < ksize; i++ {
			rowy := starty + i - kcenter
			if rowy < srcb.Min.Y {
				rowy = srcb.Min.Y
			} else if rowy > srcb.Max.Y-1 {
				rowy = srcb.Max.Y - 1
			}
			row := make([]pixel, srcb.Dx())
			pixGetter.getPixelRow(rowy, &row)
			rows[i] = row
		}

		for y := pmin; y < pmax; y++ {
			// calculate dst row
			for x := srcb.Min.X; x < srcb.Max.X; x++ {
				var r, g, b, a float32 = math.MaxFloat32, math.MaxFloat32, math.MaxFloat32, math.MaxFloat32
				for _, w := range masks {
					wx := x + w.u
					if wx < srcb.Min.X {
						wx = srcb.Min.X
					} else if wx > srcb.Max.X-1 {
						wx = srcb.Max.X - 1
					}
					rowsx := wx - srcb.Min.X
					rowsy := kcenter + w.v
					px := rows[rowsy][rowsx]

					r = minf32(r, px.R*w.weight)
					g = minf32(g, px.G*w.weight)
					b = minf32(b, px.B*w.weight)
					a = minf32(a, px.A*w.weight)
				}

				pixSetter.setPixel(dstb.Min.X+x-srcb.Min.X, dstb.Min.Y+y-srcb.Min.Y, pixel{r, g, b, a})
			}
			// rotate temp rows
			if y < pmax-1 {
				tmprow := rows[0]
				for i := 0; i < ksize-1; i++ {
					rows[i] = rows[i+1]
				}
				nextrowy := y + ksize/2 + 1
				if nextrowy > srcb.Max.Y-1 {
					nextrowy = srcb.Max.Y - 1
				}
				pixGetter.getPixelRow(nextrowy, &tmprow)
				rows[ksize-1] = tmprow
			}
		}
	})
}

func Erosion(kernel []float32) Filter {
	return &morphologyErosion{
		kernel: kernel,
	}
}

type morphologyDilation struct {
	kernel []float32
}

func (*morphologyDilation) Bounds(srcBounds image.Rectangle) (dstBounds image.Rectangle) {
	return image.Rect(0, 0, srcBounds.Dx(), srcBounds.Dy())
}

func (p *morphologyDilation) Draw(dst draw.Image, src image.Image, options *Options) {
	if options == nil {
		options = &defaultOptions
	}

	srcb := src.Bounds()
	dstb := dst.Bounds()

	if srcb.Dx() <= 0 || srcb.Dy() <= 0 {
		return
	}

	ksize, masks := prepareMorphologyMasks(p.kernel)
	kcenter := ksize / 2

	if ksize < 1 {
		copyimage(dst, src, options)
		return
	}

	_ = masks != nil

	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	parallelize(options.Parallelization, srcb.Min.Y, srcb.Max.Y, func(pmin, pmax int) {
		// init temp rows
		starty := pmin
		rows := make([][]pixel, ksize)
		for i := 0; i < ksize; i++ {
			rowy := starty + i - kcenter
			if rowy < srcb.Min.Y {
				rowy = srcb.Min.Y
			} else if rowy > srcb.Max.Y-1 {
				rowy = srcb.Max.Y - 1
			}
			row := make([]pixel, srcb.Dx())
			pixGetter.getPixelRow(rowy, &row)
			rows[i] = row
		}

		for y := pmin; y < pmax; y++ {
			// calculate dst row
			for x := srcb.Min.X; x < srcb.Max.X; x++ {
				var r, g, b, a float32 = -math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32
				for _, w := range masks {
					wx := x + w.u
					if wx < srcb.Min.X {
						wx = srcb.Min.X
					} else if wx > srcb.Max.X-1 {
						wx = srcb.Max.X - 1
					}
					rowsx := wx - srcb.Min.X
					rowsy := kcenter + w.v
					px := rows[rowsy][rowsx]

					r = maxf32(r, px.R)
					g = maxf32(g, px.G)
					b = maxf32(b, px.B)
					a = maxf32(a, px.A)
				}

				pixSetter.setPixel(dstb.Min.X+x-srcb.Min.X, dstb.Min.Y+y-srcb.Min.Y, pixel{r, g, b, a})
			}
			// rotate temp rows
			if y < pmax-1 {
				tmprow := rows[0]
				for i := 0; i < ksize-1; i++ {
					rows[i] = rows[i+1]
				}
				nextrowy := y + ksize/2 + 1
				if nextrowy > srcb.Max.Y-1 {
					nextrowy = srcb.Max.Y - 1
				}
				pixGetter.getPixelRow(nextrowy, &tmprow)
				rows[ksize-1] = tmprow
			}
		}
	})
}

type morphGroup byte

type groupFilter struct {
	filters []Filter
}

func (*groupFilter) Bounds(srcBounds image.Rectangle) (dstBounds image.Rectangle) {
	return image.Rect(0, 0, srcBounds.Dx(), srcBounds.Dy())
}

func (p *groupFilter) Draw(dst draw.Image, src image.Image, options *Options) {

	var tmpImage *image.NRGBA
	last := len(p.filters) - 1

	for i, f := range p.filters {

		if i == last {
			f.Draw(dst, src, options)
			continue
		}

		tmpImage = image.NewNRGBA(dst.Bounds())
		f.Draw(tmpImage, src, options)
		src = tmpImage
	}
}

func Dilation(kernel []float32) Filter {
	return &morphologyDilation{
		kernel: kernel,
	}
}

func Opening(kernel []float32) Filter {
	return &groupFilter{
		filters: []Filter{
			Erosion(kernel),
			Dilation(kernel),
		},
	}
}

func Closing(kernel []float32) Filter {
	return &groupFilter{
		filters: []Filter{
			Dilation(kernel),
			Erosion(kernel),
		},
	}
}
