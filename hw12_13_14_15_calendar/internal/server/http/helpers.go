package internalhttp

import (
	"encoding/json"
	"net/http"
)

func methodHandler(expectedMethod string, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != expectedMethod {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlerFunc(w, r)
	}
}

type wrapper map[string]interface{}

func (s *Server) writeJSON(w http.ResponseWriter, status int, data wrapper) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
func (s *Server) errorResponse(w http.ResponseWriter, status int, message interface{}) {
	wr := wrapper{"error": message}
	err := s.writeJSON(w, status, wr)
	if err != nil {
		s.logger.Info("failed error responel", "error", err)
		w.WriteHeader(500)
	}
}
