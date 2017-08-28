package tile_surge

import (
	//"fmt"
	m "github.com/murphy214/mercantile"
	pc "github.com/murphy214/polyclip"
	"math"
)

// Distance finds the length of the hypotenuse between two points.
// Forumula is the square root of (x2 - x1)^2 + (y2 - y1)^2
func Distance(p1 pc.Point, p2 pc.Point) float64 {
	first := math.Pow(float64(p2.X-p1.X), 2)
	second := math.Pow(float64(p2.Y-p1.Y), 2)
	return math.Sqrt(first + second)
}

// converts a single point from a coordinate to a tile point
func single_point(row pc.Point, bound m.Extrema) []int32 {
	deltax := (bound.E - bound.W)
	deltay := (bound.N - bound.S)

	factorx := (row.X - bound.W) / deltax
	factory := (bound.N - row.Y) / deltay

	xval := int32(factorx * 4096)
	yval := int32(factory * 4096)

	//here1 := uint32((row[0] - bound.w) / (bound.e - bound.w))
	//here2 := uint32((bound.n-row[1])/(bound.n-bound.s)) * 4096
	if xval >= 4095 {
		xval = 4095
	}

	if yval >= 4095 {
		yval = 4095
	}

	return []int32{xval, yval}
}

// makes coordinates for a line
func Make_Coords(coord []pc.Point, bound m.Extrema, tileid m.TileID) [][]int32 {
	var newlist [][]int32
	var oldpt []int32
	east := int32(0)
	west := int32(4095)
	south := int32(4095)
	north := int32(0)

	for ii, i := range coord {
		pt := single_point(i, bound)

		if ii == 0 {
			newlist = append(newlist, pt)
		} else {
			if ((oldpt[0] == pt[0]) && (oldpt[1] == pt[1])) == false {
				newlist = append(newlist, pt)
			}
		}
		oldpt = pt

		//
		if pt[0] > east {
			east = pt[0]
		}
		if pt[0] < west {
			west = pt[0]
		}

		if pt[1] < south {
			south = pt[1]
		}
		if pt[1] > north {
			north = pt[1]
		}
	}

	return newlist

}

// makes coordinates for a line that is float
func Make_Coords_Float(coord [][]float64, bound m.Extrema, tileid m.TileID) [][]int32 {
	var newlist [][]int32
	var oldpt []int32
	east := int32(0)
	west := int32(4095)
	south := int32(4095)
	north := int32(0)

	for ii, i := range coord {
		pt := single_point(pc.Point{i[0], i[1]}, bound)

		if ii == 0 {
			newlist = append(newlist, pt)
		} else {
			if ((oldpt[0] == pt[0]) && (oldpt[1] == pt[1])) == false {
				newlist = append(newlist, pt)
			}
		}
		oldpt = pt

		//
		if pt[0] > east {
			east = pt[0]
		}
		if pt[0] < west {
			west = pt[0]
		}

		if pt[1] < south {
			south = pt[1]
		}
		if pt[1] > north {
			north = pt[1]
		}
	}

	return newlist

}

// makes polygon layer for cordinate positions
func Make_Coords_Polygon(polygon pc.Polygon, bound m.Extrema) [][][]int32 {
	var newlist [][][]int32

	for _, cont := range polygon {
		newcont := [][]int32{}
		for _, i := range cont {
			newcont = append(newcont, single_point(i, bound))
		}
		newlist = append(newlist, newcont)
	}
	return newlist

}

// makes polygon layer for cordinate positions
func Make_Coords_Polygon_Float(polygon [][][]float64, bound m.Extrema) [][][]int32 {
	var newlist [][][]int32
	oldpt := []int32{0, 0}

	for _, cont := range polygon {
		newcont := [][]int32{}
		count := 0
		for _, i := range cont {
			pt := single_point(pc.Point{i[0], i[1]}, bound)
			if count == 0 {
				count = 1
				newcont = append(newcont, pt)

			} else {
				if ((oldpt[0] == pt[0]) && (oldpt[1] == pt[1])) == false {
					newcont = append(newcont, pt)
				}

			}
			oldpt = pt
		}
		newlist = append(newlist, newcont)
	}

	return newlist

}
