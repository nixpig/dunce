package site

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/nixpig/dunce/pkg/templates"
	"github.com/stretchr/testify/mock"
)

var mockLogger = new(MockLogger)
var mockSessionManager = new(MockSessionManager)
var mockErrorHandlers = new(MockErrorHandlers)
var mockService = new(MockService)
var mockTemplateCache = templates.TemplateCache{
	"pages/admin/site.tmpl": mockTemplate,
}

func TestSiteController(t *testing.T) {
	scenarios := map[string]func(t *testing.T, ctrl SiteController){
		"test create site item view": testSiteControllerViewCreateItem,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			ctrl := NewSiteController(mockService, SiteControllerConfig{
				ErrorHandlers:  mockErrorHandlers,
				Log:            mockLogger,
				SessionManager: mockSessionManager,
				TemplateCache:  mockTemplateCache,
				CsrfToken: func(r *http.Request) string {
					return "mock-token"
				},
			})

			fn(t, ctrl)
		})
	}
}

func testSiteControllerViewCreateItem(t *testing.T, ctrl SiteController) {

}

type MockErrorHandlers struct {
	mock.Mock
}

func (e *MockErrorHandlers) NotFound(w http.ResponseWriter, r *http.Request) {
	e.Called(w, r)
}

func (e *MockErrorHandlers) InternalServerError(w http.ResponseWriter, r *http.Request) {
	e.Called(w, r)
}

func (e *MockErrorHandlers) BadRequest(w http.ResponseWriter, r *http.Request) {
	e.Called(w, r)
}

type MockSessionManager struct {
	mock.Mock
}

func (s *MockSessionManager) Exists(ctx context.Context, key string) bool {
	args := s.Called(ctx, key)

	return args.Bool(0)
}

func (s *MockSessionManager) PopString(ctx context.Context, key string) string {
	args := s.Called(ctx, key)

	return args.String(0)
}

func (s *MockSessionManager) GetString(ctx context.Context, key string) string {
	args := s.Called(ctx, key)

	return args.String(0)
}

func (s *MockSessionManager) LoadAndSave(next http.Handler) http.Handler {
	args := s.Called(next)

	return args.Get(0).(http.Handler)
}

func (s *MockSessionManager) RenewToken(ctx context.Context) error {
	args := s.Called(ctx)

	return args.Error(0)
}

func (s *MockSessionManager) Put(
	ctx context.Context,
	key string,
	val interface{},
) {
	s.Called(ctx, key, val)
}

func (s *MockSessionManager) Remove(ctx context.Context, key string) {
	s.Called(ctx, key)
}

type MockLogger struct {
	mock.Mock
}

func (l *MockLogger) Info(format string, values ...any) {
	l.Called(format, values)
}

func (l *MockLogger) Error(format string, values ...any) {
	l.Called(format, values)
}

var mockTemplate = new(MockTemplate)

type MockTemplate struct {
	mock.Mock
}

func (t *MockTemplate) ExecuteTemplate(
	wr io.Writer,
	name string,
	data any,
) error {
	args := t.Called(wr, name, data)

	return args.Error(0)
}

type MockService struct {
	mock.Mock
}

func (s *MockService) Create(key, value string) (*SiteItemResponseDto, error) {
	args := s.Called(key, value)

	return args.Get(0).(*SiteItemResponseDto), args.Error(1)
}
