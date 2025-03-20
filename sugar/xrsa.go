package sugar

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"github.com/qingfeng-lian/gosugar/logger"
	"go.uber.org/zap"
	"os"
)

// ================================================================================
// RSA加密解密
// ================================================================================

type XRsa struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewXRsa(publicKey []byte, privateKey []byte) (*XRsa, error) {
	newRasObj := &XRsa{}

	if publicKey != nil || len(publicKey) > 0 {
		block, _ := pem.Decode(publicKey)
		if block == nil {
			return nil, errors.New("public key error")
		}
		pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		pub := pubInterface.(*rsa.PublicKey)
		newRasObj.publicKey = pub
	}

	if privateKey != nil || len(privateKey) > 0 {
		block, _ := pem.Decode(privateKey)
		if block == nil {
			return nil, errors.New("private key error!")
		}
		priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		pri, ok := priv.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("private key not supported")
		}
		newRasObj.privateKey = pri
	}

	return newRasObj, nil
}

func (r *XRsa) LoadRsaPublicKey(publicKeyPath string) *XRsa {
	// 读取文件
	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		logger.Info(context.Background(), "读取文件失败", zap.Error(err))
		return nil
	}
	block, _ := pem.Decode(publicKey)
	if block == nil {
		logger.Info(context.Background(), "读取文件失败", zap.Error(err))
		return nil
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		logger.Info(context.Background(), "x509.ParsePKIXPublicKey失败", zap.Error(err))
		return nil
	}
	pub := pubInterface.(*rsa.PublicKey)
	r.publicKey = pub
	return r
}

// 公钥加密
func (r *XRsa) PublicEncrypt(data string) (string, error) {
	partLen := r.publicKey.N.BitLen()/8 - 11
	chunks := split([]byte(data), partLen)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		bytes, err := rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(bytes)
	}

	return base64.RawURLEncoding.EncodeToString(buffer.Bytes()), nil
}

// 私钥解密
func (r *XRsa) PrivateDecrypt(encrypted string) (string, error) {
	partLen := r.publicKey.N.BitLen() / 8
	raw, err := base64.RawURLEncoding.DecodeString(encrypted)
	chunks := split([]byte(raw), partLen)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(decrypted)
	}

	return buffer.String(), err
}

// 数据加签
func (r *XRsa) Sign(data string) (string, error) {
	h := crypto.SHA256.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)

	sign, err := rsa.SignPKCS1v15(rand.Reader, r.privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(sign), err
}

// 数据验签
func (r *XRsa) Verify(data string, sign string) error {
	h := crypto.SHA256.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)

	decodedSign, err := base64.RawURLEncoding.DecodeString(sign)
	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(r.publicKey, crypto.SHA256, hashed, decodedSign)
}

func MarshalPKCS8PrivateKey(key *rsa.PrivateKey) []byte {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}
	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)

	k, _ := asn1.Marshal(info)
	return k
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
