package tile_surge

import (
	//l "github.com/murphy214/layersplit"
	m "github.com/murphy214/mercantile"
	//pc "github.com/murphy214/polyclip"
	"github.com/paulmach/go.geojson"
	//"strings"
	"fmt"
	//"sync"
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
func Intialize_Drill(tilemap map[m.TileID][]*geojson.Feature,config Config,totalmap map[m.TileID]Vector_Tile) map[m.TileID]Vector_Tile {
	prefix := config.Prefix
	endsize := config.Maxzoom
	// iterating through each tileid
	count2 := 0
	sizetilemap := len(tilemap)
	//var wg sync.WaitGroup
	count := 0
	c := make(chan []Vector_Tile)
	for k, v := range tilemap {
		//wg.Add(1)
		go func(k m.TileID, v []*geojson.Feature,c chan []Vector_Tile) {
			c <- Make_Zoom_Drill(k, v, prefix, endsize,config)
			fmt.Printf("\r[%d / %d] Tiles Recursively Drilled to endsize, %d", count2, sizetilemap, endsize)
			count2 += 1

			//wg.Done()
		}(k, v,c)
		
		count +=1 
	}
	//wg.Wait()

	// iterating through each value in the tilemap
	for range tilemap {
		vts := <- c
		for _,v := range vts {
			totalmap[v.Tileid] = v
		}
	}
	return totalmap
}


// vector tile struct 
type Vector_Tile struct {
	Filename string
	Data []byte
	Tileid m.TileID
}

// recursively drills until the max zoom is reached
func Make_Zoom_Drill(k m.TileID, v []*geojson.Feature, prefix string, endsize int,config Config) []Vector_Tile {
	outputsize := int(k.Z) + 1
	cc := make(chan map[m.TileID][]*geojson.Feature)
	for _, i := range v {
		go func(k m.TileID, i *geojson.Feature, cc chan map[m.TileID][]*geojson.Feature) {
			if i.Geometry.Type == "Polygon" {
				partmap := Children_Polygon(i, k) 
				cc <- partmap
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
		for kk, vv := range partmap {
			childmap[kk] = append(childmap[kk], vv...)
		}
	}

	// iterating through each value in the child map and waiting to complete
	//var wg sync.WaitGroup
	vtchan := make(chan Vector_Tile)
	for kkk, vvv := range childmap {
		//childmap = map[m.TileID][]*geojson.Feature{}
		//wg.Add(1)
		go func(kkk m.TileID, vvv []*geojson.Feature, prefix string,vtchan chan Vector_Tile) {
			vtchan <- Make_Tile(kkk, vvv, prefix,config)
				//Make_Zoom_Drill(kkk, vvv, prefix, endsize)
			//wg.Done()

		}(kkk, vvv, prefix,vtchan)
	}
	
	vector_tiles := []Vector_Tile{}
	for range childmap {
		vt := <-vtchan
		vector_tiles = append(vector_tiles,vt)
	}

	//wg.Wait()
	if endsize != outputsize {
		ccc := make(chan []Vector_Tile)
		for kkk, vvv := range childmap {
			go func(kkk m.TileID, vvv []*geojson.Feature, prefix string,ccc chan []Vector_Tile) {
				ccc <- Make_Zoom_Drill(kkk,vvv,prefix,endsize,config)
			}(kkk,vvv,prefix,ccc)
		}
		// appending to the major vector tiles shit
		for range childmap {
			vts := <- ccc
			vector_tiles = append(vector_tiles,vts...)
		}
		return vector_tiles

	} else {
		return vector_tiles
	}
}
