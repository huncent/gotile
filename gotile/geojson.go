package tile_surge

import (
	//l "github.com/murphy214/layersplit"
	"fmt"
	m "github.com/murphy214/mercantile"
	"github.com/paulmach/go.geojson"
	"io/ioutil"
	//"sync"
	"time"
	"database/sql"

)


// Lints properties within multigeometies
func Lint_Properties(props map[string]interface{}) map[string]interface{} {
	newprops := map[string]interface{}{}
	for k,v := range props {
		if k != "id" {
			newprops[k] = v
		}
	}
	return newprops
}

// Splits multiple geometries into single geoemetries
func Split_Multi(gjson *geojson.FeatureCollection) *geojson.FeatureCollection {
	// splitting multi geometriess
	c := make(chan []*geojson.Feature)
	for _,i := range gjson.Features {
		go func(i *geojson.Feature,c chan []*geojson.Feature) {
			if i.Geometry.Type == "MultiLineString" {
				props := i.Properties
				props = Lint_Properties(props)
				newfeats := []*geojson.Feature{}
				for _,newline := range i.Geometry.MultiLineString {
					newfeats = append(newfeats,&geojson.Feature{Geometry:&geojson.Geometry{LineString:newline,Type:"LineString"},Properties:props})
				}

				c <- newfeats
			} else if i.Geometry.Type == "MultiPolygon" {
				props := i.Properties
				props = Lint_Properties(props)

				newfeats := []*geojson.Feature{}
				for _,newline := range i.Geometry.MultiPolygon {
					newfeats = append(newfeats,&geojson.Feature{Geometry:&geojson.Geometry{Polygon:newline,Type:"Polygon"},Properties:props})
				}
				c <- newfeats
			} else if i.Geometry.Type == "MultiPoint" {
				props := i.Properties

				props = Lint_Properties(props)

				newfeats := []*geojson.Feature{}
				for _,newline := range i.Geometry.MultiPoint {
					newfeats = append(newfeats,&geojson.Feature{Geometry:&geojson.Geometry{Point:newline,Type:"Point"},Properties:props})
				}
				c <- newfeats
			} else {
				i.Properties = Lint_Properties(i.Properties)

				c <- []*geojson.Feature{i}
			}
		}(i,c)
	}
	newfeats := []*geojson.Feature{}
	for range gjson.Features {
		newfeats = append(newfeats,<-c...)
	}	
	return &geojson.FeatureCollection{Features:newfeats}
} 



// reads geojson feature collection into memory
func Read_Geojson(filename string) *geojson.FeatureCollection {
	e, _ := ioutil.ReadFile(filename)
	fc1, _ := geojson.UnmarshalFeatureCollection(e)
	return fc1
}

// upper zoom configuration for sql shit
type Upper_Zoom_Config struct {
	Unique string // unique field
	Zoom_Map map[int]float64
}

// Configuration shit
type Config struct {
	Type string // json or mbtiles
	Minzoom int // minimum zoom
	Maxzoom int // maximum zoom
	Number_Features int // number of features (not needed)
	Prefix string // prefix
	Zooms []int // zooms (not needed)
	Currentzoom int // current zoom (not needed)
	Outputjsonfilename string // output json filename (not needed)
	Outputmbtilesfilename string // output mbtiles filename (not needed)
	Memory float64 // memory (not needed)
	Zoom_Config Upper_Zoom_Config // zoom config (currently not used)
	New_Output bool // output whether to delete the old output or keep it
	Json_Meta string // json metadata
}

// epands the configuration structure
func Expand_Config(config Config) Config {
	count := config.Minzoom
	zooms := []int{}
	for count <= config.Maxzoom {
		zooms = append(zooms,count)
		count += 1
	}

	if config.New_Output != false {
		config.New_Output = true
	}

	config.Memory = 2.5
	config.Zooms = zooms
	config.Outputjsonfilename = config.Prefix + ".json"
	if len(config.Outputmbtilesfilename) == 0 {
		config.Outputmbtilesfilename = config.Prefix + ".mbtiles"
	}
	return config
}

