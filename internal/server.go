package internal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type ServerOpts struct {
}

type Server struct {
	config        BornoConfig
	logger        *logrus.Logger
	templateStore TemplateStore
	configOpts    ConfigOpts
}

func NewServer(opts *ConfigOpts) (*Server, error) {
	templateStore, err := NewInBuiltTemplateStore()
	if err != nil {
		return nil, err
	}

	return &Server{config: opts.Config, logger: opts.Logger, templateStore: templateStore, configOpts: *opts}, nil
}

func (s Server) Run() error {
	mux := http.NewServeMux()

	// fs := http.FileServer(http.Dir(s.configOpts.AssetsRoot))
	// mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", s.handlerFunction)
	return http.ListenAndServe(":8000", mux)
}

func (s Server) handlerFunction(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)

	path := r.URL.Path

	if path == "/" {
		s.handleIndex(w, r)
		return
	}

	if s.isSlideURL(r.URL.Path) {
		s.handleSlideURL(w, r)
		return
	}

	http.NotFound(w, r)
}

func (s Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	responseData := struct {
		Title  string
		Author string
		Links  []Link
		Groups []TalkGroup
	}{
		Title:  s.config.PageTitle,
		Author: s.config.Author,
		Links:  s.config.Links,
		Groups: groupByYear(s.config.Talks),
	}

	template, err := s.templateStore.GetIndex()
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	if err := template.Execute(w, responseData); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
}

func (s Server) handleSlideURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	slugGroups := groupBySlug(s.config.Talks)
	requested := r.URL.Path
	if _, ok := slugGroups[requested]; !ok {
		http.NotFound(w, r)
	}

	talk := slugGroups[requested]

	if !strings.HasPrefix(talk.SlideURL, "https://docs.google.com/presentation/d/e/") {
		http.Redirect(w, r, talk.SlideURL, http.StatusFound)
		return
	}

	responseData := struct {
		Title    string
		Location string
		Slides   string
	}{
		Title:    talk.Title,
		Location: talk.Location,
		Slides:   talk.SlideURL,
	}

	template, err := s.templateStore.GetSlide()
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	if err := template.Execute(w, responseData); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
}

func (s Server) isSlideURL(path string) bool {
	if strings.HasPrefix(path, "/slides") {
		return true
	}

	return false
}

func (s Server) logRequest(r *http.Request) {
	s.logger.Printf("Requested: %v", r.URL.Path)
}
