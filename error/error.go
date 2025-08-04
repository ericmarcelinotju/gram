package domain

import (
	"errors"
	"fmt"
)

const (
	// NotFound error indicates a missing / not found record
	NotFoundError        = "NotFound"
	notFoundErrorMessage = "Record not found"

	// ValidationError indicates an error in payload validation
	ValidationError        = "ValidationError"
	validationErrorMessage = "Validation error"

	// ResourceAlreadyExistsError indicates a duplicate / already existing record
	ResourceAlreadyExistsError = "ResourceAlreadyExists"
	alreadyExistsErrorMessage  = "Resource already exists"

	// RepositoryError indicates a reposity error
	RepositoryError        = "RepositoryError"
	repositoryErrorMessage = "Error in repository operation"

	// CacheError indicates a cache error
	CacheError        = "CacheError"
	cacheErrorMessage = "Error in cache operation"

	// DatabaseError indicates a database error
	DatabaseError        = "DatabaseError"
	databaseErrorMessage = "Error in database operation"

	// HTTPClientError indicates a http client error
	HTTPClientError        = "HTTPClientError"
	httpClientErrorMessage = "Error in http client operation"

	// NotAuthenticated indicates an authentication error
	NotAuthenticated             = "NotAuthenticated"
	notAuthenticatedErrorMessage = "Not Authenticated"

	// TokenGeneratorError indicates an token generation error
	TokenGeneratorError        = "TokenGeneratorError"
	tokenGeneratorErrorMessage = "Error in token generation"

	// NotAuthorized indicates an authorization error
	NotAuthorized             = "NotAuthorized"
	notAuthorizedErrorMessage = "Not Authorized"

	// NotAuthorized indicates an authorization error
	UnsupportedError        = "Unsupported"
	unsupportedErrorMessage = "payload not supported"

	// DismissedError indicates a dismissed process error
	DismissedError        = "Dismissed"
	dismissedErrorMessage = "Operation dismissed"

	// StorageError indicates an error in storage operation
	StorageError        = "StorageError"
	storageErrorMessage = "Storage error"

	// UnknownError indicates an error that the app cannot find the cause for
	UnknownError        = "UnknownError"
	unknownErrorMessage = "Something went wrong"
)

var (
	// ErrNotFound not found error
	ErrNotFound = errors.New(NotFoundError)
	// ErrValidationError invalid payload
	ErrValidationError = errors.New(ValidationError)
	// ErrResourceAlreadyExists duplicate
	ErrResourceAlreadyExists = errors.New(ResourceAlreadyExistsError)
	// ErrRepositoryError general repositoty error
	ErrRepositoryError = errors.New(RepositoryError)
	// ErrCacheError cache error
	ErrCacheError = errors.New(CacheError)
	// ErrDatabaseError query error
	ErrDatabaseError = errors.New(DatabaseError)
	// ErrHTTPClientError HTTP Client request
	ErrHTTPClientError = errors.New(HTTPClientError)
	// ErrNotAuthenticated authentication required
	ErrNotAuthenticated = errors.New(NotAuthenticated)
	// ErrNotAuthorized authorization invalid
	ErrNotAuthorized = errors.New(NotAuthorized)
	// ErrTokenGeneratorError token generation error
	ErrTokenGeneratorError = errors.New(TokenGeneratorError)
	// ErrUnsupported indicates an error because the app have not support the client request
	ErrUnsupported = errors.New(UnsupportedError)
	// ErrDismissed indicates an error because the app have dismissed the request
	ErrDismissed = errors.New(DismissedError)
	// ErrStorage indicates an error when doing storage operation
	ErrStorage = errors.New(StorageError)
	// ErrUnknownError indicates an error that the app cannot find the cause for
	ErrUnknownError = errors.New(UnknownError)
)

// AppError defines an application (domain) error
type AppError struct {
	Err  error
	Type string
}

// NewAppError initializes a new domain error using an error and its type.
func NewAppError(err error, errType string) *AppError {
	return &AppError{
		Err:  err,
		Type: errType,
	}
}

// NewAppErrorWithType initializes a new default error for a given type.
func NewAppErrorWithType(errType string) *AppError {
	var err error

	switch errType {
	case NotFoundError:
		err = ErrNotFound
	case ValidationError:
		err = ErrValidationError
	case ResourceAlreadyExistsError:
		err = ErrResourceAlreadyExists
	case CacheError:
		err = ErrCacheError
	case DatabaseError:
		err = ErrDatabaseError
	case HTTPClientError:
		err = ErrHTTPClientError
	case NotAuthenticated:
		err = ErrNotAuthenticated
	case NotAuthorized:
		err = ErrNotAuthorized
	case TokenGeneratorError:
		err = ErrTokenGeneratorError
	case UnsupportedError:
		err = ErrUnsupported
	case DismissedError:
		err = ErrDismissed
	case StorageError:
		err = ErrStorage
	default:
		err = ErrUnknownError
	}

	return &AppError{
		Err:  err,
		Type: errType,
	}
}

// String converts the app error to a human-readable string.
func (appErr *AppError) Error() string {
	return fmt.Sprintf("%s : %s", appErr.Type, appErr.Err.Error())
}

// Is check error type.
func (appErr *AppError) Is(err error) bool {
	return appErr.Type == err.Error()
}
