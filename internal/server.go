package internal

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type ServerOpts struct {
	Config           BornoConfig
	TemplateLocation string
	Logger           *logrus.Logger
	Host             string
	Port             int
}

type Server struct {
	config        BornoConfig
	logger        *logrus.Logger
	templateStore TemplateStore
}

func NewServer(opts *ServerOpts) (*Server, error) {
	templateStore, err := NewInBuiltTemplateStore()
	if err != nil {
		return nil, err
	}

	return &Server{config: opts.Config, logger: opts.Logger, templateStore: templateStore}, nil
}

func (s Server) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handlerFunction)
	return http.ListenAndServe(":8000", mux)
}

func (s Server) handlerFunction(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)

	if s.isSlideURL(r.URL.Path) {
		s.handleSlideURL(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")

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
		s.processError(w, r, err)
		return
	}

	if err := template.Execute(w, responseData); err != nil {
		s.processError(w, r, err)
		return
	}
}

func (s Server) processError(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Errorln(err)

	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("ERROR"))

	return
}

func (s Server) processNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("NOT FOUND"))

	return
}

func (s Server) processRedirect(w http.ResponseWriter, r *http.Request, location string) {
	http.Redirect(w, r, location, http.StatusFound)
	return
}

func (s Server) handleSlideURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	slugGroups := groupBySlug(s.config.Talks)
	requested := r.URL.Path
	if _, ok := slugGroups[requested]; !ok {
		s.processNotFound(w, r)
	}

	talk := slugGroups[requested]

	if !strings.HasPrefix(talk.SlideURL, "https://docs.google.com/presentation/d/e/") {
		s.processRedirect(w, r, talk.SlideURL)
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
		s.processError(w, r, err)
		return
	}

	if err := template.Execute(w, responseData); err != nil {
		s.processError(w, r, err)
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
