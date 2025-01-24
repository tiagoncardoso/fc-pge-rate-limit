package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strings"
)

func GetRequestIp(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}

	return ip
}

func GenerateMD5Hash(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	hashSum := hash.Sum(nil)
	return hex.EncodeToString(hashSum)
}
