package run

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockServer struct {
	mock.Mock
}

func (m *MockServer) ListenAndServe() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockServer) Shutdown(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestNewApp(t *testing.T) {
	app := NewApp()

	assert.NotNil(t, app, "NewApp() should not return nil")

	assert.IsType(t, &App{}, app, "NewApp() should return an instance of *App")
}

func TestApp_Run(t *testing.T) {
	mockServer := new(MockServer)
	app := &App{srv: mockServer}

	mockServer.On("ListenAndServe").Return(nil)
	mockServer.On("Shutdown", mock.Anything).Return(nil)

	done := make(chan struct{})

	go func() {
		app.Run()
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not finish in time")
	}

	mockServer.AssertCalled(t, "ListenAndServe")
	mockServer.AssertCalled(t, "Shutdown", mock.Anything)
}
