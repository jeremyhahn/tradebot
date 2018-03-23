package service

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/viewmodel"
)

type JsonWebTokenServiceImpl struct {
	ctx             common.Context
	databaseManager common.DatabaseManager
	expiration      time.Duration
	authService     AuthService
	rsaKeyPair      common.KeyPair
	jsonWriter      common.HttpWriter
	contexts        map[uint]common.Context
	tokens          map[uint]*jwt.Token
	JsonWebTokenService
	Middleware
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JsonWebTokenClaims struct {
	Id            uint   `json:"user_id"`
	Username      string `json:"username"`
	LocalCurrency string `json:"local_currency"`
	Etherbase     string `json:"etherbase"`
	Keystore      string `json:"keystore"`
	jwt.StandardClaims
}

var JWT_CONTEXT_LOCK sync.Mutex
var JWT_TOKEN_LOCK sync.Mutex

func NewJsonWebTokenService(ctx common.Context, databaseManager common.DatabaseManager,
	authService AuthService, jsonWriter common.HttpWriter) (JsonWebTokenService, error) {
	keypair, err := common.NewRsaKeyPair(ctx)
	if err != nil {
		return nil, err
	}
	return CreateJsonWebTokenService(ctx, databaseManager, authService, jsonWriter, 60, keypair), nil
}

func CreateJsonWebTokenService(ctx common.Context, databaseManager common.DatabaseManager, authService AuthService,
	jsonWriter common.HttpWriter, expiration int64, rsaKeyPair common.KeyPair) JsonWebTokenService {
	return &JsonWebTokenServiceImpl{
		ctx:             ctx,
		databaseManager: databaseManager,
		authService:     authService,
		jsonWriter:      jsonWriter,
		expiration:      time.Duration(expiration),
		contexts:        make(map[uint]common.Context),
		tokens:          make(map[uint]*jwt.Token),
		rsaKeyPair:      rsaKeyPair}
}

func (service *JsonWebTokenServiceImpl) CreateContext(w http.ResponseWriter, r *http.Request) (common.Context, error) {
	_, claims, err := service.ParseToken(r)
	if err != nil {
		return nil, err
	}
	service.ctx.GetLogger().Debugf("[JsonWebTokenService.CreateContext] Claims: %+v", claims)
	return &common.Ctx{
		User: &dto.UserContextDTO{
			Id:            claims.Id,
			Username:      claims.Username,
			LocalCurrency: claims.LocalCurrency,
			Etherbase:     claims.Etherbase,
			Keystore:      claims.Keystore},
		AppRoot:      service.ctx.GetAppRoot(),
		Logger:       service.ctx.GetLogger(),
		CoreDB:       service.databaseManager.ConnectCoreDB(),
		PriceDB:      service.databaseManager.ConnectPriceDB(),
		Debug:        service.ctx.GetDebug(),
		SSL:          service.ctx.GetSSL(),
		IPC:          service.ctx.GetIPC(),
		Keystore:     service.ctx.GetKeystore(),
		EthereumMode: service.ctx.GetEthereumMode()}, nil
}

func (service *JsonWebTokenServiceImpl) GetContext(userID uint) common.Context {
	return service.contexts[userID]
}

func (service *JsonWebTokenServiceImpl) ParseToken(r *http.Request) (*jwt.Token, *JsonWebTokenClaims, error) {
	service.ctx.GetLogger().Debugf("[JsonWebTokenService.ParseToken]")
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return service.rsaKeyPair.GetPublicKey(), nil
		})
	if err != nil {
		return nil, nil, err
	}
	claims := &JsonWebTokenClaims{}
	_, err = jwt.ParseWithClaims(token.Raw, claims,
		func(token *jwt.Token) (interface{}, error) {
			return service.rsaKeyPair.GetPublicKey(), nil
		})
	if err != nil {
		return nil, nil, err
	}
	return token, claims, nil
}

