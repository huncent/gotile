package tile_surge

import (
	m "github.com/murphy214/mercantile"
	"github.com/paulmach/go.geojson"

	"testing"

)

// testing the envelope lines function
func Test_Env_Polygon(t *testing.T) {
	testcases := []struct {
			Feature *geojson.Feature
			Zoom int
			Expected map[m.TileID][]*geojson.Feature
	}{	
		{
			Feature:&geojson.Feature{Geometry:&geojson.Geometry{Polygon:[][][]float64{{{-79.438103,38.216235},{-79.308692,38.382706999999996},{-79.3112960655144,38.418454601455},{-79.297758,38.416438},{-79.23161999999999,38.474041},{-79.22825588138188,38.4800395510852},{-78.91249599999999,38.303486},{-78.888025,38.303596999999996},{-78.749354,38.206621},{-78.781078,38.080757},{-78.83921099999999,38.047565},{-78.90727799999999,37.945958999999995},{-79.005129,37.88169},{-79.062454,37.9176},{-79.157423,37.890995},{-79.183978,37.914193999999995},{-79.482405,38.086104999999996},{-79.436678,38.162800000000004},{-79.512158,38.180419},{-79.438103,38.216235}},{{-79.11339699999999,38.154047},{-79.049779,38.121112},{-79.014845,38.157129999999995},{-79.02301,38.195777},{-79.095818,38.185691999999996},{-79.11339699999999,38.154047}},{{-78.921908,38.031569},{-78.86479299999999,38.04882},{-78.864708,38.095693},{-78.880213,38.094711},{-78.950277,38.069486},{-78.921908,38.031569}}},Type:"Polygon"},Properties:map[string]interface{}{"AREA":"51015","area":"51015","index":2605,"colorkey":"#BDDE07"}},
			Zoom:10,
			Expected:map[m.TileID][]*geojson.Feature{{286,395,10}:[]*geojson.Feature{&geojson.Feature{Geometry:&geojson.Geometry{Polygon:[][][]float64{{{-79.1015625,37.996162679728116},{-79.1015625,37.90664398653771},{-79.157423,37.890995},{-79.183978,37.914193999999995},{-79.32627062342272,37.996162679728116}}},Type:"Polygon"},Properties:map[string]interface{}{"index":2605,"colorkey":"#BDDE07","AREA":"51015","area":"51015"}}},{288,394,10}:[]*geojson.Feature{&geojson.Feature{Geometry:&geojson.Geometry{Polygon:[][][]float64{{{-78.75,38.20707276349778},{-78.75,38.20405801475222},{-78.749354,38.206621}}},Type:"Polygon"},Properties:map[string]interface{}{"AREA":"51015","area":"51015","index":2605,"colorkey":"#BDDE07"}}},{287,393,10}:[]*geojson.Feature{&geojson.Feature{Geometry:&geojson.Geometry{Polygon:[][][]float64{{{-78.888025,38.303596999999996},{-78.91249599999999,38.303486},{-79.1015625,38.40920038594468},{-79.1015625,38.272688535980954},{-78.84382738807555,38.272688535980954}}},Type:"Polygon"},Properties:map[string]interface{}{"area":"51015","index":2605,"colorkey":"#BDDE07","AREA":"51015"}}},{285,394,10}:[]*geojson.Feature{&geojson.Feature{Geometry:&geojson.Geometry{Polygon:[][][]float64{{{-79.453125,38.135214488923395},{-79.453125,38.06923804734491},{-79.482405,38.086104999999996}}},Type:"Polygon"},Properties:map[string]interface{}{"AREA":"51015","area":"51015","index":2605,"colorkey":"#BDDE07"}},&geojson.Feature{Geometry:&geojson.Geometry{Polygon:[][][]float64{{{-79.453125,38.208969751846595},{-79.453125,38.166639158624804},{-79.512158,38.180419}}},Type:"Polygon"},Properties:map[string]interface{}{"area":"51015","index":2605,"colorkey":"#BDDE07","AREA":"51015"}}},{287,394,10}:[]*geojson.Feature{&geojson.Feature{Geometry:&geojson.Geometry{Polygon:[][][]float64{{{-78.84382738807555,38.272688535980954},{-79.1015625,38.272688535980954},{-79.1015625,38.17535098501051},{-79.095818,38.185691999999996},{-79.02301,38.195777},{-79.014845,38.157129999999995},{-79.049779,38.121112},{-79.1015625,38.14792028653054},{-79.1015625,37.996162679728116},{-78.87364599137793,37.996162679728116},{-78.83921099999999,38.047565},{-78.781078,38.080757},{-78.75,38.20405801475222},{-78.75,38.20707276349778}},{{-78.864708,38.095693},{-78.86479299999999,38.04882},{-78.921908,38.031569},{-78.950277,38.069486},{-78.880213,38.094711}}},Type:"Polygon"},Properties:map[string]interface{}{"AREA":"51015","area":"51015","index":2605,"colorkey":"#BDDE07"}}},{287,395,10}:[]*geojson.Feature{&geojson.Feature{Geometry:&geojson.Geometry{Polygon:[][][]float64{{{-79.1015625,37.996162679728116},{-79.1015625,37.90664398653771},{-79.062454,37.9176},{-79.005129,37.88169},{-78.90727799999999,37.945958999999995},{-78.87364599137793,37.996162679728116}}},Type:"Polygon"},Properties:map[string]interface{}{"AREA":"51015","area":"51015","index":2605,"colorkey":"#BDDE07"}}},{286,393,10}:[]*geojson.Feature{&geojson.Feature{Geometry:&geojson.Geometry{Polygon:[][][]float64{{{-79.22825588138188,38.4800395510852},{-79.23161999999999,38.474041},{-79.297758,38.416438},{-79.3112960655144,38.418454601455},{-79.308692,38.382706999999996},{-79.39421749045586,38.272688535980954},{-79.1015625,38.272688535980954},{-79.1015625,38.40920038594468}}},Type:"Polygon"},Properties:map[string]interface{}{"area":"51015","index":2605,"colorkey":"#BDDE07","AREA":"51015"}}},{286,394,10}:[]*geojson.Feature{&geojson.Feature{Geometry:&geojson.Geometry{Polygon:[][][]float64{{{-79.39421749045586,38.272688535980954},{-79.438103,38.216235},{-79.453125,38.208969751846595},{-79.453125,38.166639158624804},{-79.436678,38.162800000000004},{-79.453125,38.135214488923395},{-79.453125,38.06923804734491},{-79.32627062342272,37.996162679728116},{-79.1015625,37.996162679728116},{-79.1015625,38.14792028653054},{-79.11339699999999,38.154047},{-79.1015625,38.17535098501051},{-79.1015625,38.272688535980954}}},Type:"Polygon"},Properties:map[string]interface{}{"AREA":"51015","area":"51015","index":2605,"colorkey":"#BDDE07"}}}},
		},
	}

	for _, tcase := range testcases {
		tilemap := Env_Polygon(tcase.Feature,tcase.Zoom)
		for k := range tcase.Expected {
			featsexpected := tcase.Expected[k]
			vals,boolval := tilemap[k]
			if boolval == false {
				t.Errorf("Env_Polygon Error, Key not found %s",k)				
			} else {
				// getting the total number of expected values
				totalexpected := 0
				for _,i := range featsexpected {
					totalexpected += len(i.Geometry.Polygon)
				}

				// getting the total number of found values
				totalfound := 0
				for _,i := range vals {
					totalfound += len(i.Geometry.Polygon)
				}
				if totalfound != totalexpected {
					t.Errorf("Env_Polygon Error, Geometries not same number of points %s %s",totalfound,totalexpected)				

				}				
			}
		}
	}
}


