package user

import (
	"context"
	"github.com/anthonysyk/go-rest-api-mongo-template/internal/auth"
	"github.com/anthonysyk/go-rest-api-mongo-template/internal/fs"
	"github.com/anthonysyk/go-rest-api-mongo-template/internal/password"
	"log"
	"sync"
	"time"
)

type Service struct {
	fs             fs.Client
	userRepository Repository
}

func NewService(userRepository Repository, fsClient fs.Client) *Service {
	return &Service{userRepository: userRepository, fs: fsClient}
}

func (s Service) Login(ctx context.Context, id, pwd, secret string) (string, error) {
	user, err := s.userRepository.GetUser(ctx, id)
	if err != nil {
		return "", err
	}

	err = password.VerifyPassword(user.Password, pwd)
	if err != nil {
		return "", err
	}

	return auth.CreateToken(secret, id, 24*time.Hour)
}

func (s Service) GetUser(ctx context.Context, id string) (*User, error) {
	return s.userRepository.GetUser(ctx, id)
}

func (s Service) ListUsers(ctx context.Context) ([]User, error) {
	return s.userRepository.ListUsers(ctx)
}

func (s Service) CreateUsers(ctx context.Context, users []*User) int {
	wg := sync.WaitGroup{}
	total := 0
	for _, u := range users {
		wg.Add(1)
		go func(user *User) {
			defer wg.Done()
			err := s.CreateUser(ctx, user)
			if err != nil {
				log.Print(err.Error())
				return
			}
			total++
		}(u)
	}
	wg.Wait()
	return total
}

func (s Service) CreateUser(ctx context.Context, user *User) error {
	// Create user un database
	err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	// Create data file
	err = s.fs.Write(user.ID, []byte(user.Data))
	if err != nil {
		return err
	}

	return nil
}

func (s Service) DeleteUser(ctx context.Context, id string) error {
	// delete mongodb user
	err := s.userRepository.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	// delete file associated
	err = s.fs.Remove(id)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) UpdateUser(ctx context.Context, id string, user *User) error {
	dbUser, err := s.userRepository.GetUser(ctx, id)
	if err != nil {
		return err
	}

	err = s.userRepository.UpdateUser(ctx, id, user)
	if err != nil {
		return err
	}

	if dbUser.Data != user.Data {
		err := s.fs.Write(id, []byte(user.Data))
		if err != nil {
			return err
		}
	}

	return nil
}
