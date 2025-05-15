package mongo_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rshelekhov/go-db/mongo"
	"github.com/rshelekhov/go-db/mongo/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testDB *testutil.TestDB
	conn   mongo.ConnectionManager
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	testDB, err = testutil.NewTestDB(ctx)
	if err != nil {
		panic(err)
	}

	conn, err = mongo.NewConnection(ctx, testDB.URI(), "testdb",
		mongo.WithTimeout(time.Second*5),
		mongo.WithServerAPI("1"),
	)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	conn.Close(ctx)
	testDB.Close(ctx)

	os.Exit(code)
}

type TestDoc struct {
	Name  string `bson:"name" fake:"{firstname}"`
	Value int    `bson:"value" fake:"{number:1,100}"`
}

func TestInsertOne(t *testing.T) {
	ctx := context.Background()
	coll := "test_insert_one"

	doc := TestDoc{}
	gofakeit.Struct(&doc)

	result, err := conn.InsertOne(ctx, coll, doc)
	require.NoError(t, err)
	assert.NotNil(t, result.InsertedID)

	// With options
	doc2 := TestDoc{}
	gofakeit.Struct(&doc2)
	opts := options.InsertOne().SetBypassDocumentValidation(true)
	result, err = conn.InsertOne(ctx, coll, doc2, opts)
	require.NoError(t, err)
	assert.NotNil(t, result.InsertedID)
}

func TestInsertMany(t *testing.T) {
	ctx := context.Background()
	coll := "test_insert_many"

	docs := make([]TestDoc, 2)
	for i := range docs {
		gofakeit.Struct(&docs[i])
	}
	docsAny := make([]any, len(docs))
	for i, doc := range docs {
		docsAny[i] = doc
	}

	result, err := conn.InsertMany(ctx, coll, docsAny)
	require.NoError(t, err)
	assert.Len(t, result.InsertedIDs, 2)

	// With options
	docs2 := make([]TestDoc, 2)
	for i := range docs2 {
		gofakeit.Struct(&docs2[i])
	}
	docs2Any := make([]any, len(docs2))
	for i, doc := range docs2 {
		docs2Any[i] = doc
	}
	opts := options.InsertMany().SetOrdered(false)
	result, err = conn.InsertMany(ctx, coll, docs2Any, opts)
	require.NoError(t, err)
	assert.Len(t, result.InsertedIDs, 2)
}

func TestFindOne(t *testing.T) {
	ctx := context.Background()
	coll := "test_find_one"

	doc := TestDoc{}
	gofakeit.Struct(&doc)
	_, err := conn.InsertOne(ctx, coll, doc)
	require.NoError(t, err)

	var result TestDoc
	err = conn.FindOne(ctx, coll, bson.M{"name": doc.Name}, &result)
	require.NoError(t, err)
	assert.Equal(t, doc.Name, result.Name)
	assert.Equal(t, doc.Value, result.Value)

	// With options
	opts := options.FindOne().SetMaxTime(time.Second)
	var result2 TestDoc
	err = conn.FindOne(ctx, coll, bson.M{"name": doc.Name}, &result2, opts)
	require.NoError(t, err)
	assert.Equal(t, doc.Name, result2.Name)
}

func TestFind(t *testing.T) {
	ctx := context.Background()
	coll := "test_find"

	docs := make([]TestDoc, 2)
	for i := range docs {
		gofakeit.Struct(&docs[i])
	}
	docsAny := make([]any, len(docs))
	for i, doc := range docs {
		docsAny[i] = doc
	}
	_, err := conn.InsertMany(ctx, coll, docsAny)
	require.NoError(t, err)

	cursor, err := conn.Find(ctx, coll, bson.M{})
	require.NoError(t, err)
	defer cursor.Close(ctx)

	var results []TestDoc
	err = cursor.All(ctx, &results)
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// With options
	opts := options.Find().SetLimit(1)
	cursor, err = conn.Find(ctx, coll, bson.M{}, opts)
	require.NoError(t, err)
	defer cursor.Close(ctx)

	var results2 []TestDoc
	err = cursor.All(ctx, &results2)
	require.NoError(t, err)
	assert.Len(t, results2, 1)
}

func TestUpdateOne(t *testing.T) {
	ctx := context.Background()
	coll := "test_update_one"

	doc := TestDoc{}
	gofakeit.Struct(&doc)
	_, err := conn.InsertOne(ctx, coll, doc)
	require.NoError(t, err)

	newValue := gofakeit.Number(1, 100)
	update := bson.M{"$set": bson.M{"value": newValue}}
	result, err := conn.UpdateOne(ctx, coll, bson.M{"name": doc.Name}, update)
	require.NoError(t, err)
	assert.Equal(t, int64(1), result.ModifiedCount)

	// With options
	opts := options.Update().SetUpsert(true)
	doc2 := TestDoc{}
	gofakeit.Struct(&doc2)
	update2 := bson.M{"$set": bson.M{"value": newValue}}
	result, err = conn.UpdateOne(ctx, coll, bson.M{"name": doc2.Name}, update2, opts)
	require.NoError(t, err)
	assert.Equal(t, int64(0), result.ModifiedCount)
	assert.Equal(t, int64(1), result.UpsertedCount)
}

