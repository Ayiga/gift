package gift

import (
	"image"
	"image/draw"
	"math"
)

type bilateralFilter struct {
	kernelSize    int
	gaussianSigma float32
	colorSigma    float32
}

func (*bilateralFilter) Bounds(srcBounds image.Rectangle) (dstBounds image.Rectangle) {
	dstBounds = image.Rect(0, 0, srcBounds.Dx(), srcBounds.Dy())
	return
}

func (p *bilateralFilter) Draw(dst draw.Image, src image.Image, options *Options) {
	if options == nil {
		options = &defaultOptions
	}

	srcb := src.Bounds()
	dstb := dst.Bounds()

	if srcb.Dx() <= 0 || srcb.Dy() <= 0 {
		return
	}

	ksize := p.kernelSize
	kcenter := ksize / 2

	if ksize < 1 {
		copyimage(dst, src, options)
		return
	}

	gaussianWeights := make([]uvweight, 0, ksize*ksize)

	for i := 0; i < ksize; i++ {
		for j := 0; j < ksize; j++ {
			gaussianWeights = append(gaussianWeights, uvweight{
				u:      i - kcenter,
				v:      j - kcenter,
				weight: float32(math.Hypot(float64(i-kcenter), float64(j-kcenter))) / p.gaussianSigma,
			})
		}
	}

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
				var r, g, b, a float32
				var sumWeight float64
				for _, w := range gaussianWeights {
					wx := x + w.u
					if wx < srcb.Min.X {
						wx = srcb.Min.X
					} else if wx > srcb.Max.X-1 {
						wx = srcb.Max.X - 1
					}
					rowsx := wx - srcb.Min.X
					rowsy := kcenter + w.v
					px := rows[rowsy][rowsx]
					cpx := rows[kcenter][x]

					imageDist := w.weight
					colorDist := math.Sqrt(float64(((px.R-cpx.R)*(px.R-cpx.R))+((px.G-cpx.G)*(px.G-cpx.G))+((px.B-cpx.B)*(px.B-cpx.B)))) / float64(p.colorSigma)

					weight := 1.0 / (math.Exp(float64(imageDist*imageDist*0.5)) * math.Exp(colorDist*colorDist*0.5))

					sumWeight += weight
					r += float32(weight) * px.R
					g += float32(weight) * px.G
					b += float32(weight) * px.B
					a = cpx.A
				}

				r /= float32(sumWeight)
				g /= float32(sumWeight)
				b /= float32(sumWeight)

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

func BilateralFilter(kernelSize int, gaussianSigma, colorSigma float32) Filter {
	return &bilateralFilter{
		kernelSize:    kernelSize,
		gaussianSigma: gaussianSigma,
		colorSigma:    colorSigma,
	}
}
