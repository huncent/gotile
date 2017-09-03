package tile_surge

import (
	"io/ioutil"
	//"sync"
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	m "github.com/murphy214/mercantile"
	"encoding/json"
	"reflect"
	"bytes"
	"compress/gzip"
	"github.com/paulmach/go.geojson"
	"os"
	"log"

)

var jsondata = `{
        "vector_layers": [
            {
                "id": "county",
                "description": "",
                "minzoom": 5,
                "maxzoom": 13,
                "fields": {
                    "area": "String",
                    "colorkey":"String"
                }
            }]
    }`


// vector layer json
type Vector_Layer struct {
	ID string `json:"id"`
	Description string `json:"description"`
	Minzoom int `json:"minzoom"`
	Maxzoom int `json:"maxzoom"`
	Fields map[string]string `json:"fields"`
}


type Vector_Layers struct {
	Vector_Layers []Vector_Layer `json:"vector_layers"`
}

// returns the string of the json meta data
func Make_Json_Meta(config Config,feat *geojson.Feature) string {
	layer := Vector_Layer{ID:config.Prefix,Description:"",Minzoom:config.Minzoom,Maxzoom:config.Maxzoom}

	fields := Reflect_Fields(feat.Properties)
	layer.Fields = fields

	vector_layers := Vector_Layers{Vector_Layers:[]Vector_Layer{layer}}

	b,_ := json.Marshal(vector_layers)
	return string(b)
}

// creating the slice that will be used to create the metadata table
func Make_Metadata_Slice(config Config,feat *geojson.Feature) [][]string {
	// getting the json blob metadata
	jsondata := Make_Json_Meta(config,feat)

	// creating values 
	values := [][]string{{"name",config.Outputmbtilesfilename},{"type","overlay"},{"version","2"},{"description",config.Outputmbtilesfilename},{"format","pbf"},{"json",jsondata}}

	return values
}

// creates the sqllite database and inserts metadata 
func Create_Database_Meta(config Config,feat *geojson.Feature) *sql.DB {
	os.Remove(config.Outputmbtilesfilename)

	db, err := sql.Open("sqlite3", config.Outputmbtilesfilename)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()
	fmt.Printf("Creating and opening %s.\n",config.Outputmbtilesfilename)


	sqlStmt := `
	CREATE TABLE metadata (name text, value text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
	fmt.Printf("Created metadata table: %s.\n",config.Outputmbtilesfilename)

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into metadata(name, value) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	// creating metadata slice string
	values := Make_Metadata_Slice(config,feat)


	defer stmt.Close()
	for _,i := range values {
		_, err = stmt.Exec(i[0],i[1])
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()


	sqlStmt = `
	CREATE TABLE tiles (zoom_level integer, tile_column integer, tile_row integer, tile_data blob);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}

	tx, err = db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted appropriate metadata: %s.\n",config.Outputmbtilesfilename)
	return db 
}



// reflects a tile value back and stuff
func Reflect_Fields(mymap map[string]interface{}) map[string]string {
	newmap := map[string]string{}
	for k,v := range mymap {

		vv := reflect.ValueOf(v)
		kd := vv.Kind()
		if (reflect.Float64 == kd) || (reflect.Float32 == kd) {
			//fmt.Print(v, "float", k)
			newmap[k] = "Float"
			//hash = Hash_Tv(tv)
		} else if (reflect.Int == kd) || (reflect.Int8 == kd) || (reflect.Int16 == kd) || (reflect.Int32 == kd) || (reflect.Int64 == kd) || (reflect.Uint8 == kd) || (reflect.Uint16 == kd) || (reflect.Uint32 == kd) || (reflect.Uint64 == kd) {
			//fmt.Print(v, "int", k)
			newmap[k] = "Integer"
			//hash = Hash_Tv(tv)
		} else if reflect.String == kd {
			//fmt.Print(v, "str", k)
			newmap[k] = "String"
			//hash = Hash_Tv(tv)

		} else {
			fmt.Print(k,v,"\n")
		}
	}
	return newmap
}


