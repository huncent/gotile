package tile_surge

import (
	m "github.com/murphy214/mercantile"
	"testing"
	"github.com/paulmach/go.geojson"
	h "github.com/mitchellh/hashstructure"
	"vector-tile/2.1"
)

// hashs a given tv structure
func Hash_Tv(tv *vector_tile.Tile_Value) uint64 {
	hash, _ := h.Hash(tv, nil)
	return hash
}

// test for Make_Tv_String
func Test_Make_Tv_String(t *testing.T) {
	testcases := []struct {
			Valin string
			Expected_Hash uint64
	}{	
		{
			Valin:"shit",
			Expected_Hash:13591876025604954094,
		},
	}

	for _, tcase := range testcases {
		val := Make_Tv_String(tcase.Valin)
		val_hash := Hash_Tv(val)
		if tcase.Expected_Hash != val_hash {
			t.Errorf("Make_Tv_String Error Expected %d Got %d",tcase.Expected_Hash,val_hash)
		}
	}
}

// test for Make_Tv_Float
func Test_Make_Tv_Float(t *testing.T) {
	testcases := []struct {
			Valin float64
			Expected_Hash uint64
	}{	
		{
			Valin:42.32323,
			Expected_Hash:17337609332515500439,
		},
	}

	for _, tcase := range testcases {
		val := Make_Tv_Float(tcase.Valin)
		val_hash := Hash_Tv(val)
		if tcase.Expected_Hash != val_hash {
			t.Errorf("Make_Tv_Float Error Expected %d Got %d",tcase.Expected_Hash,val_hash)
		}
	}
}


// test for Make_Tv_Int
func Test_Make_Tv_Int(t *testing.T) {
	testcases := []struct {
			Valin int
			Expected_Hash uint64
	}{	
		{
			Valin:2324332132,
			Expected_Hash:8061578276763486459,
		},
	}

	for _, tcase := range testcases {
		val := Make_Tv_Int(tcase.Valin)
		val_hash := Hash_Tv(val)
		if tcase.Expected_Hash != val_hash {
			t.Errorf("Make_Tv_Int Error Expected %d Got %d",tcase.Expected_Hash,val_hash)
		}
	}
}

// test for Make_Tv_Int
func Test_Reflect_Value(t *testing.T) {
	testcases := []struct {
			Valin interface{}
			Expected_Hash uint64
	}{	
		{
			Valin:2324332132,
			Expected_Hash:8061578276763486459,
		},
		{
			Valin:42.32323,
			Expected_Hash:17337609332515500439,
		},
		{
			Valin:"shit",
			Expected_Hash:13591876025604954094,
		},
	}

	for _, tcase := range testcases {
		val := Reflect_Value(tcase.Valin)
		val_hash := Hash_Tv(val)
		if tcase.Expected_Hash != val_hash {
			t.Errorf("Reflect_Value Error Expected %d Got %d",tcase.Expected_Hash,val_hash)
		}
	}
}



// test for Make_Zoom_Drill
func Test_Make_Tile(t *testing.T) {
	testcases := []struct {
			K m.TileID
			V []*geojson.Feature
			Prefix string
			Con Config
			Expected_Size int
	}{	
		{
			K:m.TileID{X:524,Y:841,Z:11},
			V:[]*geojson.Feature{&geojson.Feature{Geometry:&geojson.Geometry{Polygon:[][][]float64{{{-87.71484375,30.600093873550065},{-87.890625,30.600093873550065},{-87.890625,30.44867367928756},{-87.71484375,30.44867367928756}}},Type:"Polygon"},Properties:map[string]interface{}{"area":"1003","index":2846,"colorkey":"#DDBD07","AREA":"1003"}}},
			Prefix:"shit",
			Con:Config{Prefix:"shit",Minzoom:5,Maxzoom:13},
			Expected_Size:105,
		},
	}

	for _, tcase := range testcases {
		vts := Make_Tile(tcase.K,tcase.V,tcase.Prefix,tcase.Con)
		sizevts := len(vts.Data)
		if tcase.Expected_Size != sizevts {
			t.Errorf("Make_Zoom_Drill Error Expected tilemap size %d Got %d",tcase.Expected_Size,sizevts)
		}
	}
}