package encryption

//go:generate mockgen -destination=../mocks/mockEncryptorService.go -package=mocks gitlab.connectwisedev.com/platform/platform-tasking-service/src/encryption EncryptorService

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"

	agentModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
)

// EncPrefix is a prefix, added to encrypted using Rijndael-AES algorithm data
const EncPrefix = "aes256:"

// Encryptor is encryption service var
var Encryptor EncryptorService

// EncryptorService is encryption interface
type EncryptorService interface {
	Encrypt(creds agentModels.Credentials) (agentModels.Credentials, error)
	Decrypt(creds agentModels.Credentials) (agentModels.Credentials, error)
}

// NewService creates AES encryption implementation
func NewService(encKey string) EncryptorService {
	return AESEncryptor{
		Key: encKey,
	}
}

// AESEncryptor holds encryption key
type AESEncryptor struct {
	Key string
}

// Encrypt encrypts given agentModels.Credentials with AES256 using service key from config
func (e AESEncryptor) Encrypt(creds agentModels.Credentials) (encrypted agentModels.Credentials, err error) {
	encrypted.UseCurrentUser = creds.UseCurrentUser

	encrypted.Password, err = e.encrypt(creds.Password, e.Key)
	if err != nil {
		return encrypted, fmt.Errorf("Encrypt: encryption of password failed. err %s", err.Error())
	}

	encrypted.Username, err = e.encrypt(creds.Username, e.Key)
	if err != nil {
		return encrypted, fmt.Errorf("Encrypt: encryption of username failed. err %s", err.Error())
	}

	encrypted.Domain, err = e.encrypt(creds.Domain, e.Key)
	if err != nil {
		return encrypted, fmt.Errorf("Encrypt: encryption of domain failed. err %s", err.Error())
	}

	return
}

// Decrypt decrypts given agentModels.Credentials with AES256 using service key from config
func (e AESEncryptor) Decrypt(creds agentModels.Credentials) (decrypted agentModels.Credentials, err error) {
	decrypted.UseCurrentUser = creds.UseCurrentUser

	decrypted.Password, err = e.decrypt(creds.Password, e.Key)
	if err != nil {
		return decrypted, fmt.Errorf("Decrypt: decryption of password failed. err %s", err.Error())
	}

	decrypted.Username, err = e.decrypt(creds.Username, e.Key)
	if err != nil {
		return decrypted, fmt.Errorf("Decrypt: decryption of username failed. err %s", err.Error())
	}

	decrypted.Domain, err = e.decrypt(creds.Domain, e.Key)
	if err != nil {
		return decrypted, fmt.Errorf("Decrypt: decryption of domain failed. err %s", err.Error())
	}

	return
}

func (AESEncryptor) encrypt(data string, key string) (string, error) {
	if len(data) == 0 {
		return "", nil
	}
	keyHash := sha256.Sum256([]byte(key))
	block, _ := aes.NewCipher(keyHash[:])
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "cannot create GCM")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.Wrap(err, "cannot create nonce")
	}

	cipherText := base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(data), nil))

	return EncPrefix + cipherText, err
}

func (AESEncryptor) decrypt(data string, key string) (string, error) {
	if len(data) == 0 {
		return "", nil
	}
	keyHash := sha256.Sum256([]byte(key))

	if !strings.HasPrefix(data, EncPrefix) {
		return data, nil
	}

	data = strings.Replace(data, EncPrefix, "", 1)
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return data, errors.Wrapf(err, "cannot decode string: %s", data)
	}

	block, err := aes.NewCipher(keyHash[:])
	if err != nil {
		return data, errors.Wrap(err, "cannot create block")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return data, errors.Wrap(err, "cannot create GCM")
	}

	nonceSize := gcm.NonceSize()

	nonce, cipherText := decoded[:nonceSize], decoded[nonceSize:]
	result, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return data, errors.Wrap(err, "cannot open")
	}

	return string(result), nil
}
