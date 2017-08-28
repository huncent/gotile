package tile_surge

import (
	//l "github.com/murphy214/layersplit"
	"fmt"
	m "github.com/murphy214/mercantile"
	"github.com/paulmach/go.geojson"
	"io/ioutil"
	"sync"
	"time"
)

// reads geojson feature collection into memory
func Read_Geojson(filename string) *geojson.FeatureCollection {
	e, _ := ioutil.ReadFile(filename)
	fc1, _ := geojson.UnmarshalFeatureCollection(e)
	return fc1
}

// creates the tiles from a given configuration
func Make_Tiles(gjson *geojson.FeatureCollection, prefix string, zooms []int) {
	fmt.Print("Writing Layers ", zooms, "\n")
	// reading geojson
	s := time.Now()

	// iterating through each zoom
	// creating tilemap
	tilemap := Make_Tilemap(gjson, zooms[0])

	// iterating through each tileid in the tilemap
	sizetilemap := len(tilemap)
	count := 0
	var wg sync.WaitGroup
	for k, v := range tilemap {
		wg.Add(1)
		go func(k m.TileID, v []*geojson.Feature, prefix string) {
			Make_Tile(k, v, prefix)
			fmt.Printf("\r[%d / %d] Tiles Complete of Size %d", count, sizetilemap, zooms[0])
			count += 1
			wg.Done()
		}(k, v, prefix)
	}
	wg.Wait()

	// drilling if needed
	// sending the tilemap into the driller
	Intialize_Drill(tilemap, prefix, zooms[len(zooms)-1])

	fmt.Printf("\nCompleted in %s.\n", time.Now().Sub(s))
}
