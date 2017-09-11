# What is it?

Gotile (which is a temporary name) ingests a data source in geojson format either from a file or postgis database and cuts the features within the geojson into vector tiles for a given zoom range in a manner that focuses primarily on performance. All properties within the geojson are carried throughout the pipeline to the tile. Currently it supports to output formats, a json "directory":buf key value store for each tile and a mbtiles file which is essentially a sqlite db with a little metadata. However with that output you can use mapbox's mbview <file> to view output immediately. 

# How fast is it?
While I don't have definitive benchmarks yet for most datasets I've tested against the closest output the project tippecannoe can provide its anywhere from 3x-6x faster. Take into consideration I've really only thrown a couple data sets, but still performance is pretty good I would estimate. 

#### Why don't you support a normal directory output of tiles??

Go language's (and most programming languages) file creation is basically a wrapped os / unix operation which means you subject to unix / mac os x constraints. At least in Mac Os X the maximum amount of file contexts your allowed to have open is a little over 12k and by default is 256, the nature in which I produce tiles with recursion and drilling means I could have hundreds of thousands of contexts I potentially need to open to output to file. So its an unncessary bottleneck that really isn't worth supporting. I say that out of pragmatism having a directory with potentially a million files / directories takes your os like 10 minutes to do a simple rm -rf dir. Most importantly you can replicate the same functionality with the json output (as long as its not like 6 gb) and a 40 line implementation of a http server in go. I'll probably build that as a command when I get around to doing cli interfaces. 

### File Structure 
- **DB_Interface.go** - this file handles interfacing with postgis for raw data ingestion in bulk using the function db_Interface() and Make_Bounds_Sql() handles iteratively incrementing over each tile in a specific zoom and drilling down turning each tile into a unique query reducing the load on the ram signicantly. 
- **envelope.go** - envelope handles the concurrent processing of the mapping of features to tiles at a specific zoom it also houses the Make_Zoom_Drill() function where most of the work is done.
- **geometry.go** - handles 3 things really the change in projection from the expected srid 4326 to the speriodal espg 3857 so to eliminate distoration, it then converts said points into tile x tile y mapbox tile coords (y = 0 at the top of the tile) for polygon it also ensures winding order is correct post tile xy conversion. 
- **line_envelope** - handles the enveloping or mapping / splitting of single line features to tile ids concurrently
- **poly_envelope** - handles the enveloping or mapping / splitting of single polygon features to tile ids concurrently (also lints the ring order and splits into separate polygons if necessary) 
- **output.go** - handles alot of the procedural stuff for generating the mbtiles metadata and mbtiles transactions
- **tile.go** - handles the creation of a single tile and outputs a structure that contains the raw tile byte data the file directory associated with it and the actual tileid. 
- **util.go** - takes a few important objects within this pipeline as a golang object and returns a string that can be copy and pasted as a test case into a struct in other words does all the parsing to generate a structure raw from the stucture itself. 

# Example 

```go
package main 

import (
	t "github.com/murphy214/gotile/gotile"
)


func main() {
	gjson := t.Read_Geojson("county.geojson")
	config := t.Config{Minzoom:0,Maxzoom:13,Prefix:"county",Type:"mbtiles",New_Output:true,Outputmbtilesfilename:"county.mbtiles"}
	t.Make_Tiles(gjson,config)
}
```

**To view output install mbview with node and export a mapbox key to global env**

```
npm install mbview
mbview county.mbtiles
```

# Output View

![](https://user-images.githubusercontent.com/10904982/30285563-81ecf430-96ec-11e7-9778-52141b18d742.png)

### TODO
* ~~Support multi-geometries~~ (Drills multi geoms to single geometries currently) 
* ~~Write more compresensive unit test cases~~ (Added a few more tests however coverage is still shit) 
* Clean up code, evidence of hacking out different methods or implementations everywhere lol. (still needs done)
  - Also I think to clean up logging kind of hard to decide how to long as the process is recursive and concurrent as well, sort of hard to decide what to log.
* Consider / test supporting geobuf in some capacity within the codebase, not sure how much I want to support currently.
  - If this can be done effeciently memory usage will drop drastically of course. 
* Show a few examles of the different apis etc.
