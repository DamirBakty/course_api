package services

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"
	"web/config"
	"web/models"
	"web/repos"
	"web/schemas"
)

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

func (s *AuthService) ValidateToken(tokenString string) (*KeycloakClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &KeycloakClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid header not found in token")
		}

		key, err := s.getPublicKey(kid)
		if err != nil {
			return nil, err
		}

		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*KeycloakClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	expectedIssuerPrefix := fmt.Sprintf("%s/realms/%s", s.config.KeycloakURL, s.config.KeycloakRealm)
	if claims.Issuer == "" || !strings.HasPrefix(claims.Issuer, expectedIssuerPrefix) {
		return nil, fmt.Errorf("invalid token issuer: expected issuer to start with %s", expectedIssuerPrefix)
	}

	return claims, nil
}

func (s *AuthService) IntrospectToken(tokenString string) (bool, error) {
	if tokenString == "" {
		return false, errors.New("token is required")
	}

	introspectionURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token/introspect",
		s.config.KeycloakURL, s.config.KeycloakRealm)
	formData := url.Values{}
	formData.Set("token", tokenString)
	formData.Set("client_id", s.config.KeycloakClientID)
	formData.Set("client_secret", s.config.KeycloakClientSecret)

	req, err := http.NewRequest("POST", introspectionURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return false, fmt.Errorf("failed to create introspection request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to send introspection request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("token introspection failed: %s (status code: %d)", string(body), resp.StatusCode)
	}

	var introspectionResponse struct {
		Active bool `json:"active"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&introspectionResponse); err != nil {
		return false, fmt.Errorf("failed to parse introspection response: %w", err)
	}

	return introspectionResponse.Active, nil
}

func (s *AuthService) ExtractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}

func (s *AuthService) getPublicKey(kid string) (interface{}, error) {

	if key, ok := s.keysCache[kid]; ok && time.Since(s.keysCacheTime) < time.Hour {
		return key, nil
	}

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

	for _, key := range jwks.Keys {
		if key.Kid == kid {

			if key.Kty != "RSA" {
				return nil, fmt.Errorf("key type %s not supported", key.Kty)
			}

			n, err := base64.RawURLEncoding.DecodeString(key.N)
			if err != nil {
				return nil, err
			}

			e, err := base64.RawURLEncoding.DecodeString(key.E)
			if err != nil {
				return nil, err
			}

			modulus := new(big.Int)
			modulus.SetBytes(n)

			var exponent int
			for i := 0; i < len(e); i++ {
				exponent = exponent*256 + int(e[i])
			}

			publicKey := &rsa.PublicKey{
				N: modulus,
				E: exponent,
			}

			s.keysCache[kid] = publicKey
			s.keysCacheTime = time.Now()

			return publicKey, nil
		}
	}

	return nil, fmt.Errorf("key with ID %s not found", kid)
}

func (s *AuthService) HasRole(claims *KeycloakClaims, role string) bool {

	for _, r := range claims.RealmAccess.Roles {
		if r == role {
			return true
		}
	}

	if clientRoles, ok := claims.ResourceAccess[s.config.KeycloakClientID]; ok {
		for _, r := range clientRoles.Roles {
			if r == role {
				return true
			}
		}
	}

	return false
}

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

func (s *AuthService) GetUserBySub(sub string) (models.User, error) {
	if sub == "" {
		return models.User{}, errors.New("sub is required")
	}

	return s.userRepo.GetBySub(sub)
}

func (s *AuthService) CreateUser(user models.User) (models.User, error) {
	return s.userRepo.Create(user)
}

func (s *AuthService) RegisterUserInKeycloak(username, email, password string, roles []string) error {

	tokenURL := fmt.Sprintf("%s/realms/master/protocol/openid-connect/token", s.config.KeycloakURL)

	formData := url.Values{}
	formData.Set("grant_type", "password")
	formData.Set("client_id", "admin-cli")
	formData.Set("username", s.config.KeycloakAdminUsername)
	formData.Set("password", s.config.KeycloakAdminPassword)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create admin token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send admin token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("admin authentication failed: %s (status code: %d)", string(body), resp.StatusCode)
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return fmt.Errorf("failed to parse admin token response: %w", err)
	}

	userURL := fmt.Sprintf("%s/admin/realms/%s/users", s.config.KeycloakURL, s.config.KeycloakRealm)

	userData := map[string]interface{}{
		"username": username,
		"email":    email,
		"enabled":  true,
		"credentials": []map[string]interface{}{
			{
				"type":      "password",
				"value":     password,
				"temporary": false,
			},
		},
	}

	userDataJSON, err := json.Marshal(userData)
	if err != nil {
		return fmt.Errorf("failed to marshal user data: %w", err)
	}

	req, err = http.NewRequest("POST", userURL, bytes.NewBuffer(userDataJSON))
	if err != nil {
		return fmt.Errorf("failed to create user creation request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenResponse.AccessToken)

	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send user creation request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {

		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("user creation failed: %s (status code: %d)", string(body), resp.StatusCode)
	}

	if len(roles) > 0 {

		getUserURL := fmt.Sprintf("%s/admin/realms/%s/users?username=%s", s.config.KeycloakURL, s.config.KeycloakRealm, username)

		req, err = http.NewRequest("GET", getUserURL, nil)
		if err != nil {
			return fmt.Errorf("failed to create get user request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+tokenResponse.AccessToken)

		resp, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send get user request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("get user failed: %s (status code: %d)", string(body), resp.StatusCode)
		}

		var users []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
			return fmt.Errorf("failed to parse get user response: %w", err)
		}

		if len(users) == 0 {
			return fmt.Errorf("user not found after creation")
		}

		userID := users[0]["id"].(string)

		getRolesURL := fmt.Sprintf("%s/admin/realms/%s/roles", s.config.KeycloakURL, s.config.KeycloakRealm)

		req, err = http.NewRequest("GET", getRolesURL, nil)
		if err != nil {
			return fmt.Errorf("failed to create get roles request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+tokenResponse.AccessToken)

		resp, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send get roles request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("get roles failed: %s (status code: %d)", string(body), resp.StatusCode)
		}

		var availableRoles []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&availableRoles); err != nil {
			return fmt.Errorf("failed to parse get roles response: %w", err)
		}

		var rolesToAssign []map[string]interface{}
		for _, role := range availableRoles {
			roleName := role["name"].(string)
			for _, requestedRole := range roles {
				if roleName == requestedRole {
					rolesToAssign = append(rolesToAssign, role)
					break
				}
			}
		}

		if len(rolesToAssign) > 0 {

			assignRolesURL := fmt.Sprintf("%s/admin/realms/%s/users/%s/role-mappings/realm", s.config.KeycloakURL, s.config.KeycloakRealm, userID)

			rolesJSON, err := json.Marshal(rolesToAssign)
			if err != nil {
				return fmt.Errorf("failed to marshal roles data: %w", err)
			}

			req, err = http.NewRequest("POST", assignRolesURL, bytes.NewBuffer(rolesJSON))
			if err != nil {
				return fmt.Errorf("failed to create assign roles request: %w", err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tokenResponse.AccessToken)

			resp, err = client.Do(req)
			if err != nil {
				return fmt.Errorf("failed to send assign roles request: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("assign roles failed: %s (status code: %d)", string(body), resp.StatusCode)
			}
		}
	}

	return nil
}

func (s *AuthService) GetUserRepo() repos.UserRepositoryInterface {
	return s.userRepo
}

func (s *AuthService) Login(username, password string, service UserService) (*schemas.LoginResponse, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
		s.config.KeycloakURL, s.config.KeycloakRealm)
	fmt.Println(s.config.KeycloakClientID)

	formData := url.Values{}
	formData.Set("grant_type", "password")
	formData.Set("client_id", s.config.KeycloakClientID)
	formData.Set("client_secret", s.config.KeycloakClientSecret)
	formData.Set("username", username)
	formData.Set("password", password)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("authentication failed: %s (status code: %d)", string(body), resp.StatusCode)
	}

	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	claims, err := s.ValidateToken(tokenResponse.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	_, err = service.ClaimUserUserFromToken(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to validate claims: %w", err)
	}

	loginResponse := &schemas.LoginResponse{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		TokenType:    tokenResponse.TokenType,
		ExpiresIn:    tokenResponse.ExpiresIn,
	}

	return loginResponse, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*schemas.LoginResponse, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
		s.config.KeycloakURL, s.config.KeycloakRealm)

	formData := url.Values{}
	formData.Set("grant_type", "refresh_token")
	formData.Set("client_id", s.config.KeycloakClientID)
	formData.Set("client_secret", s.config.KeycloakClientSecret)
	formData.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token refresh failed: %s (status code: %d)", string(body), resp.StatusCode)
	}

	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	loginResponse := &schemas.LoginResponse{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		TokenType:    tokenResponse.TokenType,
		ExpiresIn:    tokenResponse.ExpiresIn,
	}

	return loginResponse, nil
}
