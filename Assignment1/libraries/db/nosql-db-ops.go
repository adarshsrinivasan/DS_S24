package db

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/common"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	NoSQLClient *Client
)

type Client struct {
	client   *mongo.Client
	dbClient *mongo.Database
}

const (
	MongoHostEnv     = "MONGO_HOST"
	MongoPortEnv     = "MONGO_PORT"
	MongoUsernameEnv = "MONGO_USERNAME"
	MongoPasswordEnv = "MONGO_PASSWORD"
	MongoDbEnv       = "MONGO_DB"
)

func NewNoSQLClient(ctx context.Context, applicationName string) (*Client, error) {
	host := common.GetEnv(MongoHostEnv, "localhost")
	port := common.GetEnv(MongoPortEnv, "27017")
	username := common.GetEnv(MongoUsernameEnv, "admin")
	password := common.GetEnv(MongoPasswordEnv, "admin")
	dbName := common.GetEnv(MongoDbEnv, "marketplace")

	credential := options.Credential{
		Username: username,
		Password: password,
	}
	client, err := mongo.Connect(ctx,
		options.Client().
			ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)),
		options.Client().
			SetAuth(credential).
			SetAppName(applicationName))
	if err != nil {
		err = fmt.Errorf("exception while connecting to mongo DB: %v", err)
		logrus.Errorf("NewNoSQLClient: %v\n", err)
		return nil, err
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		err = fmt.Errorf("exception while pinging mongo DB: %v", err)
		logrus.Errorf("NewNoSQLClient: %v\n", err)
		return nil, err
	}

	//return client.Database(dbName).Collection(main.main.ProductTableName), nil
	return &Client{
		client:   client,
		dbClient: client.Database(dbName),
	}, nil

}

func VerifyNOSQLDatabaseConnection(ctx context.Context, client *Client) error {
	if client.dbClient == nil || client.client == nil {
		return fmt.Errorf("database connection not initialized")
	}
	if err := client.client.Ping(ctx, readpref.Primary()); err != nil {
		err = fmt.Errorf("exception while pinging mongo DB: %v", err)
		logrus.Errorf("NewNoSQLClient: %v\n", err)
		return err
	}
	return nil
}

func (client *Client) CreateCollection(ctx context.Context, collectionName string) error {
	if !client.isCollectionPresent(ctx, collectionName) {
		if err := client.dbClient.CreateCollection(ctx, collectionName); err != nil {
			err = fmt.Errorf("exception while creating collection in mongo DB: %v", err)
			logrus.Errorf("CreateCollection: %v\n", err)
			return err
		}
	}
	return nil
}

func (client *Client) isCollectionPresent(ctx context.Context, collectionName string) bool {
	coll, _ := client.dbClient.ListCollectionNames(ctx, bson.D{{"name", collectionName}})
	return len(coll) == 1
}

