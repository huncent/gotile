package tile_surge

import (
	//l "github.com/murphy214/layersplit"
	m "github.com/murphy214/mercantile"
	//pc "github.com/murphy214/polyclip"
	"github.com/paulmach/go.geojson"
	//"strings"
	"fmt"
	"database/sql"
	"math"
	//"log"
	//"sync"
)

// makes a tilemap and returns
func Make_Tilemap(feats *geojson.FeatureCollection, size int) (map[m.TileID][]*geojson.Feature,int) {
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

	// getting size of total number of features within the tilemap
	totalsize := 0
	for _,v := range totalmap {
		totalsize += len(v)
	}



	return totalmap,totalsize
}

// rounds a number to an adequate value 
func Round(val float64, roundOn float64, places int ) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

// calculates the maximum amount of concurrent actions that can be performed in memory 
func Size_Stovepipe(config Config) int {
	// getting the delta between the zooms 
	delta := config.Maxzoom - config.Currentzoom - 1

	// getting the number of go routines
	number_go_routines := math.Pow(4.0,float64(delta))

	// assuming 10 kb per go routine
	size_each_routine := float64(4) // kb

	// estimated or maximum memory within stovepipe
	size_stovepipe := number_go_routines * size_each_routine + float64(config.Number_Features) * float64(4)

	// getting the size of a gb
	size_gb := size_stovepipe / 1000.0 / 1000.0

	// getting the number of sem things to make
	size_sem := int(Round(config.Memory / size_gb,.5,0))
	return size_sem

}

// makes children and returns tilemap of a first intialized tilemap
func Make_Tilemap_Children(tilemap map[m.TileID][]*geojson.Feature, prefix string) (map[m.TileID][]*geojson.Feature,int) {

	// iterating through each tileid
	ccc := make(chan map[m.TileID][]*geojson.Feature)
	newmap := map[m.TileID][]*geojson.Feature{}
	count2 := 0
	counter := 0
	sizetilemap := len(tilemap)
	buffer := 100000

	// iterating through each tielmap
	for k, v := range tilemap {
		go func(k m.TileID, v []*geojson.Feature, ccc chan map[m.TileID][]*geojson.Feature) {
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
				tempmap := <-cc
				for k, v := range tempmap {
					childmap[k] = append(childmap[k], v...)
				}
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


	// getting size of total number of features within the tilemap
	totalsize := 0
	for _,v := range newmap {
		totalsize += len(v)
	}



	return newmap,totalsize
}


// makes children and returns tilemap of a first intialized tilemap
func Intialize_Drill(tilemap map[m.TileID][]*geojson.Feature,config Config, db *sql.DB) []Vector_Tile {
	// getting size sema (i.e. the limitation on how many go functions are called )
	// in the routine below, calculated by teh memory input config
	size_sem := Size_Stovepipe(config)
	if size_sem == 0 {
		size_sem = 1
	}
	fmt.Printf("Max Make_Zoom_Drill Go Routines: %d\n",size_sem)

	// creating sema
	var sema = make(chan struct{}, size_sem)

	// intializing values	
	prefix := config.Prefix
	endsize := config.Maxzoom
	count2 := 0
	sizetilemap := len(tilemap)
	//var wg sync.WaitGroup
	count := 0
	c := make(chan []Vector_Tile)
	for k, v := range tilemap {

		go func(k m.TileID, v []*geojson.Feature,c chan []Vector_Tile) {

			sema <- struct{}{}        // acquire token
			defer func() { <-sema }() // release token
			c <- Make_Zoom_Drill(k, v, prefix, endsize,config)
			fmt.Printf("[%d / %d] Tiles Recursively Drilled to endsize, %d\n", count2, sizetilemap, endsize)
			count2 += 1
		}(k, v,c)



		count +=1 
	}
	//wg.Wait()

	total := 0
	// iterating through each value in the tilemap
	totalvts := []Vector_Tile{}
	for range tilemap {
		vts := <- c
		total += len(vts)
		if config.Type == "mbtiles" {
			Insert_Data3(vts,db,config)

		} else if config.Type == "json" {
			totalvts = append(totalvts,vts...)
		}
		fmt.Printf("\nTotal Number of Tiles %d\n",total)

	}
	return totalvts
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
