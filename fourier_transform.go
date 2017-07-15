package gift

import (
	"image"
	"image/draw"
	"math"
	"math/cmplx"

	giftimage "github.com/disintegration/gift/image"
	// giftcolor "github.com/disintegration/gift/image/color"
)

type discreteFourierTransform struct{}

func (*discreteFourierTransform) Bounds(src image.Rectangle) image.Rectangle {
	return image.Rect(0, 0, src.Dx(), src.Dy())
}

func (h *discreteFourierTransform) Draw(dst draw.Image, src image.Image, options *Options) {
	if options == nil {
		options = &defaultOptions
	}

	srcb := src.Bounds()
	dstb := dst.Bounds()

	if srcb.Dx() <= 0 || srcb.Dy() <= 0 {
		return
	}

	cSrc := giftimage.NewC64RGBA(srcb)
	cDst := giftimage.NewC64RGBA(dstb)

	N := srcb.Dx() * srcb.Dy()
	invN := 1 / float64(N)

	c := cmplx.Exp(complex(0, -2*math.Pi*invN))

	// cmplx.Exp(complex(0, -2*math.PI*))

	copyimage(cSrc, src, options)

	cAcc := complex(1, 0)
	for i, k := 0, 0; i < len(cDst.Pix); i, k = i+4, k+1 {
		cAcc *= c
		ck := cAcc
		ckAcc := complex(1, 0)
		var r, g, b complex128

		for j, n := 0, 0; j < len(cSrc.Pix); j, n = j+4, n+1 {
			ckAcc *= ck
			r += complex128(cSrc.Pix[j+0]) * ckAcc
			g += complex128(cSrc.Pix[j+1]) * ckAcc
			b += complex128(cSrc.Pix[j+2]) * ckAcc
		}

		cDst.Pix[i+0] = complex64(r)
		cDst.Pix[i+1] = complex64(g)
		cDst.Pix[i+2] = complex64(b)
		cDst.Pix[i+3] = complex(1, 0)
	}

	copyimage(dst, cDst, options)
}

var _ Filter = &discreteFourierTransform{}

func DiscreteFourierTransform() Filter {
	return new(discreteFourierTransform)
}
