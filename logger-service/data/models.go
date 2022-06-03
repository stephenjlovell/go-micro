package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func NewLogEntry(entry LogEntry) LogEntry {
	t := time.Now()
	return LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: t,
		UpdatedAt: t,
	}
}

type Models struct {
	LogEntry LogEntry
}

func New(c *mongo.Client) Models {
	client = c
	return Models{
		LogEntry: LogEntry{},
	}

}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), NewLogEntry(entry))
	if err != nil {
		log.Println("error inserting into logs:", err)
		return err
	}
	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFunc() // cancel if request times out
	collection := client.Database("logs").Collection("logs")

	opts := options.Find().SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("unable to find all docs:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry
	for cursor.Next(ctx) {
		item := &LogEntry{}
		err := cursor.Decode(item)
		if err != nil {
			log.Println("error decoding log entry:", err)
			return nil, err
		}
		// FIXME inefficient
		logs = append(logs, item)
	}
	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	// TODO: duplicated in LogEntry.All()
	ctx, cancelFunc := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFunc() // cancel if request times out
	collection := client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	entry := &LogEntry{}
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (l *LogEntry) DropCollection() error {
	// TODO: duplicated in LogEntry.All()
	ctx, cancelFunc := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFunc() // cancel if request times out
	collection := client.Database("logs").Collection("logs")

	return collection.Drop(ctx)
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	// TODO: duplicated in All()
	ctx, cancelFunc := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFunc() // cancel if request times out
	collection := client.Database("logs").Collection("logs")
	// TODO: duplicated in GetOne()
	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}
	result, err := collection.UpdateOne(ctx, bson.M{"_id": docID}, bson.D{{"$set", bson.D{{"name", l.Name}, {"data", l.Data}}}})
	if err != nil {
		return nil, err
	}
	return result, nil
}
