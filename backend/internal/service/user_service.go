package service

import (
	"errors"
	"time"

	"github.com/dimasrizkyfebrian/coursify/internal/config"
	"github.com/dimasrizkyfebrian/coursify/internal/dto"
	"github.com/dimasrizkyfebrian/coursify/internal/model"
	"github.com/dimasrizkyfebrian/coursify/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo   *repository.UserRepository
	Config *config.Config
}

func NewUserService(repo *repository.UserRepository, cfg *config.Config) *UserService {
	return &UserService{Repo: repo, Config: cfg}
}

func (s *UserService) Register(req *model.User) error {
	return s.Repo.CreateUser(req)
}

func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	if user.Status != "active" {
		return "", errors.New("account is not active")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.Config.JWTSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *UserService) GetUserByID(id string) (*model.User, error) {
	return s.Repo.GetUserByID(id)
}

func (s *UserService) GetPendingUsers() ([]model.User, error) {
	return s.Repo.GetUsersByStatus("pending")
}

func (s *UserService) ApproveUser(id string) error {
	return s.Repo.UpdateUserStatus(id, "active")
}

func (s *UserService) RejectUser(id string) error {
	return s.Repo.UpdateUserStatus(id, "rejected")
}

func (s *UserService) GetPendingUserCount() (int, error) {
	return s.Repo.GetPendingUserCount()
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	return s.Repo.GetAllUsers()
}

func (s *UserService) UpdateUser(user *model.User) error {
	return s.Repo.UpdateUser(user)
}

func (s *UserService) DeleteUser(id string) error {
	return s.Repo.DeleteUser(id)
}

func (s *UserService) GetUserStats() (map[string]int, error) {
	return s.Repo.GetUserStats()
}

func (s *UserService) UpdateUserAvatarURL(id, url string) error {
	return s.Repo.UpdateUserAvatarURL(id, url)
}

// MapUserToResponse helper
func (s *UserService) MapUserToResponse(user *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		AvatarURL: &user.AvatarURL.String,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (s *UserService) MapUsersToResponse(users []model.User) []dto.UserResponse {
	var response []dto.UserResponse
	for _, user := range users {
		response = append(response, s.MapUserToResponse(&user))
	}
	return response
}
