package tile_surge

import (
	m "github.com/murphy214/mercantile"

	"testing"

)

// test for Make_Coords_Float
func Test_Make_Line_Geom(t *testing.T) {
	testcases := []struct {
			Coord [][]int32
			Row []int32
			Expected []uint32
	}{	
		{
			Coord:[][]int32{{1551,4005},{1546,4062},{1545,4068},{1545,4071},{1545,4087},{1545,4095}},
			Row:[]int32{0,0},
			Expected:[]uint32{9,3102,8010,42,9,114,1,12,0,6,0,32,0,16},
		},
	}

	for _, tcase := range testcases {
		coords,_ := Make_Line_Geom(tcase.Coord,tcase.Row)
		for i := range coords {
			valmine := coords[i]
			valtcase := tcase.Expected[i]
			if valmine != valtcase {
				t.Errorf("Make_Line_Geom Error, Expected %s Got %s",valtcase,valmine)
			}
		}
	}
}



// test for Make_Coords_Float
func Test_Make_Polygon(t *testing.T) {
	testcases := []struct {
			Coord [][][]int32
			Row []int32
			Expected  []uint32

	}{	
		{	
			Expected:[]uint32{9,8044,5782,82,971,294,95,1197,114,1957,891,1037,979,1003,1117,877,4099,0,0,8190,8190,0,0,2405,15},
			Coord:[][][]int32{{{4022,2891},{3536,3038},{3488,2439},{3545,1460},{3099,941},{2609,439},{2050,0},{0,0},{0,4095},{4095,4095},{4095,2892}}},
			Row:[]int32{0,0},
		},
	}

	for _, tcase := range testcases {
		coords,_ := Make_Polygon(tcase.Coord,tcase.Row)
		for i := range coords {
			valmine := coords[i]
			valtcase := tcase.Expected[i]
			if valmine != valtcase {
				t.Errorf("Make_Coords_Polygon_Float Error, Expected %s Got %s",valtcase,valmine)
			}
		}
	}
}


// test for Make_Coords_Float
func Test_Make_Point(t *testing.T) {
	testcases := []struct {
			Row []float64
			OldRow []int32
			Bd m.Extrema
			Expected []uint32
	}{	
		{
			OldRow:[]int32{0,0},
			Row:[]float64{-82.324,40.0},
			Bd: m.Extrema{-82.6171875, -82.265625, 40.17887331434696, 39.909736234537185},
			Expected:[]uint32{9, 6830, 5444},
		},
	}

	for _, tcase := range testcases {
		coords,_ := Make_Point(tcase.Row,tcase.OldRow,tcase.Bd)
		for i := range coords {
			valmine := coords[i]
			valtcase := tcase.Expected[i]
			if valmine != valtcase {
				t.Errorf("Make_Line_Geom Error, Expected %s Got %s",valtcase,valmine)
			}
		}
	}
}