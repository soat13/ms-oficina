package jwt

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    app "github.com/soat13/fase-1-oficina/internal/auth/application"
)

type Service struct {
    secret []byte
    issuer string
    ttl    time.Duration
}

func NewService(secret, issuer string, ttl time.Duration) *Service {
    return &Service{secret: []byte(secret), issuer: issuer, ttl: ttl}
}

type customClaims struct {
    Email string `json:"email"`
    jwt.RegisteredClaims
}

func (s *Service) GenerateToken(userID uuid.UUID, email string) (string, error) {
    now := time.Now()
    claims := customClaims{
        Email: email,
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    s.issuer,
            Subject:   userID.String(),
            IssuedAt:  jwt.NewNumericDate(now),
            ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.secret)
}

func (s *Service) ParseToken(tokenStr string) (app.Claims, error) {
    var claims customClaims
    token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return s.secret, nil
    })
    if err != nil {
        return app.Claims{}, err
    }
    if !token.Valid {
        return app.Claims{}, errors.New("invalid token")
    }
    uid, err := uuid.Parse(claims.Subject)
    if err != nil {
        return app.Claims{}, err
    }
    return app.Claims{UserID: uid, Email: claims.Email}, nil
}

var _ interface{ GenerateToken(uuid.UUID, string) (string, error); ParseToken(string) (app.Claims, error) } = (*Service)(nil)

