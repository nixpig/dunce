package tag

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/nixpig/dunce/pkg"
	"github.com/stretchr/testify/mock"
)

var mockSessionManager = new(MockSessionManager)

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

func (s *MockSessionManager) Put(ctx context.Context, key string, val interface{}) {
	s.Called(ctx, key, val)
}

func (s *MockSessionManager) Remove(ctx context.Context, key string) {
	s.Called(ctx, key)
}

type MockTagService struct {
	mock.Mock
}

func (s *MockTagService) Create(tag *TagNewRequestDto) (*TagResponseDto, error) {
	args := s.Called(tag)

	return args.Get(0).(*TagResponseDto), args.Error(1)
}

func (s *MockTagService) DeleteById(id int) error {
	args := s.Called(id)

	return args.Error(0)
}

func (s *MockTagService) GetAll() (*[]TagResponseDto, error) {
	args := s.Called()

	return args.Get(0).(*[]TagResponseDto), args.Error(1)
}

func (s *MockTagService) GetByAttribute(attr, slug string) (*TagResponseDto, error) {
	args := s.Called(slug)

	return args.Get(0).(*TagResponseDto), args.Error(1)
}

func (s *MockTagService) Update(tag *TagUpdateRequestDto) (*TagResponseDto, error) {
	args := s.Called(tag)

	return args.Get(0).(*TagResponseDto), args.Error(1)
}

var mockService = new(MockTagService)

type MockLogger struct {
	mock.Mock
}

func (l *MockLogger) Info(format string, values ...any) {
	l.Called(format, values)
}

func (l *MockLogger) Error(format string, values ...any) {
	l.Called(format, values)
}

var mockLogger = new(MockLogger)

func mockTemplate() *template.Template {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	ts, err := template.ParseFiles(
		// FIXME: less than ideal arbitrarily jumping up two levels 😒
		path.Join(pwd, "..", "..", "test", "templates", "admin.tmpl"),
	)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	return ts
}

var mockTemplateCache = map[string]*template.Template{
	"pages/admin/admin-new-tag.tmpl": mockTemplate(),
}

func TestTagsControllerNewHandler(t *testing.T) {
	scenarios := map[string]func(t *testing.T, ctrl TagController){
		"test handle get new tag": testGetAdminTagsNewHandler,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			config := pkg.ControllerConfig{
				Log:            mockLogger,
				TemplateCache:  mockTemplateCache,
				SessionManager: mockSessionManager,
				CsrfToken: func(r *http.Request) string {
					return ""
				},
			}

			ctrl := NewTagController(mockService, config)

			fn(t, ctrl)
		})
	}
}

func testGetAdminTagsNewHandler(t *testing.T, ctrl TagController) {
	req, err := http.NewRequest("GET", "/admin/tags/create", nil)
	if err != nil {
		t.Fatal("failed to construct request", err)
	}

	mockSessionManagerExists := mockSessionManager.
		On("Exists", mock.Anything, "logged_in_username").
		Return(true)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ctrl.GetAdminTagsNewHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Error("status not ok")
	}

	mockSessionManagerExists.Unset()
	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("unmet expectations")
	}
}
