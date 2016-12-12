package api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

// Handler is base object to attach all REST methods
type Handler struct {
	*mux.Router
	// RequestID returns ID for the current request
	RequestID func(req *http.Request) string
	// Send details for 500 erorrs or just generic "internal error"
	Enable500WithDetails bool
}

// ErrorRepr is JSON serialization for error
type ErrorRepr struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

// ServeHTTP wraps every request with logging, error serialization, rollbar
func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var (
		statusCode string
		requestID  string
	)

	start := time.Now().UTC()
	statusCode = "200"

	if h.RequestID != nil {
		requestID = h.RequestID(req)
	} else {
		requestID = UUIDRequestID(req)
	}

	defer func() {
		if p := recover(); p != nil {
			switch e := p.(type) {
			case *Error:
				statusCode = strconv.Itoa(e.Code)

				rw.Header().Set("Content-Type", "application/json")
				rw.WriteHeader(e.Code)

				json.NewEncoder(rw).Encode(&ErrorRepr{Code: e.ExtendedCode, Message: e.Message})
			case *RawError:
				statusCode = strconv.Itoa(e.Code)

				if e.ContentType != "" {
					rw.Header().Set("Content-Type", e.ContentType)
				}

				for key, value := range e.Headers {
					rw.Header().Set(key, value)
				}
				rw.WriteHeader(e.Code)
				rw.Write(e.Body)
			default:
				statusCode = "500"

				rw.Header().Set("Content-Type", "application/json")
				rw.WriteHeader(500)

				if h.Enable500WithDetails {
					json.NewEncoder(rw).Encode(&ErrorRepr{
						Message: fmt.Sprintf("%s", p),
					})
				} else {
					rw.Write([]byte("{\"Message\":\"internal error\"}"))
				}

				stack := make([]byte, 1024*1024)
				n := runtime.Stack(stack, false)

				log.Printf("[%s] [%s] panic:\n%s", requestID, p, stack[:n])

			}
		}

		responseTime := time.Since(start)

		log.Printf("[%s] << %s %s %s %.3fs", requestID, statusCode, req.Method, req.RequestURI, responseTime.Seconds())
	}()

	log.Printf("[%s] >> %s %s ...", requestID, req.Method, req.RequestURI)
	rw.Header().Set("X-Request-ID", requestID)

	h.Router.ServeHTTP(rw, req)
}

// UUIDRequestID generates random UUID for each request
func UUIDRequestID(req *http.Request) string {
	return uuid.NewV4().String()
}

// JSON sends application/json response
func (h *Handler) JSON(rw http.ResponseWriter, response interface{}) {
	h.JSONStatus(rw, 200, response)
}

// JSONStatus sends response with code
func (h *Handler) JSONStatus(rw http.ResponseWriter, status int, response interface{}) {
	if status != 204 {
		rw.Header().Set("Content-Type", "application/json")
	}
	rw.WriteHeader(status)

	if status == 204 {
		return
	}

	e := json.NewEncoder(rw)
	Check(e.Encode(response))
}

// XML sends application/xml response
func (h *Handler) XML(rw http.ResponseWriter, response interface{}) {
	h.XMLStatus(rw, 200, response)
}

// XMLStatus sends response with code
func (h *Handler) XMLStatus(rw http.ResponseWriter, status int, response interface{}) {
	if status != 204 {
		rw.Header().Set("Content-Type", "application/xml")
	}
	rw.WriteHeader(status)

	if status == 204 {
		return
	}

	e := xml.NewEncoder(rw)
	Check(e.EncodeToken(xml.ProcInst{Target: "xml", Inst: []byte("version=\"1.0\" encoding=\"UTF-8\"")}))
	Check(e.Encode(response))
}

// EmptyStatus sends empty response with specified status
func (h *Handler) EmptyStatus(rw http.ResponseWriter, status int) {
	rw.WriteHeader(status)
}
