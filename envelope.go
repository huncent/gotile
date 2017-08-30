package tile_surge

import (
	//l "github.com/murphy214/layersplit"
	m "github.com/murphy214/mercantile"
	//pc "github.com/murphy214/polyclip"
	"github.com/paulmach/go.geojson"
	//"strings"
	"fmt"
	"sync"
)

// makes a tilemap and returns
func Make_Tilemap(feats *geojson.FeatureCollection, size int) map[m.TileID][]*geojson.Feature {
	c := make(chan map[m.TileID][]*geojson.Feature)
	for _, i := range feats.Features {
		partmap := map[m.TileID][]*geojson.Feature{}

		go func(i *geojson.Feature, size int, c chan map[m.TileID][]*geojson.Feature) {
			//partmap := map[m.TileID][]*geojson.Feature{}

			if i.Geometry.Type == "Polygon" {
				partmap = Env_Polygon(i, size)
			} else if i.Geometry.Type == "LineString" {
				partmap = Env_Line(i, size)
			} else if i.Geometry.Type == "Point" {
				pt := i.Geometry.Point
				tileid := m.Tile(pt[0], pt[1], size)
				partmap[tileid] = append(partmap[tileid], i)
			}
			c <- partmap
		}(i, size, c)
	}

	// collecting channel shit
	totalmap := map[m.TileID][]*geojson.Feature{}
	for range feats.Features {
		partmap := <-c
		for k, v := range partmap {
			totalmap[k] = append(totalmap[k], v...)
		}

	}

	return totalmap
}

// makes children and returns tilemap of a first intialized tilemap
func Make_Tilemap_Children(tilemap map[m.TileID][]*geojson.Feature, prefix string) map[m.TileID][]*geojson.Feature {

	// iterating through each tileid
	ccc := make(chan map[m.TileID][]*geojson.Feature)
	newmap := map[m.TileID][]*geojson.Feature{}
	count2 := 0
	counter := 0
	sizetilemap := len(tilemap)
	buffer := 100000
	//if sizetilemap > 10000 {
	//	buffer = sizetilemap / 5
	//}
	for k, v := range tilemap {
		go func(k m.TileID, v []*geojson.Feature, ccc chan map[m.TileID][]*geojson.Feature) {
			cc := make(chan map[m.TileID][]*geojson.Feature)
			for _, i := range v {
				go func(k m.TileID, i *geojson.Feature, cc chan map[m.TileID][]*geojson.Feature) {
					if i.Geometry.Type == "Polygon" {
						cc <- Children_Polygon(i, k)
					} else if i.Geometry.Type == "LineString" {
						partmap := Env_Line(i, int(k.Z+1))
						cc <- partmap
					} else if i.Geometry.Type == "Point" {
						partmap := map[m.TileID][]*geojson.Feature{}
						pt := i.Geometry.Point
						tileid := m.Tile(pt[0], pt[1], int(k.Z+1))
						partmap[tileid] = append(partmap[tileid], i)
						cc <- partmap
					}
				}(k, i, cc)
			}

			// collecting all into child map
			childmap := map[m.TileID][]*geojson.Feature{}
			for range v {
				tempmap := <-cc
				for k, v := range tempmap {
					childmap[k] = append(childmap[k], v...)
				}
			}

			// making each value in the created childmap and
			// waiting to complete
			var wg sync.WaitGroup
			for k, v := range childmap {
				wg.Add(1)
				go func(k m.TileID, v []*geojson.Feature, prefix string) {
					Make_Tile(k, v, prefix)
					wg.Done()
				}(k, v, prefix)
			}
			wg.Wait()

			ccc <- childmap
		}(k, v, ccc)

		counter += 1
		// collecting shit
		if (counter == buffer) || (sizetilemap-1 == count2) {
			count := 0

			for count < counter {
				tempmap := <-ccc
				for k, v := range tempmap {
					newmap[k] = append(newmap[k], v...)
				}
				count += 1
			}
			counter = 0
			fmt.Printf("\r[%d / %d] Tiles Complete, Size: %d       ", count2, sizetilemap, int(k.Z)+1)

		}
		count2 += 1

	}

	return newmap
}

// makes children and returns tilemap of a first intialized tilemap
func Intialize_Drill(tilemap map[m.TileID][]*geojson.Feature, prefix string, endsize int) {

	// iterating through each tileid
	count2 := 0
	sizetilemap := len(tilemap)
	var wg sync.WaitGroup

	for k, v := range tilemap {
		wg.Add(1)
		go func(k m.TileID, v []*geojson.Feature) {
			Make_Zoom_Drill(k, v, prefix, endsize)
			fmt.Printf("\r[%d / %d] Tiles Recursively Drilled to endsize, %d", count2, sizetilemap, endsize)
			count2 += 1

			wg.Done()
		}(k, v)
	}
	wg.Wait()
}

// recursively drills until the max zoom is reached
func Make_Zoom_Drill(k m.TileID, v []*geojson.Feature, prefix string, endsize int) {
	outputsize := int(k.Z) + 1
	cc := make(chan map[m.TileID][]*geojson.Feature)
	for _, i := range v {
		go func(k m.TileID, i *geojson.Feature, cc chan map[m.TileID][]*geojson.Feature) {
			if i.Geometry.Type == "Polygon" {
				cc <- Children_Polygon(i, k)
			} else if i.Geometry.Type == "LineString" {
				partmap := Env_Line(i, int(k.Z+1))
				partmap = Lint_Children_Lines(partmap, k)
				cc <- partmap
			} else if i.Geometry.Type == "Point" {
				partmap := map[m.TileID][]*geojson.Feature{}
				pt := i.Geometry.Point
				tileid := m.Tile(pt[0], pt[1], int(k.Z+1))
				partmap[tileid] = append(partmap[tileid], i)
				cc <- partmap
			}
		}(k, i, cc)
	}

	// collecting all into child map
	childmap := map[m.TileID][]*geojson.Feature{}
	for range v {
		partmap := <-cc
		for k, v := range partmap {
			childmap[k] = append(childmap[k], v...)
		}
	}

	// iterating through each value in the child map and waiting to complete
	var wg sync.WaitGroup
	for k, v := range childmap {
		//childmap = map[m.TileID][]*geojson.Feature{}
		wg.Add(1)
		go func(k m.TileID, v []*geojson.Feature, prefix string) {
			Make_Tile(k, v, prefix)
			if endsize != outputsize {
				Make_Zoom_Drill(k, v, prefix, endsize)
			}
			wg.Done()

		}(k, v, prefix)
	}
	wg.Wait()

}
