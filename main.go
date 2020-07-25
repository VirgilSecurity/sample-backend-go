package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/VirgilSecurity/virgil-sdk-go/v6/crypto"
	"github.com/VirgilSecurity/virgil-sdk-go/v6/session"
)

var (
	userStorage sync.Map

	cryptoInstance   *crypto.Crypto
	cryptoPrivateKey crypto.PrivateKey
)

func main() {
	cryptoInstance := &crypto.Crypto{}
	cryptoPrivateKey, _ = cryptoInstance.ImportPrivateKey([]byte(os.Getenv("APP_KEY")))

	http.HandleFunc("/authenticate", auth)
	http.HandleFunc("/virgil-jwt", provideJWT)

	http.ListenAndServe("localhost:3000", nil)
}

func auth(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	type packet struct {
		Identity string `json:"identity"`
	}
	packetAuth := packet{}
	err := json.Unmarshal(body, &packetAuth)
	if err != nil || packetAuth.Identity == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authToken := newToken()
	storageSet(authToken, packetAuth.Identity)

	resp, _ := json.Marshal(struct {
		AuthToken string `json:"authToken"`
	}{
		AuthToken: authToken,
	})
	w.Write(resp)
}

func provideJWT(w http.ResponseWriter, r *http.Request) {
	authHeaders := strings.Fields(r.Header.Get("Authorization"))

	if 2 < len(authHeaders) || "Bearer" != authHeaders[0] || !tokenExists(authHeaders[1]) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jwt, err := generateJWT(storageGet(authHeaders[1]))
	if err != nil || nil == jwt {
		w.WriteHeader(http.StatusInternalServerError)
	}

	resp, _ := json.Marshal(struct {
		JWT string `json:"virgilToken"`
	}{
		JWT: jwt.String(),
	})

	w.Write(resp)
}

func newToken() string {
	var symbols = []byte("abcdefghijklmnopqrstuvwxyz0123456789")
	randomstring := make([]byte, 32)
	for i := 0; i < 32; i++ {
		randomstring[i] = symbols[rand.Intn(len(symbols)-1)]
	}

	return base64.StdEncoding.EncodeToString(randomstring)
}

func storageSet(token string, identity string) {
	userStorage.Store(token, identity)
}

func storageGet(token string) string {
	identity, ok := userStorage.Load(token)
	if ok {
		return identity.(string)
	}
	return ""
}

func tokenExists(token string) bool {
	_, ok := userStorage.Load(token)
	return ok
}

func generateJWT(identity string) (*session.Jwt, error) {
	tokenSigner := &session.VirgilAccessTokenSigner{}
	generator := &session.JwtGenerator{AppKey: cryptoPrivateKey, AppKeyID: os.Getenv("APP_KEY_ID"), AppID: os.Getenv("APP_ID"), AccessTokenSigner: tokenSigner, TTL: time.Hour}
	return generator.GenerateToken(identity, nil)
}
