package tile_surge

import (
	pc "github.com/murphy214/polyclip"
	m "github.com/murphy214/mercantile"

	"testing"

)

// test for distance
func Test_Distance(t *testing.T) {
	testcases := []struct {
			Expected  float64
			Pt1 pc.Point
			Pt2 pc.Point
	}{	
		{
			Expected:0.3255764119219919,
			Pt1:pc.Point{-82.324,43.232},
			Pt2:pc.Point{-82.0,43.2},
		},
	}

	for _, tcase := range testcases {
		val := Distance(tcase.Pt1,tcase.Pt2)
		if val != tcase.Expected {
			t.Errorf("Distance Error, Expected %f Got %f",tcase.Expected,val)
		}
	}
}

// test for distance
func Test_Single_Point(t *testing.T) {
	testcases := []struct {
			Expected  []int32
			Pt pc.Point
			Bd m.Extrema
	}{	
		{
			Expected:[]int32{1877, 3750},
			Pt:pc.Point{-82.324,40.0},
			Bd:m.Bounds(m.TileID{X:69,Y:96,Z:8}),
		},
	}

	for _, tcase := range testcases {
		val := single_point(tcase.Pt,tcase.Bd)
		if (val[0] != tcase.Expected[0]) || (val[1] != tcase.Expected[1]) {
			t.Errorf("Single Point Error, Expected %s Got %s",tcase.Expected,val)
		}
	}
}


// test for Make_Coords_Float
func Test_Make_Coords_Float(t *testing.T) {
	testcases := []struct {
			Expected  [][]int32
			Coord [][]float64
			Bd m.Extrema
			Tileid m.TileID
	}{	
		{
			Expected:[][]int32{{1551,4005},{1546,4062},{1545,4068},{1545,4071},{1545,4087},{1545,4095}},
			Coord:[][]float64{{-77.913702,39.64553},{-77.913702,39.64553},{-77.914165,39.641742},{-77.91422,39.641348},{-77.914224,39.641122},{-77.914236,39.640089},{-77.914241110187,39.6395375643667}},
			Bd:m.Bounds(m.TileID{X:290,Y:388,Z:10}),
			Tileid:m.TileID{X:290,Y:388,Z:10},
		},
	}

	for _, tcase := range testcases {
		coords := Make_Coords_Float(tcase.Coord,tcase.Bd,tcase.Tileid)
		for i := range coords {
			valmine := coords[i]
			valtcase := tcase.Expected[i]
			if (valmine[0] != valtcase[0]) || (valmine[1] != valtcase[1]) {
				t.Errorf("Make_Coords_Float Error, Expected %s Got %s",valtcase,valmine)
			}
		}
	}
}

// test for Make_Coords_Float
func Test_Make_Coords_Polygon_Float(t *testing.T) {
	testcases := []struct {
			Expected  [][][]int32
			Coord [][][]float64
			Bd m.Extrema
	}{	
		{
			Expected:[][][]int32{{{4022,2891},{3536,3038},{3488,2439},{3545,1460},{3099,941},{2609,439},{2050,0},{0,0},{0,4095},{4095,4095},{4095,2892}}},
			Coord:[][][]float64{{{-155.045382,19.739824},{-155.087118,19.728013},{-155.091216,19.776367999999998},{-155.086341,19.855399},{-155.124618,19.897288},{-155.166625,19.93789},{-155.21460251667767,19.9733487861106},{-155.390625,19.9733487861106},{-155.390625,19.642587534013032},{-155.0390625,19.642587534013032},{-155.0390625,19.73973673156395}}},
			Bd:m.Extrema{-155.390625, -155.0390625, 19.9733487861106, 19.642587534013032},
		},
	}

	for _, tcase := range testcases {
		rings := Make_Coords_Polygon_Float(tcase.Coord,tcase.Bd)
		for coordsi,coords := range rings {
			for i := range coords {
				valmine := rings[coordsi][i]
				valtcase := tcase.Expected[coordsi][i]
				if (valmine[0] != valtcase[0]) || (valmine[1] != valtcase[1]) {
					t.Errorf("Make_Coords_Polygon_Float Error, Expected %s Got %s",valtcase,valmine)
				}
			}
		}
	}
}

	