package tests

import (
	"../Image"
	"container/list"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math/rand"
	"os"
)

func Test(point list.List, img image.Image) {

	new_img := image.NewRGBA(img.Bounds())
	w := img.Bounds().Dx()
	l := img.Bounds().Dy()
	triangle_list := Image.Delaunay(point, w, l)
	m_color := make(map[Image.Triangle]color.RGBA)

	for e := triangle_list.Front(); e != nil; e = e.Next() {
		m_color[*(e.Value.(*Image.Triangle))] = color.RGBA{uint8(rand.Int() % 255), uint8(rand.Int() % 255), uint8(rand.Int() % 255), 255}
	}

	Render(w, l, new_img, m_color)
	new_file, err := os.Create("new_test.jpg")
	if err != nil {
		panic(err)
	}
	jpeg.Encode(new_file, new_img, nil)
	fmt.Print("OK\n")
}

func Render(x int, y int, new_img *image.RGBA, m_color map[Image.Triangle]color.RGBA) {

	isRender := make(map[image.Point]bool)
	for i := 0; i < x; i++ {
		for k := 0; k < y; k++ {
			new_img.Set(i, k, color.White)
		}
	}
	for triangles := range m_color {
		for i := 0; i < x; i++ {
			for k := 0; k < y; k++ {
				if !isRender[image.Point{i, k}] && triangles.IsInTriangle(image.Point{i, k}) {
					isRender[image.Point{i, k}] = true
					new_img.Set(i, k, m_color[triangles])
				}
			}
		}
	}
}
