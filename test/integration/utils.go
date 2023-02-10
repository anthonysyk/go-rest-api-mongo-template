package integration

import "github.com/anthonysyk/go-rest-api-mongo-template/api/user"

func getUser1() user.User {
	return user.User{
		ID:       "123",
		Password: "pwd1",
		IsActive: false,
		Age:      31,
		Name:     "Jane Doe",
		Data:     "data",
	}
}

func getUser2() user.User {
	return user.User{
		ID:       "321",
		Password: "pwd2",
		IsActive: false,
		Age:      41,
		Name:     "John Doe",
		Data:     "data",
	}
}