// InsertOne inserts a document into the specified collection.
func (client *Client) InsertOne(ctx context.Context, collectionName string, document interface{}) (int, error) {
	collection := client.dbClient.Collection(collectionName)
	_, err := collection.InsertOne(ctx, document)
	if err != nil {
		err = fmt.Errorf("exception while Inserting document in mongo DB: %v", err)
		logrus.Errorf("InsertOne: %v\n", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// FindOne finds a document in the specified collection based on the filter.
func (client *Client) FindOne(ctx context.Context, collectionName string, whereClauses []WhereClauseType, result interface{}) (int, error) {
	filter := whereClausesToFilter(whereClauses)
	collection := client.dbClient.Collection(collectionName)
	if err := collection.FindOne(ctx, filter).Decode(result); err != nil {
		err = fmt.Errorf("exception while Reading document in mongo DB: %v", err)
		logrus.Errorf("FindOne: %v\n", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// FindMany finds documents in the specified collection based on the filter.
func (client *Client) FindMany(ctx context.Context, collectionName string, whereClauses []WhereClauseType, result interface{}) (int, error) {
	filter := whereClausesToFilter(whereClauses)
	collection := client.dbClient.Collection(collectionName)
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		err = fmt.Errorf("exception while Reading document in mongo DB: %v", err)
		logrus.Errorf("FindMany: %v\n", err)
		return http.StatusInternalServerError, err
	}
	if err := cursor.All(ctx, result); err != nil {
		err = fmt.Errorf("exception while Parsing document List result in mongo DB: %v", err)
		logrus.Errorf("FindMany: %v\n", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, err
}

// UpdateOne updates a document in the specified collection based on the filter.
func (client *Client) UpdateOne(ctx context.Context, collectionName string, whereClauses []WhereClauseType, data interface{}) (int, error) {
	filter := whereClausesToFilter(whereClauses)

	bsonDoc, _ := toDoc(data)
	update := bson.D{{"$set", bsonDoc}}
	//update, err := buildUpdateModel(data)
	//if err != nil {
	//	err = fmt.Errorf("exception while build Update query in mongo DB: %v", err)
	//	logrus.Errorf("UpdateOne: %v\n", err)
	//	return http.StatusInternalServerError, err
	//}
	collection := client.dbClient.Collection(collectionName)
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		err = fmt.Errorf("exception while Updating document in mongo DB: %v", err)
		logrus.Errorf("UpdateOne: %v\n", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// DeleteOne deletes a document from the specified collection based on the filter.
func (client *Client) DeleteOne(ctx context.Context, collectionName string, whereClauses []WhereClauseType) (int, error) {
	filter := whereClausesToFilter(whereClauses)
	collection := client.dbClient.Collection(collectionName)
	if _, err := collection.DeleteOne(ctx, filter); err != nil {
		err = fmt.Errorf("exception while Deleting document in mongo DB: %v", err)
		logrus.Errorf("UpdateOne: %v\n", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func whereClausesToFilter(whereClauses []WhereClauseType) bson.D {
	filter := bson.D{}

	for _, wc := range whereClauses {
		fieldName := wc.ColumnName

		switch wc.RelationType {
		case EQUAL:
			filter = append(filter, bson.E{Key: fieldName, Value: wc.ColumnValue})
		case NOT_EQUAL:
			filter = append(filter, bson.E{Key: fieldName, Value: bson.D{{"$ne", wc.ColumnValue}}})
		case IN:
			filter = append(filter, bson.E{Key: fieldName, Value: bson.D{{"$in", wc.ColumnValue}}})
		case NOT_IN:
			filter = append(filter, bson.E{Key: fieldName, Value: bson.D{{"$nin", wc.ColumnValue}}})
		case IS:
			filter = append(filter, bson.E{Key: fieldName, Value: bson.D{{"$exists", wc.ColumnValue != nil}}})
		case LIKE:
			filter = append(filter, bson.E{Key: fieldName, Value: bson.D{{"$regex", wc.ColumnValue}, {"$options", "i"}}})
		case GT:
			filter = append(filter, bson.E{Key: fieldName, Value: bson.D{{"$gt", wc.ColumnValue}}})
		case LT:
			filter = append(filter, bson.E{Key: fieldName, Value: bson.D{{"$lt", wc.ColumnValue}}})
		default:
			// Unsupported relation type, ignore or handle accordingly
		}
	}

	return filter
}

// BuildUpdateModel builds a BSON update model based on the provided interface.
func buildUpdateModel(data interface{}) (bson.D, error) {
	updateModel := bson.D{}

	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct")
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i).Name
		value := val.Field(i).Interface()

		if reflect.DeepEqual(value, reflect.Zero(val.Field(i).Type()).Interface()) {
			continue
		}

		updateModel = append(updateModel, bson.E{Key: field, Value: value})
	}

	return updateModel, nil
}

//	func (product *ProductModel) updateProductByID() (int, error) {
//		if err := db.VerifyNOSQLDatabaseConnection(db.Client); err != nil {
//			return http.StatusInternalServerError, fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Update", ProductTableName, err)
//		}
//		product.UpdatedAt = time.Now()
//		bsonCourt, _ := toDoc(product)
//		update := bson.D{{"$set", bsonCourt}}
//		if _, err := db.Client.UpdateByID(common.Ctx, product.ID, update); err != nil {
//			return http.StatusInternalServerError, fmt.Errorf("exception while performing %s Operation on Table: %s. %v", "Update", ProductTableName, err)
//		}
//		return http.StatusOK, nil
//	}
func toDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}
