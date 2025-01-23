package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGormDB struct {
	mock.Mock
}

func (m *MockGormDB) AutoMigrate(models ...interface{}) error {
	args := m.Called(models)
	return args.Error(0)
}

func TestMigrateDB(t *testing.T) {
	mockDB := new(MockGormDB)

	t.Run("Success", func(t *testing.T) {
		mockDB.On("AutoMigrate", mock.Anything).Return(nil)

		err := MigrateDB(mockDB)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		mockDB.AssertCalled(t, "AutoMigrate", mock.Anything)
		mockDB.AssertExpectations(t)
	})
}

func TestNewDbConf(t *testing.T) {
	envs := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"DB_USER":     "test_user",
		"DB_PASSWORD": "test_password",
		"DB_NAME":     "test_db",
		"DB_DRIVER":   "postgres",
	}

	for key, value := range envs {
		assert.NoError(t, os.Setenv(key, value), "Failed to set environment variable")
	}

	conf := NewDbConf()

	assert.Equal(t, envs["DB_HOST"], conf.Host, "DB_HOST mismatch")
	assert.Equal(t, envs["DB_PORT"], conf.Port, "DB_PORT mismatch")
	assert.Equal(t, envs["DB_USER"], conf.User, "DB_USER mismatch")
	assert.Equal(t, envs["DB_PASSWORD"], conf.Password, "DB_PASSWORD mismatch")
	assert.Equal(t, envs["DB_NAME"], conf.Name, "DB_NAME mismatch")
	assert.Equal(t, envs["DB_DRIVER"], conf.Driver, "DB_DRIVER mismatch")

	for key := range envs {
		assert.NoError(t, os.Unsetenv(key), "Failed to unset environment variable")
	}
}

func TestNewDbConf_EmptyEnvironment(t *testing.T) {
	os.Clearenv()

	conf := NewDbConf()

	assert.Empty(t, conf.Host, "Expected Host to be empty")
	assert.Empty(t, conf.Port, "Expected Port to be empty")
	assert.Empty(t, conf.User, "Expected User to be empty")
	assert.Empty(t, conf.Password, "Expected Password to be empty")
	assert.Empty(t, conf.Name, "Expected Name to be empty")
	assert.Empty(t, conf.Driver, "Expected Driver to be empty")
}
