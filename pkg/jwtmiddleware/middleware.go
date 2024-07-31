package jwtmiddleware

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strings"
	"time"
)

type UserClaims struct {
	jwt.StandardClaims
	UserID   string              `json:"user_id"`
	DeviceID string              `json:"device_id"`
	Email    string              `json:"email"`
	Role     models.UserCategory `json:"role"`
}

type Manager struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

type TokenManager interface {
	ParseToken(signedTokenString string) (*UserClaims, error)
	putClaimsOnContext(ctx context.Context, claims *UserClaims) (context.Context, error)
	ExtractUserClaims(ctx context.Context) (*UserClaims, error)
}

var _ TokenManager = &Manager{}

var AuthenticatedUserMetadataKey = "AuthenticatedUser"

func New(publicKey, privateKey string) (*Manager, error) {
	return generateKey(publicKey, privateKey)
}

func generateKey(publicKey, privateKey string) (*Manager, error) {
	tokenGeneratorPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	tokenGeneratorPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	return &Manager{
		publicKey:  tokenGeneratorPublicKey,
		privateKey: tokenGeneratorPrivateKey,
	}, nil
}

func (handler *Manager) GenerateTokenWithExpiration(claims *UserClaims) (string, error) {
	claims.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(handler.privateKey)
}

// BuildAuthResponse Set user details and generate token
func (handler *Manager) BuildAuthResponse(email, userID, deviceID string, role models.UserCategory) (string, error) {
	claims := UserClaims{
		Email:    email,
		UserID:   userID,
		DeviceID: deviceID,
		Role:     role,
	}
	return handler.GenerateTokenWithExpiration(&claims)
}

func (claims *UserClaims) Valid() error {
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return fmt.Errorf("token has expired")
	}
	return nil
}

func (handler *Manager) ParseToken(signedTokenString string) (*UserClaims, error) {
	t, err := jwt.ParseWithClaims(signedTokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, errors.New("invalid signing algorithm")
		}
		return handler.publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("invalid token")
	}
	return t.Claims.(*UserClaims), nil
}

// PutClaimsOnContext put user token on context
func (handler *Manager) putClaimsOnContext(ctx context.Context, claims *UserClaims) (context.Context, error) {
	jsonClaims, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}

	return metadata.AppendToOutgoingContext(ctx, AuthenticatedUserMetadataKey, string(jsonClaims)), nil
}

// ValidateMiddleware middleware required endpoints: verify claims and put claims on context
func (handler *Manager) ValidateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader != "" {
			handler.validateHeaderToken(authorizationHeader, next, w, r, false)
		} else {
			errMsg := errors.New("no token in authorization header")
			WriteJSONResponse(w, errs.Body(errs.ErrorUnauthorized, errMsg), http.StatusUnauthorized)
			return
		}

	})
}

// ValidateRestrictedAccessMiddleware middleware required endpoints: verify claims
// extensively check if they have superior access to these endpoints
// and put claims on context
func (handler *Manager) ValidateRestrictedAccessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader != "" {
			handler.validateHeaderToken(authorizationHeader, next, w, r, true)
		} else {
			WriteJSONResponse(w, errs.Body(errs.ErrorUnauthorized, errors.New("no token in authorization header")), http.StatusUnauthorized)
			return
		}

	})
}

func (handler *Manager) validateHeaderToken(authorizationHeader string, next http.Handler, w http.ResponseWriter, r *http.Request, isAdminPrivileged bool) {
	bearerToken := strings.Split(authorizationHeader, " ")
	if len(bearerToken) == 1 {
		WriteJSONResponse(w, errs.Body(errs.ErrorUnauthorized, errors.New("malformed token in authorization header")), http.StatusUnauthorized)
	}

	if len(bearerToken) == 2 {
		if bearerToken[1] == "" {
			log.Error().Msg("bearer token is empty")
			next.ServeHTTP(w, r)
			return
		}

		claims, err := handler.ParseToken(bearerToken[1])
		if err != nil {
			log.Error().Msgf("unable to parse token string: %v", err)
			WriteJSONResponse(w, errs.Body(errs.ErrorUnauthorized, err), http.StatusUnauthorized)
			return
		}

		if isAdminPrivileged {
			// validate that user if user has permission to access the endpoint
			if claims.Role == models.CustomerCategory {
				WriteJSONResponse(w, errs.Body(errs.RestrictedAccessError, err), http.StatusUnauthorized)
				return
			}
		}

		ctx, _ := handler.putClaimsOnContext(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))

	}

}

// ExtractUserClaims returns claims from an authenticated user
func (handler *Manager) ExtractUserClaims(ctx context.Context) (*UserClaims, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return nil, errors.New("unable to parse authenticated user")
	}
	jsonClaims := md.Get(AuthenticatedUserMetadataKey)
	if len(jsonClaims) == 0 {
		return nil, errors.New("fail decode authenticated user claims")
	}

	var claims UserClaims
	err := json.Unmarshal([]byte(jsonClaims[0]), &claims)
	if err != nil {
		return nil, errors.New("fail to unmarshall claims")
	}

	return &claims, nil
}

func WriteJSONResponse(w http.ResponseWriter, result any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	data := struct {
		Data any `json:"data"`
	}{
		Data: result,
	}

	if result != nil {
		err := json.NewEncoder(w).Encode(&data)
		if err != nil {
			log.Err(err).Msg("fail to encode result")
			return
		}
	}

}

func WriteJSONErrorResponse(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	data := struct {
		Data any `json:"data"`
	}{
		Data: err,
	}

	newErr := json.NewEncoder(w).Encode(&data)
	if newErr != nil {
		log.Err(newErr).Msg("fail to encode result")
		return
	}
}
