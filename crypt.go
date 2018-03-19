package main

import (
	"crypto/aes"
	"io"
	"crypto/cipher"
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha256"
	"crypto/rand"
)

//
// See: https://gist.github.com/manishtpatel/8222606
//
// Encrypted format:
// encrypted[:SaltBytes] == salt
// encrypted[SaltBytes:] == cryptodata
// cryptodata[:aes.BlockSize] == iv
// cryptodata[aes.BlockSize:] == encryptedData
//
func encrypt(password, data []byte) ([]byte, error) {
	salt, err := makeSalt(SaltBytes)

	if err != nil {
		return nil, err
	}

	key := makeKey(password, salt)

	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err // @todo FailedToEncrypt
	}

	cipherText := make([]byte, aes.BlockSize+len(data))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err // @todo FailedToEncrypt
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], data)

	return append(salt, cipherText...), nil
}

//
// See: https://gist.github.com/manishtpatel/8222606
//
// Encrypted format:
// encrypted[:SaltBytes] == salt
// encrypted[SaltBytes:] == cryptodata
// cryptodata[:aes.BlockSize] == iv
// cryptodata[aes.BlockSize:] == encryptedData
//
func decrypt(password []byte, encrypted []byte) ([]byte, error) {
	salt := encrypted[:SaltBytes]
	encrypted = encrypted[SaltBytes:]
	key := makeKey(password, salt)

	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err // @todo FailedToDecrypt
	}

	if len(encrypted) < aes.BlockSize {
		// Cipher is too short
		return nil, err // @todo FailedToDecrypt
	}

	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(encrypted, encrypted)

	return encrypted, nil
}

//
// Generate a random salt for use in a key derivation function
//
func makeSalt(length int) ([]byte, error) {
	salt := make([]byte, length)

	_, err := io.ReadFull(rand.Reader, salt)

	if err != nil {
		return nil, err
	}

	return salt, nil
}

//
// Stretch a short, human readable key into a more secure key using a key derivation function
//
func makeKey(password, salt []byte) []byte {
	return pbkdf2.Key(password, salt, Pbkdf2Iters, KeySize, sha256.New)
}