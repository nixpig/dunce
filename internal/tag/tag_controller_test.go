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

	"github.com/nixpig/dunce/pkg/session"
	"github.com/nixpig/dunce/pkg/templates"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var mockTemplateCache = templates.TemplateCache{
	"pages/admin/new-tag.tmpl": mockTemplate,
	"pages/admin/tags.tmpl":    mockTemplate,
	"pages/admin/tag.tmpl":     mockTemplate,
}

var mockLogger = new(MockLogger)
var mockSessionManager = new(MockSessionManager)
var mockErrorHandlers = new(MockErrorHandlers)

func TestTagsControllerNewHandler(t *testing.T) {
	scenarios := map[string]func(t *testing.T, ctrl TagController){
		"test handle get new tag (success)":              testGetAdminTagsNewHandler,
		"test handle get new tag (error - template)":     testGetAdminTagsNewHandlerTemplateError,
		"test handle create new tag (success)":           testPostAdminTagsHandler,
		"test handle create new tag (error - service)":   testPostAdminTagsHandlerServiceError,
		"test handle delete tag (success)":               testPostAdminTagsDeleteHandler,
		"test handle delete tag (error - bad id)":        testPostAdminTagsDeleteHandlerErrorBadId,
		"test handle delete tag (error - service error)": testPostAdminTagsDeleteHandlerServiceError,
		"test get tags (success)":                        testGetAdminTagsHandler,
		"test get tags (error - service error)":          testGetAdminTagsHandlerServiceError,
		"test get tags (error - template error)":         testGetAdminTagsHandlerTemplateError,
		"test get tags by slug (success)":                testGetAdminTagsBySlugHandler,
		"test get tags by slug (error - service error)":  testGetAdminTagsBySlugHandlerServiceError,
		"test get tags by slug (error - template error)": testGetAdminTagsBySlugHandlerTemplateError,
		"test post tag by slug (success)":                testPostTagBySlugToUpdateHandler,
		"test post tag by slug (error - bad form id)":    testPostTagBySlugToUpdateHandlerBadFormIdError,
		"test post tag by slug (error - service error)":  testPostTagBySlugToUpdateHandlerServiceError,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			config := TagControllerConfig{
				Log:            mockLogger,
				TemplateCache:  mockTemplateCache,
				SessionManager: mockSessionManager,
				CsrfToken: func(r *http.Request) string {
					return "mock-token"
				},
				ErrorHandlers: mockErrorHandlers,
			}

			ctrl := NewTagController(mockService, config)

			fn(t, ctrl)
		})
	}
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

type MockTagService struct {
	mock.Mock
}

func (s *MockTagService) Create(
	tag *TagNewRequestDto,
) (*TagResponseDto, error) {
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

func (s *MockTagService) GetByAttribute(
	attr, slug string,
) (*TagResponseDto, error) {
	args := s.Called(attr, slug)

	return args.Get(0).(*TagResponseDto), args.Error(1)
}

func (s *MockTagService) Update(
	tag *TagUpdateRequestDto,
) (*TagResponseDto, error) {
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

func (t *MockTemplate) ExecuteTemplate(
	wr io.Writer,
	name string,
	data any,
) error {
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

	require.Equal(
		t,
		http.StatusOK,
		rr.Result().StatusCode,
		"should return status code ok",
	)

	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should evaluate if session exists")
	}

	if res := mockTemplate.AssertExpectations(t); !res {
		t.Error("should execute template")
	}

	mockTemplateExecuteTemplate.Unset()
	mockSessionManagerExists.Unset()
}

func testGetAdminTagsNewHandlerTemplateError(t *testing.T, ctrl TagController) {
	req, err := http.NewRequest("GET", "/admin/tags/new", nil)
	if err != nil {
		t.Error("failed to construct request", err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.GetAdminTagsNewHandler)

	mockSessionManagerExists := mockSessionManager.
		On("Exists", mock.Anything, session.LOGGED_IN_USERNAME).
		Return(true)

	mockTemplateExecuteTemplate := mockTemplate.
		On("ExecuteTemplate", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("template_error"))

	mockErrorHandlersInternalServerError := mockErrorHandlers.
		On("InternalServerError", rr, req).
		Run(func(args mock.Arguments) {
			rr.WriteHeader(http.StatusInternalServerError)
		})

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusInternalServerError,
		rr.Result().StatusCode,
		"should return status code internal server error",
	)

	if res := mockErrorHandlers.AssertExpectations(t); !res {
		t.Error("should call error handler")
	}

	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should evaluate session exists")
	}

	if res := mockTemplate.AssertExpectations(t); !res {
		t.Error("should execute template")
	}

	mockTemplateExecuteTemplate.Unset()
	mockSessionManagerExists.Unset()
	mockErrorHandlersInternalServerError.Unset()
}

