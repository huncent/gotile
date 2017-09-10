package main

import (
  //"fmt"
  "os"
  t "github.com/murphy214/gotile/gotile"
  "github.com/urfave/cli"
  "strings"
)

func get_config(filename string) t.Config {
	prefix := strings.Split(filename,".")[0]
	return t.Config{Minzoom:0,Maxzoom:14,Prefix:prefix,Type:"mbtiles"}
}


func main() {
  app := cli.NewApp()

  app.Action = func(c *cli.Context) error {
    filename := c.Args().Get(0)
    config := get_config(filename)
    gjson := t.Read_Geojson(filename)
    t.Make_Tiles(gjson,config)
    //fmt.Printf("Hello %q", c.Args().Get(0))
    return nil
  }

  app.Run(os.Args)
}