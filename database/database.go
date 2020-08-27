package database

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client struct
type Client struct {
	Context context.Context

	*mongo.Client
	cancel     context.CancelFunc
	connection Config
}

// Config struct
type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func NewConnection(conn Config) (*Client, error) {
	opts := options.Client()
	opts.SetHosts([]string{
		fmt.Sprintf(`%s:%d`, conn.Host, conn.Port),
	})
	if len(conn.User) > 0 || len(conn.Password) > 0 {
		opts.SetAuth(options.Credential{
			AuthSource:  conn.Database,
			Username:    conn.User,
			Password:    conn.Password,
			PasswordSet: true,
		})
	}

	ctx, canceller := context.WithCancel(context.Background())
	mCli, err := mongo.Connect(ctx, opts)

	go func() {
		select {
		case <-ctx.Done():
			mCli.Disconnect(ctx)
		}
	}()

	cli := &Client{
		Context:    ctx,
		Client:     mCli,
		cancel:     canceller,
		connection: conn,
	}

	return cli, err
}

// Disconnect function
func (s *Client) Disconnect() error {
	if s.cancel != nil {
		s.cancel()
	}

	return nil
}

// Operate function
func (s *Client) Operate(operator func(db *mongo.Database) (interface{}, error)) (interface{}, error) {
	db := s.Database(s.connection.Database)

	return operator(db)
}

// FindOne function
func (s *Client) FindOne(collection string, filter bson.D, data interface{}) error {
	db := s.Database(s.connection.Database)

	if filter == nil {
		filter = bson.D{}
	}

	cur, err := db.Collection(collection).Find(s.Context, filter)
	if err != nil {
		return err
	}

	if cur.Next(s.Context) {
		return cur.Decode(data)
	}

	return nil
}

// Find function
func (s *Client) Find(collection string, filter bson.D, data interface{}) error {
	db := s.Database(s.connection.Database)

	if filter == nil {
		filter = bson.D{}
	}

	cur, err := db.Collection(collection).Find(s.Context, filter)
	if err != nil {
		return err
	}

	return cur.All(s.Context, data)
}

// Create function
func (s *Client) Create(collection string, model interface{}) (*mongo.InsertOneResult, error) {
	db := s.Database(s.connection.Database)

	return db.Collection(collection).InsertOne(s.Context, model)
}

// CreateMany function
func (s *Client) CreateMany(collection string, models interface{}) (*mongo.InsertManyResult, error) {
	var docs []interface{}
	switch reflect.TypeOf(models).Kind() {
	case reflect.Slice, reflect.Array, reflect.Map: // TODO: need to verify wether map can be boxed as slice of interface
		ms := reflect.ValueOf(models)
		len := ms.Len()
		docs = make([]interface{}, len)
		for i := 0; i < len; i++ {
			docs[i] = ms.Index(i).Interface()
		}
		break
	default:
		docs = []interface{}{models}
	}

	db := s.Database(s.connection.Database)

	return db.Collection(collection).InsertMany(s.Context, docs)
}

// Upsert function
func (s *Client) Upsert(collection string, filter interface{}, model interface{}) error {
	b, _ := bson.Marshal(model)
	var doc primitive.D
	_ = bson.Unmarshal(b, &doc)
	m := doc.Map()

	if filter == nil || 0 == reflect.ValueOf(filter).Len() {
		id, ok := m["_id"]
		if !ok {
			id = primitive.NewObjectID()
			m["_id"] = id
		}
		filter = bson.M{"_id": id}
	}

	delete(m, "created_at")
	delete(m, "updated_at")

	update := bson.D{
		{"$set", m},
		{"$setOnInsert", bson.M{
			"created_at": time.Now(),
		}},
		{"$currentDate", bson.M{
			"updated_at": true,
		}},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)

	db := s.Database(s.connection.Database)

	rst := db.Collection(collection).FindOneAndUpdate(s.Context, filter, update, opts)
	if rst.Err() != nil {
		return rst.Err()
	}

	return rst.Decode(model)
}
