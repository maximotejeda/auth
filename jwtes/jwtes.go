// package in chage of manage keys and create, refresh and validate JWT
// the package is called from main to make accessible globally
package jwtes

import (
	"errors"
	"fmt"
	"maps"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// jwtType
type jwtType struct {
	*keyManager
}

// NewJwtType
func NewJwtType() (*jwtType, error) {
	keys, err := NewKeys("jwtsign")
	if err != nil {
		return nil, fmt.Errorf("creating jwtType: %w", err)
	}
	keys.writeToDisk()
	return &jwtType{keyManager: keys}, nil
}

// CreateSignToken
// need claims and time to live for the token
//
//	claims := map[string]interface{}{
//			"id": 1,
//			"username": "juan",
//			"email": "juan@example.com",
//			"rol": "user,admin",
//			"sub": "pasword reset",
//		}
func (j *jwtType) CreateSignToken(claim map[string]any, ttl time.Duration) (string, error) {
	rid := uuid.NewString()

	jwtClaims := jwt.MapClaims{
		"iss": "auth JWT",      //  who is the issuer(creator)
		"sub": "authorization", // to whom the token was created
		"exp": jwt.NewNumericDate(time.Now().Add(ttl)),
		"nbf": jwt.NewNumericDate(time.Now()),
		"iat": jwt.NewNumericDate(time.Now()),
		"jti": rid,
	}
	maps.Copy(jwtClaims, claim)
	token := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, jwtClaims)
	ss, err := token.SignedString(j.privKey)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}
	return ss, nil
}

// ValidateJWTString
// here only jwt.MapClaims default keys are verified and algo from sign
func (j *jwtType) ValidateJWTString(tokenSTR string) (map[string]any, error) {
	token, err := jwt.Parse(tokenSTR, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected singning method: %v", token.Header["alg"])
		}
		return j.pubKey, nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	} else if errors.Is(err, jwt.ErrTokenExpired) {
		return claims, err
	}
	return nil, fmt.Errorf("invalid token: %w", err)
}

// RefreshToken
// if Validate fails for token expiration we can issue a new one
// will fail in any other context, bad alg, or still valid token
func (j *jwtType) RefreshExpiredToken(tokenSTR string) (string, error) {
	claims, err := j.ValidateJWTString(tokenSTR)
	if !errors.Is(err, jwt.ErrTokenExpired) {
		return "", fmt.Errorf("error is not expired: %w", err)
	}
	now := time.Now()

	issuedflt, ok := claims["iat"].(float64)
	if !ok {
		return "", fmt.Errorf("unknown issued type")
	}
	issued := time.Unix(int64(issuedflt), 0)
	expirationflt, ok := claims["exp"].(float64)
	expiration := time.Unix(int64(expirationflt), 0)
	if !ok {
		return "", fmt.Errorf("unknown expiration type")
	}
	ttl := expiration.Sub(issued) //time to live of the token

	if now.Sub(expiration) > ttl {
		return "", fmt.Errorf("time expired more than permited time\n\texpired: %#v\n\tissued: %#v\n\tpermited time:%v\n\telapsed Time:%v", expiration, issued, ttl, time.Since(expiration)) // if expired with more time than the previous ttl must issue a normal token
	}
	delete(claims, "exp")
	delete(claims, "nbf")
	delete(claims, "jti")
	delete(claims, "iat")
	claims["sub"] = "refresh token"
	jwtSTR, err := j.CreateSignToken(claims /*ttl*/, 3*time.Second)
	if err != nil {
		return "", fmt.Errorf("creating token: %w", err)
	}
	return jwtSTR, err
}
