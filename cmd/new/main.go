package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

type module struct {
	Name      string
	LowerName string
	Path      string
	tp        string
}

func main() {
	name := flag.String("module", "foo", "new module")
	flag.Parse()

	var tasks []module
	tasks = append(tasks, module{strings.Title(*name), *name, *name + ".go", "base"})
	tasks = append(tasks, module{strings.Title(*name), *name, "mysql/" + *name + ".go", "model"})
	tasks = append(tasks, module{strings.Title(*name), *name, "http/" + *name + ".go", "http"})

	for _, task := range tasks {
		if _, err := os.Stat(task.Path); err == nil {
			fmt.Printf("file %s is existed, skip...\n", task.Path)
			continue
		}

		f, err := os.Create(task.Path)
		if err != nil {
			log.Fatal(err)
		}

		var t *template.Template
		switch task.tp {
		case "base":
			t = template.Must(template.New(task.LowerName).Parse(tpl_base))
		case "model":
			t = template.Must(template.New(task.LowerName).Parse(tpl_model))
		case "http":
			t = template.Must(template.New(task.LowerName).Parse(tpl_http))
		}

		if err = t.Execute(f, task); err != nil {
			log.Fatal(err)
		}
	}
}

var (
	tpl_base = `package api

// {{ .Name }}
type {{ .Name }} struct {
	Model

	// your field
}

// {{ .Name }}Service
type {{ .Name }}Service interface {
	Find{{ .Name }}ByKV(key string, val interface{}) (*{{ .Name }}, error)
	Find{{ .Name }}s(filter {{ .Name }}Filter) ([]*{{ .Name }}, int, error)
	Create{{ .Name }}({{ .LowerName }} *{{ .Name }}) error
}

type {{ .Name }}Filter struct {
}`

	tpl_model = `package mysql

import (
	"github.com/lukedever/api"
	"gorm.io/gorm"
)

var _ api.{{ .Name }}Service = (*{{ .Name }}Service)(nil)

// {{ .Name }}Service
type {{ .Name }}Service struct {
	db *gorm.DB
}

// New{{ .Name }}Service
func New{{ .Name }}Service(db *gorm.DB) *{{ .Name }}Service {
	return &{{ .Name }}Service{db: db}
}

// Find{{ .Name }}ByKV
func (s *{{ .Name }}Service) Find{{ .Name }}ByKV(key string, val interface{}) (*api.{{ .Name }}, error) {
	var {{ .LowerName}} api.{{ .Name }}
	r := s.db.Where(key+" = ?", val).First(&{{ .LowerName}})
	return &{{ .LowerName}}, r.Error
}

// Find{{ .Name }}s
func (s *{{ .Name }}Service) Find{{ .Name }}s(filter api.{{ .Name }}Filter) ([]*api.{{ .Name }}, int, error) {
	return nil, 0, nil
}

// Create{{ .Name }}
func (s *{{ .Name }}Service) Create{{ .Name }}({{ .LowerName}} *api.{{ .Name }}) error {
	r := s.db.Create({{ .LowerName}})
	return r.Error
}`

	tpl_http = `package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handle{{ .Name }}
func (s *Server) Handle{{ .Name }}(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}`
)
