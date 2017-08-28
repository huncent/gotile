package tile_surge

import (
	m "github.com/murphy214/mercantile"
	pc "github.com/murphy214/polyclip"
)

func Pos() []int32 {
	return []int32{0, 0}
}

func moverow(row []int32, geometry []uint32) []uint32 {
	geometry = append(geometry, moveTo(1))
	geometry = append(geometry, uint32(paramEnc(row[0])))
	geometry = append(geometry, uint32(paramEnc(row[1])))
	return geometry

}
func linerow(row []int32, geometry []uint32) []uint32 {
	geometry = append(geometry, uint32(paramEnc(row[0])))
	geometry = append(geometry, uint32(paramEnc(row[1])))
	return geometry

}

func cmdEnc(id uint32, count uint32) uint32 {
	return (id & 0x7) | (count << 3)
}

func moveTo(count uint32) uint32 {
	return cmdEnc(1, count)
}

func lineTo(count uint32) uint32 {
	return cmdEnc(2, count)
}

func closePath(count uint32) uint32 {
	return cmdEnc(7, count)
}

func paramEnc(value int32) int32 {
	return (value << 1) ^ (value >> 31)
}

func Make_Point(row []float64, oldrow []int32, bound m.Extrema) ([]uint32, []int32) {
	deltax := (bound.E - bound.W)
	deltay := (bound.N - bound.S)

	factorx := (row[0] - bound.W) / deltax
	factory := (bound.N - row[1]) / deltay

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

	coords := []int32{xval, yval}
	geometry := []uint32{moveTo(uint32(1))}
	geometry = linerow([]int32{coords[0] - oldrow[0], coords[1] - oldrow[1]}, geometry)

	return geometry, coords

}

func Make_Line_Geom(coords [][]int32, position []int32) ([]uint32, []int32) {
	var count uint32
	count = 0
	var geometry []uint32
	var oldrow []int32
	//total := map[uint32][]int32{}
	//var linetocount uint32
	linetocount := uint32(len(coords) - 1)

	for _, row := range coords {
		if count == 0 {
			geometry = moverow([]int32{row[0] - position[0], row[1] - position[1]}, geometry)
			geometry = append(geometry, lineTo(linetocount))

			count = 1
		} else {
			geometry = linerow([]int32{row[0] - oldrow[0], row[1] - oldrow[1]}, geometry)
		}
		oldrow = row
	}

	return geometry, oldrow
}

// reverses the coord list
func reverse(coord [][]int32) [][]int32 {
	current := len(coord) - 1
	newlist := [][]int32{}
	for current != -1 {
		newlist = append(newlist, coord[current])
		current = current - 1
	}
	return newlist
}

// reverses the coord list
func reverse_float(coord [][]float64) [][]float64 {
	current := len(coord) - 1
	newlist := [][]float64{}
	for current != -1 {
		newlist = append(newlist, coord[current])
		current = current - 1
	}
	return newlist
}

// asserts a winding order
func assert_winding_order(coord [][]int32, exp_orient string) [][]int32 {
	count := 0
	firstpt := coord[0]
	weight := 0.0
	var oldpt []int32
	for _, pt := range coord {
		if count == 0 {
			count = 1
		} else {
			weight += float64((pt[0] - oldpt[0]) * (pt[1] + oldpt[1]))
		}
		oldpt = pt
	}

	weight += float64((firstpt[0] - oldpt[0]) * (firstpt[1] + oldpt[1]))
	var orientation string
	if weight > 0 {
		orientation = "clockwise"
	} else {
		orientation = "counter"
	}

	if orientation != exp_orient {
		return reverse(coord)
	} else {
		return coord
	}
	return coord
}

// asserts a winding order
func assert_winding_order_float(coord [][]float64, exp_orient string) [][]float64 {
	count := 0
	firstpt := coord[0]
	weight := 0.0
	var oldpt []float64
	for _, pt := range coord {
		if count == 0 {
			count = 1
		} else {
			weight += float64((pt[0] - oldpt[0]) * (pt[1] + oldpt[1]))
		}
		oldpt = pt
	}

	weight += float64((firstpt[0] - oldpt[0]) * (firstpt[1] + oldpt[1]))
	var orientation string
	if weight > 0 {
		orientation = "clockwise"
	} else {
		orientation = "counter"
	}

	if orientation != exp_orient {
		return reverse_float(coord)
	} else {
		return coord
	}
	return coord
}

// makes a polygon for a list of polygon geometries.
func Make_Polygon(coords [][][]int32, position []int32) ([]uint32, []int32) {
	var count uint32
	count = 0
	var geometry []uint32
	var oldrow []int32
	//total := map[uint32][]int32{}
	//var linetocount uint32

	for i, coord := range coords {
		if i == 0 {
			coord = assert_winding_order(coord, "clockwise")
			linetocount := uint32(len(coord) - 1)

			for _, row := range coord {
				if count == 0 {
					geometry = moverow([]int32{row[0] - position[0], row[1] - position[1]}, geometry)
					geometry = append(geometry, lineTo(linetocount))

					count = 1
				} else {
					geometry = linerow([]int32{row[0] - oldrow[0], row[1] - oldrow[1]}, geometry)
				}
				oldrow = row

			}
			geometry = append(geometry, closePath(1))
		} else {
			count = 0
			coord = assert_winding_order(coord, "counter")
			linetocount := uint32(len(coord) - 1)

			for _, row := range coord {
				if count == 0 {
					geometry = moverow([]int32{row[0] - oldrow[0], row[1] - oldrow[1]}, geometry)
					geometry = append(geometry, lineTo(linetocount))

					count = 1
				} else {
					geometry = linerow([]int32{row[0] - oldrow[0], row[1] - oldrow[1]}, geometry)
				}
				oldrow = row

			}
			geometry = append(geometry, closePath(1))
		}

	}

	return geometry, oldrow
}

// makes a polygon for a list of polygon geometries.
func Make_Polygon_Float(coords [][][]float64, position []int32, bds m.Extrema) ([]uint32, []int32) {
	var count uint32
	count = 0
	var geometry []uint32
	var oldrow []int32
	//total := map[uint32][]int32{}
	//var linetocount uint32

	for i, coord := range coords {
		if i == 0 {
			coord = assert_winding_order_float(coord, "clockwise")
			linetocount := uint32(len(coord) - 1)

			for _, rowi := range coord {
				row := single_point(pc.Point{rowi[0], rowi[1]}, bds)
				if count == 0 {
					geometry = moverow([]int32{row[0] - position[0], row[1] - position[1]}, geometry)
					geometry = append(geometry, lineTo(linetocount))

					count = 1
				} else {
					geometry = linerow([]int32{row[0] - oldrow[0], row[1] - oldrow[1]}, geometry)
				}
				oldrow = row

			}
			geometry = append(geometry, closePath(1))
		} else {
			count = 0
			coord = assert_winding_order_float(coord, "counter")
			linetocount := uint32(len(coord) - 1)

			for _, rowi := range coord {
				row := single_point(pc.Point{rowi[0], rowi[1]}, bds)
				if count == 0 {
					geometry = moverow([]int32{row[0] - oldrow[0], row[1] - oldrow[1]}, geometry)
					geometry = append(geometry, lineTo(linetocount))

					count = 1
				} else {
					geometry = linerow([]int32{row[0] - oldrow[0], row[1] - oldrow[1]}, geometry)
				}
				oldrow = row

			}
			geometry = append(geometry, closePath(1))
		}

	}

	return geometry, oldrow
}
