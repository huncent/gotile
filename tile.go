package tile_surge

import (
	//l "github.com/murphy214/layersplit"
	m "github.com/murphy214/mercantile"
	//pc "github.com/murphy214/polyclip"
	"github.com/paulmach/go.geojson"
	"os"
	"strconv"
	"vector-tile/2.1"
	//"strings"
	//"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"reflect"
)

var total = 0
var total2 = 0

func Reflect_Value(v interface{}) *vector_tile.Tile_Value {
	var tv *vector_tile.Tile_Value
	//fmt.Print(v)
	vv := reflect.ValueOf(v)
	kd := vv.Kind()
	if (reflect.Float64 == kd) || (reflect.Float32 == kd) {
		//fmt.Print(v, "float", k)
		tv = Make_Tv_Float(float64(vv.Float()))
		//hash = Hash_Tv(tv)
	} else if (reflect.Int == kd) || (reflect.Int8 == kd) || (reflect.Int16 == kd) || (reflect.Int32 == kd) || (reflect.Int64 == kd) || (reflect.Uint8 == kd) || (reflect.Uint16 == kd) || (reflect.Uint32 == kd) || (reflect.Uint64 == kd) {
		//fmt.Print(v, "int", k)
		tv = Make_Tv_Int(int(vv.Int()))
		//hash = Hash_Tv(tv)
	} else if reflect.String == kd {
		//fmt.Print(v, "str", k)
		tv = Make_Tv_String(string(vv.String()))
		//hash = Hash_Tv(tv)

	} else {
		tv := new(vector_tile.Tile_Value)
		t := ""
		tv.StringValue = &t
	}
	return tv
}

// makes a tile_value string
func Make_Tv_String(stringval string) *vector_tile.Tile_Value {
	tv := new(vector_tile.Tile_Value)
	t := string(stringval)
	tv.StringValue = &t
	return tv
}

// makes a tile value float
func Make_Tv_Float(val float64) *vector_tile.Tile_Value {
	tv := new(vector_tile.Tile_Value)
	t := float64(val)
	tv.DoubleValue = &t
	return tv
}

// makes a tile value int
func Make_Tv_Int(val int) *vector_tile.Tile_Value {
	tv := new(vector_tile.Tile_Value)
	t := int64(val)
	tv.SintValue = &t
	return tv
}

// updates all values and tags
func Update_Properties(properties map[string]interface{}, keys []string, values []*vector_tile.Tile_Value, keysmap map[string]uint32, valuesmap map[*vector_tile.Tile_Value]uint32) ([]uint32, []string, []*vector_tile.Tile_Value, map[string]uint32, map[*vector_tile.Tile_Value]uint32) {
	tags := []uint32{}
	// iterating through each property
	for k, v := range properties {
		value := Reflect_Value(v)

		// logic for keys
		keyint, keybool := keysmap[k]
		// logic for if it isnt in keymap
		if keybool == false {
			keys = append(keys, k)
			keysmap[k] = uint32(len(keys) - 1)
			tags = append(tags, uint32(len(keys)-1))
		} else {
			tags = append(tags, keyint)
		}

		// logic for keys
		valueint, valuebool := valuesmap[value]
		// logic for if it isnt in keymap
		if valuebool == false {
			values = append(values, value)
			valuesmap[value] = uint32(len(values) - 1)
			tags = append(tags, uint32(len(values)-1))
		} else {
			tags = append(tags, valueint)
		}

	}

	return tags, keys, values, keysmap, valuesmap
}

// makes a single tile for a given polygon
func Make_Tile(tileid m.TileID, feats []*geojson.Feature, prefix string) {
	filename := prefix + "/" + strconv.Itoa(int(tileid.Z)) + "/" + strconv.Itoa(int(tileid.X)) + "/" + strconv.Itoa(int(tileid.Y))
	dir := prefix + "/" + strconv.Itoa(int(tileid.Z)) + "/" + strconv.Itoa(int(tileid.X))
	os.MkdirAll(dir, os.ModePerm)
	bound := m.Bounds(tileid)
	var keys []string
	var values []*vector_tile.Tile_Value
	keysmap := map[string]uint32{}
	valuesmap := map[*vector_tile.Tile_Value]uint32{}

	// iterating through each feature
	features := []*vector_tile.Tile_Feature{}
	//position := []int32{0, 0}
	for _, i := range feats {
		var tags, geometry []uint32
		var feat vector_tile.Tile_Feature
		tags, keys, values, keysmap, valuesmap = Update_Properties(i.Properties, keys, values, keysmap, valuesmap)

		// logic for point feature
		if i.Geometry.Type == "Point" {
			geometry, _ = Make_Point(i.Geometry.Point, []int32{0, 0}, bound)
			feat_type := vector_tile.Tile_POINT
			feat = vector_tile.Tile_Feature{Tags: tags, Type: &feat_type, Geometry: geometry}
			features = append(features, &feat)

		} else if i.Geometry.Type == "LineString" {
			eh := Make_Coords_Float(i.Geometry.LineString, bound, tileid)
			if len(eh) > 0 {
				geometry, _ = Make_Line_Geom(eh, []int32{0, 0})
				feat_type := vector_tile.Tile_LINESTRING
				feat = vector_tile.Tile_Feature{Tags: tags, Type: &feat_type, Geometry: geometry}
				features = append(features, &feat)
			}

		} else if i.Geometry.Type == "Polygon" {
			geometry, _ = Make_Polygon(Make_Coords_Polygon_Float(i.Geometry.Polygon, bound), []int32{0, 0})
			feat_type := vector_tile.Tile_POLYGON
			feat = vector_tile.Tile_Feature{Tags: tags, Type: &feat_type, Geometry: geometry}
			features = append(features, &feat)

		}

	}

	layerVersion := uint32(15)
	extent := vector_tile.Default_Tile_Layer_Extent
	//var bound []Bounds
	layername := prefix
	layer := vector_tile.Tile_Layer{
		Version:  &layerVersion,
		Name:     &layername,
		Extent:   &extent,
		Values:   values,
		Keys:     keys,
		Features: features,
	}

	tile := vector_tile.Tile{}
	tile.Layers = append(tile.Layers, &layer)

	bytevals, _ := proto.Marshal(&tile)
	ioutil.WriteFile(filename, bytevals, 0666)

	//fmt.Printf("\r%s", filename)
}
