package smux

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"

	"golang.org/x/crypto/pbkdf2"
)

var SigKey = ""

const (
	KeyLength   = 32
	NonceLength = 16
	HMACLength  = 32
	PBKDF2Iter  = 100000
	MagicHeader = "ENC"
)

type Signature struct{}

// Encrypt 加密数据
func (s Signature) Encrypt(data []byte) ([]byte, error) {
	salt := make([]byte, 16)
	nonce := make([]byte, NonceLength)

	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	if SigKey == "" {
		SigKey = "hello world"
	}
	key := deriveKey(append(salt, nonce...), SigKey)
	ciphertext := streamCipher(data, key)

	macData := append(nonce, ciphertext...)
	mac := hmac.New(sha256.New, key)
	mac.Write(macData)
	hmacSum := mac.Sum(nil)

	var buf bytes.Buffer
	buf.WriteString(MagicHeader)
	buf.Write(salt)
	buf.Write(nonce)
	buf.Write(ciphertext)
	buf.Write(hmacSum)

	return buf.Bytes(), nil
}

// Decrypt 解密数据
func (s Signature) Decrypt(data []byte) ([]byte, error) {
	if len(data) < 3+16+NonceLength+HMACLength {
		return nil, errors.New("数据长度不合法")
	}
	if !bytes.HasPrefix(data, []byte(MagicHeader)) {
		return nil, errors.New("数据头部无效")
	}

	salt := data[3 : 3+16]
	nonce := data[3+16 : 3+16+NonceLength]
	macOffset := len(data) - HMACLength
	ciphertext := data[3+16+NonceLength : macOffset]
	mac := data[macOffset:]

	if SigKey == "" {
		SigKey = "hello world"
	}
	key := deriveKey(append(salt, nonce...), SigKey)

	expectedMac := hmac.New(sha256.New, key)
	expectedMac.Write(append(nonce, ciphertext...))
	expectedSum := expectedMac.Sum(nil)

	if !hmac.Equal(mac, expectedSum) {
		return nil, errors.New("HMAC 校验失败，数据可能被篡改")
	}

	plaintext := streamCipher(ciphertext, key)
	return plaintext, nil
}

// deriveKey 使用 PBKDF2 派生密钥
func deriveKey(salt []byte, sigKey string) []byte {
	return pbkdf2.Key([]byte(sigKey), salt, PBKDF2Iter, KeyLength, sha256.New)
}

// streamCipher 简单的伪随机异或加密
func streamCipher(data, key []byte) []byte {
	var output []byte
	keystream := sha256.Sum256(key)
	stream := keystream[:]
	i := 0
	for _, b := range data {
		if i >= len(stream) {
			sum := sha256.Sum256(stream)
			stream = sum[:]
			i = 0
		}
		output = append(output, b^stream[i])
		i++
	}
	return output
}
