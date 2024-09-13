package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"tg-backend/pkg/log"
	"tg-backend/server/types"
	"tg-backend/server/util"
	"time"
)

func NewTelegramAuthMiddleware(
	telegramBotToken string,
) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tgUser, err := checkUser(r, telegramBotToken)
			if err != nil {
				errMsg := "Authentication error!"
				http.Error(w, errMsg, http.StatusForbidden)
				log.Error(err.Error())
				return
			}
			r = r.WithContext(util.NewContext(r.Context(), tgUser))
			next.ServeHTTP(w, r)
		})
	}
}

func checkUser(r *http.Request, token string) (*types.TelegramUser, error) {
	var query url.Values

	// auth header
	authHeader := r.Header.Get("Authorization")
	var err error
	//log.Info("Authorization: %s", authHeader)
	// base64 decode
	//authHeaderBytes, err := base64.StdEncoding.DecodeString(authHeader)
	//authHeader = string(authHeaderBytes)
	//if len(authHeader) == 0 || err != nil {
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	query, err = url.ParseQuery(authHeader)
	if err != nil {
		return nil, err
	}

	expire := checkExpire(query)
	if expire {
		return nil, errors.New("authorization expired")
	}

	res, err := checkAuth(token, query)
	if err != nil || !res {
		return nil, err
	}

	tgUser := &types.TelegramUser{}
	err = json.Unmarshal([]byte(query.Get("user")), tgUser)
	if err != nil {
		return nil, err
	}
	return tgUser, nil
}

func checkExpire(query url.Values) bool {
	authDate := query.Get("auth_date")
	if len(authDate) == 0 {
		return false
	}
	authUnix, err := strconv.Atoi(authDate)
	if err != nil {
		return false
	}

	return time.Now().Unix() > int64(authUnix+1800)
}

func checkAuth(secret string, query url.Values) (bool, error) {
	hash := query.Get("hash")
	if len(hash) == 0 {
		return false, errors.New("hash empty")
	}

	authCheckString, err := getAuthCheckString(query)
	if err != nil {
		return false, err
	}

	secretKey := getHmac256Signature([]byte("WebAppData"), []byte(secret))
	expectedHash := getHmac256Signature(secretKey, []byte(authCheckString))
	expectedHashString := hex.EncodeToString(expectedHash)

	if expectedHashString != hash {
		return false, errors.New("hash incorrect")
	}
	return true, nil
}

// get alphabetic sorted query string
func getAuthCheckString(values url.Values) (string, error) {
	paramKeys := make([]string, 0)
	for key, v := range values {
		if key == "hash" {
			continue
		}
		if len(v) != 1 {
			return "", errors.New("is not a valid auth query")
		}
		paramKeys = append(paramKeys, key)
	}

	// sort keys
	sort.Strings(paramKeys)

	dataCheckArr := make([]string, len(paramKeys))
	for i, key := range paramKeys {
		dataCheckArr[i] = key + "=" + values.Get(key)
	}

	return strings.Join(dataCheckArr, "\n"), nil
}

func getHmac256Signature(secretKey []byte, data []byte) []byte {
	mac := hmac.New(sha256.New, secretKey)
	mac.Write(data)
	sum := mac.Sum(nil)
	return sum
}

func EnableCORS(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Allow requests from any origin

		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Allow specified HTTP methods

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")

		// Allow specified headers

		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

		// Continue with the next handler

		next.ServeHTTP(w, r)

	})

}
