package user

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
	"pages/admin/login.tmpl":    mockTemplate,
	"pages/admin/new-user.tmpl": mockTemplate,
	"pages/admin/users.tmpl":    mockTemplate,
	"pages/admin/user.tmpl":     mockTemplate,
}

var mockLogger = new(MockLogger)
var mockSessionManager = new(MockSessionManager)
var mockService = new(MockUserService)

func TestUserController(t *testing.T) {
	scenarios := map[string]func(t *testing.T, ctrl UserController){
		"is authenticated helper":                            testIsAuthenticatedHelper,
		"get user login screen (already logged in)":          testGetUserLoginScreenHandlerIsLoggedIn,
		"get user login screen (not logged in)":              testGetUserLoginScreenHandlerNotLoggedIn,
		"get user login screen (error - template rendering)": testGetUserLoginScreenHandlerTemplateError,
		"post user login (success)":                          testPostUserLogin,
		"post user login (error - login failed)":             testPostUserLoginUsernamePasswordFailed,
		"post user login (error - renew token failed)":       testPostUserLoginRenewTokenFailed,
		"post user logout (success)":                         testPostUserLogout,
		"get create user page (success)":                     testGetCreateUserPage,
		"get create user page (error - template)":            testGetCreateUserPageTemplateError,
		"post create user (success)":                         testPostCreateUser,
		"post create user (error - service error)":           testPostCreateUserServiceError,
		"get all users (success)":                            testGetAllUsers,
		"get all users (error - service error)":              testGetAllUsersServiceError,
		"get all users (error - template error)":             testGetAllUsersTemplateError,
		"get user by username (success)":                     testGetUserByUsername,
		"get user by username (error - service error)":       testGetUserByUsernameServiceError,
		"get user by username (error - template error)":      testGetUserByUsernameTemplateError,
		"post delete user (success)":                         testPostDeleteUser,
		"post delete user (error - form error)":              testPostDeleteUserFormError,
		"post delete user (error - service error)":           testPostDeleteUserServiceError,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			config := UserControllerConfig{
				Log:            mockLogger,
				TemplateCache:  mockTemplateCache,
				SessionManager: mockSessionManager,
				CsrfToken: func(r *http.Request) string {
					return "mock-token"
				},
			}

			ctrl := NewUserController(mockService, config)
			fn(t, ctrl)
		})

	}

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

type MockUserService struct {
	mock.Mock
}

func (u *MockUserService) Create(
	user *UserNewRequestDto,
) (*UserResponseDto, error) {
	args := u.Called(user)

	return args.Get(0).(*UserResponseDto), args.Error(1)
}

func (u *MockUserService) DeleteById(id uint) error {
	args := u.Called(id)

	return args.Error(0)
}

func (u *MockUserService) Exists(username string) (bool, error) {
	args := u.Called(username)

	return args.Bool(0), args.Error(1)
}

func (u *MockUserService) GetAll() (*[]UserResponseDto, error) {
	args := u.Called()

	return args.Get(0).(*[]UserResponseDto), args.Error(1)
}

func (u *MockUserService) GetByAttribute(
	attr, value string,
) (*UserResponseDto, error) {
	args := u.Called(attr, value)

	return args.Get(0).(*UserResponseDto), args.Error(1)
}

func (u *MockUserService) Update(
	user *User,
) (*UserResponseDto, error) {
	args := u.Called(user)

	return args.Get(0).(*UserResponseDto), args.Error(1)
}

func (u *MockUserService) LoginWithUsernamePassword(
	username, password string,
) error {
	args := u.Called(username, password)

	return args.Error(0)
}

