package tag

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/nixpig/dunce/pkg"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var mockTemplateCache = map[string]pkg.Template{
	"pages/admin/admin-new-tag.tmpl": mockTemplate,
}

var mockLogger = new(MockLogger)
var mockSessionManager = new(MockSessionManager)

func TestTagsControllerNewHandler(t *testing.T) {
	scenarios := map[string]func(t *testing.T, ctrl TagController){
		"test handle get new tag (success)":              testGetAdminTagsNewHandler,
		"test handle get new tag (error - template)":     testGetAdminTagsNewHandlerTemplateError,
		"test handle create new tag (success)":           testPostAdminTagsHandler,
		"test handle create new tag (error - service)":   testPostAdminTagsHandlerServiceError,
		"test handle delete tag (success)":               testPostAdminTagsDeleteHandler,
		"test handle delete tag (error - bad id)":        testPostAdminTagsDeleteHandlerErrorBadId,
		"test handle delete tag (error - service error)": testPostAdminTagsDeleteHandlerServiceError,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			config := pkg.ControllerConfig{
				Log:            mockLogger,
				TemplateCache:  mockTemplateCache,
				SessionManager: mockSessionManager,
				CsrfToken: func(r *http.Request) string {
					return "mock-token"
				},
			}

			ctrl := NewTagController(mockService, config)

			fn(t, ctrl)
		})
	}
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

var mockTemplate = new(MockTemplate)

type MockTemplate struct {
	mock.Mock
}

func (t *MockTemplate) ExecuteTemplate(wr io.Writer, name string, data any) error {
	args := t.Called(wr, name, data)

	return args.Error(0)
}

func testGetAdminTagsNewHandler(t *testing.T, ctrl TagController) {
	req, err := http.NewRequest("GET", "/admin/tags/new", nil)
	if err != nil {
		t.Error("failed to construct request", err)
	}

	mockSessionManagerExists := mockSessionManager.
		On("Exists", mock.Anything, "logged_in_username").
		Return(true)

	mockTemplateExecuteTemplate := mockTemplate.
		On("ExecuteTemplate", mock.Anything, "admin", TagCreateView{
			CsrfToken:       "mock-token",
			IsAuthenticated: true,
		}).
		Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ctrl.GetAdminTagsNewHandler)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Result().StatusCode, "should return status code ok")

	mockSessionManagerExists.Unset()
	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should evaluate if session exists")
	}

	mockTemplateExecuteTemplate.Unset()
	if res := mockTemplate.AssertExpectations(t); !res {
		t.Error("should execute template")
	}
}

func testGetAdminTagsNewHandlerTemplateError(t *testing.T, ctrl TagController) {
	req, err := http.NewRequest("GET", "/admin/tags/new", nil)
	if err != nil {
		t.Error("failed to construct request", err)
	}

	mockSessionManagerExists := mockSessionManager.
		On("Exists", mock.Anything, "logged_in_username").
		Return(true)

	mockTemplateExecuteTemplate := mockTemplate.
		On("ExecuteTemplate", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("template_error"))

	mockLoggerError := mockLogger.On("Error", "template_error", mock.Anything).Return()

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.GetAdminTagsNewHandler)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode, "should return status code internal server error")

	mockSessionManagerExists.Unset()
	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should evaluate session exists")
	}

	mockTemplateExecuteTemplate.Unset()
	if res := mockTemplate.AssertExpectations(t); !res {
		t.Error("should execute template")
	}

	mockLoggerError.Unset()
	if res := mockLogger.AssertExpectations(t); !res {
		t.Error("should log error")
	}
}

