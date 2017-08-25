package tile_surge

import (
	//l "github.com/murphy214/layersplit"
	"fmt"
	m "github.com/murphy214/mercantile"
	"github.com/paulmach/go.geojson"
	"io/ioutil"
	"time"
)

// reads geojson feature collection into memory
func Read_Geojson(filename string) *geojson.FeatureCollection {
	e, _ := ioutil.ReadFile(filename)
	fc1, _ := geojson.UnmarshalFeatureCollection(e)
	return fc1
}

// creates the tiles from a given configuration
func Make_Tiles_Geojson(gjson *geojson.FeatureCollection, prefix string, zooms []int) {
	// reading geojson
	s := time.Now()
	// iterating through each zoom
	for _, i := range zooms {
		// creating tilemap
		tilemap := Make_Tilemap(gjson, i)

		// iterating through each tileid in the tilemap
		c := make(chan string)
		for k, v := range tilemap {
			go func(k m.TileID, v []*geojson.Feature, prefix string, c chan string) {
				Make_Tile(k, v, prefix)
				c <- ""
			}(k, v, prefix, c)
		}
		count := 0
		sizetilemap := len(tilemap)
		for range tilemap {
			msg1 := <-c
			fmt.Printf("\r[%d / %d] Tiles Complete of Size %d%s", count, sizetilemap, i, msg1)
			count += 1
		}

	}
	fmt.Printf("\nCompleted in %s.\n", time.Now().Sub(s))
}

// creates the tiles from a given configuration
func Make_Tiles_Geojson2(gjson *geojson.FeatureCollection, prefix string, zooms []int) {
	// reading geojson
	s := time.Now()
	// iterating through each zoom
	// creating tilemap
	tilemap := Make_Tilemap(gjson, zooms[0])

	// iterating through each tileid in the tilemap
	c := make(chan string)
	for k, v := range tilemap {
		go func(k m.TileID, v []*geojson.Feature, prefix string, c chan string) {
			Make_Tile(k, v, prefix)
			c <- ""
		}(k, v, prefix, c)
	}
	count := 0
	sizetilemap := len(tilemap)
	for range tilemap {
		msg1 := <-c
		fmt.Printf("\r[%d / %d] Tiles Complete of Size %d%s", count, sizetilemap, zooms[0], msg1)
		count += 1
	}

	for range zooms[1:] {
		tilemap = Make_Tilemap_Children(tilemap, prefix)
	}

	fmt.Printf("\nCompleted in %s.\n", time.Now().Sub(s))
}

// creates the tiles from a given configuration
func Make_Tiles_Geojson3(gjson *geojson.FeatureCollection, prefix string, zooms []int) {
	fmt.Print("Writing Layers ", zooms, "\n")
	// reading geojson
	s := time.Now()

	// iterating through each zoom
	// creating tilemap
	tilemap := Make_Tilemap(gjson, zooms[0])

	// iterating through each tileid in the tilemap
	c := make(chan string)
	for k, v := range tilemap {
		go func(k m.TileID, v []*geojson.Feature, prefix string, c chan string) {
			Make_Tile(k, v, prefix)
			c <- ""
		}(k, v, prefix, c)
	}
	count := 0
	sizetilemap := len(tilemap)
	for range tilemap {
		msg1 := <-c
		fmt.Printf("\r[%d / %d] Tiles Complete of Size %d%s", count, sizetilemap, zooms[0], msg1)
		count += 1
	}
	tilemap = Make_Tilemap_Children(tilemap, prefix)
	tilemap = Make_Tilemap_Children(tilemap, prefix)

	Make_Tilemap_Children2(tilemap, prefix, zooms[len(zooms)-1])

	fmt.Printf("\nCompleted in %s.\n", time.Now().Sub(s))
}
