package mdb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type MongoServer struct {
	Uri      string
	Username string
	Password string
	Client   *mongo.Client
	DBNames []string
	DB       map[string]*mongo.Database
	Log *log.Logger
}

// the Client object returned then can initiate specific databases
// the database is not included in the connect string
func (ms *MongoServer) ClientConnect() {
	clientOptions := options.Client().ApplyURI(ms.Uri)
	clientOptions.SetAuth(options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		Username:      ms.Username,
		Password:      ms.Password,
	})
	//clientOptions.SetTimeout(time.Duration(5))
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		ms.Log.Printf("connecting to %s - %s\n", ms.Uri, err)
		return
	}
	if err = client.Ping(context.Background(), nil); err != nil {
		ms.Log.Printf("pinging %s - %s\n", ms.Uri, err)
		return
	}
	ms.Client = client
	ms.DB = make(map[string]*mongo.Database)
	ms.Log.Printf("connected to host '%s'\n", ms.Uri)
}

func (ms *MongoServer) ClientDisconnect() {
	err := ms.Client.Disconnect(context.TODO())
	if err != nil {
		ms.Log.Printf("disconnecting from '%s' - %s\n", ms.Uri, err)
	}
}

func (ms *MongoServer) DBConnectAll() {
	if len(ms.DBNames) < 1 {
		return
	}
	for _, db := range ms.DBNames {
		ms.DBConnect(db)
	}
}

func (ms *MongoServer) DBConnect(db string) {
	ms.DB[db] = ms.Client.Database(db)
	err := ms.DB[db].Client().Ping(context.TODO(), nil)
	if err != nil {
		ms.Log.Printf("pinging (%s) - %s\n", ms.Uri, err)
	} else {
		ms.Log.Printf("connected to db '%s'\n", db)
	}
}

func (ms *MongoServer) InsertOne(db, collection string, data interface{}) error {
	err := ms.DB[db].Client().Ping(context.TODO(), nil)
	if err != nil {
		// db must already be connected
		return err
	}
	_, err = ms.DB[db].Collection(collection).InsertOne(context.TODO(), data)
	return err
}

func (ms *MongoServer) InsertMany(db, collection string, data []interface{}) error {
	err := ms.DB[db].Client().Ping(context.TODO(), nil)
	if err != nil {
		// db must already be connected
		return err
	}
	_, err = ms.DB[db].Collection(collection).InsertMany(context.TODO(), data)
	return err
}

func (ms *MongoServer) RetrieveAll(db, collection string) (results []interface{}) {
	err := ms.DB[db].Client().Ping(context.TODO(), nil)
	if err != nil {
		ms.Log.Printf("pinging '%s' - %s", db, err)
		return
	}
	cursor, err := ms.DB[db].Collection(collection).Find(context.TODO(), bson.M{})
	if err != nil {
		ms.Log.Printf("retrieving all from '%s.%s' - %s", db, collection, err)
		return
	}
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		ms.Log.Printf("extracting collection documents - %s", err)
		return
	}
	/*
		for i := 0; i < len(res); i++ {
			results = append(results, res[i])
		}
	*/
	return
}

func (ms *MongoServer) ExampleUpdateOne(db, collection, findKey, findValue, updateKeyOne, updateKeyTwo string, updateValueOne, updateValueTwo interface{}) {
	filter := bson.M{findKey: findValue}
	update := bson.D{{"$set", bson.M{updateKeyOne: updateValueOne, updateKeyTwo: updateValueTwo}}}
	err := ms.updateOne(db, collection, filter, update)
	if err != nil {
		ms.Log.Printf("updating one of (%s)/(%s)/(%s)&(%s) - %s\n", db, collection, updateKeyOne, updateKeyTwo, err)
	} else {
		ms.Log.Printf("updated one (%s)/(%s)/(%s)&(%s) to (%v)&(%v) on matches of (%s:%s)\n", db, collection, updateKeyOne, updateKeyTwo, updateValueOne, updateValueTwo, findKey, findValue)
	}
}

func (ms *MongoServer) updateOne(db, collection string, filter bson.M, update bson.D) error {
	err := ms.DB[db].Client().Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	_, err = ms.DB[db].Collection(collection).UpdateOne(context.Background(), filter, update)
	return err
}

// simplified filter & update builder only takes one field as an example
func (ms *MongoServer) ExampleUpdateMany(db, collection, findKey, findValue, updateKey string, updateValue interface{}) {
	filter := bson.M{findKey: findValue}
	update := bson.D{{"$set", bson.M{updateKey: updateValue}}}
	err := ms.updateMany(db, collection, filter, update)
	if err != nil {
		ms.Log.Printf("updating many of (%s)/(%s)/(%s) - %s\n", db, collection, updateKey, err)
	} else {
		ms.Log.Printf("updated many (%s)/(%s)/(%s) to (%v) on matches of (%s:%s)\n", db, collection, updateKey, updateValue, findKey, findValue)
	}
}

// requires custom filters and should be extended for each use-case
func (ms *MongoServer) updateMany(db, collection string, filter bson.M, update bson.D) error {
	err := ms.DB[db].Client().Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	_, err = ms.DB["network"].Collection("access_ports").UpdateMany(context.Background(), filter, update)
	return err
}
