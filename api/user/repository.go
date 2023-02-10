package user

import (
	"context"
	"github.com/anthonysyk/go-rest-api-mongo-template/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	usersCollection *mongo.Collection
}

func NewRepository(db *mongo.Database) Repository {
	return Repository{usersCollection: db.Collection("users")}
}

func (r Repository) CreateUser(ctx context.Context, user *User) error {
	createBSON, err := user.Create()
	if err != nil {
		return err
	}

	// Insert user into MongoDB
	if _, err := r.usersCollection.InsertOne(ctx, createBSON); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return db.ErrIsDuplicateKey
		}
		return err
	}
	return nil
}

func (r Repository) DeleteUser(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := r.usersCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) UpdateUser(ctx context.Context, id string, user *User) error {
	updateBSON, err := user.Update()
	if err != nil {
		return err
	}

	_, err = r.usersCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.D{{Key: "$set", Value: updateBSON}})
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) GetUser(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.usersCollection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r Repository) ListUsers(ctx context.Context) ([]User, error) {
	cursor, err := r.usersCollection.Find(ctx, bson.D{{}}, nil)
	if err != nil {
		return nil, err
	}

	var users []User
	for cursor.Next(ctx) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
