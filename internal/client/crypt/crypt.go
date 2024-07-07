package crypt

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"github.com/zenazn/pkcs7pad"
)

// getMD5Hash md5 as cipher key length 32
func getMD5Hash(text string) []byte {
	hash := md5.Sum([]byte(text))
	return hash[:]
}

// AES256CBCEncode encrypt plain text into cipher text
func AES256CBCEncode(plainText []byte, key string) (cipherText []byte, err error) {
	md5key := getMD5Hash(key)

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

	block, err = aes.NewCipher(md5key)
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

// AES256CBCDecode decrypt cipher text into plain text
func AES256CBCDecode(cipherText []byte, key string) (plainText []byte, err error) {
	bKey := getMD5Hash(key)

	var block cipher.Block
	block, err = aes.NewCipher(bKey)
	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("cipherText too short")
		return
	}
	iv := cipherText[:aes.BlockSize]
	if len(cipherText[aes.BlockSize:])%aes.BlockSize != 0 {
		err = errors.New("plainText has wrong block size")
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
