package tile_surge

import (
	//l "github.com/murphy214/layersplit"
	"fmt"
	m "github.com/murphy214/mercantile"
	"github.com/paulmach/go.geojson"
	"io/ioutil"
	//"sync"
	"time"
)

// reads geojson feature collection into memory
func Read_Geojson(filename string) *geojson.FeatureCollection {
	e, _ := ioutil.ReadFile(filename)
	fc1, _ := geojson.UnmarshalFeatureCollection(e)
	return fc1
}

type Config struct {
	Type string // json, mbtiles, or files
	Minzoom int
	Maxzoom int
	Prefix string
	Zooms []int
	Outputjsonfilename string
	Outputmbtilesfilename string
}

func Expand_Config(config Config) Config {
	count := config.Minzoom
	zooms := []int{}
	for count <= config.Maxzoom {
		zooms = append(zooms,count)
		count += 1
	}

	config.Zooms = zooms
	config.Outputjsonfilename = config.Prefix + ".json"
	config.Outputmbtilesfilename = config.Prefix + ".mbtiles"
	return config
}

// creates the tiles from a given configuration
func Make_Tiles(gjson *geojson.FeatureCollection, config Config) {
	config = Expand_Config(config)

	fmt.Print("Writing Layers ", config.Zooms, "\n")
	// reading geojson
	s := time.Now()

	// iterating through each zoom
	// creating tilemap
	tilemap := Make_Tilemap(gjson, config.Minzoom)
	prefix := config.Prefix
	// iterating through each tileid in the tilemap
	//sizetilemap := len(tilemap)
	count := 0
	totalmap := map[m.TileID]Vector_Tile{}
	c := make(chan Vector_Tile)
	for k, v := range tilemap {
		go func(k m.TileID, v []*geojson.Feature, prefix string,c chan Vector_Tile) {
			c <- Make_Tile(k, v, prefix,config)
			//fmt.Printf("\r[%d / %d] Tiles Complete of Size %d", count, sizetilemap, zooms[0])
			count += 1
		}(k, v, prefix,c)
	}

	for range tilemap {
		v := <- c
		totalmap[v.Tileid] = v
	}



	// drilling if needed
	// sending the tilemap into the driller
	totalmap = Intialize_Drill(tilemap,config,totalmap)

	if config.Type == "json" {
		Write_Json(totalmap,config.Outputjsonfilename)
	} else if config.Type == "mbtiles" {
		db := Create_Database_Meta(config,gjson.Features[0])
		db = Insert_Data2(totalmap,db)
		Make_Index(db)
	}




	//fmt.Print(len(totalmap),"\n")

	fmt.Printf("\nCompleted in %s.\n", time.Now().Sub(s))
}

