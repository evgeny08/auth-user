package httpserver

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/evgeny08/auth-user/types"
)

// service manages HTTP server methods.
type service interface {
	createUser(ctx context.Context, user *types.User) error
}

type basicService struct {
	logger  log.Logger
	storage Storage
}

// createUser creates a new User
func (s *basicService) createUser(ctx context.Context, user *types.User) error {
	// Validate
	passwordLength := 8
	if strings.TrimSpace(user.Login) == "" {
		return errorf(ErrBadParams, "empty login")
	}
	if len(user.Password) < passwordLength {
		return errorf(ErrBadParams, "password must be 8 symbols min")
	}
	_, err := s.storage.FindUserByLogin(ctx, user.Login)
	if err == nil {
		return errorf(ErrConflict, "user with this login already exist")
	}

	err = s.storage.CreateUser(ctx, user)
	if err != nil {
		return errorf(ErrInternal, "failed to insert user: %v", err)
	}
	return nil
}

// storageErrIsNotFound checks if the storage error is "not found".
func storageErrIsNotFound(err error) bool {
	type notFound interface {
		NotFound() bool
	}
	e, ok := err.(notFound)
	return ok && e.NotFound()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

//encrypt userPassword
//var (
//	initialVector = "1374653317652235"
//	passphrase    = "Impassphrasegood"
//)
//
//func (s*basicService)encryptPass(user *types.User) () {
//	var userPass = user.Password
//
//	encryptedData := EncryptAES(userPass, []byte(passphrase))
//	encryptedString := base64.StdEncoding.EncodeToString(encryptedData)
//
//	encryptedData, _ = base64.StdEncoding.DecodeString(encryptedString)
//	decryptedText := DecryptAES(encryptedData, []byte(passphrase))
//	fmt.Println(string(decryptedText))
//}
//
//func EncryptAES(src string, key []byte) []byte {
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		fmt.Println("key error1", err)
//	}
//	if src == "" {
//		fmt.Println("plain content empty")
//	}
//	ecb := cipher.NewCBCEncrypter(block, []byte(initialVector))
//	content := []byte(src)
//	content = PaddingPKCS5(content, block.BlockSize())
//	crypted := make([]byte, len(content))
//	ecb.CryptBlocks(crypted, content)
//
//	return crypted
//}
//
//func DecryptAES(crypt []byte, key []byte) []byte {
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		fmt.Println("key error1", err)
//	}
//	if len(crypt) == 0 {
//		fmt.Println("plain content empty")
//	}
//	ecb := cipher.NewCBCDecrypter(block, []byte(initialVector))
//	decrypted := make([]byte, len(crypt))
//	ecb.CryptBlocks(decrypted, crypt)
//
//	return TrimmingPKCS5(decrypted)
//}
//
//func PaddingPKCS5(cipherText []byte, blockSize int) []byte {
//	padding := blockSize - len(cipherText)%blockSize
//	padText := bytes.Repeat([]byte{byte(padding)}, padding)
//	return append(cipherText, padText...)
//}
//
//func TrimmingPKCS5(encrypt []byte) []byte {
//	padding := encrypt[len(encrypt)-1]
//	return encrypt[:len(encrypt)-int(padding)]
//}
