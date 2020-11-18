package module

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
    "errors"
    "fmt"
)

var keyCrypt = []byte("1234567890ABCDEF1234567890ABCDEF")
var ivCrypt = []byte("1234567890ABCDEF")

func PKCS5Padding(data []byte, blockSize int) []byte {
    padding := blockSize - len(data)%blockSize
    padtext := bytes.Repeat([]byte{byte(0)}, padding)
    return append(data, padtext...)
}

func PKCS5Trimming(data []byte, blockSize int) ([]byte, error) {
    dataLen := len(data)
    paddingLen := int(data[dataLen-1])

    if paddingLen > dataLen || paddingLen >= blockSize {
        return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
    }

    return data[:(dataLen - paddingLen)], nil
}

func AESEncrypt(text string) (string, error) {
    block, err := aes.NewCipher(keyCrypt)
    if err != nil {
        return "", err
    }

    msg := PKCS5Padding([]byte(text), aes.BlockSize)
    ciphertext := make([]byte, len(msg))

    cbc := cipher.NewCBCEncrypter(block, ivCrypt)
    cbc.CryptBlocks(ciphertext, msg)
    
    encrypted := base64.StdEncoding.EncodeToString(ciphertext)
    return encrypted, nil
}

func AESDecrypt(text string) (string, error) {
    block, err := aes.NewCipher(keyCrypt)
    if err != nil {
        return "", err
    }

    decodedMsg, err := base64.StdEncoding.DecodeString(text)
    
    if err != nil {
        return "", err
    }

    if (len(decodedMsg) % aes.BlockSize) != 0 {
        return "", errors.New("blocksize must be multipe of decoded message length")
    }

    msg := make([]byte, len(decodedMsg))

    cbc := cipher.NewCBCDecrypter(block, ivCrypt)
    cbc.CryptBlocks(msg, decodedMsg)

    unpadMsg, err := PKCS5Trimming(msg, aes.BlockSize)
    if err != nil {
        return "", err
    }

    return string(unpadMsg), nil
}

func runDemo() {
    fmt.Printf("key = %x\n", keyCrypt)
    fmt.Printf("iv = %x\n", ivCrypt)

    encrypted, _ := AESEncrypt("http://cdn.test.vn/channel/test/playlist.m3u8")
    fmt.Println("encrypted = " + encrypted)
    msg, _ := AESDecrypt(encrypted)
    fmt.Println("decrypted = " + msg)
}