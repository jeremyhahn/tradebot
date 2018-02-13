package webservice

import (
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/service"
)

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JsonWebToken struct {
	ctx         *common.Context
	expiration  time.Duration
	token       *jwt.Token
	authService service.AuthService
	rsaKeyPair  *common.RsaKeyPair
	jsonWriter  common.HttpWriter
}

type JsonWebTokenDTO struct {
	Value string `json:"token"`
}

func NewJsonWebToken(ctx *common.Context, authService service.AuthService, jsonWriter common.HttpWriter) (*JsonWebToken, error) {
	keypair, err := common.NewRsaKeyPair(ctx)
	if err != nil {
		return nil, err
	}
	return CreateJsonWebToken(ctx, authService, jsonWriter, 10, keypair), nil
}

func CreateJsonWebToken(ctx *common.Context, authService service.AuthService, jsonWriter common.HttpWriter,
	expiration int64, rsaKeyPair *common.RsaKeyPair) *JsonWebToken {
	return &JsonWebToken{
		ctx:         ctx,
		authService: authService,
		jsonWriter:  jsonWriter,
		expiration:  time.Duration(expiration),
		token:       jwt.New(jwt.SigningMethodRS256),
		rsaKeyPair:  rsaKeyPair}
}

func (_jwt *JsonWebToken) GetToken() *jwt.Token {
	return _jwt.token
}

func (_jwt *JsonWebToken) ParseToken(r *http.Request) (*jwt.Token, error) {
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return _jwt.rsaKeyPair.PublicKey, nil
		})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (_jwt *JsonWebToken) Generate(w http.ResponseWriter, req *http.Request) {

	_jwt.ctx.Logger.Debugf("[JWT.createJsonWebToken] url: %s, method: %s, remoteAddress: %s, requestUri: %s ",
		req.URL.Path, req.Method, req.RemoteAddr, req.RequestURI)

	var user UserCredentials
	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		_jwt.ctx.Logger.Errorf("%v", "[WebServerRequest error")
		http.Error(w, "Request has error", http.StatusForbidden)
		return
	}

	err = _jwt.authService.Login(user.Username, user.Password)
	if err != nil {
		_jwt.ctx.Logger.Errorf("%v", "Invalid credentials")
		http.Error(w, "Invalid credentials", http.StatusForbidden)
		return
	}

	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * _jwt.expiration).Unix()
	claims["iat"] = time.Now().Unix()
	claims["username"] = user.Username
	_jwt.token.Claims = claims

	tokenString, err := _jwt.token.SignedString(_jwt.rsaKeyPair.PrivateKey)
	if err != nil {
		_jwt.ctx.Logger.Errorf("%v", "Error signing JWT token")
		http.Error(w, "Error signing the token", http.StatusInternalServerError)
		return
	}

	tokenDTO := JsonWebTokenDTO{tokenString}
	_jwt.jsonWriter.Write(w, http.StatusOK, tokenDTO)
}

func (_jwt *JsonWebToken) MiddlewareValidator(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, err := _jwt.ParseToken(r)
	if err == nil {
		if token.Valid {
			next(w, r)
		} else {
			_jwt.ctx.Logger.Errorf("[JsonWebToken.MiddlewareValidator] Invalid token")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
		}
	} else {

		if token == nil {
			errmsg := "Invalid request. JSON Web Token required"
			_jwt.ctx.Logger.Errorf("[JsonWebToken.MiddlewareValidator] %s", errmsg)
			http.Error(w, errmsg, http.StatusBadRequest)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if _, ok := claims["username"]; !ok {
			errmsg := "Invalid request. Username parameter required"
			_jwt.ctx.Logger.Errorf("[JsonWebToken.MiddlewareValidator] %s", errmsg)
			http.Error(w, errmsg, http.StatusBadRequest)
			return
		}

		_jwt.ctx.Logger.Errorf("Unauthorized access to %s by %s", r.RequestURI, claims["username"])
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
	}
}
