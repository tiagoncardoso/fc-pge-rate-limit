package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/tiagoncardoso/fc-pge-rate-limit/pkg/fcrl/rllog"
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

func ParseStructToString(data interface{}) string {
	emptyData, _ := json.Marshal(struct{}{})

	jsonData, err := json.Marshal(data)
	if err != nil {
		rllog.Error("Failed to marshal CacheData: " + err.Error())
		return string(emptyData)
	}

	return string(jsonData)
}
