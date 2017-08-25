package tile_surge

import (
	//l "github.com/murphy214/layersplit"
	m "github.com/murphy214/mercantile"
	//pc "github.com/murphy214/polyclip"
	"github.com/paulmach/go.geojson"
	//"strings"
	"fmt"
)

func Make_Tilemap(feats *geojson.FeatureCollection, size int) map[m.TileID][]*geojson.Feature {
	c := make(chan map[m.TileID][]*geojson.Feature)
	for _, i := range feats.Features {
		go func(i *geojson.Feature, size int, c chan map[m.TileID][]*geojson.Feature) {
			partmap := map[m.TileID][]*geojson.Feature{}

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
						cc <- Env_Line(i, int(k.Z+1))
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

			c := make(chan string)
			for k, v := range childmap {
				go func(k m.TileID, v []*geojson.Feature, prefix string, c chan string) {
					Make_Tile(k, v, prefix)
					c <- ""
				}(k, v, prefix, c)
			}

			count := 0
			for range childmap {
				msg1 := <-c
				fmt.Printf("%s", msg1)
				count += 1
			}

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
func Make_Tilemap_Children2(tilemap map[m.TileID][]*geojson.Feature, prefix string, endsize int) {

	// iterating through each tileid
	ccc := make(chan string)
	//newmap := map[m.TileID][]*geojson.Feature{}
	count2 := 0
	//counter := 0
	sizetilemap := len(tilemap)
	//buffer := 100000
	//if sizetilemap > 10000 {
	//	buffer = sizetilemap / 5
	//}
	for k, v := range tilemap {
		go func(k m.TileID, v []*geojson.Feature, ccc chan string) {
			Make_Zoom_Drill(k, v, prefix, endsize)
			ccc <- ""
		}(k, v, ccc)
	}

	for range tilemap {
		fmt.Printf("\r[%d / %d] Tiles Recursively Drilled to endsize, %d", count2, sizetilemap, endsize)
		fmt.Print(<-ccc)
		count2 += 1
	}
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
				cc <- Env_Line(i, int(k.Z+1))
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

	c := make(chan string)
	for k, v := range childmap {
		go func(k m.TileID, v []*geojson.Feature, prefix string, c chan string) {
			Make_Tile(k, v, prefix)
			c <- ""
		}(k, v, prefix, c)
	}

	count := 0
	for range childmap {
		msg1 := <-c
		fmt.Printf("%s", msg1)
		count += 1
	}

	if endsize != outputsize {
		chans := make(chan string)
		for k, v := range childmap {
			go func(k m.TileID, v []*geojson.Feature, endsize int, chans chan string) {
				Make_Zoom_Drill(k, v, prefix, endsize)
				c <- ""
			}(k, v, endsize, chans)
		}

		for range childmap {
			fmt.Print(<-c)
		}
	}
}