func testGetUserLoginScreenHandlerIsLoggedIn(
	t *testing.T,
	ctrl UserController,
) {
	req, err := http.NewRequest("GET", "/admin/login", nil)
	if err != nil {
		t.Error("unable to create request")
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.UserLoginGet)

	ctx := context.WithValue(req.Context(), session.IS_LOGGED_IN_CONTEXT_KEY, true)

	handler.ServeHTTP(
		rr,
		req.WithContext(ctx),
	)

	require.Equal(
		t,
		http.StatusSeeOther,
		rr.Result().StatusCode,
		"should return status code see other",
	)

	require.Equal(
		t,
		"/admin/articles",
		rr.Result().Header.Get("Location"),
		"should redirect to articles page",
	)
}

func testGetUserLoginScreenHandlerNotLoggedIn(
	t *testing.T,
	ctrl UserController,
) {
	req, err := http.NewRequest("GET", "/admin/login", nil)
	if err != nil {
		t.Error("unable to create request")
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.UserLoginGet)

	ctx := context.WithValue(req.Context(), session.IS_LOGGED_IN_CONTEXT_KEY, false)

	mockSessionManagerPopString := mockSessionManager.
		On("PopString", ctx, session.SESSION_KEY_MESSAGE).
		Return("msg")

	mockTemplateExecuteTemplate := mockTemplate.On(
		"ExecuteTemplate",
		mock.Anything,
		"admin",
		UserLoginView{
			Message:         "msg",
			CsrfToken:       "mock-token",
			IsAuthenticated: false,
		},
	).Return(nil)

	handler.ServeHTTP(
		rr,
		req.WithContext(ctx),
	)

	require.Equal(
		t,
		http.StatusOK,
		rr.Result().StatusCode,
		"should return status code ok",
	)

	mockTemplateExecuteTemplate.Unset()
	if res := mockTemplate.AssertExpectations(t); !res {
		t.Error("should execute template with view struct")
	}

	mockSessionManagerPopString.Unset()
	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should pop message from session context")
	}
}

func testGetUserLoginScreenHandlerTemplateError(
	t *testing.T,
	ctrl UserController,
) {
	req, err := http.NewRequest("GET", "/admin/login", nil)
	if err != nil {
		t.Error("unable to create request")
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.UserLoginGet)

	ctx := context.WithValue(req.Context(), session.IS_LOGGED_IN_CONTEXT_KEY, false)

	mockSessionManagerPopString := mockSessionManager.
		On("PopString", ctx, session.SESSION_KEY_MESSAGE).
		Return("msg")

	mockTemplateExecuteTemplate := mockTemplate.On(
		"ExecuteTemplate",
		mock.Anything,
		"admin",
		UserLoginView{
			Message:         "msg",
			CsrfToken:       "mock-token",
			IsAuthenticated: false,
		},
	).Return(errors.New("template_error"))

	mockLoggerError := mockLogger.On("Error", "template_error", mock.Anything)

	handler.ServeHTTP(
		rr,
		req.WithContext(ctx),
	)

	require.Equal(
		t,
		http.StatusInternalServerError,
		rr.Result().StatusCode,
		"should return status code internal server error",
	)

	mockTemplateExecuteTemplate.Unset()
	if res := mockTemplate.AssertExpectations(t); !res {
		t.Error("should execute template with view struct")
	}

	mockSessionManagerPopString.Unset()
	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should pop message from session context")
	}

	mockLoggerError.Unset()
	if res := mockLogger.AssertExpectations(t); !res {
		t.Error("should log error")
	}
}

func testPostUserLogin(t *testing.T, ctrl UserController) {
	form := url.Values{}
	form.Add("username", "janedoe")
	form.Add("password", "p4ssw0rd")

	req, err := http.NewRequest(
		"POST",
		"/admin/login",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		t.Error("unable to construct request")
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.UserLoginPost)

	mockServiceLoginUsernamePassword := mockService.On("LoginWithUsernamePassword", "janedoe", "p4ssw0rd").
		Return(nil)

	mockSessionManagerRenewToken := mockSessionManager.On("RenewToken", req.Context()).
		Return(nil)

	mockSessionManagerPut := mockSessionManager.On(
		"Put",
		req.Context(),
		session.LOGGED_IN_USERNAME,
		"janedoe",
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
		"/admin/articles",
		rr.Result().Header.Get("Location"),
		"should redirect to articles page",
	)

	mockServiceLoginUsernamePassword.Unset()
	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call login service")
	}

	mockSessionManagerRenewToken.Unset()
	mockSessionManagerPut.Unset()
	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should call session manager")
	}
}

