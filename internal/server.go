package internal

import (
	"html/template"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ServerOpts struct {
	Config   BornoConfig
	Template *template.Template
	Logger   *logrus.Logger
	Host     string
	Port     int
}

type Server struct {
	config   BornoConfig
	logger   *logrus.Logger
	template *template.Template
}

func NewServer(opts *ServerOpts) (*Server, error) {
	return &Server{config: opts.Config, logger: opts.Logger, template: opts.Template}, nil
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

	if err := s.template.Execute(w, responseData); err != nil {
		s.logger.Errorln(err)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("ERROR"))

		return
	}
}

func (s Server) logRequest(r *http.Request) {
	s.logger.Printf("Requested: %v", r.URL.Path)
}
