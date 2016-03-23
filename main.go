package main

import (
	"./tests"
	"container/list"
	"image"
	"image/jpeg"
	"math/rand"
	"os"
	"time"
)

func main() {
	delaunay();

}

func delaunay()  {

	list := list.New()

	fp, err := os.Open("test.jpg")
	if err != nil {
		panic(err)
	}

	img, err := jpeg.Decode(fp)
	fp.Close()

	w := img.Bounds().Dx()
	l := img.Bounds().Dy()


	list.PushBack(image.Point{0,0})
	list.PushBack(image.Point{w,l})
	list.PushBack(image.Point{0,l})
	list.PushBack(image.Point{w,0})

	p := make(map[image.Point]bool)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tp := image.Point{r.Int() % w, r.Int() % l}
		if !p[tp] {
			p[tp] = true
			list.PushBack(tp)
		}
	}
	tests.Delaunay(*list, img)
}
