package integration_tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/authorizerdev/authorizer/internal/graph/model"
)

// TestLogout tests the logout functionality of the Authorizer application.
func TestLogout(t *testing.T) {
	cfg := getTestConfig()
	ts := initTestSetup(t, cfg)
	_, ctx := createContext(ts)

	// Create a test user and login to get tokens
	email := "logout_test_" + uuid.New().String() + "@authorizer.dev"
	password := "Password@123"

	// Signup the user
	signupReq := &model.SignUpInput{
		Email:           &email,
		Password:        password,
		ConfirmPassword: password,
	}
	signupRes, err := ts.GraphQLProvider.SignUp(ctx, signupReq)
	assert.NoError(t, err)
	assert.NotNil(t, signupRes)
	assert.Equal(t, email, *signupRes.User.Email)
	assert.NotEmpty(t, *signupRes.AccessToken)

	// Login to get fresh tokens
	loginReq := &model.LoginInput{
		Email:    &email,
		Password: password,
	}
	loginRes, err := ts.GraphQLProvider.Login(ctx, loginReq)
	assert.NoError(t, err)
	assert.NotNil(t, loginRes)
	assert.NotEmpty(t, *loginRes.AccessToken)

	// Test cases
	t.Run("should fail logout without access token", func(t *testing.T) {
		// Clear any existing authorization header
		ts.GinContext.Request.Header.Set("Authorization", "")

		logoutRes, err := ts.GraphQLProvider.Logout(ctx)
		assert.Error(t, err)
		assert.Nil(t, logoutRes)
	})

	t.Run("should fail logout with invalid access token", func(t *testing.T) {
		// Set an invalid token
		ts.GinContext.Request.Header.Set("Authorization", "Bearer invalid_token")

		logoutRes, err := ts.GraphQLProvider.Logout(ctx)
		assert.Error(t, err)
		assert.Nil(t, logoutRes)
	})

	t.Run("should successfully logout with valid access token", func(t *testing.T) {
		// Set the valid access token
		ts.GinContext.Request.Header.Set("Authorization", "Bearer "+*loginRes.AccessToken)
		logoutRes, err := ts.GraphQLProvider.Logout(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, logoutRes)
		assert.NotNil(t, logoutRes.Message)

		// Try to access a protected endpoint with the same token
		profile, err := ts.GraphQLProvider.Profile(ctx)
		assert.Error(t, err)
		assert.Nil(t, profile)
	})
}
