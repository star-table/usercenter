package random

import (
	"encoding/base64"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/martinlindhe/base36"
	"github.com/star-table/usercenter/pkg/util/times"
	"github.com/star-table/usercenter/pkg/util/uuid"
)

func Token() string {
	//return uuid.NewUuid()
	token, _ := GenerateRandomStringAsBase36(32)
	return token
}

func RandomString(l int) string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func RandomFileName() string {
	return uuid.NewUuid() + strconv.FormatInt(times.GetNowMillisecond(), 10)
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomStringAsBase64(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), err
}

func GenerateRandomStringAsBase36(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	if err != nil {
		return "", err
	}
	return strings.ToLower(base36.EncodeBytes(b)), err
}
