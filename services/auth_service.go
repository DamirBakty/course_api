package services

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/big"
	"net/http"
	"strings"
	"time"
	"web/config"
	"web/repos"
)

// KeycloakClaims represents the claims in a Keycloak JWT token
type KeycloakClaims struct {
	jwt.RegisteredClaims
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAccess map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	Sub               string `json:"sub"`
}

// AuthService handles authentication and authorization with Keycloak
type AuthService struct {
	config        *config.AppConfig
	jwksURL       string
	keysCache     map[string]interface{}
	keysCacheTime time.Time
	userRepo      repos.UserRepositoryInterface
}

func NewAuthService(config *config.AppConfig, userRepo repos.UserRepositoryInterface) *AuthService {
	jwksURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs",
		config.KeycloakURL, config.KeycloakRealm)

	return &AuthService{
		config:    config,
		jwksURL:   jwksURL,
		keysCache: make(map[string]interface{}),
		userRepo:  userRepo,
	}
}

// ValidateToken validates a JWT token from Keycloak
func (s *AuthService) ValidateToken(tokenString string) (*KeycloakClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &KeycloakClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get the key ID from the token header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid header not found in token")
		}

		// Get the public key for this kid
		key, err := s.getPublicKey(kid)
		if err != nil {
			return nil, err
		}

		return key, nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Get the claims
	claims, ok := token.Claims.(*KeycloakClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

// ExtractToken extracts the token from the Authorization header
func (s *AuthService) ExtractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	// Check if the Authorization header has the Bearer prefix
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}

// getPublicKey gets the public key for the given key ID
func (s *AuthService) getPublicKey(kid string) (interface{}, error) {
	// Check if we have the key in cache and if the cache is still valid (1 hour)
	if key, ok := s.keysCache[kid]; ok && time.Since(s.keysCacheTime) < time.Hour {
		return key, nil
	}

	// Fetch the JWKS from Keycloak
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", s.jwksURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get JWKS: %s", resp.Status)
	}

	// Parse the JWKS
	var jwks struct {
		Keys []struct {
			Kid string   `json:"kid"`
			Kty string   `json:"kty"`
			Alg string   `json:"alg"`
			Use string   `json:"use"`
			N   string   `json:"n"`
			E   string   `json:"e"`
			X5c []string `json:"x5c"`
		} `json:"keys"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	// Find the key with the matching kid
	for _, key := range jwks.Keys {
		if key.Kid == kid {
			// Convert the key to a *rsa.PublicKey
			if key.Kty != "RSA" {
				return nil, fmt.Errorf("key type %s not supported", key.Kty)
			}

			// Decode the modulus and exponent
			n, err := base64.RawURLEncoding.DecodeString(key.N)
			if err != nil {
				return nil, err
			}

			e, err := base64.RawURLEncoding.DecodeString(key.E)
			if err != nil {
				return nil, err
			}

			// Convert the modulus to a big int
			modulus := new(big.Int)
			modulus.SetBytes(n)

			// Convert the exponent to an int
			var exponent int
			for i := 0; i < len(e); i++ {
				exponent = exponent*256 + int(e[i])
			}

			// Create the public key
			publicKey := &rsa.PublicKey{
				N: modulus,
				E: exponent,
			}

			// Cache the key
			s.keysCache[kid] = publicKey
			s.keysCacheTime = time.Now()

			return publicKey, nil
		}
	}

	return nil, fmt.Errorf("key with ID %s not found", kid)
}

// HasRole checks if the user has the specified role
func (s *AuthService) HasRole(claims *KeycloakClaims, role string) bool {
	// Check realm roles
	for _, r := range claims.RealmAccess.Roles {
		if r == role {
			return true
		}
	}

	// Check client roles
	if clientRoles, ok := claims.ResourceAccess[s.config.KeycloakClientID]; ok {
		for _, r := range clientRoles.Roles {
			if r == role {
				return true
			}
		}
	}

	return false
}

// ValidateSession checks if a user with the given sub exists in the database
func (s *AuthService) ValidateSession(sub string) (bool, error) {
	if sub == "" {
		return false, errors.New("sub is required")
	}

	_, err := s.userRepo.GetBySub(sub)
	if err != nil {
		if err.Error() == "user not found" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
