package common

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"

	jwt "github.com/dgrijalva/jwt-go"
)

type RsaKeyPair struct {
	Directory    string
	PrivateKey   *rsa.PrivateKey
	PrivateBytes []byte
	PublicKey    *rsa.PublicKey
	PublicBytes  []byte
}

func NewRsaKeyPair(ctx *Context) (*RsaKeyPair, error) {
	return CreateRsaKeyPair(ctx, "./keys")
}

func CreateRsaKeyPair(ctx *Context, directory string) (*RsaKeyPair, error) {
	privateKeyBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", directory, "rsa.key"))
	if err != nil {
		ctx.Logger.Errorf("[RsaKeyPair] %s", err.Error())
		return nil, err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		ctx.Logger.Errorf("[RsaKeyPair] %s", err.Error())
		return nil, err
	}
	publicKeyBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", directory, "rsa.pub"))
	if err != nil {
		ctx.Logger.Errorf("[RsaKeyPair] %s", err.Error())
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		ctx.Logger.Errorf("[RsaKeyPair] %s", err.Error())
		return nil, err
	}
	return &RsaKeyPair{
		Directory:    directory,
		PrivateKey:   privateKey,
		PrivateBytes: privateKeyBytes,
		PublicKey:    publicKey,
		PublicBytes:  publicKeyBytes}, nil
}
