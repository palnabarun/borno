package internal

import (
	"net/http"

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

func (s Server) logRequest(r *http.Request) {
	s.logger.Printf("Requested: %v", r.URL.Path)
}
