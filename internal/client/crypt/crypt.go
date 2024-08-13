package crypt

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	"github.com/zenazn/pkcs7pad"
)

// Encode encrypt plain text into cipher text
func Encode(plainText []byte, key string) (cipherText []byte, err error) {
	bKey := sha256.Sum256([]byte(key))

	compressed := new(bytes.Buffer)
	w := gzip.NewWriter(compressed)
	_, err = w.Write(plainText)
	if err != nil {
		err = fmt.Errorf("error writing to gzip: %v", err)
		return
	}
	err = w.Close()
	if err != nil {
		err = fmt.Errorf("error close writing to gzip: %v", err)
		return
	}

	paddedText := pkcs7pad.Pad(compressed.Bytes(), aes.BlockSize)

	if len(paddedText)%aes.BlockSize != 0 {
		err = fmt.Errorf("after pkcs7pad it has the wrong block size")
		return
	}
	var block cipher.Block

	block, err = aes.NewCipher(bKey[:])
	if err != nil {
		return
	}

	cipherText = make([]byte, aes.BlockSize+len(paddedText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	encrypt := cipher.NewCBCEncrypter(block, iv)
	encrypt.CryptBlocks(cipherText[aes.BlockSize:], paddedText)

	return
}

// Decode decrypt cipher text into plain text
func Decode(cipherText []byte, key string) (plainText []byte, err error) {
	bKey := sha256.Sum256([]byte(key))

	var block cipher.Block
	block, err = aes.NewCipher(bKey[:])
	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("cipherText too short")
		return
	}
	iv := cipherText[:aes.BlockSize]
	if len(cipherText[aes.BlockSize:])%aes.BlockSize != 0 {
		err = errors.New("cipherText has wrong block size")
		return
	}
	paddedText := make([]byte, len(cipherText)-aes.BlockSize)
	decrypt := cipher.NewCBCDecrypter(block, iv)
	decrypt.CryptBlocks(paddedText, cipherText[aes.BlockSize:])

	var compressed []byte
	compressed, err = pkcs7pad.Unpad(paddedText)

	if err != nil {
		err = fmt.Errorf("unpad error %w", err)
		return
	}

	var r *gzip.Reader
	b := bytes.NewBuffer(compressed)
	r, err = gzip.NewReader(b)
	if err != nil {
		err = fmt.Errorf("ungzip newReader error %w", err)
		return
	}
	plainText, err = io.ReadAll(r)
	if err != nil {
		err = fmt.Errorf("ungzip readAll error %w", err)
		return
	}
	err = r.Close()
	if err != nil {
		err = fmt.Errorf("ungzip close error %w", err)
		return
	}

	return
}
