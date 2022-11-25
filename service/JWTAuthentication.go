package service

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTService interface {
	GenerateToken(email string, isUser bool) string
	ValidateToken(token string) (*jwt.Token, error)

}

type authCustomClaims struct {
	Name string `json:"name"`
	User bool   `json:"user"`
	jwt.StandardClaims
}

type jwtServices struct {
	secretKey string
	issuer    string
}


func getSecretKey() string {
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "secret"
	}
	return secret
}

func NewJWTAuthService() JWTService {
	return &jwtServices{
		secretKey: getSecretKey(),
		issuer:    "Bikash",
	}
}

func (info *jwtServices) GenerateToken(email string, isUser bool) string {
claims := &authCustomClaims{
	Name: email,
	User: isUser,
	StandardClaims: jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
		Issuer:    info.issuer,
		IssuedAt:  time.Now().Unix(),
	},
	
}
token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims);

//encode string
t, err := token.SignedString([]byte(info.secretKey))

if err!=nil{
	panic(err)
}

return t
}

func (info *jwtServices) ValidateToken(encodedToken string) (*jwt.Token, error){
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("Invalid token", token.Header["alg"])
		}
		return []byte(info.secretKey), nil
	})
}