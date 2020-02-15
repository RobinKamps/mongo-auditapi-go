package db

import (
	"context"
	"strconv"
	"strings"
	"time"

	"mongo-auditapi/pkg/config"
	"mongo-auditapi/pkg/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DataAccess provides data access functions for the persistent audit records.
type DataAccess struct {
	DbClient *mongo.Client
}

// InitializeDataAccess initializes a connection to the indicated MongoDB database.
func InitializeDataAccess(dbURL string) (*DataAccess, error) {
	dao := DataAccess{}

	c, err := mongo.NewClient(options.Client().ApplyURI(dbURL))
	if err != nil {
		return &dao, err
	}

	err = c.Connect(context.Background())
	if err != nil {
		return &dao, err
	}

	dao.DbClient = c

	return &dao, nil
}

// AuditFetcher fetches field audit trails from stored change events.
type AuditFetcher interface {
	GetFieldAuditTrail(database string, collection string, documentKey primitive.ObjectID, fieldID string) ([]model.FieldAuditRecord, error)
}

// MongoDBAuditFetcher fetches field audit trails from stored change events in a MongoDB audit collection.
type MongoDBAuditFetcher struct {
	Config config.Configuration
	Dao    *DataAccess
}

// GetFieldAuditTrail fetches audit records for a given field of a specific record in a specific MongoDB collection.
func (m *MongoDBAuditFetcher) GetFieldAuditTrail(database string, collection string, documentKey primitive.ObjectID, fieldID string) ([]model.FieldAuditRecord, error) {
	ctx := context.Background()
	coll := m.Dao.DbClient.Database(m.Config.AuditDatabase).Collection(m.Config.AuditCollection)

	var opts options.FindOptions
	opts.Sort = map[string]int{"timestamp": -1}

	cur, err := coll.Find(ctx, bson.D{{Key: "database", Value: database}, {Key: "collection", Value: collection}, {Key: "documentKey", Value: documentKey}}, &opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var recs []model.FieldAuditRecord
	for cur.Next(ctx) {
		ce := model.ChangeEvent{}
		err := cur.Decode(&ce)
		if err != nil {
			return nil, err
		}

		switch ce.OperationType {
		case "insert":
			recs = append(recs, model.FieldAuditRecord{FieldID: fieldID, FieldValue: TraverseForFieldValue(strings.Split(fieldID, "."), ce.FullDocument), UpdatedBy: ce.User, UpdatedAt: time.Unix(int64(ce.Timestamp.T), 0)})
		case "update":
			if _, ok := ce.UpdateDescription.UpdatedFields[fieldID]; ok {
				recs = append(recs, model.FieldAuditRecord{FieldID: fieldID, FieldValue: ce.UpdateDescription.UpdatedFields[fieldID], UpdatedBy: ce.User, UpdatedAt: time.Unix(int64(ce.Timestamp.T), 0)})
			}
		}
	}

	return recs, nil
}

// TraverseForFieldValue walks through the supplied document of type primitive.D to locate and fetch the value of the
// a specific field, whose field path is supplied as an array (for example, {"arr", "0", "field1"} denotes
// the field "field1" of the first element of the array "arr").
func TraverseForFieldValue(f []string, payload primitive.D) interface{} {
	f1 := f[0]
	v1 := payload.Map()[f1]
	if len(f) == 1 {
		return v1
	}

	f2 := f[1]
	if i, err := strconv.ParseInt(f2, 10, 64); err == nil {
		arr := v1.(primitive.A)
		v1 = arr[i]
	}

	if len(f) == 2 {
		return v1
	}

	return TraverseForFieldValue(f[2:], v1.(primitive.D))
}
