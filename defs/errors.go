package defs

const (
	ErrUnauthorizedAdmin             = "UNAUTHORIZED"
	ErrDuplicateEntry                = "DUPLICATE_ENTRY"
	ErrFailedQuery                   = "BAD_QUERY"
	ErrUserManagerUnauthorizedDomain = "unauthorized-domain"
	ErrUserManagerDuplicate          = "duplicate-user"

	ErrGoogleBadAuthCode           = "BAD_AUTH_CODE"
	ErrGoogleNoClientAssociated    = "NO_ASSOICATED_CLIENT"
	ErrGoogleMissingClientRedirect = "NO_REDIRECT_URI"
	ErrGoogleMissingAuthEndpoint   = "BAD_AUTH_ENDPOINT"
	ErrGoogleInvalidGoogleResponse = "BAD_GOOGLE_RESPONSE"
)
