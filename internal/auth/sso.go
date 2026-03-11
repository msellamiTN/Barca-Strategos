package auth

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/barca-strategos/phoenix/pkg/types"
)

// SAMLRequest represents a SAML AuthnRequest.
type SAMLRequest struct {
	XMLName xml.Name `xml:"AuthnRequest"`
	ID      string   `xml:"ID,attr"`
	IssueInstant string `xml:"IssueInstant,attr"`
	Version string   `xml:"Version,attr"`
	Destination string `xml:"Destination,attr"`
}

// SSOHandler manages SAML/OIDC SSO flows.
type SSOHandler struct {
	AuthSvc *Service
	SPEntityID string
	SSOURL     string
}

// NewSSOHandler creates an SSO handler.
func NewSSOHandler(authSvc *Service, spEntityID, ssoURL string) *SSOHandler {
	return &SSOHandler{
		AuthSvc:    authSvc,
		SPEntityID: spEntityID,
		SSOURL:     ssoURL,
	}
}

// InitiateSAML starts SAML SSO.
func (s *SSOHandler) InitiateSAML(w http.ResponseWriter, r *http.Request) {
	// Redirect to IdP with SAML request
	redirectURL := fmt.Sprintf("%s?SAMLRequest=%s", s.SSOURL, url.QueryEscape("base64-encoded-request"))
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// ConsumeSAML processes SAML response.
func (s *SSOHandler) ConsumeSAML(w http.ResponseWriter, r *http.Request) {
	// Parse SAMLResponse from POST
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	samlResponse := r.FormValue("SAMLResponse")
	// In production, validate signature and decrypt assertion
	// For demo, extract attributes (mock)
	userID := "user123"
	role := types.RoleAnalyst
	token, err := s.AuthSvc.GenerateToken(userID, role)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

// InitiateOIDC starts OIDC flow.
func (s *SSOHandler) InitiateOIDC(w http.ResponseWriter, r *http.Request) {
	// Build OIDC auth URL
	authURL := fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s&scope=openid profile",
		s.SSOURL, "phoenix", url.QueryEscape("http://localhost:3000/auth/oidc/callback"))
	http.Redirect(w, r, authURL, http.StatusFound)
}

// OIDCCallback handles OIDC callback.
func (s *SSOHandler) OIDCCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	// Exchange code for token (mock)
	userID := "user123"
	role := types.RoleAnalyst
	token, _ := s.AuthSvc.GenerateToken(userID, role)
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

// ValidateJWTFromCookie extracts and validates JWT from cookie.
func (s *SSOHandler) ValidateJWTFromCookie(r *http.Request) (*Claims, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return nil, fmt.Errorf("no jwt cookie")
	}
	token := strings.TrimPrefix(cookie.Value, "Bearer ")
	return s.AuthSvc.ValidateToken(token)
}
