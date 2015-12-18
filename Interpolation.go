package main

import (
	"image"
	"image/color"
	"math"
)

func bicubicInterpolate(p [4][4]float64, x, y float64) float64 {
	var a00, a01, a02, a03, a10, a11, a12, a13, a20, a21, a22, a23, a30, a31, a32, a33 float64
	a00 = p[1][1]
	a01 = -.5*p[1][0] + .5*p[1][2]
	a02 = p[1][0] - 2.5*p[1][1] + 2*p[1][2] - .5*p[1][3]
	a03 = -.5*p[1][0] + 1.5*p[1][1] - 1.5*p[1][2] + .5*p[1][3]
	a10 = -.5*p[0][1] + .5*p[2][1]
	a11 = .25*p[0][0] - .25*p[0][2] - .25*p[2][0] + .25*p[2][2]
	a12 = -.5*p[0][0] + 1.25*p[0][1] - p[0][2] + .25*p[0][3] + .5*p[2][0] - 1.25*p[2][1] + p[2][2] - .25*p[2][3]
	a13 = .25*p[0][0] - .75*p[0][1] + .75*p[0][2] - .25*p[0][3] - .25*p[2][0] + .75*p[2][1] - .75*p[2][2] + .25*p[2][3]
	a20 = p[0][1] - 2.5*p[1][1] + 2*p[2][1] - .5*p[3][1]
	a21 = -.5*p[0][0] + .5*p[0][2] + 1.25*p[1][0] - 1.25*p[1][2] - p[2][0] + p[2][2] + .25*p[3][0] - .25*p[3][2]
	a22 = p[0][0] - 2.5*p[0][1] + 2*p[0][2] - .5*p[0][3] - 2.5*p[1][0] + 6.25*p[1][1] - 5*p[1][2] + 1.25*p[1][3] + 2*p[2][0] - 5*p[2][1] + 4*p[2][2] - p[2][3] - .5*p[3][0] + 1.25*p[3][1] - p[3][2] + .25*p[3][3]
	a23 = -.5*p[0][0] + 1.5*p[0][1] - 1.5*p[0][2] + .5*p[0][3] + 1.25*p[1][0] - 3.75*p[1][1] + 3.75*p[1][2] - 1.25*p[1][3] - p[2][0] + 3*p[2][1] - 3*p[2][2] + p[2][3] + .25*p[3][0] - .75*p[3][1] + .75*p[3][2] - .25*p[3][3]
	a30 = -.5*p[0][1] + 1.5*p[1][1] - 1.5*p[2][1] + .5*p[3][1]
	a31 = .25*p[0][0] - .25*p[0][2] - .75*p[1][0] + .75*p[1][2] + .75*p[2][0] - .75*p[2][2] - .25*p[3][0] + .25*p[3][2]
	a32 = -.5*p[0][0] + 1.25*p[0][1] - p[0][2] + .25*p[0][3] + 1.5*p[1][0] - 3.75*p[1][1] + 3*p[1][2] - .75*p[1][3] - 1.5*p[2][0] + 3.75*p[2][1] - 3*p[2][2] + .75*p[2][3] + .5*p[3][0] - 1.25*p[3][1] + p[3][2] - .25*p[3][3]
	a33 = .25*p[0][0] - .75*p[0][1] + .75*p[0][2] - .25*p[0][3] - .75*p[1][0] + 2.25*p[1][1] - 2.25*p[1][2] + .75*p[1][3] + .75*p[2][0] - 2.25*p[2][1] + 2.25*p[2][2] - .75*p[2][3] - .25*p[3][0] + .75*p[3][1] - .75*p[3][2] + .25*p[3][3]
	x2 := x * x
	x3 := x2 * x
	y2 := y * y
	y3 := y2 * y
	return (a00 + a01*y + a02*y2 + a03*y3) +
	(a10+a11*y+a12*y2+a13*y3)*x +
	(a20+a21*y+a22*y2+a23*y3)*x2 +
	(a30+a31*y+a32*y2+a33*y3)*x3
}

func safe(x float64) uint8 {
	if x > 255 {
		return 255
	} else if x < 0 {
		return 0
	} else {
		return (uint8)(x)
	}
}

func getpix(p [4][4][4]uint8, x, y float64) []uint8 {
	var ans [4]uint8
	var mat [4][4]float64
	for k := 0; k < 4; k++ {
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				mat[i][j] = float64(p[i][j][k])
			}
		}
		ans[k] = safe(bicubicInterpolate(mat, x, y))
	}
	return ans[:]
}

func bound(v, d, u int) int {
	if v > u {
		return u
	}
	if v < d {
		return d
	}
	return v
}

func GetPix(m image.Image, u, v float64, back color.Color) color.Color {
	if u < 0 || v < 0 || u > 1 || v > 1 {
		return back
	}
	ww := m.Bounds().Dx()
	hh := m.Bounds().Dy()
	u = u*(float64)(ww) - 0.5
	v = v*(float64)(hh) - 0.5

	xx := math.Floor(u)
	yy := math.Floor(v)
	x := int(xx)
	y := int(yy)

	u_ratio := (u - (float64)(x))
	v_ratio := (v - (float64)(y))

	xb := ww
	yb := hh
	var p [4][4][4]uint8

	switch m := m.(type) {
	case *image.RGBA:
		for i := 0; i <= 3; i++ {
			for j := 0; j <= 3; j++ {
				off := (bound(x+i-1, 0, xb-1) + bound(y+j-1, 0, yb-1)*ww) * 4
				for k := 0; k <= 3; k++ {
					p[i][j][k] = m.Pix[off+k]
				}
			}
		}
	}
	ans := getpix(p, u_ratio, v_ratio)
	return color.RGBA{ans[0], ans[1], ans[2], ans[3]}
}
