package tag

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"testing"

	"github.com/nixpig/dunce/pkg"
	"github.com/stretchr/testify/mock"
)

type MockTagService struct {
	mock.Mock
}

func (s *MockTagService) Create(tag *TagData) (*Tag, error) {
	args := s.Called(tag)

	return args.Get(0).(*Tag), args.Error(1)
}

func (s *MockTagService) DeleteById(id int) error {
	args := s.Called(id)

	return args.Error(0)
}

func (s *MockTagService) GetAll() (*[]Tag, error) {
	args := s.Called()

	return args.Get(0).(*[]Tag), args.Error(1)
}

func (s *MockTagService) GetByAttribute(attr, slug string) (*Tag, error) {
	args := s.Called(slug)

	return args.Get(0).(*Tag), args.Error(1)
}

func (s *MockTagService) Update(tag *Tag) (*Tag, error) {
	args := s.Called(tag)

	return args.Get(0).(*Tag), args.Error(1)
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
		path.Join(pwd, "..", "..", "test", "templates", "base.tmpl"),
	)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	return ts
}

var mockTemplateCache = map[string]*template.Template{
	"admin-new-tag.tmpl": mockTemplate(),
	"admin-tags.tmpl":    mockTemplate(),
}

func TestTagsControllerNewHandler(t *testing.T) {
	scenarios := map[string]func(t *testing.T, ctrl TagController){
		"test handle get new tag":  testGetAdminTagsNewHandler,
		"test handle get all tags": testGetAdminTagsHandler,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			config := pkg.ControllerConfig{
				Log:           mockLogger,
				TemplateCache: mockTemplateCache,
			}

			ctrl := NewTagController(mockService, config)

			fn(t, ctrl)
		})
	}

}

func testGetAdminTagsNewHandler(t *testing.T, ctrl TagController) {
	// req, err := http.NewRequest("GET", "/admin/tags/create", nil)
	// if err != nil {
	// 	t.Fatal("failed to construct request", err)
	// }
	//
	// rr := httptest.NewRecorder()
	// handler := http.HandlerFunc(ctrl.GetAdminTagsNewHandler)
	//
	// handler.ServeHTTP(rr, req)
	//
	// if status := rr.Code; status != http.StatusOK {
	// 	t.Errorf("status not ok")
	// }
	//
}

// FIXME: this doesn't really test anything
func testGetAdminTagsHandler(t *testing.T, ctrl TagController) {
	// mockService.On("GetAll").Return(&[]Tag{
	// 	{
	// 		Id: 23,
	// 		TagData: TagData{
	// 			Name: "Go",
	// 			Slug: "golang",
	// 		},
	// 	},
	// 	{
	// 		Id: 69,
	// 		TagData: TagData{
	// 			Name: "Rust",
	// 			Slug: "rust-lang",
	// 		},
	// 	},
	// }, nil)
	//
	// req, err := http.NewRequest("GET", "/admin/tags", nil)
	// if err != nil {
	// 	t.Fatal("failed to construct request")
	// }
	//
	// rr := httptest.NewRecorder()
	// handler := http.HandlerFunc(ctrl.GetAdminTagsHandler)
	//
	// handler.ServeHTTP(rr, req)
}
