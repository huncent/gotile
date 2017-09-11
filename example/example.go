package main 

import (
	t "github.com/murphy214/gotile/gotile"
)


func main() {

	gjson := t.Read_Geojson("county.geojson")
	config := t.Config{Minzoom:0,Maxzoom:13,Prefix:"county",Type:"mbtiles",New_Output:true,Outputmbtilesfilename:"county.mbtiles"}
	t.Make_Tiles(gjson,config)

}