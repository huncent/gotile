package tile_surge

import (
	m "github.com/murphy214/mercantile"
	"testing"
)

// test for Make_Tilemap
func Test_Make_Tilelist(t *testing.T) {
	testcases := []struct {
			Extrema m.Extrema
			Size int
			Expected_Tilelist  []m.TileID
	}{	
		{
			Extrema:m.Extrema{W:-176.684744, E:145.830418, N:71.341223, S:-14.37374},
			Size:5,
			Expected_Tilelist:[]m.TileID{m.TileID{X:0,Y:18,Z:5},m.TileID{X:1,Y:18,Z:5},m.TileID{X:2,Y:18,Z:5},m.TileID{X:3,Y:18,Z:5},m.TileID{X:4,Y:18,Z:5},m.TileID{X:5,Y:18,Z:5},m.TileID{X:6,Y:18,Z:5},m.TileID{X:7,Y:18,Z:5},m.TileID{X:8,Y:18,Z:5},m.TileID{X:9,Y:18,Z:5},m.TileID{X:10,Y:18,Z:5},m.TileID{X:11,Y:18,Z:5},m.TileID{X:12,Y:18,Z:5},m.TileID{X:13,Y:18,Z:5},m.TileID{X:14,Y:18,Z:5},m.TileID{X:15,Y:18,Z:5},m.TileID{X:16,Y:18,Z:5},m.TileID{X:17,Y:18,Z:5},m.TileID{X:18,Y:18,Z:5},m.TileID{X:19,Y:18,Z:5},m.TileID{X:20,Y:18,Z:5},m.TileID{X:21,Y:18,Z:5},m.TileID{X:22,Y:18,Z:5},m.TileID{X:23,Y:18,Z:5},m.TileID{X:24,Y:18,Z:5},m.TileID{X:25,Y:18,Z:5},m.TileID{X:26,Y:18,Z:5},m.TileID{X:27,Y:18,Z:5},m.TileID{X:28,Y:18,Z:5},m.TileID{X:0,Y:17,Z:5},m.TileID{X:1,Y:17,Z:5},m.TileID{X:2,Y:17,Z:5},m.TileID{X:3,Y:17,Z:5},m.TileID{X:4,Y:17,Z:5},m.TileID{X:5,Y:17,Z:5},m.TileID{X:6,Y:17,Z:5},m.TileID{X:7,Y:17,Z:5},m.TileID{X:8,Y:17,Z:5},m.TileID{X:9,Y:17,Z:5},m.TileID{X:10,Y:17,Z:5},m.TileID{X:11,Y:17,Z:5},m.TileID{X:12,Y:17,Z:5},m.TileID{X:13,Y:17,Z:5},m.TileID{X:14,Y:17,Z:5},m.TileID{X:15,Y:17,Z:5},m.TileID{X:16,Y:17,Z:5},m.TileID{X:17,Y:17,Z:5},m.TileID{X:18,Y:17,Z:5},m.TileID{X:19,Y:17,Z:5},m.TileID{X:20,Y:17,Z:5},m.TileID{X:21,Y:17,Z:5},m.TileID{X:22,Y:17,Z:5},m.TileID{X:23,Y:17,Z:5},m.TileID{X:24,Y:17,Z:5},m.TileID{X:25,Y:17,Z:5},m.TileID{X:26,Y:17,Z:5},m.TileID{X:27,Y:17,Z:5},m.TileID{X:28,Y:17,Z:5},m.TileID{X:0,Y:16,Z:5},m.TileID{X:1,Y:16,Z:5},m.TileID{X:2,Y:16,Z:5},m.TileID{X:3,Y:16,Z:5},m.TileID{X:4,Y:16,Z:5},m.TileID{X:5,Y:16,Z:5},m.TileID{X:6,Y:16,Z:5},m.TileID{X:7,Y:16,Z:5},m.TileID{X:8,Y:16,Z:5},m.TileID{X:9,Y:16,Z:5},m.TileID{X:10,Y:16,Z:5},m.TileID{X:11,Y:16,Z:5},m.TileID{X:12,Y:16,Z:5},m.TileID{X:13,Y:16,Z:5},m.TileID{X:14,Y:16,Z:5},m.TileID{X:15,Y:16,Z:5},m.TileID{X:16,Y:16,Z:5},m.TileID{X:17,Y:16,Z:5},m.TileID{X:18,Y:16,Z:5},m.TileID{X:19,Y:16,Z:5},m.TileID{X:20,Y:16,Z:5},m.TileID{X:21,Y:16,Z:5},m.TileID{X:22,Y:16,Z:5},m.TileID{X:23,Y:16,Z:5},m.TileID{X:24,Y:16,Z:5},m.TileID{X:25,Y:16,Z:5},m.TileID{X:26,Y:16,Z:5},m.TileID{X:27,Y:16,Z:5},m.TileID{X:28,Y:16,Z:5},m.TileID{X:0,Y:15,Z:5},m.TileID{X:1,Y:15,Z:5},m.TileID{X:2,Y:15,Z:5},m.TileID{X:3,Y:15,Z:5},m.TileID{X:4,Y:15,Z:5},m.TileID{X:5,Y:15,Z:5},m.TileID{X:6,Y:15,Z:5},m.TileID{X:7,Y:15,Z:5},m.TileID{X:8,Y:15,Z:5},m.TileID{X:9,Y:15,Z:5},m.TileID{X:10,Y:15,Z:5},m.TileID{X:11,Y:15,Z:5},m.TileID{X:12,Y:15,Z:5},m.TileID{X:13,Y:15,Z:5},m.TileID{X:14,Y:15,Z:5},m.TileID{X:15,Y:15,Z:5},m.TileID{X:16,Y:15,Z:5},m.TileID{X:17,Y:15,Z:5},m.TileID{X:18,Y:15,Z:5},m.TileID{X:19,Y:15,Z:5},m.TileID{X:20,Y:15,Z:5},m.TileID{X:21,Y:15,Z:5},m.TileID{X:22,Y:15,Z:5},m.TileID{X:23,Y:15,Z:5},m.TileID{X:24,Y:15,Z:5},m.TileID{X:25,Y:15,Z:5},m.TileID{X:26,Y:15,Z:5},m.TileID{X:27,Y:15,Z:5},m.TileID{X:28,Y:15,Z:5},m.TileID{X:0,Y:14,Z:5},m.TileID{X:1,Y:14,Z:5},m.TileID{X:2,Y:14,Z:5},m.TileID{X:3,Y:14,Z:5},m.TileID{X:4,Y:14,Z:5},m.TileID{X:5,Y:14,Z:5},m.TileID{X:6,Y:14,Z:5},m.TileID{X:7,Y:14,Z:5},m.TileID{X:8,Y:14,Z:5},m.TileID{X:9,Y:14,Z:5},m.TileID{X:10,Y:14,Z:5},m.TileID{X:11,Y:14,Z:5},m.TileID{X:12,Y:14,Z:5},m.TileID{X:13,Y:14,Z:5},m.TileID{X:14,Y:14,Z:5},m.TileID{X:15,Y:14,Z:5},m.TileID{X:16,Y:14,Z:5},m.TileID{X:17,Y:14,Z:5},m.TileID{X:18,Y:14,Z:5},m.TileID{X:19,Y:14,Z:5},m.TileID{X:20,Y:14,Z:5},m.TileID{X:21,Y:14,Z:5},m.TileID{X:22,Y:14,Z:5},m.TileID{X:23,Y:14,Z:5},m.TileID{X:24,Y:14,Z:5},m.TileID{X:25,Y:14,Z:5},m.TileID{X:26,Y:14,Z:5},m.TileID{X:27,Y:14,Z:5},m.TileID{X:28,Y:14,Z:5},m.TileID{X:0,Y:13,Z:5},m.TileID{X:1,Y:13,Z:5},m.TileID{X:2,Y:13,Z:5},m.TileID{X:3,Y:13,Z:5},m.TileID{X:4,Y:13,Z:5},m.TileID{X:5,Y:13,Z:5},m.TileID{X:6,Y:13,Z:5},m.TileID{X:7,Y:13,Z:5},m.TileID{X:8,Y:13,Z:5},m.TileID{X:9,Y:13,Z:5},m.TileID{X:10,Y:13,Z:5},m.TileID{X:11,Y:13,Z:5},m.TileID{X:12,Y:13,Z:5},m.TileID{X:13,Y:13,Z:5},m.TileID{X:14,Y:13,Z:5},m.TileID{X:15,Y:13,Z:5},m.TileID{X:16,Y:13,Z:5},m.TileID{X:17,Y:13,Z:5},m.TileID{X:18,Y:13,Z:5},m.TileID{X:19,Y:13,Z:5},m.TileID{X:20,Y:13,Z:5},m.TileID{X:21,Y:13,Z:5},m.TileID{X:22,Y:13,Z:5},m.TileID{X:23,Y:13,Z:5},m.TileID{X:24,Y:13,Z:5},m.TileID{X:25,Y:13,Z:5},m.TileID{X:26,Y:13,Z:5},m.TileID{X:27,Y:13,Z:5},m.TileID{X:28,Y:13,Z:5},m.TileID{X:0,Y:12,Z:5},m.TileID{X:1,Y:12,Z:5},m.TileID{X:2,Y:12,Z:5},m.TileID{X:3,Y:12,Z:5},m.TileID{X:4,Y:12,Z:5},m.TileID{X:5,Y:12,Z:5},m.TileID{X:6,Y:12,Z:5},m.TileID{X:7,Y:12,Z:5},m.TileID{X:8,Y:12,Z:5},m.TileID{X:9,Y:12,Z:5},m.TileID{X:10,Y:12,Z:5},m.TileID{X:11,Y:12,Z:5},m.TileID{X:12,Y:12,Z:5},m.TileID{X:13,Y:12,Z:5},m.TileID{X:14,Y:12,Z:5},m.TileID{X:15,Y:12,Z:5},m.TileID{X:16,Y:12,Z:5},m.TileID{X:17,Y:12,Z:5},m.TileID{X:18,Y:12,Z:5},m.TileID{X:19,Y:12,Z:5},m.TileID{X:20,Y:12,Z:5},m.TileID{X:21,Y:12,Z:5},m.TileID{X:22,Y:12,Z:5},m.TileID{X:23,Y:12,Z:5},m.TileID{X:24,Y:12,Z:5},m.TileID{X:25,Y:12,Z:5},m.TileID{X:26,Y:12,Z:5},m.TileID{X:27,Y:12,Z:5},m.TileID{X:28,Y:12,Z:5},m.TileID{X:0,Y:11,Z:5},m.TileID{X:1,Y:11,Z:5},m.TileID{X:2,Y:11,Z:5},m.TileID{X:3,Y:11,Z:5},m.TileID{X:4,Y:11,Z:5},m.TileID{X:5,Y:11,Z:5},m.TileID{X:6,Y:11,Z:5},m.TileID{X:7,Y:11,Z:5},m.TileID{X:8,Y:11,Z:5},m.TileID{X:9,Y:11,Z:5},m.TileID{X:10,Y:11,Z:5},m.TileID{X:11,Y:11,Z:5},m.TileID{X:12,Y:11,Z:5},m.TileID{X:13,Y:11,Z:5},m.TileID{X:14,Y:11,Z:5},m.TileID{X:15,Y:11,Z:5},m.TileID{X:16,Y:11,Z:5},m.TileID{X:17,Y:11,Z:5},m.TileID{X:18,Y:11,Z:5},m.TileID{X:19,Y:11,Z:5},m.TileID{X:20,Y:11,Z:5},m.TileID{X:21,Y:11,Z:5},m.TileID{X:22,Y:11,Z:5},m.TileID{X:23,Y:11,Z:5},m.TileID{X:24,Y:11,Z:5},m.TileID{X:25,Y:11,Z:5},m.TileID{X:26,Y:11,Z:5},m.TileID{X:27,Y:11,Z:5},m.TileID{X:28,Y:11,Z:5},m.TileID{X:0,Y:10,Z:5},m.TileID{X:1,Y:10,Z:5},m.TileID{X:2,Y:10,Z:5},m.TileID{X:3,Y:10,Z:5},m.TileID{X:4,Y:10,Z:5},m.TileID{X:5,Y:10,Z:5},m.TileID{X:6,Y:10,Z:5},m.TileID{X:7,Y:10,Z:5},m.TileID{X:8,Y:10,Z:5},m.TileID{X:9,Y:10,Z:5},m.TileID{X:10,Y:10,Z:5},m.TileID{X:11,Y:10,Z:5},m.TileID{X:12,Y:10,Z:5},m.TileID{X:13,Y:10,Z:5},m.TileID{X:14,Y:10,Z:5},m.TileID{X:15,Y:10,Z:5},m.TileID{X:16,Y:10,Z:5},m.TileID{X:17,Y:10,Z:5},m.TileID{X:18,Y:10,Z:5},m.TileID{X:19,Y:10,Z:5},m.TileID{X:20,Y:10,Z:5},m.TileID{X:21,Y:10,Z:5},m.TileID{X:22,Y:10,Z:5},m.TileID{X:23,Y:10,Z:5},m.TileID{X:24,Y:10,Z:5},m.TileID{X:25,Y:10,Z:5},m.TileID{X:26,Y:10,Z:5},m.TileID{X:27,Y:10,Z:5},m.TileID{X:28,Y:10,Z:5},m.TileID{X:0,Y:8,Z:5},m.TileID{X:1,Y:8,Z:5},m.TileID{X:2,Y:8,Z:5},m.TileID{X:3,Y:8,Z:5},m.TileID{X:4,Y:8,Z:5},m.TileID{X:5,Y:8,Z:5},m.TileID{X:6,Y:8,Z:5},m.TileID{X:7,Y:8,Z:5},m.TileID{X:8,Y:8,Z:5},m.TileID{X:9,Y:8,Z:5},m.TileID{X:10,Y:8,Z:5},m.TileID{X:11,Y:8,Z:5},m.TileID{X:12,Y:8,Z:5},m.TileID{X:13,Y:8,Z:5},m.TileID{X:14,Y:8,Z:5},m.TileID{X:15,Y:8,Z:5},m.TileID{X:16,Y:8,Z:5},m.TileID{X:17,Y:8,Z:5},m.TileID{X:18,Y:8,Z:5},m.TileID{X:19,Y:8,Z:5},m.TileID{X:20,Y:8,Z:5},m.TileID{X:21,Y:8,Z:5},m.TileID{X:22,Y:8,Z:5},m.TileID{X:23,Y:8,Z:5},m.TileID{X:24,Y:8,Z:5},m.TileID{X:25,Y:8,Z:5},m.TileID{X:26,Y:8,Z:5},m.TileID{X:27,Y:8,Z:5},m.TileID{X:28,Y:8,Z:5},m.TileID{X:0,Y:6,Z:5},m.TileID{X:1,Y:6,Z:5},m.TileID{X:2,Y:6,Z:5},m.TileID{X:3,Y:6,Z:5},m.TileID{X:4,Y:6,Z:5},m.TileID{X:5,Y:6,Z:5},m.TileID{X:6,Y:6,Z:5},m.TileID{X:7,Y:6,Z:5},m.TileID{X:8,Y:6,Z:5},m.TileID{X:9,Y:6,Z:5},m.TileID{X:10,Y:6,Z:5},m.TileID{X:11,Y:6,Z:5},m.TileID{X:12,Y:6,Z:5},m.TileID{X:13,Y:6,Z:5},m.TileID{X:14,Y:6,Z:5},m.TileID{X:15,Y:6,Z:5},m.TileID{X:16,Y:6,Z:5},m.TileID{X:17,Y:6,Z:5},m.TileID{X:18,Y:6,Z:5},m.TileID{X:19,Y:6,Z:5},m.TileID{X:20,Y:6,Z:5},m.TileID{X:21,Y:6,Z:5},m.TileID{X:22,Y:6,Z:5},m.TileID{X:23,Y:6,Z:5},m.TileID{X:24,Y:6,Z:5},m.TileID{X:25,Y:6,Z:5},m.TileID{X:26,Y:6,Z:5},m.TileID{X:27,Y:6,Z:5},m.TileID{X:28,Y:6,Z:5}},
		},
	}

	for _, tcase := range testcases {
		tilelist := Make_Tilelist(tcase.Extrema,tcase.Size)
		if len(tcase.Expected_Tilelist) != len(tilelist) {
			t.Errorf("Make_Tilelist Error Expected tilemap size %d Got %d",len(tcase.Expected_Tilelist),len(tilelist))
		}

	}
	
}


