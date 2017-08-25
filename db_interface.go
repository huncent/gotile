package tile_surge

import (
	//"github.com/golang/protobuf/proto"
	"fmt"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
	m "github.com/murphy214/mercantile"
	"github.com/paulmach/go.geojson"
	"reflect"
	"strconv"
	"strings"
)

func Add_BBox(tablename string, tileid m.TileID) string {
	bds := m.Bounds(tileid)

	return fmt.Sprintf("(%s.geom && ST_MakeEnvelope(%f, %f, %f, %f, 4326))", tablename, bds.S, bds.W, bds.N, bds.E)

}

func get_type(text string) string {
	val := strings.Split(text, "(")[0]
	val = strings.Replace(val, " ", "", -1)
	return val
}

func get_polygon(polystring string) [][][]float64 {
	polystring = polystring[7:]
	polyvals := strings.Split(polystring, "),")
	coords := [][][]float64{}

	for _, ring := range polyvals {
		ring = strings.Replace(ring, "(", "", -1)
		ring = strings.Replace(ring, ")", "", -1)
		ringfloat := [][]float64{}
		is := strings.Split(ring, ",")
		for _, i := range is {
			vals := strings.Split(i, " ")
			x, _ := strconv.ParseFloat(vals[0], 64)
			y, _ := strconv.ParseFloat(vals[1], 64)
			ringfloat = append(ringfloat, []float64{x, y})
		}

		coords = append(coords, ringfloat)

	}
	//fmt.Print(coords, polystring, "\n")

	return coords
	//fmt.Print(polystring[7:], "\n")
}

func get_linestring(polystring string) [][]float64 {
	ring := polystring[10:]
	//fmt.Print(ring, "\n")
	//ring := strings.Split(polystring, "),")
	//coords := [][]float64{}

	ring = strings.Replace(ring, "(", "", -1)
	ring = strings.Replace(ring, ")", "", -1)
	ringfloat := [][]float64{}
	is := strings.Split(ring, ",")
	for _, i := range is {
		vals := strings.Split(i, " ")
		x, _ := strconv.ParseFloat(vals[0], 64)
		y, _ := strconv.ParseFloat(vals[1], 64)
		ringfloat = append(ringfloat, []float64{x, y})
	}

	return ringfloat
	//fmt.Print(polystring[7:], "\n")
}

func get_point(polystring string) []float64 {
	ring := polystring[6:]
	//fmt.Print(ring, "\n")
	//ring := strings.Split(polystring, "),")
	//coords := [][]float64{}
	ring = ring[1 : len(ring)-1]
	vals := strings.Split(ring, " ")
	x, _ := strconv.ParseFloat(vals[0], 64)
	y, _ := strconv.ParseFloat(vals[1], 64)
	return []float64{x, y}

	//fmt.Print(polystring[7:], "\n")
}

func DB_Interface(database string, query string) *geojson.FeatureCollection {

	a := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Port:     5432,
			Database: database,
			User:     "postgres",
		},
		MaxConnections: 1,
	}
	p, _ := pgx.NewConnPool(a)

	rows, _ := p.Query(query)
	var keys []string
	fdescs := rows.FieldDescriptions()
	for _, i := range fdescs {
		keys = append(keys, i.Name)
	}
	pos := len(keys) - 1
	//fmt.Print(keys, rows, "\n")
	//var ppp *geo.Line
	featcollection := &geojson.FeatureCollection{}
	for rows.Next() {
		vals, _ := rows.Values()
		tempmap := map[string]interface{}{}
		for ii, i := range vals[:pos] {
			tempmap[keys[ii]] = i
		}
		//fmt.Print(vals[pos], "\n")
		text := reflect.ValueOf(vals[pos]).String()

		geomtype := get_type(text)

		// getting the right geometry from string
		if geomtype == "POLYGON" {
			geom := get_polygon(text)
			geomnew := geojson.Geometry{Polygon: geom, Type: "Polygon"}
			feature := geojson.Feature{Geometry: &geomnew, Properties: tempmap}
			featcollection.Features = append(featcollection.Features, &feature)

		} else if geomtype == "LINESTRING" {
			geom := get_linestring(text)
			geomnew := geojson.Geometry{LineString: geom, Type: "LineString"}
			feature := geojson.Feature{Geometry: &geomnew, Properties: tempmap}
			featcollection.Features = append(featcollection.Features, &feature)

		} else if geomtype == "POINT" {
			geom := get_point(text)
			geomnew := geojson.Geometry{Point: geom, Type: "Point"}
			feature := geojson.Feature{Geometry: &geomnew, Properties: tempmap}
			featcollection.Features = append(featcollection.Features, &feature)

		}

		//fmt.Print(get_type(text), text, "\n")

		//feature.Properties = tempmap

		//featcollection.Features = append(featcollection.Features, feature)

		//fmt.Print(ppp.Scan(vals[pos]), "\n")
		//fmt.Print(vals[:pos], , "\n")
		//totalvals = append(totalvals, vals)
	}

	return featcollection
}
