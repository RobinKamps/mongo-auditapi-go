package api

import (
	"encoding/json"
	"log"
	"net/http"

	"mongo-auditapi/pkg/config"
	"mongo-auditapi/pkg/db"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FieldAuditService provides the API function for fetching field-level audit records for records that
// are audited.
type FieldAuditService struct {
	DataAccess db.AuditFetcher
	Config     config.Configuration
}

// InitializeRoutes sets up the URL routes supported by the FieldAuditService API.
func (s *FieldAuditService) InitializeRoutes() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/health", s.Health).Methods(http.MethodGet)
	r.Handle("/auditrecords/{documentKey}/{fieldPath}", s.GetFieldAudit()).Methods(http.MethodGet)

	return r
}

// GetFieldAudit fetches field-level audit records for records that are audited.
func (s *FieldAuditService) GetFieldAudit() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var documentKey = vars["documentKey"]
		id, err := primitive.ObjectIDFromHex(documentKey)
		if err != nil {
			log.Printf("ERROR, supplied document key %s is invalid, encountered error: %v\n", documentKey, err)
			respondWithError(w, http.StatusBadRequest, "supplied document key is invalid")
			return
		}
		var fieldPath = vars["fieldPath"]
		atArr, err := s.DataAccess.GetFieldAuditTrail(s.Config.AppDatabase, s.Config.AppCollection, id, fieldPath)
		if err != nil {
			log.Printf("ERROR, failed to fetch audit trail for field %s for document key %s, encountered error: %v\n", fieldPath, documentKey, err)
			respondWithError(w, http.StatusInternalServerError, "failed to fetch audit trail")
			return
		}

		respondWithJSON(w, http.StatusOK, atArr)
	})
}

// Health verifies if the API is in a healthy state.
func (s *FieldAuditService) Health(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, "Status: 200 OK, Version: "+s.Config.Version)
}

// respondWithError converts an error message into a JSON response.
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

// respondWithJson converts the payload into a JSON response.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