func testPostAdminTagsHandler(t *testing.T, ctrl TagController) {
	form := url.Values{}
	form.Add("name", "tag name")
	form.Add("slug", "tag-slug")

	req, err := http.NewRequest(
		"POST",
		"/admin/tags",
		strings.NewReader(form.Encode()),
	)
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

	require.Equal(
		t,
		http.StatusSeeOther,
		rr.Result().StatusCode,
		"should return status code see other",
	)
	require.Equal(
		t,
		"/admin/tags",
		rr.Result().Header.Get("Location"),
		"should set redirect location",
	)

	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call through to tag service to create tag")
	}

	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should put response message in session")
	}

	mockSessionManagerPut.Unset()
	mockServiceCreate.Unset()
}

func testPostAdminTagsHandlerServiceError(t *testing.T, ctrl TagController) {
	form := url.Values{}
	form.Add("name", "tag name")
	form.Add("slug", "tag-slug")

	req, err := http.NewRequest(
		"POST",
		"/admin/tags",
		strings.NewReader(form.Encode()),
	)
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

	mockErrorHandlersInternalServerError := mockErrorHandlers.
		On("InternalServerError", rr, req).
		Run(func(args mock.Arguments) {
			rr.WriteHeader(http.StatusInternalServerError)
		})

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusInternalServerError,
		rr.Result().StatusCode,
		"should return status code internal server error",
	)

	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call through to tag service to create")
	}

	if res := mockErrorHandlers.AssertExpectations(t); !res {
		t.Error("should call error handler")
	}

	mockServiceCreate.Unset()
	mockErrorHandlersInternalServerError.Unset()
}

func testPostAdminTagsDeleteHandler(t *testing.T, ctrl TagController) {
	form := url.Values{}
	form.Add("id", "23")
	form.Add("name", "tag name")

	req, err := http.NewRequest(
		"POST",
		"/admin/tags/tag-slug/delete",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		t.Error("unable to construct request")
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.DeleteAdminTagsSlugHandler)

	mockServiceDeleteById := mockService.On("DeleteById", 23).Return(nil)

	mockSessionManagerPut := mockSessionManager.On(
		"Put",
		mock.Anything,
		"message",
		"Deleted tag 'tag name'.",
	)

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusSeeOther,
		rr.Result().StatusCode,
		"should return status see other",
	)
	require.Equal(
		t,
		"/admin/tags",
		rr.Result().Header.Get("Location"),
		"should redirect to tags page",
	)

	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call through to delete in tag service")
	}

	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should put message into session context")
	}

	mockSessionManagerPut.Unset()
	mockServiceDeleteById.Unset()
}

func testPostAdminTagsDeleteHandlerErrorBadId(
	t *testing.T,
	ctrl TagController,
) {
	form := url.Values{}
	form.Add("id", "nonsense")
	form.Add("name", "tag name")

	req, err := http.NewRequest(
		"POST",
		"/admin/tags/tag-slug/delete",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		t.Error("unable to create request")
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.DeleteAdminTagsSlugHandler)

	mockErrorHandlersBadRequest := mockErrorHandlers.
		On("BadRequest", rr, req).
		Run(func(args mock.Arguments) {
			rr.WriteHeader(http.StatusBadRequest)
		})

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusBadRequest,
		rr.Result().StatusCode,
		"should return status code bad request",
	)

	if res := mockErrorHandlers.AssertExpectations(t); !res {
		t.Error("should call error handler")
	}

	mockErrorHandlersBadRequest.Unset()
}

func testPostAdminTagsDeleteHandlerServiceError(
	t *testing.T,
	ctrl TagController,
) {
	form := url.Values{}
	form.Add("id", "23")
	form.Add("name", "tag name")

	req, err := http.NewRequest(
		"POST",
		"/admin/tags/tag-slug/delete",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		t.Error("unable to create request")
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.DeleteAdminTagsSlugHandler)

	mockServiceDeleteById := mockService.On("DeleteById", 23).
		Return(errors.New("service_error"))

	mockErrorHandlersInternalServerError := mockErrorHandlers.
		On("InternalServerError", rr, req).
		Run(func(args mock.Arguments) {
			rr.WriteHeader(http.StatusInternalServerError)
		})

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusInternalServerError,
		rr.Result().StatusCode,
		"should return status internal server error",
	)

	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call tag service to delete")
	}

	if res := mockErrorHandlers.AssertExpectations(t); !res {
		t.Error("should call error handlers")
	}

	mockServiceDeleteById.Unset()
	mockErrorHandlersInternalServerError.Unset()
}

