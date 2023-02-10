//go:build integration

package integration

import (
	"context"
	"github.com/anthonysyk/go-rest-api-mongo-template/api/user"
	"github.com/anthonysyk/go-rest-api-mongo-template/internal/db"
	"github.com/anthonysyk/go-rest-api-mongo-template/internal/password"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository(t *testing.T) {
	ctx := context.Background()
	client, err := db.NewMongoDB(ctx, "mongodb://root:example@localhost:27017")
	if err != nil {
		t.Fatal(err)
	}
	dbClient := client.Database(t.Name())
	repo := user.NewRepository(dbClient)
	defer dbClient.Drop(ctx)

	user1 := getUser1()
	user2 := getUser2()

	// Test Create User
	err = repo.CreateUser(ctx, &user1)
	assert.NoError(t, err)
	err = repo.CreateUser(ctx, &user2)
	assert.NoError(t, err)

	// Test List Users
	users, err := repo.ListUsers(ctx)
	assert.NoError(t, err)
	assert.Len(t, users, 2)

	// Test Get User
	dbUser1, err := repo.GetUser(ctx, user1.ID)
	assert.NoError(t, err)
	assert.Equal(t, user1, *dbUser1)
	dbUser2, err := repo.GetUser(ctx, user2.ID)
	assert.NoError(t, err)
	assert.Equal(t, user2, *dbUser2)

	// Test Update User
	updatedName := "John Smith"
	err = repo.UpdateUser(ctx, user2.ID, &user.User{Name: updatedName})
	assert.NoError(t, err)
	dbUser2, err = repo.GetUser(ctx, user2.ID)
	assert.NoError(t, err)
	assert.Equal(t, dbUser2.Name, updatedName)
	assert.Equal(t, dbUser2.Age, user2.Age)

	// Test Delete User
	err = repo.DeleteUser(ctx, user2.ID)
	assert.NoError(t, err)
	users, err = repo.ListUsers(ctx)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
}

func TestUserRepository_Password(t *testing.T) {
	ctx := context.Background()
	client, err := db.NewMongoDB(ctx, "mongodb://root:example@localhost:27017")
	if err != nil {
		t.Fatal(err)
	}
	dbClient := client.Database(t.Name())
	repo := user.NewRepository(dbClient)
	defer dbClient.Drop(ctx)

	user1 := getUser1()
	pwd := user1.Password

	// Create user
	err = repo.CreateUser(ctx, &user1)
	assert.NoError(t, err)

	// Verify password stored is hashed with bcrypt
	dbUser1, err := repo.GetUser(ctx, user1.ID)
	assert.NoError(t, err)
	assert.NoError(t, password.VerifyPassword(dbUser1.Password, pwd))

	// Verify password update is still hashed with bcrypt
	updatedPassword := "newpassword"
	err = repo.UpdateUser(ctx, user1.ID, &user.User{Password: updatedPassword})
	dbUser1, err = repo.GetUser(ctx, user1.ID)
	assert.NoError(t, err)
	assert.NoError(t, password.VerifyPassword(dbUser1.Password, updatedPassword))
}
