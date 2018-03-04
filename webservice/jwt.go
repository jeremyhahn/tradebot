package webservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/service"
)

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JsonWebToken struct {
	ctx         common.Context
	expiration  time.Duration
	token       *jwt.Token
	authService service.AuthService
	rsaKeyPair  common.KeyPair
	jsonWriter  common.HttpWriter
	common.JsonWebToken
}

type JsonWebTokenDTO struct {
	Value string `json:"token"`
	Error string `json:"error"`
}

func NewJsonWebToken(ctx *common.Context, authService service.AuthService,
	jsonWriter common.HttpWriter) (*JsonWebToken, error) {
	keypair, err := common.NewRsaKeyPair(ctx)
	if err != nil {
		return nil, err
	}
	return CreateJsonWebToken(ctx, authService, jsonWriter, 60, keypair), nil
}

func CreateJsonWebToken(ctx common.Context, authService service.AuthService,
	jsonWriter common.HttpWriter, expiration int64, rsaKeyPair *common.RsaKeyPair) *JsonWebToken {
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

func (_jwt *JsonWebToken) GetClaims() jwt.MapClaims {
	return _jwt.GetToken().Claims.(jwt.MapClaims)
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

func (_jwt *JsonWebToken) GenerateToken(w http.ResponseWriter, req *http.Request) {

	_jwt.ctx.Logger.Debugf("[JWT.createJsonWebToken] url: %s, method: %s, remoteAddress: %s, requestUri: %s ",
		req.URL.Path, req.Method, req.RemoteAddr, req.RequestURI)

	var user UserCredentials
	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		_jwt.jsonWriter.Write(w, http.StatusBadRequest, JsonWebTokenDTO{Error: "Bad request"})
		return
	}

	userDTO, err := _jwt.authService.Login(user.Username, user.Password)
	if err != nil {
		_jwt.jsonWriter.Write(w, http.StatusForbidden, JsonWebTokenDTO{Error: "Invalid credentials"})
		return
	}

	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * _jwt.expiration).Unix()
	claims["iat"] = time.Now().Unix()
	claims["user_id"] = userDTO.GetId()
	claims["username"] = userDTO.GetUsername()
	claims["local_currency"] = userDTO.GetLocalCurrency()
	claims["etherbase"] = userDTO.GetEtherbase()
	claims["keystore"] = userDTO.GetKeystore()
	_jwt.token.Claims = claims

	tokenString, err := _jwt.token.SignedString(_jwt.rsaKeyPair.PrivateKey)
	if err != nil {
		_jwt.jsonWriter.Write(w, http.StatusInternalServerError, JsonWebTokenDTO{Error: "Error signing token"})
		return
	}

	if err = _jwt.setApplicationStateFromToken(); err != nil {
		_jwt.ctx.Logger.Errorf("[JsonWebToken.GenerateToken] Failed to set application state using generated claims: %s", err.Error())
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
	}

	tokenDTO := JsonWebTokenDTO{Value: tokenString}
	_jwt.jsonWriter.Write(w, http.StatusOK, tokenDTO)
}

func (_jwt *JsonWebToken) MiddlewareValidator(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, err := _jwt.ParseToken(r)
	if err == nil {
		if token.Valid {
			if err := _jwt.setApplicationStateFromToken(); err != nil {
				_jwt.ctx.Logger.Errorf("[JsonWebToken.MiddlewareValidator] Failed to set application state using JWT claims: %s", err.Error())
				http.Error(w, "Invalid claims", http.StatusInternalServerError)
			}
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

func (_jwt *JsonWebToken) setApplicationStateFromToken() error {
	claims := _jwt.token.Claims.(jwt.MapClaims)
	_jwt.ctx.Logger.Debugf("[JsonWebToken.setApplicationState] Setting claims: %+v", claims)
	claimIds := []string{"user_id", "username", "local_currency", "etherbase", "keystore"}
	for _, id := range claimIds {
		if _, ok := claims[id]; !ok {
			return errors.New(fmt.Sprintf("missing %s field", id))
		}
	}
	userDTO := &dto.UserDTO{
		Id:            claims["user_id"].(uint),
		Username:      claims["username"].(string),
		LocalCurrency: claims["local_currency"].(string),
		Etherbase:     claims["etherbase"].(string),
		Keystore:      claims["keystore"].(string)}

	newContext := &common.AppContext{
		Logger: *_jwt.ctx.GetLogger(),
		CoreDB: *_jwt.ctx.GetCoreDB(),
		PriceDB: *_jwt.ctx.GetPriceDB(),
		Debug: *_jwt.ctx.GetDebug().
		SSL: *_jwt.ctx.GetSSL()}

	_jwt.ctx = newContext
	return nil
}