func testGetAdminTagsHandler(t *testing.T, ctrl TagController) {
	req, err := http.NewRequest("GET", "/admin/tags", nil)
	if err != nil {
		t.Error("unable to construct request")
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.GetAdminTagsHandler)

	mockServiceGetAll := mockService.On("GetAll").Return(&[]TagResponseDto{
		{
			Id:   1,
			Name: "tag one",
			Slug: "tag-one",
		},
		{
			Id:   2,
			Name: "tag two",
			Slug: "tag-two",
		},
	}, nil)

	mockSessionManagerPopString := mockSessionManager.
		On("PopString", mock.Anything, "message").
		Return("session_message")

	mockSessionManagerExists := mockSessionManager.
		On(
			"Exists",
			mock.Anything,
			"logged_in_username",
		).Return(true)

	mockTemplateExecuteTemplate := mockTemplate.
		On("ExecuteTemplate", mock.Anything, "admin", TagsView{
			Message: "session_message",
			Tags: &[]TagResponseDto{
				{
					Id:   1,
					Name: "tag one",
					Slug: "tag-one",
				},
				{
					Id:   2,
					Name: "tag two",
					Slug: "tag-two",
				},
			},
			CsrfToken:       "mock-token",
			IsAuthenticated: true,
		}).Return(nil)

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusOK,
		rr.Result().StatusCode,
		"should return status code ok",
	)

	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call tag service to get all tags")
	}

	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should call session manager")
	}

	if res := mockTemplate.AssertExpectations(t); !res {
		t.Error("should execute template with tags")
	}

	mockServiceGetAll.Unset()
	mockSessionManagerPopString.Unset()
	mockSessionManagerExists.Unset()
	mockTemplateExecuteTemplate.Unset()
}

func testGetAdminTagsHandlerServiceError(t *testing.T, ctrl TagController) {
	req, err := http.NewRequest("GET", "/admin/tags", nil)
	if err != nil {
		t.Error("unable to construct request")
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.GetAdminTagsHandler)

	mockServiceGetAll := mockService.On("GetAll").
		Return(&[]TagResponseDto{}, errors.New("service_error"))

	mockErrorHandlersInternalServerError := mockErrorHandlers.
		On("InternalServerError", rr, req).
		Run(func(args mock.Arguments) {
			rr.WriteHeader(http.StatusInternalServerError)
		})

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusInternalServerError,
		rr.Result().StatusCode,
		"should return status code internal server error",
	)

	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call tag service to get all tags")
	}

	if res := mockErrorHandlers.AssertExpectations(t); !res {
		t.Error("should call error handler")
	}

	mockErrorHandlersInternalServerError.Unset()
	mockServiceGetAll.Unset()
}

func testGetAdminTagsHandlerTemplateError(t *testing.T, ctrl TagController) {
	req, err := http.NewRequest("GET", "/admin/tags", nil)
	if err != nil {
		t.Error("unable to construct request")
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.GetAdminTagsHandler)

	mockServiceGetAll := mockService.On("GetAll").Return(&[]TagResponseDto{
		{
			Id:   1,
			Name: "tag one",
			Slug: "tag-one",
		},
		{
			Id:   2,
			Name: "tag two",
			Slug: "tag-two",
		},
	}, nil)

	mockSessionManagerPopString := mockSessionManager.On("PopString", mock.Anything, "message").
		Return("mock_message")
	mockSessionManagerExists := mockSessionManager.On("Exists", mock.Anything, "logged_in_username").
		Return(true)

	mockTemplateExecuteTemplate := mockTemplate.On("ExecuteTemplate", mock.Anything, "admin", TagsView{
		Message: "mock_message",
		Tags: &[]TagResponseDto{
			{
				Id:   1,
				Name: "tag one",
				Slug: "tag-one",
			},
			{
				Id:   2,
				Name: "tag two",
				Slug: "tag-two",
			},
		},
		CsrfToken:       "mock-token",
		IsAuthenticated: true,
	}).
		Return(errors.New("template_error"))

	mockErrorHandlersInternalServerError := mockErrorHandlers.
		On("InternalServerError", rr, req).
		Run(func(args mock.Arguments) {
			rr.WriteHeader(http.StatusInternalServerError)
		})

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusInternalServerError,
		rr.Result().StatusCode,
		"should return status code internal server error",
	)

	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call service to get all tags")
	}

	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should call session methods")
	}

	if res := mockTemplate.AssertExpectations(t); !res {
		t.Error("should execute template")
	}

	if res := mockErrorHandlers.AssertExpectations(t); !res {
		t.Error("should call error handler")
	}

	mockErrorHandlersInternalServerError.Unset()
	mockServiceGetAll.Unset()
	mockSessionManagerPopString.Unset()
	mockSessionManagerExists.Unset()
	mockTemplateExecuteTemplate.Unset()
}

