package Image

import (
	"fmt"
	"image"
	"math"
)

type Triangle struct {
	p [3]image.Point
}

func NewTriangle(x1 int, y1 int, x2 int, y2 int, x3 int, y3 int) *Triangle {
	if (x1 == x2 && x2 == x3) || float64(y1-y2)/float64(x1-x2) == float64(y3-y2)/float64(x3-x2) {
		fmt.Print(image.Point{x1, y1}, image.Point{x2, y2}, image.Point{x3, y3}, "\n")
		panic("This is not a triangle.")
	}
	var p = [3]image.Point{image.Point{x1, y1}, image.Point{x2, y2}, image.Point{x3, y3}}
	return &Triangle{p}
}

func (this *Triangle) GetPoint(i int) image.Point {
	return this.p[i]
}

func (this *Triangle) GetLine(i int) Line {
	a := this.GetPoint(i % 3)
	b := this.GetPoint((i + 1) % 3)

	vsa := a.X
	vsb := b.X
	if vsa == vsb {
		vsa = a.Y
		vsb = b.Y
	}

	if vsa > vsb {
		tmp := a
		a = b
		b = tmp
	}

	return *NewLine(a.X, a.Y, b.X, b.Y)
}

func (this *Triangle) IsInCircumcircle(ip image.Point) bool {
	int_rand_x := ip.X
	int_rand_y := ip.Y
	r, p := this.Circumcircle()
	if r > (float64(int_rand_x-p.X)*float64(int_rand_x-p.X) + float64(int_rand_y-p.Y)*float64(int_rand_y-p.Y)) {
		return true
	}
	return false
}
func (this *Triangle) Circumcircle() (float64, image.Point) {
	var b [2]float64
	var k [2]float64
	ai := 0
	bi := 1
	if this.p[ai].Y-this.p[bi].Y == 0 {
		ai = 2
		bi = 0
	}
	co := 0

	k[co] = -(float64(this.p[ai].X-this.p[bi].X) / float64(this.p[ai].Y-this.p[bi].Y))
	x1 := float64(this.p[ai].X+this.p[bi].X) / 2
	y1 := float64(this.p[ai].Y+this.p[bi].Y) / 2
	b[co] = y1 - k[co]*x1

	ai = 1
	bi = 2
	if this.p[ai].Y-this.p[bi].Y == 0 {
		ai = 2
		bi = 0
	}
	co = 1
	k[co] = -(float64(this.p[ai].X-this.p[bi].X) / float64(this.p[ai].Y-this.p[bi].Y))
	x1 = float64(this.p[ai].X+this.p[bi].X) / 2
	y1 = float64(this.p[ai].Y+this.p[bi].Y) / 2
	b[co] = y1 - k[co]*x1

	ax := float64((b[1] - b[0]) / (k[0] - k[1]))
	ay := float64(float64(k[0])*ax + float64(b[0]))
	p := image.Point{int(ax), int(ay)}
	r := ((ax-float64(this.p[0].X))*(ax-float64(this.p[0].X)) + (ay-float64(this.p[0].Y))*(ay-float64(this.p[0].Y)))
	return r, p
}

func (this *Triangle) IsInSide(o image.Point) bool {

	for i := 0; i < 3; i++ {

		line := this.GetLine(i)
		dxa := o.X - line.Get(0).X
		dxb := o.X - line.Get(1).X
		dxc := line.Get(0).X - line.Get(1).X
		dya := o.Y - line.Get(0).Y
		dyb := o.Y - line.Get(1).Y
		dyc := line.Get(0).Y - line.Get(1).Y
		da := math.Sqrt(float64(dxa*dxa + dya*dya))
		db := math.Sqrt(float64(dxb*dxb + dyb*dyb))
		dc := math.Sqrt(float64(dxc*dxc + dyc*dyc))
		if dya*dyb <= 0 || dxa*dxb <= 0 {
			if line.Get(0).X == line.Get(1).X {
				if o.X == line.Get(0).X {
					if dya*dyb <= 0 {
						return true
					}
				}

			}
			if line.Get(0).Y == line.Get(1).Y {
				if o.Y == line.Get(0).Y {
					if dxa*dxb <= 0 {
						return true
					}
				}
			}
			if (da+db-dc)*(da+db-dc) < 4 {
				return true
			}

		}
	}
	return false
}

func (this *Triangle) IsInTriangle(o image.Point) bool {
	if this.p[0].X == this.p[1].X && this.p[1].X == this.p[2].X || this.p[0].Y == this.p[1].Y && this.p[1].Y == this.p[2].Y {
		return false
	}
	var x0, y0, x1, y1, x2, y2, x3, y3 float64
	x0 = float64(o.X)
	y0 = float64(o.Y)

	x1 = float64(this.p[0].X)
	x2 = float64(this.p[1].X)
	x3 = float64(this.p[2].X)
	y1 = float64(this.p[0].Y)
	y2 = float64(this.p[1].Y)
	y3 = float64(this.p[2].Y)

	if (y2-y1)/(x2-x1) == (y3-y1)/(x3-x1) {
		return false
	}
	var k, cb, tya, tyb float64

	if (x1 - x2) != 0 {
		k = (y1 - y2) / (x1 - x2)
		cb = y1 - k*x1
		tya = y3 - (k*x3 + cb)
		tyb = y0 - (k*x0 + cb)
		if tya*tyb < 0 {
			return false
		}
	} else {
		if (x3-x1)*(x0-x1) < 0 {
			return false
		}
	}

	if (x2 - x3) != 0 {
		k = (y2 - y3) / (x2 - x3)
		cb = y2 - k*x2

		tya = y1 - (k*x1 + cb)
		tyb = y0 - (k*x0 + cb)
		if tya*tyb < 0 {
			return false
		}
	} else {
		if (x1-x2)*(x0-x2) < 0 {
			return false
		}
	}
	if (x3 - x1) != 0 {
		k = (y3 - y1) / (x3 - x1)
		cb = y1 - k*x1

		tya = y2 - (k*x2 + cb)
		tyb = y0 - (k*x0 + cb)
		if tya*tyb < 0 {
			return false
		}
	} else {
		if (x2-x1)*(x0-x1) < 0 {
			return false
		}
	}

	return true

}

type Line struct {
	p [2]image.Point
}

func NewLine(x1 int, y1 int, x2 int, y2 int) *Line {

	var p = [2]image.Point{image.Point{x1, y1}, image.Point{x2, y2}}
	return &Line{p}
}

func (this *Line) Get(i int) image.Point {
	return this.p[i]
}