// creates the tiles from a given configuration
func Make_Tiles(gjson *geojson.FeatureCollection, config Config) {
	// Splittng multiple geometries
	gjson = Split_Multi(gjson)


	// creating config expansion
	config = Expand_Config(config)
	fmt.Print("Writing Layers ", config.Zooms, "\n")

	// reading geojson
	var db *sql.DB
	if config.Type == "mbtiles" {
		db = Create_Database_Meta(config,gjson.Features[0])
	}


	s := time.Now()

	// iterating through each zoom
	// creating tilemap
	// getting prefix and min zooom 
	prefix := config.Prefix
	config.Currentzoom = config.Minzoom

	// creating totalmap for tiles under 5 
	// any tiles under 5 arent worth recursively drilling
	totalmap := map[m.TileID]Vector_Tile{}
	tilemap := map[m.TileID][]*geojson.Feature{}
	totalsize := 0
	for (config.Currentzoom <= 5) || (config.Currentzoom == config.Minzoom) {
		// creating tile map for current layer
		if config.Currentzoom == config.Minzoom {
			tilemap,totalsize = Make_Tilemap(gjson, config.Currentzoom)
		} else {
			tilemap,totalsize = Make_Tilemap_Children(tilemap, prefix)
		}

		c := make(chan Vector_Tile)
		for k, v := range tilemap {
			go func(k m.TileID, v []*geojson.Feature, prefix string,c chan Vector_Tile) {
				c <- Make_Tile(k, v, prefix,config)
			}(k, v, prefix,c)
		}

		// iterating through tile map
		for range tilemap {
			v := <- c
			totalmap[v.Tileid] = v
		}

		// incrementing the current zoom 
		config.Currentzoom = config.Currentzoom + 1
	}

	// number of features
	config.Number_Features = totalsize

	// drilling if needed
	// sending the tilemap into the driller
	if config.Type == "mbtiles" {
		db = Insert_Data2(totalmap,db,config)
	}

	// drilling down tilemap
	var totalvts []Vector_Tile
	if config.Currentzoom < config.Maxzoom {
		totalvts = Intialize_Drill(tilemap,config,db)
	}

	// sending into the correct function for output type\
	if config.Type == "json" {
		// adding totalvts to the totalmap
		for _,i := range totalvts {
			totalmap[i.Tileid] = i
		}

		Write_Json(totalmap,config.Outputjsonfilename)
	} else if config.Type == "mbtiles" {
		Make_Index(db)
	}

	fmt.Printf("\nCompleted in %s.\n", time.Now().Sub(s))
}

// makes a sql query for a tile specific query at the first zoom level
// this function will then iteratively go through each query based on a rough memory calculation
// to estimate how many top level routines can be throughput at once
func Make_Bounds_Sql(database string,tablename string,basesql string,config Config) {
	s := time.Now()
	// getting the extent and number of rows
	ext,num_b := Get_Extent(database,tablename)

	// getting the config shit
	config = Expand_Config(config)
	config.Currentzoom = config.Minzoom
	fmt.Print("Writing Layers ", config.Zooms, "\n")

	// getting json sample
	gjson := DB_Interface(database,basesql + " LIMIT 1;")

	// reading geojson
	var db *sql.DB
	if config.Type == "mbtiles" {
		db = Create_Database_Meta(config,gjson.Features[0])
	}

	// getting the total map for the upper zomos
	//totalmap := map[m.TileID]Vector_Tile{}


	config.Number_Features = num_b

	// getting the sema size
	sema_size_sql := Size_Stovepipe(config) 
	fmt.Printf("Max Make_Tiles_Sql Go Routines: %d\n",sema_size_sql)

	// creating sema
	var sema_sql = make(chan struct{}, sema_size_sql)

	// getting the tile list
	tilelist := Make_Tilelist(ext,config.Currentzoom)
	count := 0
	// iterating through each tile in the tilelist
	if config.Maxzoom > config.Currentzoom {
		c := make(chan []Vector_Tile) 
		sizetilelist := len(tilelist)
		for _,i := range tilelist {
			go func(i m.TileID,c chan []Vector_Tile) {
				sema_sql <- struct{}{}        // acquire token
				defer func() { <-sema_sql }() // release token
				count += 1
				fmt.Printf("\n[%d/%d] Sql Tiles Started.\n",count,sizetilelist)

				// getting query logic
				bbox_logic := Add_BBox(tablename,i)
				one_query := fmt.Sprintf("%s WHERE %s",basesql,bbox_logic)
				
				// selecint data and piping to channel
				//DB_Interface(database,one_query).Features
				c <-Make_Zoom_Drill(i, DB_Interface(database,one_query).Features, config.Prefix, config.Maxzoom,config)
			}(i,c)
		}

		// iterating over tilelist
		for range tilelist {
			vtmap := <-c
			Insert_Data3(vtmap,db,config)
		}
	}

	// finishing creation of output type
	if config.Type == "json" {
		//Write_Json(totalmap,config.Outputjsonfilename)
	} else if config.Type == "mbtiles" {
		Make_Index(db)
	}
	fmt.Printf("Time creating mbtiles %s.\n",time.Now().Sub(s))
}

