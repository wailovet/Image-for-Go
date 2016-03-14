package main

import
(
	"./Image"
)
func main() {

	gb := Image.NewGaussianBlur("test.jpg")
	gb.Radius(5)
	gb.To("new_test.jpg")

}