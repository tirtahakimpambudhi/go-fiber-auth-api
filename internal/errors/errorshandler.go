package errors

import (
	"net/http"
	"strings"

	"github.com/tirtahakimpambudhi/restful_api/internal/model/response"
)

// TypeErr represents the type of error with a string code.
type TypeErr string

// Error codes for different types of errors.
const (
	INTERNAL_SERVER_ERROR TypeErr = "INTERNAL_SERVER_ERROR"
	CONFLICT              TypeErr = "CONFLICT"
	FORBIDEN              TypeErr = "FORBIDDEN"
	UNAUTHORIZE           TypeErr = "UNAUTHORIZED"
	UNPROCESS_ENITITY     TypeErr = "UNPROCESSABLE_ENTITY"
	BAD_REQUEST           TypeErr = "BAD_REQUEST"
	REQUEST_TIMEOUT       TypeErr = "REQUEST_TIMEOUT"
	NOT_FOUND             TypeErr = "NOT_FOUND"
	TO_MANY_REQUEST       TypeErr = "TO_MANY_REQUEST"
)

// NewError creates a new response.Error based on the provided code and detail message.
func NewError(code TypeErr, detail string) *response.Error {
	// Initialize a new response.Error object.
	err := new(response.Error)

	// Set the title of the error based on the error code.
	err.Title = string(code)
	err.Code = string(code)
	// Assign HTTP status codes based on the error code.
	err.Status = ConvertStringToStatusCode(string(code))
	// Set the detailed error message.
	err.Detail = detail

	// Return the constructed error object.
	return err
}

// ConvertHTTPStatusCodeToString converts an HTTP status code to its corresponding string description.
func ConvertStatusCodeToString(statusCode int) string {
	switch statusCode {
	case http.StatusOK:
		return "OK"
	case http.StatusCreated:
		return "CREATED"
	case http.StatusAccepted:
		return "ACCEPTED"
	case http.StatusNoContent:
		return "NO_CONTENT"
	case http.StatusMovedPermanently:
		return "MOVED_PERMANENTLY"
	case http.StatusFound:
		return "FOUND"
	case http.StatusSeeOther:
		return "SEE_OTHER"
	case http.StatusNotModified:
		return "NOT_MODIFIED"
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusPaymentRequired:
		return "PAYMENT_REQUIRED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusMethodNotAllowed:
		return "METHOD_NOT_ALLOWED"
	case http.StatusNotAcceptable:
		return "NOT_ACCEPTABLE"
	case http.StatusProxyAuthRequired:
		return "PROXY_AUTHENTICATION_REQUIRED"
	case http.StatusRequestTimeout:
		return "REQUEST_TIMEOUT"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusGone:
		return "GONE"
	case http.StatusLengthRequired:
		return "LENGTH_REQUIRED"
	case http.StatusPreconditionFailed:
		return "PRECONDITION_FAILED"
	case http.StatusRequestEntityTooLarge:
		return "REQUEST_ENTITY_TOO_LARGE"
	case http.StatusRequestURITooLong:
		return "REQUEST_URI_TOO_LONG"
	case http.StatusUnsupportedMediaType:
		return "UNSUPPORTED_MEDIA_TYPE"
	case http.StatusRequestedRangeNotSatisfiable:
		return "REQUESTED_RANGE_NOT_SATISFIABLE"
	case http.StatusExpectationFailed:
		return "EXPECTATION_FAILED"
	case http.StatusTeapot:
		return "I_AM_A_TEAPOT"
	case http.StatusMisdirectedRequest:
		return "MISDIRECTED_REQUEST"
	case http.StatusUnprocessableEntity:
		return "UNPROCESSABLE_ENTITY"
	case http.StatusLocked:
		return "LOCKED"
	case http.StatusFailedDependency:
		return "FAILED_DEPENDENCY"
	case http.StatusTooEarly:
		return "TOO_EARLY"
	case http.StatusUpgradeRequired:
		return "UPGRADE_REQUIRED"
	case http.StatusPreconditionRequired:
		return "PRECONDITION_REQUIRED"
	case http.StatusTooManyRequests:
		return "TOO_MANY_REQUESTS"
	case http.StatusRequestHeaderFieldsTooLarge:
		return "REQUEST_HEADER_FIELDS_TOO_LARGE"
	case http.StatusUnavailableForLegalReasons:
		return "UNAVAILABLE_FOR_LEGAL_REASONS"
	case http.StatusInternalServerError:
		return "INTERNAL_SERVER_ERROR"
	case http.StatusNotImplemented:
		return "NOT_IMPLEMENTED"
	case http.StatusBadGateway:
		return "BAD_GATEWAY"
	case http.StatusServiceUnavailable:
		return "SERVICE_UNAVAILABLE"
	case http.StatusGatewayTimeout:
		return "GATEWAY_TIMEOUT"
	case http.StatusHTTPVersionNotSupported:
		return "HTTP_VERSION_NOT_SUPPORTED"
	case http.StatusVariantAlsoNegotiates:
		return "VARIANT_ALSO_NEGOTIATES"
	case http.StatusInsufficientStorage:
		return "INSUFFICIENT_STORAGE"
	case http.StatusLoopDetected:
		return "LOOP_DETECTED"
	case http.StatusNotExtended:
		return "NOT_EXTENDED"
	case http.StatusNetworkAuthenticationRequired:
		return "NETWORK_AUTHENTICATION_REQUIRED"
	default:
		return "UNKNOWN_STATUS"
	}
}