func testPostAdminTagsHandler(t *testing.T, ctrl TagController) {
	form := url.Values{}
	form.Add("name", "tag name")
	form.Add("slug", "tag-slug")

	req, err := http.NewRequest("POST", "/admin/tags", strings.NewReader(form.Encode()))
	if err != nil {
		t.Error("unable to construct request")
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	mockServiceCreate := mockService.On("Create", &TagNewRequestDto{
		Name: "tag name",
		Slug: "tag-slug",
	}).Return(&TagResponseDto{
		Id:   23,
		Name: "tag name",
		Slug: "tag-slug",
	}, nil)

	mockSessionManagerPut := mockSessionManager.
		On("Put", mock.Anything, "message", "Created tag 'tag name'.")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.PostAdminTagsHandler)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusSeeOther, rr.Result().StatusCode, "should return status code see other")
	require.Equal(t, "/admin/tags", rr.Result().Header.Get("Location"), "should set redirect location")

	mockServiceCreate.Unset()
	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call through to tag service to create tag")
	}

	mockSessionManagerPut.Unset()
	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should put response message in session")
	}
}

func testPostAdminTagsHandlerServiceError(t *testing.T, ctrl TagController) {
	form := url.Values{}
	form.Add("name", "tag name")
	form.Add("slug", "tag-slug")

	req, err := http.NewRequest("POST", "/admin/tags", strings.NewReader(form.Encode()))
	if err != nil {
		t.Error("unable to construct request")
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.PostAdminTagsHandler)

	mockServiceCreate := mockService.On("Create", &TagNewRequestDto{
		Name: "tag name",
		Slug: "tag-slug",
	}).Return(&TagResponseDto{}, errors.New("service_error"))

	mockLoggerError := mockLogger.On("Error", "service_error", mock.Anything).Return()

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode, "should return status code internal server error")

	mockServiceCreate.Unset()
	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call through to tag service to create")
	}

	mockLoggerError.Unset()
	if res := mockLogger.AssertExpectations(t); !res {
		t.Error("should log error")
	}
}

func testPostAdminTagsDeleteHandler(t *testing.T, ctrl TagController) {
	form := url.Values{}
	form.Add("id", "23")
	form.Add("name", "tag name")

	req, err := http.NewRequest("POST", "/admin/tags/tag-slug/delete", strings.NewReader(form.Encode()))
	if err != nil {
		t.Error("unable to construct request")
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.DeleteAdminTagsSlugHandler)

	mockServiceDeleteById := mockService.On("DeleteById", 23).Return(nil)

	mockSessionManagerPut := mockSessionManager.On("Put", mock.Anything, "message", "Deleted tag 'tag name'.")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusSeeOther, rr.Result().StatusCode, "should return status see other")
	require.Equal(t, "/admin/tags", rr.Result().Header.Get("Location"), "should redirect to tags page")

	mockServiceDeleteById.Unset()
	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call through to delete in tag service")
	}

	mockSessionManagerPut.Unset()
	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should put message into session context")
	}
}

func testPostAdminTagsDeleteHandlerErrorBadId(t *testing.T, ctrl TagController) {
	form := url.Values{}
	form.Add("id", "nonsense")
	form.Add("name", "tag name")

	req, err := http.NewRequest("POST", "/admin/tags/tag-slug/delete", strings.NewReader(form.Encode()))
	if err != nil {
		t.Error("unable to create request")
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.DeleteAdminTagsSlugHandler)

	mockLoggerError := mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Result().StatusCode, "should return status code bad request")

	mockLoggerError.Unset()
	if res := mockLogger.AssertExpectations(t); !res {
		t.Error("should log the error")
	}
}

func testPostAdminTagsDeleteHandlerServiceError(t *testing.T, ctrl TagController) {
	form := url.Values{}
	form.Add("id", "23")
	form.Add("name", "tag name")

	req, err := http.NewRequest("POST", "/admin/tags/tag-slug/delete", strings.NewReader(form.Encode()))
	if err != nil {
		t.Error("unable to create request")
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.DeleteAdminTagsSlugHandler)

	mockServiceDeleteById := mockService.On("DeleteById", 23).Return(errors.New("service_error"))

	mockLoggerError := mockLogger.On("Error", "service_error", mock.Anything).Return()

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode, "should return status internal server error")

	mockServiceDeleteById.Unset()
	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call tag service to delete")
	}

	mockLoggerError.Unset()
	if res := mockLogger.AssertExpectations(t); !res {
		t.Error("should log error")
	}
}
