package mdb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	BTdb = "bullet_trains"
)

type BulletTrains struct {
	Route map[string][]*BulletTrain
}

func (bts *BulletTrains) Tabulate(route string) (values [][]interface{}) {
	for i := 0; i < len(bts.Route[route]); i++ {
		values = append(values, bts.Route[route][i].Flatten())
	}
	return values
}

type BulletTrain struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Velocity int `bson:"velocity"`
	Capacity int `bson:"capacity"`
	Altitude int `bson:"altitude"`
	LastLat float64 `bson:"last_lat"`
	LastLong float64 `bson:"last_long"`
}

func (bt *BulletTrain) Flatten() []interface{} {
	return []interface{}{
		bt.ID.Hex(),
		bt.Velocity,
		bt.Capacity,
		bt.Altitude,
		bt.LastLat,
		bt.LastLong,
	}
}

var BulletTrainHeaders = []string{
	"ID",
	"Velocity",
	"Capacity",
	"Altitude",
	"LastLat",
	"LastLong",
}

func (ms *MongoServer) RetrieveAllBulletTrains(collection string) (results []*BulletTrain) {
	err := ms.DB[BTdb].Client().Ping(context.TODO(), nil)
	if err != nil {
		ms.Log.Printf("pinging '%s' - %s", BTdb, err)
		return
	}
	cursor, err := ms.DB[BTdb].Collection(collection).Find(context.TODO(), bson.M{})
	if err != nil {
		ms.Log.Printf("retrieving all from '%s.%s' - %s", BTdb, collection, err)
		return
	}
	var res []BulletTrain
	err = cursor.All(context.TODO(), &res)
	if err != nil {
		ms.Log.Printf("extracting '%s' collection documents - %s", collection, err)
		return
	}
	for i := 0; i < len(res); i++ {
		results = append(results, &res[i])
	}
	return
}