func testPostUserLoginUsernamePasswordFailed(
	t *testing.T,
	ctrl UserController,
) {
	form := url.Values{}
	form.Add("username", "janedoe")
	form.Add("password", "p4ssw0rd")

	req, err := http.NewRequest(
		"POST",
		"/admin/login",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		t.Error("unable to construct request")
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.UserLoginPost)

	mockServiceLoginUsernamePassword := mockService.On("LoginWithUsernamePassword", "janedoe", "p4ssw0rd").
		Return(errors.New("username_password_error"))

	mockLoggerError := mockLogger.On(
		"Error",
		"username_password_error",
		mock.Anything,
	)

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusUnauthorized,
		rr.Result().StatusCode,
		"should return status code unauthorised",
	)

	mockServiceLoginUsernamePassword.Unset()
	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call login service")
	}

	mockLoggerError.Unset()
	if res := mockLogger.AssertExpectations(t); !res {
		t.Error("should log error")
	}

	mockSessionManager.AssertNotCalled(t, "RenewToken")
}

func testPostUserLoginRenewTokenFailed(t *testing.T, ctrl UserController) {
	form := url.Values{}
	form.Add("username", "janedoe")
	form.Add("password", "p4ssw0rd")

	req, err := http.NewRequest(
		"POST",
		"/admin/login",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		t.Error("unable to construct request")
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.UserLoginPost)

	mockServiceLoginUsernamePassword := mockService.On("LoginWithUsernamePassword", "janedoe", "p4ssw0rd").
		Return(nil)

	mockSessionManagerRenewToken := mockSessionManager.On("RenewToken", req.Context()).
		Return(errors.New("renew_token_error"))

	mockLoggerError := mockLogger.On(
		"Error",
		"renew_token_error",
		mock.Anything,
	)

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusUnauthorized,
		rr.Result().StatusCode,
		"should return status code unauthorised",
	)

	mockServiceLoginUsernamePassword.Unset()
	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call login service")
	}

	mockSessionManagerRenewToken.Unset()
	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should call session manager")
	}

	mockLoggerError.Unset()
	if res := mockLogger.AssertExpectations(t); !res {
		t.Error("should log error")
	}
}

func testPostUserLogout(t *testing.T, ctrl UserController) {
	req, err := http.NewRequest("POST", "/admin/logout", nil)
	if err != nil {
		t.Error("unable to construct request")
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.UserLogoutPost)

	mockSessionManagerRemove := mockSessionManager.On(
		"Remove",
		req.Context(),
		session.LOGGED_IN_USERNAME,
	)

	mockSessionManagerPut := mockSessionManager.On(
		"Put",
		req.Context(),
		session.SESSION_KEY_MESSAGE,
		"You've been logged out.",
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
		"/admin",
		rr.Result().Header.Get("Location"),
		"should redirect to main admin page",
	)

	mockSessionManagerRemove.Unset()
	mockSessionManagerPut.Unset()

	if res := mockSessionManager.AssertCalled(t, "Remove", req.Context(), session.LOGGED_IN_USERNAME); !res {
		t.Error("should remove logged in user from session context")
	}

	if res := mockSessionManager.AssertCalled(t, "Put", req.Context(), session.SESSION_KEY_MESSAGE, "You've been logged out."); !res {
		t.Error("should put message in session context")
	}
}