func TestUpdateMany(t *testing.T) {
	ctx := context.Background()
	coll := "test_update_many"

	docs := make([]TestDoc, 2)
	for i := range docs {
		gofakeit.Struct(&docs[i])
		docs[i].Value = 1
	}
	docsAny := make([]any, len(docs))
	for i, doc := range docs {
		docsAny[i] = doc
	}
	_, err := conn.InsertMany(ctx, coll, docsAny)
	require.NoError(t, err)

	newValue := gofakeit.Number(1, 100)
	update := bson.M{"$set": bson.M{"value": newValue}}
	result, err := conn.UpdateMany(ctx, coll, bson.M{"value": 1}, update)
	require.NoError(t, err)
	assert.Equal(t, int64(2), result.ModifiedCount)

	// With options
	opts := options.Update().SetUpsert(true)
	update2 := bson.M{"$set": bson.M{"value": newValue}}
	result, err = conn.UpdateMany(ctx, coll, bson.M{"value": 999}, update2, opts)
	require.NoError(t, err)
	assert.Equal(t, int64(0), result.ModifiedCount)
	assert.Equal(t, int64(1), result.UpsertedCount)
}

func TestDeleteOne(t *testing.T) {
	ctx := context.Background()
	coll := "test_delete_one"

	doc := TestDoc{}
	gofakeit.Struct(&doc)
	_, err := conn.InsertOne(ctx, coll, doc)
	require.NoError(t, err)

	result, err := conn.DeleteOne(ctx, coll, bson.M{"name": doc.Name})
	require.NoError(t, err)
	assert.Equal(t, int64(1), result.DeletedCount)

	// With options
	opts := options.Delete().SetCollation(&options.Collation{Locale: "en"})
	result, err = conn.DeleteOne(ctx, coll, bson.M{"name": "nonexistent"}, opts)
	require.NoError(t, err)
	assert.Equal(t, int64(0), result.DeletedCount)
}

func TestDeleteMany(t *testing.T) {
	ctx := context.Background()
	coll := "test_delete_many"

	docs := make([]TestDoc, 2)
	for i := range docs {
		gofakeit.Struct(&docs[i])
		docs[i].Value = 1
	}
	docsAny := make([]any, len(docs))
	for i, doc := range docs {
		docsAny[i] = doc
	}
	_, err := conn.InsertMany(ctx, coll, docsAny)
	require.NoError(t, err)

	result, err := conn.DeleteMany(ctx, coll, bson.M{"value": 1})
	require.NoError(t, err)
	assert.Equal(t, int64(2), result.DeletedCount)

	// With options
	opts := options.Delete().SetCollation(&options.Collation{Locale: "en"})
	result, err = conn.DeleteMany(ctx, coll, bson.M{"value": 999}, opts)
	require.NoError(t, err)
	assert.Equal(t, int64(0), result.DeletedCount)
}

func TestCountDocuments(t *testing.T) {
	ctx := context.Background()
	coll := "test_count"

	docs := make([]TestDoc, 3)
	for i := range docs {
		gofakeit.Struct(&docs[i])
		if i < 2 {
			docs[i].Value = 1
		}
	}
	docsAny := make([]any, len(docs))
	for i, doc := range docs {
		docsAny[i] = doc
	}
	_, err := conn.InsertMany(ctx, coll, docsAny)
	require.NoError(t, err)

	count, err := conn.CountDocuments(ctx, coll, bson.M{"value": 1})
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)

	// With options
	opts := options.Count().SetMaxTime(time.Second)
	count, err = conn.CountDocuments(ctx, coll, bson.M{"value": docs[2].Value}, opts)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestAggregate(t *testing.T) {
	ctx := context.Background()
	coll := "test_aggregate"

	docs := make([]TestDoc, 3)
	for i := range docs {
		gofakeit.Struct(&docs[i])
		if i > 0 {
			docs[i].Value = docs[1].Value
		}
	}
	docsAny := make([]any, len(docs))
	for i, doc := range docs {
		docsAny[i] = doc
	}
	_, err := conn.InsertMany(ctx, coll, docsAny)
	require.NoError(t, err)

	pipeline := []bson.M{
		{"$match": bson.M{"value": docs[1].Value}},
		{"$group": bson.M{
			"_id":   "$value",
			"count": bson.M{"$sum": 1},
		}},
	}

	cursor, err := conn.Aggregate(ctx, coll, pipeline)
	require.NoError(t, err)
	defer cursor.Close(ctx)

	var results []bson.M
	err = cursor.All(ctx, &results)
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, int32(2), results[0]["count"])

	// With options
	opts := options.Aggregate().SetMaxTime(time.Second)
	cursor, err = conn.Aggregate(ctx, coll, pipeline, opts)
	require.NoError(t, err)
	defer cursor.Close(ctx)

	results = nil
	err = cursor.All(ctx, &results)
	require.NoError(t, err)
	assert.Len(t, results, 1)
}
