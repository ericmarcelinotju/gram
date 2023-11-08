package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	mathRand "math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//Private Const
const (
	constantKey = "238497335fe5bc226b0d1b86da1a51c4"
)

const (
	// 2mb
	MaxFileSize                    = 1097152
	ErrorMessageMaxFileSize        = "Max file upload 1 MB"
	TEMPORARYFILE           string = "temp_files/" //constant for temporary directory
)

func Hash(str string) (crypted string, err error) {
	password := []byte(str)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), err
}

func CompareHash(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

//Encrypt string to base64 crypto using AES
func Encrypt(password string) (data string, err error) {
	key := []byte(constantKey)
	encryptPass := []byte(password)

	block, err := aes.NewCipher(key)
	if err != nil {
		return data, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(encryptPass))

	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return data, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], encryptPass)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

//Decrypt from base64 to decrypted string
func Decrypt(cryptoText string) (data string, err error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)
	key := []byte(constantKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return data, err
	}
	if len(ciphertext) < aes.BlockSize {
		return data, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext), nil
}

var (
	letterRunes                = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	newRand     *mathRand.Rand = mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
)

//GenerateRandomString with n length
func GenerateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[newRand.Intn(len(letterRunes))]
	}
	return string(b)
}

//SplitLatLng in string to each float
func SplitLatLng(data *string) (lat, lng float64) {
	if data != nil {
		trimPrefix := strings.TrimPrefix(*data, "(")
		trimSuffixAfterPrefix := strings.TrimSuffix(trimPrefix, ")")
		dataCoordinate := strings.Split(trimSuffixAfterPrefix, ",")

		lat, _ = strconv.ParseFloat(dataCoordinate[0], 64)
		lng, _ = strconv.ParseFloat(dataCoordinate[1], 64)

		return lat, lng
	} else {
		return 0, 0
	}
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
