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
	KeyPair
}

func NewRsaKeyPair(ctx Context) (KeyPair, error) {
	return CreateRsaKeyPair(ctx, "./keys")
}

func CreateRsaKeyPair(ctx Context, directory string) (KeyPair, error) {
	privateKeyBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", directory, "rsa.key"))
	if err != nil {
		ctx.GetLogger().Errorf("[RsaKeyPair] %s", err.Error())
		return nil, err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		ctx.GetLogger().Errorf("[RsaKeyPair] %s", err.Error())
		return nil, err
	}
	publicKeyBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", directory, "rsa.pub"))
	if err != nil {
		ctx.GetLogger().Errorf("[RsaKeyPair] %s", err.Error())
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		ctx.GetLogger().Errorf("[RsaKeyPair] %s", err.Error())
		return nil, err
	}
	return &RsaKeyPair{
		Directory:    directory,
		PrivateKey:   privateKey,
		PrivateBytes: privateKeyBytes,
		PublicKey:    publicKey,
		PublicBytes:  publicKeyBytes}, nil
}

func (keypair *RsaKeyPair) GetDirectory() string {
	return keypair.Directory
}

func (keypair *RsaKeyPair) GetPrivateKey() *rsa.PrivateKey {
	return keypair.PrivateKey
}

func (keypair *RsaKeyPair) GetPrivateBytes() []byte {
	return keypair.PrivateBytes
}

func (keypair *RsaKeyPair) GetPublicKey() *rsa.PublicKey {
	return keypair.PublicKey
}

func (keypair *RsaKeyPair) GetPublicBytes() []byte {
	return keypair.PublicBytes
}
