package integration_tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/authorizerdev/authorizer/internal/graph/model"
	"github.com/authorizerdev/authorizer/internal/refs"
)

// TestUpdateProfile tests the update profile functionality
// using the GraphQL API.
// It creates a user, updates the profile, and verifies
// the changes in the database.
func TestUpdateProfile(t *testing.T) {
	cfg := getTestConfig()
	ts := initTestSetup(t, cfg)
	_, ctx := createContext(ts)

	// Create a test user
	email := "update_profile_test_" + uuid.New().String() + "@authorizer.dev"
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

	// Set the authorization header for authenticated requests
	ts.GinContext.Request.Header.Set("Authorization", "Bearer "+*loginRes.AccessToken)

	// Test cases
	t.Run("should fail update profile without authentication", func(t *testing.T) {
		// Clear authorization header
		ts.GinContext.Request.Header.Set("Authorization", "")
		defer func() {
			ts.GinContext.Request.Header.Set("Authorization", "Bearer "+*loginRes.AccessToken)
		}()

		updateReq := &model.UpdateProfileInput{
			GivenName: refs.NewStringRef("Test"),
		}
		updateRes, err := ts.GraphQLProvider.UpdateProfile(ctx, updateReq)
		assert.Error(t, err)
		assert.Nil(t, updateRes)
	})

	t.Run("should update basic profile information", func(t *testing.T) {
		givenName := "John"
		familyName := "Doe"
		nickname := "Johnny"
		phoneNumber := "+1234567890"

		updateReq := &model.UpdateProfileInput{
			GivenName:   refs.NewStringRef(givenName),
			FamilyName:  refs.NewStringRef(familyName),
			Nickname:    refs.NewStringRef(nickname),
			PhoneNumber: refs.NewStringRef(phoneNumber),
		}

		updateRes, err := ts.GraphQLProvider.UpdateProfile(ctx, updateReq)
		assert.NoError(t, err)
		assert.NotNil(t, updateRes)

		// Get the profile
		profile, err := ts.GraphQLProvider.Profile(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, profile)
		assert.Equal(t, givenName, *profile.GivenName)
		assert.Equal(t, familyName, *profile.FamilyName)
		assert.Equal(t, nickname, *profile.Nickname)
		assert.Equal(t, phoneNumber, *profile.PhoneNumber)
		assert.Equal(t, email, *profile.Email)
	})

	t.Run("should update profile picture", func(t *testing.T) {
		picture := "https://example.com/profile.jpg"

		updateReq := &model.UpdateProfileInput{
			Picture: refs.NewStringRef(picture),
		}

		updateRes, err := ts.GraphQLProvider.UpdateProfile(ctx, updateReq)
		assert.NoError(t, err)
		assert.NotNil(t, updateRes)

		// Get the profile
		profile, err := ts.GraphQLProvider.Profile(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, profile)
		assert.Equal(t, picture, *profile.Picture)
	})

	t.Run("should update gender", func(t *testing.T) {
		gender := "male"

		updateReq := &model.UpdateProfileInput{
			Gender: refs.NewStringRef(gender),
		}

		updateRes, err := ts.GraphQLProvider.UpdateProfile(ctx, updateReq)
		assert.NoError(t, err)
		assert.NotNil(t, updateRes)

		// Get the profile
		profile, err := ts.GraphQLProvider.Profile(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, profile)
	})

	t.Run("should update birthdate", func(t *testing.T) {
		birthdate := "1990-01-01"

		updateReq := &model.UpdateProfileInput{
			Birthdate: refs.NewStringRef(birthdate),
		}

		updateRes, err := ts.GraphQLProvider.UpdateProfile(ctx, updateReq)
		assert.NoError(t, err)
		assert.NotNil(t, updateRes)

		// Get the profile
		profile, err := ts.GraphQLProvider.Profile(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, profile)
		assert.Equal(t, birthdate, *profile.Birthdate)
	})

	t.Run("should update multiple fields at once", func(t *testing.T) {
		givenName := "Updated"
		familyName := "User"
		picture := "https://example.com/new-profile.jpg"

		updateReq := &model.UpdateProfileInput{
			GivenName:  refs.NewStringRef(givenName),
			FamilyName: refs.NewStringRef(familyName),
			Picture:    refs.NewStringRef(picture),
		}

		updateRes, err := ts.GraphQLProvider.UpdateProfile(ctx, updateReq)
		assert.NoError(t, err)
		assert.NotNil(t, updateRes)

		// Get the profile
		profile, err := ts.GraphQLProvider.Profile(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, profile)
		assert.Equal(t, givenName, *profile.GivenName)
		assert.Equal(t, familyName, *profile.FamilyName)
		assert.Equal(t, picture, *profile.Picture)
	})
}