// testing the called children polygon function
func Test_Children_Polygon(t *testing.T) {
	testcases := []struct {
			Feature *geojson.Feature
			Tileid m.TileID
			Expected map[m.TileID][]*geojson.Feature
	}{	
		{
			Feature:&geojson.Feature{Geometry:&geojson.Geometry{Polygon:[][][]float64{{{-79.453125,38.135214488923395},{-79.453125,38.06923804734491},{-79.482405,38.086104999999996}}},Type:"Polygon"},Properties:map[string]interface{}{"AREA":"51015","area":"51015","index":2605,"colorkey":"#BDDE07"}},
			Tileid:m.TileID{X:285,Y:394,Z:10},
			Expected:map[m.TileID][]*geojson.Feature{{571,788,11}:[]*geojson.Feature{&geojson.Feature{Geometry:&geojson.Geometry{Polygon:[][][]float64{{{-79.453125,38.135214488923395},{-79.453125,38.13455657705412},{-79.45351725941778,38.13455657705412}}},Type:"Polygon"},Properties:map[string]interface{}{"AREA":"51015","area":"51015","index":2605,"colorkey":"#BDDE07"}}},{571,789,11}:[]*geojson.Feature{&geojson.Feature{Geometry:&geojson.Geometry{Polygon:[][][]float64{{{-79.45351725941778,38.13455657705412},{-79.482405,38.086104999999996},{-79.453125,38.06923804734491},{-79.453125,38.13455657705412}}},Type:"Polygon"},Properties:map[string]interface{}{"AREA":"51015","area":"51015","index":2605,"colorkey":"#BDDE07"}}}},
		},
	}

	for _, tcase := range testcases {
		tilemap := Children_Polygon(tcase.Feature,tcase.Tileid)
		for k := range tcase.Expected {
			featsexpected := tcase.Expected[k]
			vals,boolval := tilemap[k]
			if boolval == false {
				t.Errorf("Children_Polygon Error, Key not found %s",k)				
			} else {
				// getting the total number of expected values
				totalexpected := 0
				for _,i := range featsexpected {
					totalexpected += len(i.Geometry.Polygon)
				}

				// getting the total number of found values
				totalfound := 0
				for _,i := range vals {
					totalfound += len(i.Geometry.Polygon)
				}
				if totalfound != totalexpected {
					t.Errorf("Children_Polygon Error, Geometries not same number of points %s %s",totalfound,totalexpected)				

				}				
			}
		}
	}
}
