package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) FindUserByEmail(email string) (*struct{}, error) {
	args := m.Called(email)
	return args.Get(0).(*struct{}), args.Error(1)
}

func (m *MockUserService) SignUp(input *dto.UserSignup) (string, error) {
	args := m.Called(input)
	return args.String(0), args.Error(1)
}

func TestUserHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		mockSetup      func(*MockUserService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful registration",
			body: dto.UserSignup{
				UserLogin: dto.UserLogin{
					Email:    "test@example.com",
					Password: "password123",
				},
				Phone: "1234567890",
			},
			mockSetup: func(m *MockUserService) {
				m.On("SignUp", mock.AnythingOfType("*dto.UserSignup")).Return("mocked-token", nil)
			},
			expectedStatus: fiber.StatusCreated,
			expectedBody: map[string]interface{}{
				"token":   "mocked-token",
				"message": "successfully signed up",
			},
		},
		{
			name: "invalid input (malformed JSON)",
			body: "invalid json",
			mockSetup: func(m *MockUserService) {
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "please provide valid input",
			},
		},
		{
			name: "missing required fields",
			body: dto.UserSignup{
				UserLogin: dto.UserLogin{
					Email:    "",
					Password: "password123",
				},
				Phone: "1234567890",
			},
			mockSetup: func(m *MockUserService) {
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "email, password, and phone are required",
			},
		},
		{
			name: "service error",
			body: dto.UserSignup{
				UserLogin: dto.UserLogin{
					Email:    "test@example.com",
					Password: "password123",
				},
				Phone: "1234567890",
			},
			mockSetup: func(m *MockUserService) {
				m.On("SignUp", mock.AnythingOfType("*dto.UserSignup")).Return("", assert.AnError)
			},
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"message": "error signing up",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			mockService := &MockUserService{}
			tt.mockSetup(mockService)

			handler := NewUserHandler(mockService)

			app.Post("/register", handler.Register)

			body, err := json.Marshal(tt.body)
			if tt.name == "invalid input (malformed JSON)" {
				body = []byte(tt.body.(string))
			} else {
				assert.NoError(t, err, "should marshal request body")
			}

			req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err, "should not return an error when making the request")

			assert.Equal(t, tt.expectedStatus, resp.StatusCode, "should return correct status code")

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			assert.NoError(t, err, "should parse response body without error")
			assert.Equal(t, tt.expectedBody, responseBody, "should return correct response body")

			mockService.AssertExpectations(t)
		})
	}
}
