package library

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/leetatech/leeta_backend/services/library/models"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"strings"
	"time"
)

type UserClaims struct {
	jwt.StandardClaims
	SessionID string              `json:"session_id"`
	UserID    string              `json:"user_id"`
	Email     string              `json:"email"`
	Role      models.UserCategory `json:"role"`
}

type TokenHandler struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

type TokenManager interface {
	ParseToken(signedTokenString string) (*UserClaims, error)
	putClaimsOnContext(ctx context.Context, claims *UserClaims) (context.Context, error)
	GetClaimsFromCtx(ctx context.Context) (*UserClaims, error)
}

var _ TokenManager = &TokenHandler{}

var AuthenticatedUserMetadataKey = "AuthenticatedUser"

func NewMiddlewares(publicKey, privateKey string) (*TokenHandler, error) {
	return generateKey(publicKey, privateKey)
}

func generateKey(publicKey, privateKey string) (*TokenHandler, error) {
	tokenGeneratorPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	tokenGeneratorPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	return &TokenHandler{
		publicKey:  tokenGeneratorPublicKey,
		privateKey: tokenGeneratorPrivateKey,
	}, nil
}

func (handler *TokenHandler) GenerateTokenWithExpiration(claims *UserClaims) (string, error) {
	claims.ExpiresAt = time.Now().Add(time.Minute * 15).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(handler.privateKey)
}

// BuildAuthResponse Set user details and generate token
func (handler *TokenHandler) BuildAuthResponse(email, userID, sessionID string, role models.UserCategory) (string, error) {
	claims := UserClaims{
		SessionID: sessionID,
		Email:     email,
		UserID:    userID,
		Role:      role,
	}
	return handler.GenerateTokenWithExpiration(&claims)
}

func (claims *UserClaims) Valid() error {
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return fmt.Errorf("token has expired")
	}
	return nil
}

func (handler *TokenHandler) ParseToken(signedTokenString string) (*UserClaims, error) {
	t, err := jwt.ParseWithClaims(signedTokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, errors.New("invalid signing algorithm")
		}

		_, ok := t.Header["kid"].(string)
		if !ok {
			return nil, errors.New("invalid key ID")
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
func (handler *TokenHandler) putClaimsOnContext(ctx context.Context, claims *UserClaims) (context.Context, error) {
	jsonClaims, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}
	return metadata.AppendToOutgoingContext(ctx, AuthenticatedUserMetadataKey, string(jsonClaims)), nil
}

// ValidateMiddleware middleware required endpoints: verify claims and put claims on context
func (handler *TokenHandler) ValidateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				if bearerToken[1] == "" {
					log.Println("no token supplied")
					next.ServeHTTP(w, r)
					return
				}

				token, err := handler.ParseToken(bearerToken[1])
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				ctx, _ := handler.putClaimsOnContext(r.Context(), token)
				next.ServeHTTP(w, r.WithContext(ctx))

			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

	})
}

// GetClaimsFromCtx returns claims from an authenticated user
func (handler *TokenHandler) GetClaimsFromCtx(ctx context.Context) (*UserClaims, error) {
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

func EncodeResult(w http.ResponseWriter, result interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	data := struct {
		Data interface{} `json:"data"`
	}{
		Data: result,
	}

	err := json.NewEncoder(w).Encode(&data)
	if err != nil {
		return
	}
}
