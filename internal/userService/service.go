package userService

import (
	"errors"
	"strings"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)

type Service struct {
	repo Repository
}

// NewService creates user service with provided repository.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Create validates and creates a user.
func (s *Service) Create(email, password string) (*User, error) {
	email = strings.TrimSpace(email)
	if !isValidEmail(email) || strings.TrimSpace(password) == "" {
		return nil, ErrInvalidInput
	}

	u := &User{
		Email:    email,
		Password: password,
	}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

// List returns all users.
func (s *Service) List() ([]User, error) {
	return s.repo.List()
}

// Patch updates user fields by id.
func (s *Service) Patch(id uint, email, password *string) (*User, error) {
	u, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	if email != nil {
		trimmedEmail := strings.TrimSpace(*email)
		if !isValidEmail(trimmedEmail) {
			return nil, ErrInvalidInput
		}
		u.Email = trimmedEmail
	}
	if password != nil {
		if strings.TrimSpace(*password) == "" {
			return nil, ErrInvalidInput
		}
		u.Password = *password
	}

	if err := s.repo.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}

// Delete removes user by id.
func (s *Service) Delete(id uint) error {
	rows, err := s.repo.DeleteByID(id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func isValidEmail(email string) bool {
	if email == "" {
		return false
	}

	at := strings.Index(email, "@")
	if at <= 0 || at >= len(email)-1 {
		return false
	}

	domain := email[at+1:]
	return strings.Contains(domain, ".")
}
