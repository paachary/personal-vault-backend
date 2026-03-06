package db

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func collectionExists(database *mongo.Database, collectionName string) error {
	cursor, err := database.ListCollections(context.TODO(), bson.D{{Key: "name", Value: collectionName}})
	if err != nil {
		return errors.New("error listing collections: " + err.Error())
	}

	defer cursor.Close(context.TODO())

	if !cursor.Next(context.TODO()) {
		return errors.New("collection " + collectionName + " does not exist")
	}
	return nil
}

func createCollection(collectionName string) error {
	err := MongoClient.Database(DatabaseName).CreateCollection(context.TODO(), collectionName)
	if err != nil {
		return errors.New("error creating collection: " + err.Error())
	}
	fmt.Printf("Collection %v created successfully.\n", collectionName)
	return nil
}

func createCollections() error {
	if err := collectionExists(MongoClient.Database(DatabaseName), CollectionName); err != nil {
		err = createCollection(CollectionName)
		if err != nil {
			return errors.New("Error creating collection " + CollectionName + ": " + err.Error())
		}
	}

	if err := collectionExists(MongoClient.Database(DatabaseName), MFACollectionName); err != nil {
		err = createCollection(MFACollectionName)
		if err != nil {
			return errors.New("Error creating collection " + MFACollectionName + ": " + err.Error())
		}
	}
	return nil
}

func itemExists(col *mongo.Collection, filter any) (bool, error) {
	err := col.FindOne(context.TODO(), filter).Decode(&bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil // Document does not exist
		}
		return false, err // An actual error occurred
	}
	return true, nil // Document exists
}

// Insert the main personal data document
func InsertOne(document any, key, value string) (*mongo.InsertOneResult, error) {

	collection := MongoClient.Database(DatabaseName).Collection(CollectionName)

	// check for username already exists in the database collection

	filter := bson.M{key: value}

	exists, err := itemExists(collection, filter)

	if err != nil {
		return nil, errors.New("Record already exists with " + key + ". :" + err.Error())

	}

	if !exists {

		_, err := collection.InsertOne(context.TODO(), document)

		if err != nil {
			return nil, errors.New("Failed to insert new record :" + err.Error())

		}
	} else {
		return nil, errors.New("Record already exists with " + key + ".")

	}

	return nil, nil

}

func RetrieveSingleRecord[T any](filter any) (*T, error) {

	var result T

	collection := MongoClient.Database(DatabaseName).Collection(CollectionName)

	// check for username already exists in the database collection

	exists, err := itemExists(collection, filter)
	if err != nil {
		return nil, errors.New("Error occurred during fetching data. Please check and resubmit the query. :" + err.Error()) // An actual error occurred
	}

	if exists {

		err = collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {

			return nil, errors.New("Error occurred during fetching data. Please check and resubmit the query.") // An actual error occurred
		}

	} else {
		return nil, errors.New("Record does not exists. Please check and resubmit the query.")

	}

	return &result, nil

}

func FindRecordCount(keyFilter bson.M, attribute string) (bool, error) {
	collection := MongoClient.Database(DatabaseName).Collection(CollectionName)

	// check for username already exists in the database collection

	exists, err := itemExists(collection, keyFilter)
	if err != nil {
		return false, errors.New("Error occurred during fetching data. Please check and resubmit the query. :" + err.Error()) // An actual error occurred
	}

	if exists {

		// Filter: Check if the array with attribute exists and has size > 0

		// Filter: Check if pran_details exists
		filterExits := bson.M{
			"$and": []bson.M{
				keyFilter,
				{attribute: bson.M{"$exists": true}}},
		}

		// Optional: Also ensure it's an array
		filterExits["$expr"] = bson.M{
			"$isArray": fmt.Sprintf("$%s", attribute),
		}

		var result bson.M
		err := collection.FindOne(context.TODO(), filterExits).Decode(&result)
		if err != nil {
			return false, nil
		}

		filter := bson.M{
			"$and": []bson.M{
				keyFilter,
				{"$expr": bson.M{
					"$gt": []any{bson.M{"$size": fmt.Sprintf("$%s", attribute)}, 0},
				}},
			}}

		// Use CountDocuments to return boolean
		count, err := collection.CountDocuments(context.TODO(), filter)

		if err != nil {
			return false, errors.New("Count operation failed for the provided condition. Please check your filter conditions and resubmit the operation.")
		}

		exists := count > 0

		return exists, nil

	} else {
		return false, errors.New("Select operation failed for the provided filter conditions. Please check your filter conditions and resubmit the operation.")

	}
}

func AddOne(keyFilter bson.M, attribute string, data bson.M, updateFilters map[string]any) error {

	collection := MongoClient.Database(DatabaseName).Collection(CollectionName)

	// check for username already exists in the database collection

	exists, err := itemExists(collection, keyFilter)
	if err != nil {
		return errors.New("Error occurred during fetching data. Please check and resubmit the update. :" + err.Error()) // An actual error occurred
	}

	if exists {

		update := bson.M{
			"$push": bson.M{attribute: data},
		}

		result, err := collection.UpdateOne(context.TODO(), updateFilters, update)
		if err != nil {
			return errors.New("Error updating the record: " + err.Error())
		}

		if result.ModifiedCount == 0 {
			return errors.New("Either the element being added already exists or the document does not exist. Please retry with proper data.")
		} else {
			return nil
		}

	} else {
		return errors.New("Update failed for the provided filter conditions. Please check your filter conditions and resubmit the update.")
	}

}

