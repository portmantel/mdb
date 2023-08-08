package main

import (
	"flag"
	"log"
	"os"
	"crypto/rand"
	//"encoding/json"

	"github.com/portmantel/mdb"
	"github.com/portmantel/rw"
)

var (
	uri = flag.String("uri", "mongodb://localhost:27017", "MongoDB connect string")
	un = flag.String("un", "", "MongoDB Username")
	pw = flag.String("pw", "", "MongoDB Password")
	db = flag.String("db", "bullet_trains", "Database to connect to")
	col = flag.String("col", "chicago_to_florida", "Collection to interact with")
	cnt = flag.Int("cnt", 10, "Number of fake documents to generate")
)

func main() {
	flag.Parse()
	if *un == "" || *pw == "" {
		log.Fatalln("please provide credentials")
	}
	ms := &mdb.MongoServer{
		Uri: *uri,
		Username: *un,
		Password: *pw,
		Log: log.New(os.Stdout, "", 0),
	}
	ms.ClientConnect()
	if ms.Client == nil {
		panic("error connecting to MongoDB")
	}
	ms.DBConnect(*db)

	for i := 0; i < *cnt; i++ {
		b :=  make([]byte, 5)
		_, err := rand.Read(b)
		if err != nil {
			ms.Log.Printf("reading random seed - %s\n", err)
			continue
		}
		bt := &mdb.BulletTrain{
			Velocity: int(b[0]),
			Capacity: int(b[1]),
			Altitude: int(b[2]),
			LastLat: float64(b[3]),
			LastLong: float64(b[4]),
		}
		err = ms.InsertOne(*db, *col, bt)
		if err != nil {
			panic(err)
		} else {
			//ms.Log.Printf("wrote '%d' no error\n", i)
		}
	}
	res := ms.RetrieveAllBulletTrains(*col)
	if len(res) < 1 {
		panic("no results")
	}
	bts := &mdb.BulletTrains{
		Route: make(map[string][]*mdb.BulletTrain),
	}
	for i := 0; i < len(res); i++ {
		/*
		bt, ok := res[i].(BulletTrain)
		if !ok {
			ms.Log.Printf("asserting '%d' not ok\n", i)
			continue
		}
		
		raw, err := json.Marshal([]byte(resbt))
		if err != nil {
			ms.Log.Printf("marshaling raw result '%d' - %s\n", i, err)
			continue
		}
		var bt *BulletTrain
		err = json.Unmarshal(raw, &bt)
		if err != nil {
			ms.Log.Printf("extracting result '%d' - %s\n", i, err)
			continue
		}
		*/
		bts.Route[*col] = append(bts.Route[*col], res[i])
	}
	v := bts.Tabulate(*col)
	rw.TabFlex(mdb.BulletTrainHeaders, v)
}

