/**
加密实现
**/
package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/asn1"
	"hash"
	"math/big"

	"golang.org/x/crypto/sha3"

	"errors"

	"github.com/learnergo/loach/config"
)

//接口
type Crypto interface {
	GenerateKey() (interface{}, error)
	Sign([]byte, interface{}) ([]byte, error)
	Hash([]byte) []byte
}

var ecCurveHalfOrders = map[elliptic.Curve]*big.Int{
	elliptic.P224(): new(big.Int).Rsh(elliptic.P224().Params().N, 1),
	elliptic.P256(): new(big.Int).Rsh(elliptic.P256().Params().N, 1),
	elliptic.P384(): new(big.Int).Rsh(elliptic.P384().Params().N, 1),
	elliptic.P521(): new(big.Int).Rsh(elliptic.P521().Params().N, 1),
}

//椭圆曲线加密实现
type ecCrypto struct {
	curve    elliptic.Curve
	key      *ecdsa.PrivateKey
	hashFunc func() hash.Hash
}

type eCDSASignature struct {
	R, S *big.Int
}

func (ec *ecCrypto) GenerateKey() (interface{}, error) {
	key, err := ecdsa.GenerateKey(ec.curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (ec *ecCrypto) Hash(data []byte) []byte {
	h := ec.hashFunc()
	h.Write(data)
	return h.Sum(nil)
}

func (ec *ecCrypto) Sign(data []byte, key interface{}) ([]byte, error) {
	privateKey, yes := key.(*ecdsa.PrivateKey)
	if !yes {
		return nil, errors.New("Error Key Type")
	}

	dataHash := ec.Hash(data)

	R, S, err := ecdsa.Sign(rand.Reader, privateKey, dataHash)
	if err != nil {
		return nil, err
	}
	preventMalleability(privateKey, S)
	sig, err := asn1.Marshal(eCDSASignature{R, S})
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func preventMalleability(k *ecdsa.PrivateKey, S *big.Int) {
	halfOrder := ecCurveHalfOrders[k.Curve]
	if S.Cmp(halfOrder) == 1 {
		S.Sub(k.Params().N, S)
	}
}

func NewCrypto(config config.CryptoConfig) (Crypto, error) {
	var ecCrypto *ecCrypto = &ecCrypto{}

	switch config.Algorithm {
	case "P256-SHA256":
		ecCrypto.curve = elliptic.P256()
	case "P384-SHA384":
		ecCrypto.curve = elliptic.P384()
	case "P521-SHA512":
		ecCrypto.curve = elliptic.P521()
	default:
		return nil, errors.New("Error Algorithm")
	}

	switch config.Hash {

	case "SHA2-256":
		ecCrypto.hashFunc = sha256.New
	case "SHA2-384":
		ecCrypto.hashFunc = sha512.New384
	case "SHA3-256":
		ecCrypto.hashFunc = sha3.New256
	case "SHA3-384":
		ecCrypto.hashFunc = sha3.New384
	default:
		return nil, errors.New("Error Hash")
	}
	return ecCrypto, nil
}