func UpdateOneForArray(keyFilter bson.M, arrayField string, matchCondition bson.M, updates map[string]any) error {

	collection := MongoClient.Database(DatabaseName).Collection(CollectionName)

	// check for username already exists in the database collection

	exists, err := itemExists(collection, keyFilter)
	if err != nil {
		return errors.New("Error updating the record: " + err.Error())
	}

	if exists {

		// Build update document dynamically
		setFields := bson.M{}
		for key, value := range updates {
			// Example: addresses.$[elem].city
			setFields[fmt.Sprintf("%s.$[elem].%s", arrayField, key)] = value
		}

		update := bson.M{"$set": setFields}

		updateOpts := options.UpdateOne().
			SetArrayFilters([]any{matchCondition}).
			SetUpsert(false)

		result, err := collection.UpdateOne(context.TODO(), keyFilter, update, updateOpts)

		if err != nil {
			return errors.New("Error updating the record: " + err.Error())
		}

		if result.ModifiedCount == 0 {
			return errors.New("No matching array element found. Nothing updated.")
		}

	} else {
		return errors.New("Update failed for the provided filter conditions. Please check your filter conditions and resubmit the update.")
	}
	return nil

}

func UpdateOne(keyFilter bson.M, update bson.M) error {
	collection := MongoClient.Database(DatabaseName).Collection(CollectionName)
	exists, err := itemExists(collection, keyFilter)
	if err != nil {
		return errors.New("Error updating the record: " + err.Error())
	}

	if exists {

		result, err := collection.UpdateOne(context.TODO(), keyFilter, update)

		if err != nil {
			return errors.New("Error updating the record: " + err.Error())
		}

		if result.ModifiedCount == 0 {
			return errors.New("No matching array element found. Nothing updated.")
		}

	} else {
		return errors.New("Update failed for the provided filter conditions. Please check your filter conditions and resubmit the update.")
	}
	return nil

}

func DeleteOne(keyFilter bson.M, arrayField string, matchCondition bson.M) error {

	collection := MongoClient.Database(DatabaseName).Collection(CollectionName)

	// check for username already exists in the database collection

	exists, err := itemExists(collection, keyFilter)
	if err != nil {
		return errors.New("Error deleting the record: " + err.Error())
	}

	if exists {

		update := bson.M{
			"$pull": bson.M{
				arrayField: matchCondition, // Remove matching object from array
			},
		}

		result, err := collection.UpdateOne(context.TODO(), keyFilter, update)

		if err != nil {
			return errors.New("Error deleting the record: " + err.Error())
		}

		if result.ModifiedCount == 0 {
			return errors.New("No matching array element found. Nothing deleted.")
		}

	} else {
		return errors.New("Delete failed for the provided filter conditions. Please check your filter conditions and resubmit the delete.")
	}
	return nil

}

func GetAllUsersData[T any](keyFilter bson.M, filter bson.M) (*[]T, error) {

	var result []T

	collection := MongoClient.Database(DatabaseName).Collection(CollectionName)

	exists, err := itemExists(collection, keyFilter)
	if err != nil {
		return nil, errors.New("No matching user data found.")
	}

	if exists {
		cursor, err := collection.Find(context.TODO(), filter)

		if err != nil {
			return nil, errors.New("Error fetching the cursor: " + err.Error())
		}

		defer cursor.Close(context.TODO())

		if err = cursor.All(context.TODO(), &result); err != nil {
			return nil, errors.New("Error fetching all users data: " + err.Error())
		}

		return &result, nil
	} else {
		return nil, errors.New("No matching user data found.")
	}

}

// Find the MFA record for the user
func GetMfaRecord(keyFilter bson.M, condition bson.M) (*bson.M, error) {
	collection := MongoClient.Database(DatabaseName).Collection(MFACollectionName)

	combinedFilter := keyFilter

	if condition != nil {
		// Combine filters using $and
		combinedFilter = bson.M{
			"$and": []bson.M{keyFilter, condition},
		}
	}

	var result bson.M

	err := collection.FindOne(context.TODO(), combinedFilter).Decode(&result)
	if err != nil {
		return nil, errors.New("Error occurred during fetching MFA record. Please check and resubmit the query. :" + err.Error()) // An actual error occurred
	}

	return &result, nil

}

// Insert / Update MFA code for the user
func UpsertMfaRecord(keyFilter bson.M, condition bson.M, document any) error {

	collection := MongoClient.Database(DatabaseName).Collection(MFACollectionName)

	combinedFilter := keyFilter

	updateOpts := options.UpdateOne().SetUpsert(true) // Default to upsert with just keyFilter

	if condition != nil {
		// Combine filters using $and
		combinedFilter = bson.M{
			"$and": []bson.M{keyFilter, condition},
		}

		updateOpts = options.UpdateOne().
			SetUpsert(false)
	}

	_, err := collection.UpdateOne(context.TODO(), combinedFilter, bson.M{"$set": document}, updateOpts)
	if err != nil {
		return errors.New("Failed to insert MFA code: " + err.Error())
	}

	return nil

}

func DeleteMfaRecord(keyFilter bson.M, condition bson.M) error {

	collection := MongoClient.Database(DatabaseName).Collection(MFACollectionName)

	exists, err := itemExists(collection, keyFilter)
	if err != nil {
		return errors.New("Error deleting the record: " + err.Error())
	}

	if exists {

		// Combine filters using $and
		combinedFilter := bson.M{
			"$and": []bson.M{keyFilter, condition},
		}

		result, err := collection.DeleteOne(context.TODO(), combinedFilter)
		if err != nil {
			return errors.New("Error deleting the record: " + err.Error())
		}

		if result.DeletedCount == 0 {
			return errors.New("No matching document found. Nothing deleted.")
		}

	} else {
		return errors.New("Delete failed for the provided filter conditions. Please check your filter conditions and resubmit the delete.")
	}
	return nil

}