func testGetCreateUserPage(t *testing.T, ctrl UserController) {
	req, err := http.NewRequest("GET", "/admin/users/new", nil)
	if err != nil {
		t.Error("unable to construct request")
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.CreateUserGet)

	mockTemplateExecuteTemplate := mockTemplate.On(
		"ExecuteTemplate",
		mock.Anything,
		"admin",
		UserCreateView{
			CsrfToken:       "mock-token",
			IsAuthenticated: true,
		},
	).Return(nil)

	ctx := context.WithValue(req.Context(), session.IS_LOGGED_IN_CONTEXT_KEY, true)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	require.Equal(
		t,
		http.StatusOK,
		rr.Result().StatusCode,
		"should return status code ok",
	)

	mockTemplateExecuteTemplate.Unset()
	if res := mockTemplate.AssertExpectations(t); !res {
		t.Error("should execute template with view struct")
	}
}

func testIsAuthenticatedHelper(t *testing.T, ctrl UserController) {
	var res bool

	req, err := http.NewRequest("GET", "/admin", nil)
	if err != nil {
		t.Error("failed to construct request")
	}

	ctx := context.WithValue(req.Context(), session.IS_LOGGED_IN_CONTEXT_KEY, true)
	res = ctrl.IsAuthenticated(req.WithContext(ctx))
	require.True(t, res, "should be authenticated")

	ctx = context.WithValue(req.Context(), session.IS_LOGGED_IN_CONTEXT_KEY, false)
	res = ctrl.IsAuthenticated(req.WithContext(ctx))
	require.False(t, res, "should not be authenticated")

	ctx = context.WithValue(
		req.Context(),
		session.IS_LOGGED_IN_CONTEXT_KEY,
		"nonsense",
	)
	res = ctrl.IsAuthenticated(req.WithContext(ctx))
	require.False(t, res, "should not be authenticated")
}

func testGetCreateUserPageTemplateError(t *testing.T, ctrl UserController) {
	req, err := http.NewRequest("GET", "/admin/users/new", nil)
	if err != nil {
		t.Error("unable to construct request")
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.CreateUserGet)

	mockTemplateExecuteTemplate := mockTemplate.On(
		"ExecuteTemplate",
		mock.Anything,
		"admin",
		UserCreateView{
			CsrfToken:       "mock-token",
			IsAuthenticated: true,
		},
	).Return(errors.New("template_error"))

	mockLoggerError := mockLogger.On("Error", "template_error", mock.Anything)

	ctx := context.WithValue(req.Context(), session.IS_LOGGED_IN_CONTEXT_KEY, true)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	require.Equal(
		t,
		http.StatusInternalServerError,
		rr.Result().StatusCode,
		"should return status code internal server error",
	)

	mockTemplateExecuteTemplate.Unset()
	if res := mockTemplate.AssertExpectations(t); !res {
		t.Error("should execute template with view struct")
	}

	mockLoggerError.Unset()
	if res := mockLogger.AssertExpectations(t); !res {
		t.Error("should log error")
	}
}

func testPostCreateUser(t *testing.T, ctrl UserController) {
	form := url.Values{}
	form.Add("username", "janedoe")
	form.Add("password", "p4ssw0rd")
	form.Add("email", "jane@example.org")

	req, err := http.NewRequest(
		"POST",
		"/admin/users",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		t.Error("failed to construct request")
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.CreateUserPost)

	mockServiceCreate := mockService.On("Create", &UserNewRequestDto{
		Username: "janedoe",
		Password: "p4ssw0rd",
		Email:    "jane@example.org",
	}).Return(&UserResponseDto{
		Id:       23,
		Username: "janedoe",
		Email:    "jane@example.org",
	}, nil)

	mockSessionManagerPut := mockSessionManager.On(
		"Put",
		req.Context(),
		session.SESSION_KEY_MESSAGE,
		"Created user 'janedoe'.",
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
		"/admin/users",
		rr.Result().Header.Get("Location"),
		"should redirect to users admin page",
	)

	mockServiceCreate.Unset()
	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call user create service")
	}

	mockSessionManagerPut.Unset()
	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should put user created message in session context")
	}
}

