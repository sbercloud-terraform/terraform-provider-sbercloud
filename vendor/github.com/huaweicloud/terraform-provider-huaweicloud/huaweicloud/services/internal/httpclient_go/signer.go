package httpclient_go

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	BasicDateFormat     = "20060102T150405Z"
	Algorithm           = "SDK-HMAC-SHA256"
	HeaderXDate         = "X-Sdk-Date"
	HeaderHost          = "host"
	HeaderAuthorization = "Authorization"
	HeaderContentSha256 = "X-Sdk-content-sha256"
)

func CanonicalRequest(r *http.Request, signedHeaders []string) (string, error) {
	var hexencode string
	var err error
	if hex := r.Header.Get(HeaderContentSha256); hex != "" {
		hexencode = hex
	} else {
		data, err := RequestPayload(r)
		if err != nil {
			return "", err
		}
		hexencode, err = HexEncodeSHA256Hash(data)
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", r.Method, CanonicalURI(r), CanonicalQueryString(r), CanonicalHeaders(r, signedHeaders), strings.Join(signedHeaders, ";"), hexencode), err
}

func HexEncodeSHA256Hash(body []byte) (string, error) {
	hash := sha256.New()
	if body == nil {
		body = []byte("")
	}
	_, err := hash.Write(body)
	return fmt.Sprintf("%x", hash.Sum(nil)), err
}

func CanonicalURI(r *http.Request) string {
	pattens := strings.Split(r.URL.Path, "/")
	var uri []string
	for _, v := range pattens {
		uri = append(uri, escape(v))
	}
	urlPath := strings.Join(uri, "/")
	if len(urlPath) == 0 || urlPath[len(urlPath)-1] != '/' {
		urlPath = urlPath + "/"
	}
	return urlPath
}

func CanonicalQueryString(r *http.Request) string {
	var keys []string
	query := r.URL.Query()
	for key := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var a []string
	for _, key := range keys {
		k := escape(key)
		sort.Strings(query[key])
		for _, v := range query[key] {
			kv := fmt.Sprintf("%s=%s", k, escape(v))
			a = append(a, kv)
		}
	}
	queryStr := strings.Join(a, "&")
	r.URL.RawQuery = queryStr
	return queryStr
}

func CanonicalHeaders(r *http.Request, signerHeaders []string) string {
	var a []string
	header := make(map[string][]string)
	for k, v := range r.Header {
		header[strings.ToLower(k)] = v
	}
	for _, key := range signerHeaders {
		value := header[key]
		if strings.EqualFold(key, HeaderHost) {
			value = []string{r.Host}
		}
		sort.Strings(value)
		for _, v := range value {
			a = append(a, key+":"+strings.TrimSpace(v))
		}
	}
	return fmt.Sprintf("%s\n", strings.Join(a, "\n"))
}

func RequestPayload(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return []byte(""), nil
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return []byte(""), err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(b))
	return b, err
}

type Signer struct {
	Key    string
	Secret string
}

func (s *Signer) Sign(r *http.Request) error {
	var t time.Time
	var err error
	var dt string
	if dt = r.Header.Get(HeaderXDate); dt != "" {
		t, err = time.Parse(BasicDateFormat, dt)
	}
	if err != nil || dt == "" {
		t = time.Now()
		r.Header.Set(HeaderXDate, t.UTC().Format(BasicDateFormat))
	}
	signedHeaders := SignedHeaders(r)
	canonicalRequest, err := CanonicalRequest(r, signedHeaders)
	if err != nil {
		return err
	}
	stringToSign, err := StringToSign(canonicalRequest, t)
	if err != nil {
		return err
	}
	signature, err := SignStringToSign(stringToSign, []byte(s.Secret))
	if err != nil {
		return err
	}
	authValue := AuthHeaderValue(signature, s.Key, signedHeaders)
	r.Header.Set(HeaderAuthorization, authValue)
	return nil
}

func AuthHeaderValue(signature, accessKey string, signedHeaders []string) string {
	return fmt.Sprintf("%s Access=%s, SignedHeaders=%s, Signature=%s", Algorithm, accessKey, strings.Join(signedHeaders, ";"), signature)
}

func SignedHeaders(r *http.Request) []string {
	var a []string
	for key := range r.Header {
		a = append(a, strings.ToLower(key))
	}
	sort.Strings(a)
	return a
}

func StringToSign(canonicalRequest string, t time.Time) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(canonicalRequest))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n%s\n%x",
		Algorithm, t.UTC().Format(BasicDateFormat), hash.Sum(nil)), nil
}

// Create the HWS Signature.
func SignStringToSign(stringToSign string, signingKey []byte) (string, error) {
	hm, err := hmacsha256(signingKey, stringToSign)
	return fmt.Sprintf("%x", hm), err
}

func hmacsha256(key []byte, data string) ([]byte, error) {
	h := hmac.New(sha256.New, key)
	if _, err := h.Write([]byte(data)); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
