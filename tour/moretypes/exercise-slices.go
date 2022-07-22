package main

import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
	p := make([][]uint8, dy)
	for u := 0; u < dy; u++ {
		p[u] = make([]uint8, dx)
		for v := 0; v < dx; v++ {
			//p[u][v] = uint8((u+v)/2)
			//p[u][v] = uint8(u*v)
			p[u][v] = uint8(u ^ v)
		}
	}
	return p
}

func main() {
	pic.Show(Pic)
}