func testGetAdminTagsBySlugHandler(t *testing.T, ctrl TagController) {
	req, err := http.NewRequest("GET", "/admin/tags/{slug}", nil)
	if err != nil {
		t.Error("unable to construct request")
	}

	req.SetPathValue("slug", "tag-slug")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.GetAdminTagsSlugHandler)

	mockServiceGetByAttribute := mockService.On("GetByAttribute", "slug", "tag-slug").
		Return(&TagResponseDto{
			Id:   23,
			Name: "tag name",
			Slug: "tag-slug",
		}, nil)

	mockSessionManagerExists := mockSessionManager.On("Exists", mock.Anything, "logged_in_username").
		Return(true)

	mockTemplateExecuteTemplate := mockTemplate.On(
		"ExecuteTemplate",
		mock.Anything,
		"admin",
		TagView{
			Tag: &TagResponseDto{
				Id:   23,
				Name: "tag name",
				Slug: "tag-slug",
			},
			CsrfToken:       "mock-token",
			IsAuthenticated: true,
		},
	).Return(nil)

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusOK,
		rr.Result().StatusCode,
		"should return status code ok",
	)

	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should check if user session exists in context")
	}

	if res := mockTemplate.AssertExpectations(t); !res {
		t.Error("should execute template")
	}

	if res := mockService.AssertExpectations(t); !res {
		t.Error("should have called through to service to get tag")
	}

	mockSessionManagerExists.Unset()
	mockTemplateExecuteTemplate.Unset()
	mockServiceGetByAttribute.Unset()
}

func testGetAdminTagsBySlugHandlerServiceError(
	t *testing.T,
	ctrl TagController,
) {
	req, err := http.NewRequest("GET", "/admin/tags/{slug}", nil)
	if err != nil {
		t.Error("unable to construct request")
	}

	req.SetPathValue("slug", "tag-slug")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.GetAdminTagsSlugHandler)

	mockServiceGetByAttribute := mockService.On("GetByAttribute", "slug", "tag-slug").
		Return(&TagResponseDto{}, errors.New("service_error"))

	mockErrorHandlersInternalServerError := mockErrorHandlers.
		On("InternalServerError", rr, req).
		Run(func(args mock.Arguments) {
			rr.WriteHeader(http.StatusInternalServerError)
		})

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusInternalServerError,
		rr.Result().StatusCode,
		"should return status code internal server error",
	)

	if res := mockService.AssertExpectations(t); !res {
		t.Error("should have called through to service to get tag")
	}

	if res := mockErrorHandlers.AssertExpectations(t); !res {
		t.Error("should call error handler")
	}

	mockErrorHandlersInternalServerError.Unset()
	mockServiceGetByAttribute.Unset()
}

