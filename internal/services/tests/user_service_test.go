package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"topup-backend/internal/config"
	"topup-backend/internal/models"
	"topup-backend/internal/services"
)

// UserServiceTestSuite defines the test suite for UserService
type UserServiceTestSuite struct {
	suite.Suite
	db          *gorm.DB
	userService *services.UserService
	cfg         *config.Config
}

// SetupSuite runs once before all tests
func (suite *UserServiceTestSuite) SetupSuite() {
	// Setup in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)

	// Run migrations
	err = db.AutoMigrate(&models.User{})
	assert.NoError(suite.T(), err)

	// Setup config
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:        "test-secret",
			AccessExpire:  15 * 60,          // 15 minutes in seconds
			RefreshExpire: 7 * 24 * 60 * 60, // 7 days in seconds
		},
	}

	suite.db = db
	suite.cfg = cfg
	suite.userService = services.NewUserService(db, cfg)
}

// TearDownSuite runs once after all tests
func (suite *UserServiceTestSuite) TearDownSuite() {
	// Cleanup if needed
}

// SetupTest runs before each test
func (suite *UserServiceTestSuite) SetupTest() {
	// Clean up database before each test
	suite.db.Exec("DELETE FROM users")
}

// TestUserRegistration tests user registration
func (suite *UserServiceTestSuite) TestUserRegistration() {
	// Test data
	req := &services.RegisterRequest{
		Name:     "John Doe",
		Phone:    "081234567890",
		Password: "password123",
		Email:    "john@example.com",
	}

	// Test registration
	response, err := suite.userService.Register(req)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), "John Doe", response.User.Name)
	assert.Equal(suite.T(), "6281234567890", response.User.Phone) // Normalized phone
	assert.NotEmpty(suite.T(), response.AccessToken)
	assert.NotEmpty(suite.T(), response.RefreshToken)
}

// TestUserLogin tests user login
func (suite *UserServiceTestSuite) TestUserLogin() {
	// First register a user
	registerReq := &services.RegisterRequest{
		Name:     "Jane Doe",
		Phone:    "081987654321",
		Password: "password123",
	}

	_, err := suite.userService.Register(registerReq)
	assert.NoError(suite.T(), err)

	// Test login
	loginReq := &services.LoginRequest{
		Phone:    "081987654321",
		Password: "password123",
	}

	response, err := suite.userService.Login(loginReq)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), "Jane Doe", response.User.Name)
	assert.NotEmpty(suite.T(), response.AccessToken)
	assert.NotEmpty(suite.T(), response.RefreshToken)
}

// TestUserLoginWithWrongPassword tests login with wrong password
func (suite *UserServiceTestSuite) TestUserLoginWithWrongPassword() {
	// First register a user
	registerReq := &services.RegisterRequest{
		Name:     "Jane Doe",
		Phone:    "081987654321",
		Password: "password123",
	}

	_, err := suite.userService.Register(registerReq)
	assert.NoError(suite.T(), err)

	// Test login with wrong password
	loginReq := &services.LoginRequest{
		Phone:    "081987654321",
		Password: "wrongpassword",
	}

	response, err := suite.userService.Login(loginReq)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.Contains(suite.T(), err.Error(), "invalid phone number or password")
}

// TestDuplicateUserRegistration tests duplicate user registration
func (suite *UserServiceTestSuite) TestDuplicateUserRegistration() {
	// First registration
	req1 := &services.RegisterRequest{
		Name:     "John Doe",
		Phone:    "081234567890",
		Password: "password123",
	}

	_, err := suite.userService.Register(req1)
	assert.NoError(suite.T(), err)

	// Attempt duplicate registration
	req2 := &services.RegisterRequest{
		Name:     "John Smith",
		Phone:    "081234567890", // Same phone number
		Password: "password456",
	}

	response, err := suite.userService.Register(req2)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.Contains(suite.T(), err.Error(), "already exists")
}

// Run the test suite
func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