// ConvertStringToStatusCode converts a string description of an HTTP status to its corresponding status code.
func ConvertStringToStatusCode(status string) int {
	switch strings.ToUpper(status) {
	case "OK":
		return http.StatusOK
	case "CREATED":
		return http.StatusCreated
	case "ACCEPTED":
		return http.StatusAccepted
	case "NO_CONTENT":
		return http.StatusNoContent
	case "MOVED_PERMANENTLY":
		return http.StatusMovedPermanently
	case "FOUND":
		return http.StatusFound
	case "SEE_OTHER":
		return http.StatusSeeOther
	case "NOT_MODIFIED":
		return http.StatusNotModified
	case "BAD_REQUEST":
		return http.StatusBadRequest
	case "UNAUTHORIZED":
		return http.StatusUnauthorized
	case "PAYMENT_REQUIRED":
		return http.StatusPaymentRequired
	case "FORBIDDEN":
		return http.StatusForbidden
	case "NOT_FOUND":
		return http.StatusNotFound
	case "METHOD_NOT_ALLOWED":
		return http.StatusMethodNotAllowed
	case "NOT_ACCEPTABLE":
		return http.StatusNotAcceptable
	case "PROXY_AUTHENTICATION_REQUIRED":
		return http.StatusProxyAuthRequired
	case "REQUEST_TIMEOUT":
		return http.StatusRequestTimeout
	case "CONFLICT":
		return http.StatusConflict
	case "GONE":
		return http.StatusGone
	case "LENGTH_REQUIRED":
		return http.StatusLengthRequired
	case "PRECONDITION_FAILED":
		return http.StatusPreconditionFailed
	case "REQUEST_ENTITY_TOO_LARGE":
		return http.StatusRequestEntityTooLarge
	case "REQUEST_URI_TOO_LONG":
		return http.StatusRequestURITooLong
	case "UNSUPPORTED_MEDIA_TYPE":
		return http.StatusUnsupportedMediaType
	case "REQUESTED_RANGE_NOT_SATISFIABLE":
		return http.StatusRequestedRangeNotSatisfiable
	case "EXPECTATION_FAILED":
		return http.StatusExpectationFailed
	case "I_AM_A_TEAPOT":
		return http.StatusTeapot
	case "MISDIRECTED_REQUEST":
		return http.StatusMisdirectedRequest
	case "UNPROCESSABLE_ENTITY":
		return http.StatusUnprocessableEntity
	case "LOCKED":
		return http.StatusLocked
	case "FAILED_DEPENDENCY":
		return http.StatusFailedDependency
	case "TOO_EARLY":
		return http.StatusTooEarly
	case "UPGRADE_REQUIRED":
		return http.StatusUpgradeRequired
	case "PRECONDITION_REQUIRED":
		return http.StatusPreconditionRequired
	case "TOO_MANY_REQUESTS":
		return http.StatusTooManyRequests
	case "REQUEST_HEADER_FIELDS_TOO_LARGE":
		return http.StatusRequestHeaderFieldsTooLarge
	case "UNAVAILABLE_FOR_LEGAL_REASONS":
		return http.StatusUnavailableForLegalReasons
	case "INTERNAL_SERVER_ERROR":
		return http.StatusInternalServerError
	case "NOT_IMPLEMENTED":
		return http.StatusNotImplemented
	case "BAD_GATEWAY":
		return http.StatusBadGateway
	case "SERVICE_UNAVAILABLE":
		return http.StatusServiceUnavailable
	case "GATEWAY_TIMEOUT":
		return http.StatusGatewayTimeout
	case "HTTP_VERSION_NOT_SUPPORTED":
		return http.StatusHTTPVersionNotSupported
	case "VARIANT_ALSO_NEGOTIATES":
		return http.StatusVariantAlsoNegotiates
	case "INSUFFICIENT_STORAGE":
		return http.StatusInsufficientStorage
	case "LOOP_DETECTED":
		return http.StatusLoopDetected
	case "NOT_EXTENDED":
		return http.StatusNotExtended
	case "NETWORK_AUTHENTICATION_REQUIRED":
		return http.StatusNetworkAuthenticationRequired
	default:
		return -1 // Return -1 for unknown status descriptions
	}
}
