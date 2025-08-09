package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand/v2"
)

func HSLtoRGB(h, s, l float64) (int, int, int) {
	chroma := (1 - math.Abs(2*l-1)) * s
	hPrime := h / 60
	x := chroma * (1 - math.Abs(float64(int(hPrime)%2-1)))
	var r1, g1, b1 float64
	if 0 <= hPrime && hPrime < 1 {
		r1, g1, b1 = chroma, x, 0
	} else if 1 <= hPrime && hPrime < 2 {
		r1, g1, b1 = x, chroma, 0
	} else if 2 <= hPrime && hPrime < 3 {
		r1, g1, b1 = 0, chroma, x
	} else if 3 <= hPrime && hPrime < 4 {
		r1, g1, b1 = 0, x, chroma
	} else if 4 <= hPrime && hPrime < 5 {
		r1, g1, b1 = x, 0, chroma
	} else if 5 <= hPrime && hPrime < 6 {
		r1, g1, b1 = chroma, 0, x
	} else {
		panic("Invalid HSL conversion")
	}
	m := int(l - (chroma / 2))
	return int(r1) + m, int(g1) + m, int(b1) + m
}

func main() {
	scale := float64(200)
	minX := -2.0 * scale
	maxX := .47 * scale
	minY := -1.12 * scale
	maxY := 1.12 * scale
	maxIter := 2000
	pallette := make([][]int, maxIter)
	for i := range maxIter {
		pallette[i] = []int{
			rand.IntN(255),
			rand.IntN(255),
			rand.IntN(255),
		}
	}
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			z := complex128(0)

			i := 0
			for i = range maxIter {
				z = cmplx.Pow(z, 2) + complex(float64(x)/scale, float64(y)/200)
				if cmplx.Abs(z) > 4 {
					break
				}
			}
			if i == maxIter-1 {
				fmt.Printf("\033[30m■")
			} else {
				// r, g, b := HSLtoRGB(float64(int(math.Pow((float64(i)/float64(maxIter))*360, 1.5))%360), 0.5, (float64(i) / float64(maxIter) * 20 * 5))
				fmt.Printf("\033[38;2;%d;%d;%dm■", pallette[i][0], pallette[i][1], pallette[i][2])
				// fmt.Printf("\033[38;2;%d;%d;%dm■", r, g, b)
			}
		}
		fmt.Println()
	}
}
