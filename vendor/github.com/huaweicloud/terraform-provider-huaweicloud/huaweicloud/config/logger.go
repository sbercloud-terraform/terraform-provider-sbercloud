package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/unknwon/com"
)

// MAXFieldLength is the maximum string length of single field when logging
const MAXFieldLength int = 256

var maxTimeout = 10 * time.Minute

// LogRoundTripper satisfies the http.RoundTripper interface and is used to
// customize the default http client RoundTripper to allow for logging.
type LogRoundTripper struct {
	Rt         http.RoundTripper
	OsDebug    bool
	MaxRetries int
}

func retryTimeout(count int) time.Duration {
	seconds := math.Pow(2, float64(count))
	timeout := time.Duration(seconds) * time.Second
	if timeout > maxTimeout { // won't wait more than maxTimeout
		timeout = maxTimeout
	}
	return timeout
}

// RoundTrip performs a round-trip HTTP request and logs relevant information about it.
func (lrt *LogRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	defer func() {
		if request.Body != nil {
			request.Body.Close()
		}
	}()

	// for future reference, this is how to access the Transport struct:
	//tlsconfig := lrt.Rt.(*http.Transport).TLSClientConfig

	var err error

	if lrt.OsDebug {
		log.Printf("[DEBUG] API Request URL: %s %s", request.Method, request.URL)
		log.Printf("[DEBUG] API Request Headers:\n%s", FormatHeaders(request.Header, "\n"))

		if request.Body != nil {
			request.Body, err = lrt.logRequest(request.Body, request.Header.Get("Content-Type"))
			if err != nil {
				return nil, err
			}
		}
	}

	response, err := lrt.Rt.RoundTrip(request)
	if response == nil {
		errMessage := err.Error()
		if strings.Contains(errMessage, "no such host") {
			return nil, err
		}
	}

	// Retrying connection
	retry := 1
	for response == nil {

		if retry > lrt.MaxRetries {
			if lrt.OsDebug {
				log.Printf("[DEBUG] connection error, retries exhausted. Aborting")
			}
			err = fmt.Errorf("connection error, retries exhausted. Aborting. Last error was: %s", err)
			return nil, err
		}

		if lrt.OsDebug {
			log.Printf("[DEBUG] connection error, retry number %d: %s", retry, err)
		}
		//lintignore:R018
		time.Sleep(retryTimeout(retry))
		response, err = lrt.Rt.RoundTrip(request)
		retry++
	}

	if lrt.OsDebug {
		log.Printf("[DEBUG] API Response Code: %d", response.StatusCode)
		log.Printf("[DEBUG] API Response Headers:\n%s", FormatHeaders(response.Header, "\n"))

		response.Body, err = lrt.logResponse(response.Body, response.Header.Get("Content-Type"))
	}

	return response, err
}

// logRequest will log the HTTP Request details.
// If the body is JSON, it will attempt to be pretty-formatted.
func (lrt *LogRoundTripper) logRequest(original io.ReadCloser, contentType string) (io.ReadCloser, error) {
	defer original.Close()

	var bs bytes.Buffer
	_, err := io.Copy(&bs, original)
	if err != nil {
		return nil, err
	}

	// Handle request contentType
	if strings.HasPrefix(contentType, "application/json") {
		debugInfo := formatJSON(bs.Bytes(), true)
		log.Printf("[DEBUG] API Request Body: %s", debugInfo)
	} else {
		log.Printf("[DEBUG] Not logging because the request body isn't JSON")
	}

	return ioutil.NopCloser(strings.NewReader(bs.String())), nil
}

// logResponse will log the HTTP Response details.
// If the body is JSON, it will attempt to be pretty-formatted.
func (lrt *LogRoundTripper) logResponse(original io.ReadCloser, contentType string) (io.ReadCloser, error) {
	if strings.HasPrefix(contentType, "application/json") {
		var bs bytes.Buffer
		defer original.Close()
		_, err := io.Copy(&bs, original)
		if err != nil {
			return nil, err
		}
		debugInfo := formatJSON(bs.Bytes(), true)
		if debugInfo != "" {
			log.Printf("[DEBUG] API Response Body: %s", debugInfo)
		}
		return ioutil.NopCloser(strings.NewReader(bs.String())), nil
	}

	log.Printf("[DEBUG] Not logging because the response body isn't JSON")
	return original, nil
}

// formatJSON will try to pretty-format a JSON body.
// It will also mask known fields which contain sensitive information.
func formatJSON(raw []byte, maskBody bool) string {
	var data map[string]interface{}

	err := json.Unmarshal(raw, &data)
	if err != nil {
		log.Printf("[DEBUG] Unable to parse JSON: %s", err)
		return string(raw)
	}

	// Mask known password fields
	if maskBody {
		maskSecurityFields(data)
	}

	// Ignore the catalog
	if _, ok := data["catalog"]; ok {
		return "{ **skipped** }"
	}
	if v, ok := data["token"].(map[string]interface{}); ok {
		if _, ok := v["catalog"]; ok {
			return ""
		}
	}

	pretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("[DEBUG] Unable to re-marshal JSON: %s", err)
		return string(raw)
	}

	return string(pretty)
}

// REDACT_HEADERS is a list of headers that need to be redacted
var REDACT_HEADERS = []string{
	"x-auth-token", "x-security-token", "x-service-token",
	"x-subject-token", "x-storage-token", "authorization",
}

// RedactHeaders processes a headers object, returning a redacted list
func RedactHeaders(headers http.Header) (processedHeaders []string) {
	for name, header := range headers {
		for _, v := range header {
			if com.IsSliceContainsStr(REDACT_HEADERS, name) {
				processedHeaders = append(processedHeaders, fmt.Sprintf("%v: %v", name, "***"))
			} else {
				processedHeaders = append(processedHeaders, fmt.Sprintf("%v: %v", name, v))
			}
		}
	}
	return
}

// FormatHeaders processes a headers object plus a deliminator, returning a string
func FormatHeaders(headers http.Header, seperator string) string {
	redactedHeaders := RedactHeaders(headers)
	sort.Strings(redactedHeaders)

	return strings.Join(redactedHeaders, seperator)
}

func maskSecurityFields(data map[string]interface{}) bool {
	for k, val := range data {
		switch val.(type) {
		case string:
			if isSecurityFields(k) {
				data[k] = "***"
			} else if len(val.(string)) > MAXFieldLength {
				data[k] = "** large string **"
			}
		case map[string]interface{}:
			subData := val.(map[string]interface{})
			if masked := maskSecurityFields(subData); masked {
				return true
			}
		}
	}
	return false
}

func isSecurityFields(field string) bool {
	// "password" is apply to the most request JSON body
	// "secret" is apply to the AK/SK response JSON body
	if strings.Contains(field, "password") || strings.Contains(field, "secret") {
		return true
	}

	// "adminPass" is apply to the ecs/bms instance request JSON body
	// "adminPwd" is apply to the css cluster request JSON body
	// 'encrypted_user_data' is apply to the function request JSON body of FunctionGraph
	securityFields := []string{"adminPass", "adminPwd", "encrypted_user_data"}
	for _, key := range securityFields {
		if key == field {
			return true
		}
	}

	return false
}