func (service *JsonWebTokenServiceImpl) GenerateToken(w http.ResponseWriter, req *http.Request) {

	service.ctx.GetLogger().Debugf("[JsonWebTokenService.GenerateToken] url: %s, method: %s, remoteAddress: %s, requestUri: %s ",
		req.URL.Path, req.Method, req.RemoteAddr, req.RequestURI)

	var user UserCredentials
	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		service.jsonWriter.Write(w, http.StatusBadRequest, viewmodel.JsonWebToken{Error: "Bad request"})
		return
	}

	userContext, err := service.authService.Login(user.Username, user.Password)
	if err != nil {
		service.jsonWriter.Write(w, http.StatusForbidden, viewmodel.JsonWebToken{Error: "Invalid credentials"})
		return
	}

	userID := userContext.GetId()

	JWT_TOKEN_LOCK.Lock()
	service.tokens[userID] = jwt.NewWithClaims(jwt.SigningMethodRS256, JsonWebTokenClaims{
		Id:            userID,
		Username:      userContext.GetUsername(),
		LocalCurrency: userContext.GetLocalCurrency(),
		Etherbase:     userContext.GetEtherbase(),
		Keystore:      userContext.GetKeystore(),
		StandardClaims: jwt.StandardClaims{
			Issuer:    common.APPNAME,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * service.expiration).Unix()}})
	JWT_TOKEN_LOCK.Unlock()

	tokenString, err := service.tokens[userID].SignedString(service.rsaKeyPair.GetPrivateKey())
	if err != nil {
		service.jsonWriter.Write(w, http.StatusInternalServerError, viewmodel.JsonWebToken{Error: "Error signing token"})
		return
	}

	tokenDTO := viewmodel.JsonWebToken{Value: tokenString}
	service.jsonWriter.Write(w, http.StatusOK, tokenDTO)
}

func (service *JsonWebTokenServiceImpl) Validate(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, claims, err := service.ParseToken(r)
	if err == nil {
		if token.Valid {
			if claims.Id <= 0 {
				errmsg := "Invalid request. user_id claim required"
				service.ctx.GetLogger().Errorf("[JsonWebTokenService.Validate] %s", errmsg)
				http.Error(w, errmsg, http.StatusBadRequest)
				return
			}
			if claims.Username == "" {
				errmsg := "Invalid request. ssername claim required"
				service.ctx.GetLogger().Errorf("[JsonWebTokenService.Validate] %s", errmsg)
				http.Error(w, errmsg, http.StatusBadRequest)
				return
			}
			if claims.LocalCurrency == "" {
				errmsg := "Invalid request. local_currency claim required"
				service.ctx.GetLogger().Errorf("[JsonWebTokenService.Validate] %s", errmsg)
				http.Error(w, errmsg, http.StatusBadRequest)
				return
			}
			if claims.Etherbase == "" {
				errmsg := "Invalid request. etherbase claim required"
				service.ctx.GetLogger().Errorf("[JsonWebTokenService.Validate] %s", errmsg)
				http.Error(w, errmsg, http.StatusBadRequest)
				return
			}
			if claims.Keystore == "" {
				errmsg := "Invalid request. keystore claim required"
				service.ctx.GetLogger().Errorf("[JsonWebTokenService.Validate] %s", errmsg)
				http.Error(w, errmsg, http.StatusBadRequest)
				return
			}
			ctx := &common.Ctx{
				User: &dto.UserContextDTO{
					Id:            claims.Id,
					Username:      claims.Username,
					LocalCurrency: claims.LocalCurrency,
					Etherbase:     claims.Etherbase,
					Keystore:      claims.Keystore},
				AppRoot:      service.ctx.GetAppRoot(),
				Logger:       service.ctx.GetLogger(),
				CoreDB:       service.databaseManager.ConnectCoreDB(),
				PriceDB:      service.databaseManager.ConnectPriceDB(),
				Debug:        service.ctx.GetDebug(),
				SSL:          service.ctx.GetSSL(),
				IPC:          service.ctx.GetIPC(),
				Keystore:     service.ctx.GetKeystore(),
				EthereumMode: service.ctx.GetEthereumMode()}

			JWT_CONTEXT_LOCK.Lock()
			service.contexts[claims.Id] = ctx
			JWT_CONTEXT_LOCK.Unlock()

			service.ctx.GetLogger().Errorf("[JsonWebTokenService.Validate] Setting user context: %+v", ctx.GetUser())
			next(w, r)
		} else {
			JWT_CONTEXT_LOCK.Lock()
			service.contexts[claims.Id] = nil
			JWT_CONTEXT_LOCK.Unlock()

			service.ctx.GetLogger().Errorf("[JsonWebTokenService.Validate] Invalid token: %s", token.Raw)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
		}
	} else {
		errmsg := "Invalid request. JSON Web Token required"
		service.ctx.GetLogger().Errorf("[JsonWebTokenService.Validate] Error: %s", err.Error())
		http.Error(w, errmsg, http.StatusBadRequest)
	}
}
