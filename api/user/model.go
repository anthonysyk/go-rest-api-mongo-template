package user

import (
	"encoding/json"
	"errors"
	"github.com/anthonysyk/go-rest-api-mongo-template/internal/jsonflex"
	"github.com/anthonysyk/go-rest-api-mongo-template/internal/password"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	ID         string       `json:"id" bson:"_id,omitempty"`
	Password   string       `json:"password,omitempty" bson:"password,omitempty"`
	IsActive   bool         `json:"isActive" bson:"isActive,omitempty"`
	Balance    string       `json:"balance" bson:"balance,omitempty"`
	Age        jsonflex.Int `json:"age" bson:"age,omitempty"`
	Name       string       `json:"name" bson:"name,omitempty"`
	Gender     string       `json:"gender" bson:"gender,omitempty"`
	Company    string       `json:"company" bson:"company,omitempty"`
	Email      string       `json:"email" bson:"email,omitempty"`
	Phone      string       `json:"phone" bson:"phone,omitempty"`
	Address    string       `json:"address" bson:"address,omitempty"`
	About      string       `json:"about" bson:"about,omitempty"`
	Registered string       `json:"registered" bson:"registered,omitempty"`
	Latitude   float64      `json:"latitude" bson:"latitude,omitempty"`
	Longitude  float64      `json:"longitude" bson:"longitude,omitempty"`
	Tags       []string     `json:"tags" bson:"tags,omitempty"`
	Friends    []struct {
		ID   int    `json:"id" bson:"id,omitempty"`
		Name string `json:"name" bson:"name,omitempty"`
	} `json:"friends" bson:"friends,omitempty"`
	Data string `json:"data" bson:"data,omitempty"`
}

func (u *User) HashPassword() error {
	if u.Password == "" {
		return errors.New("password is empty")
	}
	hash, err := password.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hash
	return nil
}

func (u *User) Create() ([]byte, error) {
	// Hash password with bcrypt
	err := u.HashPassword()
	if err != nil {
		return nil, err
	}

	// Convert User to BSON
	userBSON, err := bson.Marshal(u)
	if err != nil {
		return nil, err
	}

	return userBSON, err
}

func (u *User) Update() (bson.M, error) {
	u.ID = ""
	if u.Password != "" {
		err := u.HashPassword()
		if err != nil {
			return nil, err
		}
	}

	bytes, err := bson.Marshal(u)
	if err != nil {
		return nil, err
	}

	var update bson.M
	err = bson.Unmarshal(bytes, &update)
	if err != nil {
		return nil, err
	}

	return update, nil
}

func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		Password string `json:"password,omitempty"`
		*Alias
	}{
		Password: "",
		Alias:    (*Alias)(u),
	})
}
