package handlers

import (
	"net/http"
)

// NewLegalHandler creates a new instance of LegalHandler
func NewLegalHandler() *LegalHandler {
	return &LegalHandler{}
}

// GetTerms handles GET /terms requests
func (h *LegalHandler) GetTerms(w http.ResponseWriter, r *http.Request) {
	const terms = "./docs/legal/policy/TERMS.html"

	http.ServeFile(w, r, terms)
}

// GetLicense handles GET /license requests
func (h *LegalHandler) GetLicense(w http.ResponseWriter, r *http.Request) {
	const license = "./docs/legal/LICENSE"

	http.ServeFile(w, r, license)
}