func testGetAdminTagsBySlugHandlerTemplateError(
	t *testing.T,
	ctrl TagController,
) {
	req, err := http.NewRequest("GET", "/admin/tags/{slug}", nil)
	if err != nil {
		t.Error("unable to construct request")
	}

	req.SetPathValue("slug", "tag-slug")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.GetAdminTagsSlugHandler)

	mockServiceGetByAttribute := mockService.On("GetByAttribute", "slug", "tag-slug").
		Return(&TagResponseDto{
			Id:   23,
			Name: "tag name",
			Slug: "tag-slug",
		}, nil)

	mockSessionManagerExists := mockSessionManager.On("Exists", mock.Anything, "logged_in_username").
		Return(true)

	mockErrorHandlersInternalServerError := mockErrorHandlers.
		On("InternalServerError", rr, req).
		Run(func(args mock.Arguments) {
			rr.WriteHeader(http.StatusInternalServerError)
		})

	mockTemplateExecuteTemplate := mockTemplate.On(
		"ExecuteTemplate",
		mock.Anything,
		"admin",
		TagView{
			Tag: &TagResponseDto{
				Id:   23,
				Name: "tag name",
				Slug: "tag-slug",
			},
			CsrfToken:       "mock-token",
			IsAuthenticated: true,
		},
	).Return(errors.New("template_error"))

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusInternalServerError,
		rr.Result().StatusCode,
		"should return status code internal server error",
	)

	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should check if user session exists in context")
	}

	if res := mockTemplate.AssertExpectations(t); !res {
		t.Error("should execute template")
	}

	if res := mockService.AssertExpectations(t); !res {
		t.Error("should have called through to service to get tag")
	}

	if res := mockErrorHandlers.AssertExpectations(t); !res {
		t.Error("should call error handlers")
	}

	mockErrorHandlersInternalServerError.Unset()
	mockSessionManagerExists.Unset()
	mockTemplateExecuteTemplate.Unset()
	mockServiceGetByAttribute.Unset()
}

func testPostTagBySlugToUpdateHandler(t *testing.T, ctrl TagController) {
	form := url.Values{}
	form.Add("id", "23")
	form.Add("name", "tag name")
	form.Add("slug", "tag-slug")

	req, err := http.NewRequest(
		"POST",
		"/admin/tags/{slug}",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		t.Error("unable to construct request")
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	req.SetPathValue("slug", "tag-slug")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.PostAdminTagsSlugHandler)

	mockServiceUpdate := mockService.On("Update", &TagUpdateRequestDto{
		Id:   23,
		Name: "tag name",
		Slug: "tag-slug",
	}).Return(&TagResponseDto{}, nil)

	mockSessionManagerPut := mockSessionManager.On(
		"Put",
		mock.Anything,
		"message",
		"Updated tag 'tag name'.",
	)

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusSeeOther,
		rr.Result().StatusCode,
		"should return status code see other",
	)
	require.Equal(
		t,
		"/admin/tags",
		rr.Result().Header.Get("Location"),
		"should redirect to tags page",
	)

	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call through to service")
	}

	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should put message in session context")
	}

	mockServiceUpdate.Unset()
	mockSessionManagerPut.Unset()
}

func testPostTagBySlugToUpdateHandlerBadFormIdError(
	t *testing.T,
	ctrl TagController,
) {
	form := url.Values{}
	form.Add("id", "nonsense")
	form.Add("name", "tag name")
	form.Add("slug", "tag-slug")

	req, err := http.NewRequest(
		"POST",
		"/admin/tags/{slug}",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		t.Error("unable to construct request")
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	req.SetPathValue("slug", "tag-slug")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.PostAdminTagsSlugHandler)

	mockErrorHandlersBadRequest := mockErrorHandlers.
		On("BadRequest", rr, req).
		Run(func(args mock.Arguments) {
			rr.WriteHeader(http.StatusBadRequest)
		})

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusBadRequest,
		rr.Result().StatusCode,
		"should return status code bad request",
	)

	if res := mockErrorHandlers.AssertExpectations(t); !res {
		t.Error("should call error handlers")
	}

	mockErrorHandlersBadRequest.Unset()
}

func testPostTagBySlugToUpdateHandlerServiceError(
	t *testing.T,
	ctrl TagController,
) {
	form := url.Values{}
	form.Add("id", "23")
	form.Add("name", "tag name")
	form.Add("slug", "tag-slug")
	req, err := http.NewRequest(
		"POST",
		"/admin/tags/{slug}",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		t.Error("unable to construct request")
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.PostAdminTagsSlugHandler)

	mockServiceUpdate := mockService.On("Update", &TagUpdateRequestDto{
		Id:   23,
		Name: "tag name",
		Slug: "tag-slug",
	}).Return(&TagResponseDto{}, errors.New("service_error"))

	mockErrorHandlersInternalServerError := mockErrorHandlers.
		On("InternalServerError", rr, req).
		Run(func(args mock.Arguments) {
			rr.WriteHeader(http.StatusInternalServerError)
		})

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusInternalServerError,
		rr.Result().StatusCode,
		"should return status code internal server error",
	)

	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call through to service")
	}

	if res := mockErrorHandlers.AssertExpectations(t); !res {
		t.Error("should call error handler")
	}

	mockErrorHandlersInternalServerError.Unset()
	mockServiceUpdate.Unset()
}
