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
	Type string // json, mbtiles, or files
	Minzoom int
	Maxzoom int
	Number_Features int
	Prefix string
	Zooms []int
	Currentzoom int
	Outputjsonfilename string
	Outputmbtilesfilename string
	Memory float64
	Zoom_Config Upper_Zoom_Config
}

// epands the configuration structure
func Expand_Config(config Config) Config {
	count := config.Minzoom
	zooms := []int{}
	for count <= config.Maxzoom {
		zooms = append(zooms,count)
		count += 1
	}
	config.Memory = 2.5
	config.Zooms = zooms
	config.Outputjsonfilename = config.Prefix + ".json"
	config.Outputmbtilesfilename = config.Prefix + ".mbtiles"
	return config
}

// creates the tiles from a given configuration
func Make_Tiles(gjson *geojson.FeatureCollection, config Config) {
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
		db = Insert_Data2(totalmap,db)
	}

	// drilling down tilemap
	totalvts := Intialize_Drill(tilemap,config,db)

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

