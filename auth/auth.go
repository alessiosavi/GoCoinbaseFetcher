package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"hash"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"
)

const URL = `https://api.exchange.coinbase.com`

var client = http.DefaultClient

type Auth struct {
	timestamp int64
	Path      string
	Body      []byte
}

func init() {
	client.Timeout = time.Second * 60
}
func (a *Auth) New() {
	a.timestamp = time.Now().Unix()
}

type SendOpts struct {
	Path   string
	Body   []byte
	Method string
}

func (a Auth) Send(opts SendOpts) {

	request, err := http.NewRequest(opts.Method, URL+opts.Path, bytes.NewReader(opts.Body))
	if err != nil {
		panic(err)
	}
	for k, v := range a.GetHeaders(opts) {
		request.Header.Add(k, v)
	}

	dumpRequest, err := httputil.DumpRequest(request, true)
	if err != nil {
		panic(err)
	}
	log.Println(string(dumpRequest))

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	dumpResponse, err := httputil.DumpResponse(response, true)
	if err != nil {
		panic(err)
	}
	log.Println(string(dumpResponse))
}

func (a Auth) GetHeaders(opts SendOpts) map[string]string {
	headers := make(map[string]string)
	headers["CB-ACCESS-KEY"] = KEY
	headers["CB-ACCESS-TIMESTAMP"] = strconv.FormatInt(a.timestamp, 10)
	headers["CB-ACCESS-PASSPHRASE"] = SECRET_PASS

	message := headers["CB-ACCESS-TIMESTAMP"] + opts.Method + opts.Path
	if opts.Body != nil {
		bodyBytes, _ := json.Marshal(opts.Body)
		message += string(bodyBytes)
	}

	_hmac, err := GetHMAC(HashSHA256,
		[]byte(message),
		[]byte(SECRET))
	if err != nil {
		panic(err)
	}
	headers["CB-ACCESS-SIGN"] = Base64Encode(_hmac)
	headers["Content-Type"] = "application/json"

	return headers
}

// Base64Encode takes in a byte array then returns an encoded base64 string
func Base64Encode(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

// GetHMAC returns a keyed-hash message authentication code using the desired hashtype
func GetHMAC(hashType int, input, key []byte) ([]byte, error) {
	var hasher func() hash.Hash

	switch hashType {
	case HashSHA1:
		hasher = sha1.New
	case HashSHA256:
		hasher = sha256.New
	case HashSHA512:
		hasher = sha512.New
	case HashSHA512_384:
		hasher = sha512.New384
	case HashMD5:
		hasher = md5.New
	}

	h := hmac.New(hasher, key)
	_, err := h.Write(input)
	return h.Sum(nil), err
}

// HexEncodeToString takes in a hexadecimal byte array and returns a string

// Base64Decode takes in a Base64 string and returns a byte array and an error

// GetRandomSalt returns a random salt

// GetMD5 returns a MD5 hash of a byte array

// GetSHA512 returns a SHA512 hash of a byte array

// GetSHA256 returns a SHA256 hash of a byte array

// Sha1ToHex takes a string, sha1 hashes it and return a hex string of the
// result

// Constants
const (
	HashSHA1 = iota
	HashSHA256
	HashSHA512
	HashSHA512_384
	HashMD5
)
