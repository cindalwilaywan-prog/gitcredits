package main

import (
	"math/rand"
)

// star field background
type starField struct {
	stars []struct {
		x, y int
		ch   rune
	}
}

func newStarField(width, totalHeight int) starField {
	sf := starField{}
	density := (width * totalHeight) / 40
	for i := 0; i < density; i++ {
		ch := '·'
		bright := rand.Intn(10)
		if bright == 0 {
			ch = '✦'
		} else if bright <= 2 {
			ch = '✧'
		} else if bright <= 4 {
			ch = '⋆'
		} else if bright <= 6 {
			ch = '·'
		} else {
			ch = '.'
		}
		sf.stars = append(sf.stars, struct {
			x, y int
			ch   rune
		}{
			x:  rand.Intn(width),
			y:  rand.Intn(totalHeight),
			ch: ch,
		})
	}
	return sf
}
