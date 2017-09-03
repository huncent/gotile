package tile_surge

import (
	//"github.com/golang/protobuf/proto"
	"fmt"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
	m "github.com/murphy214/mercantile"
	pc "github.com/murphy214/polyclip"
	"github.com/paulmach/go.geojson"
	"reflect"
	"strconv"
	"strings"
	//"sync"

)

// returns the bbox logic from a table name and tileid
func Add_BBox(tablename string, tileid m.TileID) string {
	bds := m.Bounds(tileid)

	return fmt.Sprintf("(%s.geom && ST_MakeEnvelope(%f, %f, %f, %f, 4326))", tablename, bds.W, bds.S, bds.E, bds.N)

}

// gets the geometry type of geometry text
func get_type(text string) string {
	val := strings.Split(text, "(")[0]
	val = strings.Replace(val, " ", "", -1)
	return val
}

// hacky way to get a polygon
func get_polygon(polystring string) [][][]float64 {
	polystring = polystring[7:]
	polyvals := strings.Split(polystring, "),")
	coords := [][][]float64{}

	// iterating through each text ring
	for _, ring := range polyvals {
		ring = strings.Replace(ring, "(", "", -1)
		ring = strings.Replace(ring, ")", "", -1)
		ringfloat := [][]float64{}
		is := strings.Split(ring, ",")

		// iterating through each text point
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

// hacky way to get a linestring
func get_linestring(polystring string) [][]float64 {
	ring := polystring[10:]
	ring = strings.Replace(ring, "(", "", -1)
	ring = strings.Replace(ring, ")", "", -1)
	ringfloat := [][]float64{}

	// creatiing and iterating through each point
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

// hacky way to get a point
func get_point(polystring string) []float64 {
	ring := polystring[6:]
	ring = ring[1 : len(ring)-1]
	vals := strings.Split(ring, " ")
	x, _ := strconv.ParseFloat(vals[0], 64)
	y, _ := strconv.ParseFloat(vals[1], 64)
	return []float64{x, y}
}

// this function allows you interface with a postgis database
// it create a raw feature collection geojson representation
// which would be the same if you were just reading from a geojson
func DB_Interface(database string, query string) *geojson.FeatureCollection {
	// intializing the config
	a := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Port:     5432,
			Database: database,
			User:     "postgres",
		},
		MaxConnections: 1,
	}

	// creating the connection
	p, _ := pgx.NewConnPool(a)

	// executing the query
	rows, _ := p.Query(query)

	// getting keys
	var keys []string
	fdescs := rows.FieldDescriptions()
	for _, i := range fdescs {
		keys = append(keys, i.Name)
	}

	pos := len(keys) - 1
	featcollection := &geojson.FeatureCollection{}
	// iterating through each row of the queried data
	for rows.Next() {
		// creating properties map
		vals, _ := rows.Values()
		tempmap := map[string]interface{}{}
		for ii, i := range vals[:pos] {
			tempmap[keys[ii]] = i
		}

		// getting the geometry text
		text := reflect.ValueOf(vals[pos]).String()

		// getting the geometry type
		geomtype := get_type(text)

		// getting the right geometry from string
		// adding the geojson on to the feature collection
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

	}

	return featcollection
}


// getting the extent of a given database
func Get_Extent(database string,tablename string) m.Extrema {
	sqlquery := fmt.Sprintf("SELECT ST_Extent(geom) as table_extent FROM %s;",tablename)

	// intializing the config
	a := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Port:     5432,
			Database: database,
			User:     "postgres",
		},
		MaxConnections: 1,
	}

	// creating the connection
	p, _ := pgx.NewConnPool(a)

	rows, _ := p.Query(sqlquery)
	var bbox string
	for rows.Next() {
		vals, _ := rows.Values()
		bbox = vals[0].(string)
	}
	bbox = bbox[4:len(bbox) - 1]
	bbox = strings.Replace(bbox,","," ",1)
	vals := strings.Split(bbox," ")
	west,_ := strconv.ParseFloat(vals[0],64)
	south,_ := strconv.ParseFloat(vals[1],64)
	east,_ := strconv.ParseFloat(vals[2],64)
	north,_ := strconv.ParseFloat(vals[3],64)

	return m.Extrema{N:north,S:south,E:east,W:west}
}


// making the tilemap for each tileslice in a given row
func Make_Tilelist(ext m.Extrema,minzoom int) []m.TileID {
	tileid := m.Tile(ext.W,ext.S,minzoom)
	bds := m.Bounds(tileid)
	startpt := []float64{(bds.E+bds.W)/2.0,(bds.N+bds.S)/2.0}
	currenty := startpt[1]
	currentx := startpt[0]
	startx := startpt[0]
	size := pc.Point{bds.E-bds.W,bds.N-bds.S}
	tilelist := []m.TileID{}


	for currenty < ext.N {
		currentx = startx
		for currentx < ext.E {
			tileid := m.Tile(currentx,currenty,minzoom)
			tilelist = append(tilelist,tileid)
			currentx += size.X
		}
		currenty += size.Y
	}
	return tilelist

}