func testPostCreateUserServiceError(t *testing.T, ctrl UserController) {
	form := url.Values{}
	form.Add("username", "janedoe")
	form.Add("password", "p4ssw0rd")
	form.Add("email", "jane@example.org")

	req, err := http.NewRequest(
		"POST",
		"/admin/users",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		t.Error("failed to construct request")
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.CreateUserPost)

	mockServiceCreate := mockService.On("Create", &UserNewRequestDto{
		Username: "janedoe",
		Password: "p4ssw0rd",
		Email:    "jane@example.org",
	}).Return(&UserResponseDto{}, errors.New("service_error"))

	mockLoggerError := mockLogger.On("Error", "service_error", mock.Anything)

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusInternalServerError,
		rr.Result().StatusCode,
		"should return status code internal server error",
	)

	mockServiceCreate.Unset()
	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call user create service")
	}

	mockLoggerError.Unset()
	if res := mockLogger.AssertExpectations(t); !res {
		t.Error("should log error")
	}
}

func testGetAllUsers(t *testing.T, ctrl UserController) {
	req, err := http.NewRequest("GET", "/admin/users", nil)
	if err != nil {
		t.Error("unable to construct request")
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.UsersGet)

	mockServiceGetAll := mockService.On("GetAll").Return(&[]UserResponseDto{
		{
			Id:       23,
			Username: "janedoe",
		},
		{
			Id:       42,
			Username: "johndoe",
		},
	}, nil)

	mockSessionManagerPopString := mockSessionManager.
		On("PopString", req.Context(), session.SESSION_KEY_MESSAGE).
		Return("msg")

	users := mockServiceGetAll.ReturnArguments[0].(*[]UserResponseDto)

	mockTemplateExecuteTemplate := mockTemplate.On(
		"ExecuteTemplate",
		mock.Anything,
		"admin",
		UsersView{
			Users:           users,
			Message:         "msg",
			CsrfToken:       "mock-token",
			IsAuthenticated: false,
		},
	).Return(nil)

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusOK,
		rr.Result().StatusCode,
		"should return status code ok",
	)

	mockSessionManagerPopString.Unset()
	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should get message from session context")
	}

	mockTemplateExecuteTemplate.Unset()
	if res := mockTemplate.AssertExpectations(t); !res {
		t.Error("should execute template")
	}

	mockServiceGetAll.Unset()
	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call service to get users")
	}
}

func testGetAllUsersServiceError(t *testing.T, ctrl UserController) {
	req, err := http.NewRequest("GET", "/admin/users", nil)
	if err != nil {
		t.Error("unable to construct request")
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.UsersGet)

	mockServiceGetAll := mockService.On("GetAll").
		Return(&[]UserResponseDto{}, errors.New("service_error"))

	mockLoggerError := mockLogger.On("Error", "service_error", mock.Anything)

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusInternalServerError,
		rr.Result().StatusCode,
		"should return status code internal server error",
	)

	mockServiceGetAll.Unset()
	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call service to get users")
	}

	mockLoggerError.Unset()
	if res := mockLogger.AssertExpectations(t); !res {
		t.Error("should log error")
	}
}

