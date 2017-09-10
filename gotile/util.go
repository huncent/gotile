package tile_surge 

import (
	"fmt"
	"github.com/paulmach/go.geojson"
	"strings"
	"reflect"
	m "github.com/murphy214/mercantile"
)

// tilemap feature
func Feature_String(a *geojson.Feature) string {
	var geom string
	if a.Geometry.Type == "LineString" {
		eh := fmt.Sprintf("%v",a.Geometry.LineString)
		eh = strings.Replace(eh,"[","{",1000000000)
		eh = strings.Replace(eh,"]","}",1000000000)
		eh = strings.Replace(eh," ",",",1000000000)
		geom = fmt.Sprintf(`&geojson.Geometry{LineString:[][]float64%s,Type:"Linestring"}`,eh)
		
	} else if a.Geometry.Type == "Polygon" {
		eh := fmt.Sprintf("%v",a.Geometry.Polygon)
		eh = strings.Replace(eh,"[","{",1000000000)
		eh = strings.Replace(eh,"]","}",1000000000)
		eh = strings.Replace(eh," ",",",1000000000)
		geom = fmt.Sprintf(`&geojson.Geometry{Polygon:[][][]float64%s,Type:"Polygon"}`,eh)
		
	}
	newmap := map[string]interface{}{}
	for k,v := range a.Properties {
		vv := reflect.ValueOf(v)
		kd := vv.Kind()
		if kd == reflect.String {
			stringval := `"` + v.(string) + `"`
			v = reflect.ValueOf(stringval)
			//v = stringval.(interface{}) 
		}

		newmap[`"`+k+`"`] = v
	}

	shit := fmt.Sprintf("%v",newmap)
	shit = shit[4:len(shit)-1]
	shit = "{" + shit + "}"
	shit = strings.Replace(shit," ",`,`,1000000000)
	shit = fmt.Sprintf("map[string]interface{}%s",shit)
	//fmt.Print(shit,"\n")
	total := fmt.Sprintf("&geojson.Feature{Geometry:%s,Properties:%s}",geom,shit)
	//fmt.Print(total,"\n")
	return total
}

// takes a raw tilemap and returns a string that can be 
// palced within a line of code or test
func Tilemap_String(tilemap map[m.TileID][]*geojson.Feature) string {

	newmap := map[m.TileID]string{}
	for k,v := range tilemap {
		stringlist := []string{}
		for _,i := range v {
			stringlist = append(stringlist,Feature_String(i))
		}
		shit :=fmt.Sprintf("[]*geojson.Feature{%s}",strings.Join(stringlist,","))
		newmap[k] = shit
	}


	shit := fmt.Sprintf("%v",newmap)
	shit = shit[4:len(shit)-1]
	shit = "{" + shit + "}"
	shit = strings.Replace(shit," ",`,`,1000000000)
	shit = "map[m.TileID][]*geojson.Feature" + shit

	return shit
	//shit = fmt.Sprintf("map[m.TileID][]{}%s",shit)

}

// makes a feature string list
func FeatureStringList(v []*geojson.Feature) string {

	stringlist := []string{}
	for _,i := range v {
		stringlist = append(stringlist,Feature_String(i))
	}
	shit :=fmt.Sprintf("[]*geojson.Feature{%s}",strings.Join(stringlist,","))
	return shit

}

// creates a vector tile string
func VectorTileString(vt Vector_Tile) string {
	eh := fmt.Sprintf("%+v",vt.Data)
	eh = strings.Replace(eh," ",",",10000000000)
	eh = "Data:[]byte{" + eh[1:len(eh)-1] + "}"
	eh2 := fmt.Sprintf(`Filename:"%s"`,vt.Filename)

	eh3 := fmt.Sprintf("Tileid:m.TileID%+v",vt.Tileid)
	eh3 = strings.Replace(eh3," ",",",10000000)
	stringvals := []string{eh,eh2,eh3}

	totalstring := strings.Join(stringvals,",")
	totalstring = fmt.Sprintf("Vector_Tile{%s}",totalstring)
	return totalstring
}

// creates a vector tile string list
func VectorTileListString(vts []Vector_Tile) string {
	newlist := []string{}
	for _,i := range vts {
		newlist = append(newlist,VectorTileString(i))
	}

	totalstring := strings.Join(newlist,",")
	totalstring = fmt.Sprintf("[]Vector_Tile{%s}",totalstring)
	return totalstring
}