// inserting data into shit
func Insert_Data(newmap map[m.TileID]Vector_Tile,db *sql.DB) *sql.DB {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	
	stmt, err := tx.Prepare("insert into tiles(zoom_level, tile_column,tile_row,tile_data) values(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	//values := [][]string{{"name","shit"},{"type","baselayer"},{"version","1.2"},{"description","shit"},{"format","pbf"}}


	defer stmt.Close()
	count := 0
	total := 0
	counter := 0
	count3 := 0
	sizenewmap := len(newmap)
	c := make(chan Vector_Tile)
    var b bytes.Buffer
    gz := gzip.NewWriter(&b)

	for k,v := range newmap {
		go func(k m.TileID,v Vector_Tile,c chan Vector_Tile) {
			gz.Reset(&b)
			if _, err := gz.Write(v.Data); err != nil {
				panic(err)
			}
			if err := gz.Flush(); err != nil {
				panic(err)
			}

			v.Data = b.Bytes()
			c <- v
		}(k,v,c)
		counter += 1
		if counter == 1000 || (sizenewmap - 1 == count3){
			count2 := 0
			for count2 < counter {
				v := <-c
				k := v.Tileid
				k.Y = (1 << uint64(k.Z)) - 1 - k.Y 
				_, err = stmt.Exec(int(k.Z),int(k.X),int(k.Y),v.Data)
				if err != nil {
					log.Fatal(err)
				}
				count += 1
				count2 += 1
				if count == 1000 {
					count = 0
					total += 1000
					fmt.Printf("\r[%d/%d] Compressing tiles and inserting into db.",total,sizenewmap)
				}

			}
			counter = 0
		}
		count3 += 1
		//fmt.Print(count,"\n")
		//count += 1
	}



	tx.Commit()


	return db
	
}

// inserting data into shit
func Insert_Data2(newmap map[m.TileID]Vector_Tile,db *sql.DB) *sql.DB {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	
	stmt, err := tx.Prepare("insert into tiles(zoom_level, tile_column,tile_row,tile_data) values(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	//values := [][]string{{"name","shit"},{"type","baselayer"},{"version","1.2"},{"description","shit"},{"format","pbf"}}


	defer stmt.Close()
	count := 0
	total := 0
	count3 := 0


	sizenewmap := len(newmap)
    var b bytes.Buffer
    gz := gzip.NewWriter(&b)

	for k,v := range newmap {
		if _, err := gz.Write(v.Data); err != nil {
			panic(err)
		}
		if err := gz.Flush(); err != nil {
			panic(err)
		}

		v.Data = b.Bytes()
		//fmt.Print(v,"\n")
		//if err := gz.Close(); err != nil {
        //	panic(err)
    	//}	
		//fmt.Print(len(bb.Bytes()),len(v),count,"\n")
		//bb := new(bytes.Buffer)
    	b = *bytes.NewBuffer([]byte{})
    	gz.Reset(&b)
		k.Y = (1 << uint64(k.Z)) - 1 - k.Y 
		_, err = stmt.Exec(int(k.Z),int(k.X),int(k.Y),v.Data)
		if err != nil {
			log.Fatal(err)
		}
		count += 1
		if count == 1000 {
			count = 0
			total += 1000
			fmt.Printf("\r[%d/%d] Compressing tiles and inserting into db.",total,sizenewmap)
		}

		count3 += 1
		//fmt.Print(count,"\n")
		//count += 1
	}



	tx.Commit()


	return db
	
}



func Make_Index(db *sql.DB) {
	defer db.Close()

	sqlStmt := `
	CREATE UNIQUE INDEX tile_index on tiles (zoom_level, tile_column, tile_row)
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

}



// gzip bytes 
func Gzip_Bytes(in []byte) []byte {
    var b bytes.Buffer
    gz := gzip.NewWriter(&b)
    if _, err := gz.Write(in); err != nil {
        panic(err)
    }
    if err := gz.Flush(); err != nil {
        panic(err)
    }
    if err := gz.Close(); err != nil {
        panic(err)
    }
    return b.Bytes()
}


// writes a json file
func Write_Json(totalmap map[m.TileID]Vector_Tile,jsonfilename string) {
	newmap := map[string][]byte{}
	for _,i := range totalmap {
		newmap[i.Filename] = i.Data
	}

	b,_ := json.Marshal(newmap)
	ioutil.WriteFile(jsonfilename,b,0666)

}