func testGetAllUsersTemplateError(t *testing.T, ctrl UserController) {
	req, err := http.NewRequest("GET", "/admin/users", nil)
	if err != nil {
		t.Error("unable to construct request")
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.UsersGet)

	mockServiceGetAll := mockService.On("GetAll").Return(&[]UserResponseDto{
		{
			Id:       23,
			Username: "janedoe",
		},
		{
			Id:       42,
			Username: "johndoe",
		},
	}, nil)

	mockSessionManagerPopString := mockSessionManager.
		On("PopString", req.Context(), session.SESSION_KEY_MESSAGE).
		Return("msg")

	users := mockServiceGetAll.ReturnArguments[0].(*[]UserResponseDto)

	mockTemplateExecuteTemplate := mockTemplate.On(
		"ExecuteTemplate",
		mock.Anything,
		"admin",
		UsersView{
			Users:           users,
			Message:         "msg",
			CsrfToken:       "mock-token",
			IsAuthenticated: false,
		},
	).Return(errors.New("template_error"))

	mockLoggerError := mockLogger.On("Error", "template_error", mock.Anything)

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusInternalServerError,
		rr.Result().StatusCode,
		"should return status code internal server error",
	)

	mockSessionManagerPopString.Unset()
	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should get message from session context")
	}

	mockTemplateExecuteTemplate.Unset()
	if res := mockTemplate.AssertExpectations(t); !res {
		t.Error("should execute template")
	}

	mockServiceGetAll.Unset()
	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call service to get users")
	}

	mockLoggerError.Unset()
	if res := mockLogger.AssertExpectations(t); !res {
		t.Error("should log error")
	}
}

func testGetUserByUsername(t *testing.T, ctrl UserController) {
	req, err := http.NewRequest("GET", "/admin/users/{slug}", nil)
	if err != nil {
		t.Error("unable to construct request")
	}

	req.SetPathValue("slug", "janedoe")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.UserGet)

	mockServiceGetByAttribute := mockService.On("GetByAttribute", "username", "janedoe").
		Return(&UserResponseDto{
			Id:       23,
			Username: "janedoe",
		}, nil)

	mockSessionManagerPopString := mockSessionManager.
		On("PopString", req.Context(), session.SESSION_KEY_MESSAGE).
		Return("msg")

	user := mockServiceGetByAttribute.ReturnArguments[0].(*UserResponseDto)

	mockTemplateExecuteTemplate := mockTemplate.On(
		"ExecuteTemplate",
		mock.Anything,
		"admin",
		UserView{
			User:            user,
			Message:         "",
			CsrfToken:       "mock-token",
			IsAuthenticated: false,
		},
	).Return(nil)

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusOK,
		rr.Result().StatusCode,
		"should return status code ok",
	)

	mockSessionManagerPopString.Unset()
	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should get message from session context")
	}

	mockTemplateExecuteTemplate.Unset()
	if res := mockTemplate.AssertExpectations(t); !res {
		t.Error("should execute template")
	}

	mockServiceGetByAttribute.Unset()
	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call service to get users")
	}
}

func testGetUserByUsernameServiceError(t *testing.T, ctrl UserController) {
	req, err := http.NewRequest("GET", "/admin/users/{slug}", nil)
	if err != nil {
		t.Error("unable to construct request")
	}

	req.SetPathValue("slug", "janedoe")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.UserGet)

	mockServiceGetByAttribute := mockService.On("GetByAttribute", "username", "janedoe").
		Return(&UserResponseDto{}, errors.New("service_error"))

	mockLoggerError := mockLogger.On("Error", "service_error", mock.Anything)

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusInternalServerError,
		rr.Result().StatusCode,
		"should return status code internal server error",
	)

	mockServiceGetByAttribute.Unset()
	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call service to get users")
	}

	mockLoggerError.Unset()
	if res := mockLogger.AssertExpectations(t); !res {
		t.Error("should log error")
	}
}