// test for Make_Tilemap
func Test_Add_BBox(t *testing.T) {
	testcases := []struct {
			Tablename string
			Tileid m.TileID
			Expected_String  string
	}{	
		{
			Tablename:"mine",
			Tileid:m.TileID{X:0,Y:18,Z:5},
			Expected_String:"(mine.geom && ST_MakeEnvelope(-180.000000, -31.952162, -168.750000, -21.943046, 4326))",
		},
	}

	for _, tcase := range testcases {
		stringval := Add_BBox(tcase.Tablename,tcase.Tileid)
		if tcase.Expected_String != stringval {
			t.Errorf("Add_BBox Error Expected tilemap size %s Got %s",tcase.Expected_String,stringval)
		}

	}
	
}

// test for Make_Tilemap
func Test_Lint_Extrema(t *testing.T) {
	testcases := []struct {
			Extrema m.Extrema
			Size int
			Expected_Extrema m.Extrema
	}{	
		{
			Extrema:m.Extrema{W:-176.684744, E:145.830418, N:71.341223, S:-14.37374},
			Size:5,
			Expected_Extrema:m.Extrema{W:-180, E:146.25, N:74.01954331150226, S:-21.943045533438177},
		},
	}

	for _, tcase := range testcases {
		lint := Lint_Extrema(tcase.Extrema,tcase.Size)
		if (tcase.Expected_Extrema.N != lint.N) || (tcase.Expected_Extrema.S != lint.S) || (tcase.Expected_Extrema.E != lint.E) || (tcase.Expected_Extrema.W != lint.W)   {
			t.Errorf("Lint_Extrema Error Expected extrema %+v Got %+v",tcase.Expected_Extrema,lint)
		}
	}
	
}