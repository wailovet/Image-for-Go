package Image

import (
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"time"
)

type GaussianBlur struct {
	input_file_name         string
	output_file_name        string
	radius                  int
	_image                  image.Image
	all_weights             float64
	gaussian_distribution_o float64
}

func NewGaussianBlur(file_name string) *GaussianBlur {
	return &GaussianBlur{file_name, "", 1, nil, 0, 1}
}

func (this *GaussianBlur) Radius(r int) *GaussianBlur {
	this.radius = r
	this.gaussian_distribution_o = float64(r) / 3

	for i := -this.radius; i < this.radius + 1; i++ {
		for k := -this.radius; k < this.radius + 1; k++ {
			this.all_weights += gaussianDistribution(i, k, this.gaussian_distribution_o)
		}
	}
	return this
}

func (this *GaussianBlur) avg(x int, y int) color.RGBA64 {
	all_weights := this.all_weights
	rgba := color.RGBA64{0, 0, 0, 0}
	for i := -this.radius; i < this.radius + 1; i++ {
		for k := -this.radius; k < this.radius + 1; k++ {
			//x + i, y + k加权后的值
			_rgba := gaussianDistributionToRGBA(this._image, x + i, y + k, i, k, this.gaussian_distribution_o)
			rgba = mergeRGBA(rgba, _rgba)
		}
	}
	return exRGBA(rgba, all_weights)

}

func exRGBA(rgba color.RGBA64, weights float64) color.RGBA64 {
	_r := float64(rgba.R) / weights
	_g := float64(rgba.G) / weights
	_b := float64(rgba.B) / weights
	_a := float64(rgba.A) / weights

	return color.RGBA64{uint16(_r), uint16(_g), uint16(_b), uint16(_a)}
}

func mergeRGBA(a color.RGBA64, b color.RGBA64) color.RGBA64 {
	_r := a.R + b.R
	_g := a.G + b.G
	_b := a.B + b.B
	_a := a.A + b.A

	return color.RGBA64{uint16(_r), uint16(_g), uint16(_b), uint16(_a)}
}
func gaussianDistributionToRGBA(_image image.Image, x int, y int, w_x int, w_y int, gaussian_distribution_o float64) color.RGBA64 {
	if (x < 0) {
		x *= -1
	}
	if (y < 0) {
		y *= -1
	}
	if (x >= _image.Bounds().Size().X) {
		x = _image.Bounds().Size().X - 1
	}
	if (y >= _image.Bounds().Size().Y) {
		y = _image.Bounds().Size().Y - 1
	}

	r, g, b, a := _image.At(x, y).RGBA()
	_r := float64(r) * gaussianDistribution(w_x, w_y, gaussian_distribution_o)
	_g := float64(g) * gaussianDistribution(w_x, w_y, gaussian_distribution_o)
	_b := float64(b) * gaussianDistribution(w_x, w_y, gaussian_distribution_o)
	_a := float64(a) * gaussianDistribution(w_x, w_y, gaussian_distribution_o)
	return color.RGBA64{uint16(_r), uint16(_g), uint16(_b), uint16(_a)}
}

var gdb [128][128]float64

func gaussianDistribution(_x int, _y int, o float64) float64 {
	if _x < 0 {
		_x *= -1
	}
	if _y < 0 {
		_y *= -1
	}
	if gdb[_x][_y] != 0 {
		return gdb[_x][_y]
	}

	x := float64(_x)
	y := float64(_y)
	a := 1 / (2 * math.Pi * o * o)
	b := -(x * x + y * y) / (2 * o * o)
	c := a * math.Exp(b)
	gdb[_x][_y] = float64(c)
	return float64(c)
}

func (this *GaussianBlur) To(file_name string) {
	this.output_file_name = file_name

	fp, err := os.Open(this.input_file_name)
	if err != nil {
		panic(err)
	}

	img, err := jpeg.Decode(fp)
	fp.Close()
	if err != nil {
		panic(err)
	}
	this._image = img

	new_img, err := os.Create(this.output_file_name)
	if err != nil {
		panic(err)
	}
	m := image.NewRGBA(this._image.Bounds())
	x := img.Bounds().Size().X
	y := img.Bounds().Size().Y

	ok := make(chan bool)

	count := this.radius * 2
	if count > 50 {
		count = 50
	}
	one_x := x / count
	one_x_ := x % count
	run_count := 0
	for k := 0; k < count; k++ {
		go func() {
			run_count ++
			head := k * one_x
			for i := head; i < head + one_x + one_x_; i++ {
				for j := 0; j < y; j++ {
					m.Set(i, j, this.avg(i, j))
				}
			}
			run_count--
			if (run_count <= 0) {
				ok <- true
			}
		}();
		time.Sleep(1e8)
	}

	<-ok
	jpeg.Encode(new_img, m, nil)
}