func testGetUserByUsernameTemplateError(t *testing.T, ctrl UserController) {
	req, err := http.NewRequest("GET", "/admin/users/{slug}", nil)
	if err != nil {
		t.Error("unable to construct request")
	}

	req.SetPathValue("slug", "janedoe")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.UserGet)

	mockServiceGetByAttribute := mockService.On("GetByAttribute", "username", "janedoe").
		Return(&UserResponseDto{
			Id:       23,
			Username: "janedoe",
		}, nil)

	mockSessionManagerPopString := mockSessionManager.
		On("PopString", req.Context(), session.SESSION_KEY_MESSAGE).
		Return("msg")

	user := mockServiceGetByAttribute.ReturnArguments[0].(*UserResponseDto)

	mockTemplateExecuteTemplate := mockTemplate.On(
		"ExecuteTemplate",
		mock.Anything,
		"admin",
		UserView{
			User:            user,
			Message:         "",
			CsrfToken:       "mock-token",
			IsAuthenticated: false,
		},
	).Return(errors.New("template_error"))

	mockLoggerError := mockLogger.On("Error", "template_error", mock.Anything)

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusInternalServerError,
		rr.Result().StatusCode,
		"should return status code internal server error",
	)

	mockSessionManagerPopString.Unset()
	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should get message from session context")
	}

	mockTemplateExecuteTemplate.Unset()
	if res := mockTemplate.AssertExpectations(t); !res {
		t.Error("should execute template")
	}

	mockServiceGetByAttribute.Unset()
	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call service to get users")
	}

	mockLoggerError.Unset()
	if res := mockLogger.AssertExpectations(t); !res {
		t.Error("should log error")
	}
}

func testPostDeleteUser(t *testing.T, ctrl UserController) {
	form := url.Values{}
	form.Add("id", "23")
	form.Add("username", "janedoe")

	req, err := http.NewRequest(
		"POST",
		"/admin/users/{username}/delete",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		t.Error("unable to construct request")
	}

	req.SetPathValue("username", "janedoe")

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.DeleteUserPost)

	mockServiceDeleteById := mockService.On("DeleteById", uint(23)).Return(nil)

	mockSessionManagerPut := mockSessionManager.On(
		"Put",
		req.Context(),
		session.SESSION_KEY_MESSAGE,
		"Deleted user 'janedoe'.",
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
		"/admin/users",
		rr.Result().Header.Get("Location"),
		"should redirect to users admin page",
	)

	mockServiceDeleteById.Unset()
	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call service to delete user")
	}

	mockSessionManagerPut.Unset()
	if res := mockSessionManager.AssertExpectations(t); !res {
		t.Error("should put deleted message in session context")
	}
}

func testPostDeleteUserFormError(t *testing.T, ctrl UserController) {
	form := url.Values{}
	form.Add("id", "nonsense")
	form.Add("username", "janedoe")

	req, err := http.NewRequest(
		"POST",
		"/admin/users/{username}/delete",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		t.Error("unable to construct request")
	}

	req.SetPathValue("username", "janedoe")

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.DeleteUserPost)

	mockLoggerError := mockLogger.On("Error", mock.Anything, mock.Anything)

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusInternalServerError,
		rr.Result().StatusCode,
		"should return status code internal server error",
	)

	mockLoggerError.Unset()
	if res := mockLogger.AssertExpectations(t); !res {
		t.Error("should log error")
	}
}

func testPostDeleteUserServiceError(t *testing.T, ctrl UserController) {
	form := url.Values{}
	form.Add("id", "23")
	form.Add("username", "janedoe")

	req, err := http.NewRequest(
		"POST",
		"/admin/users/{username}/delete",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		t.Error("unable to construct request")
	}

	req.SetPathValue("username", "janedoe")

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctrl.DeleteUserPost)

	mockServiceDeleteById := mockService.On("DeleteById", uint(23)).
		Return(errors.New("service_error"))

	mockLoggerError := mockLogger.On("Error", "service_error", mock.Anything)

	handler.ServeHTTP(rr, req)

	require.Equal(
		t,
		http.StatusInternalServerError,
		rr.Result().StatusCode,
		"should return status code internal server error",
	)

	mockServiceDeleteById.Unset()
	if res := mockService.AssertExpectations(t); !res {
		t.Error("should call service to delete user")
	}

	mockLoggerError.Unset()
	if res := mockLogger.AssertExpectations(t); !res {
		t.Error("should log error")
	}
}
