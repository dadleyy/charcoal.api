package defs

const (
	// ErrUnauthorizedAdmin response for invalid token on admin apis.
	ErrUnauthorizedAdmin = "UNAUTHORIZED"

	// ErrDuplicateEntry response for duplicate records based on uniqueness constraints.
	ErrDuplicateEntry = "DUPLICATE_ENTRY"

	// ErrFailedQuery response for failed server requests.
	ErrFailedQuery = "BAD_QUERY"

	// ErrUserManagerUnauthorizedDomain response for invalid login attempts.
	ErrUserManagerUnauthorizedDomain = "UNAUTHORIZED_DOMAIN"

	// ErrUserManagerDuplicate response when user records are duplicated.
	ErrUserManagerDuplicate = "DUPLICATE_USER"

	// ErrGoogleBadAuthCode bad response from google during oauth.
	ErrGoogleBadAuthCode = "BAD_AUTH_CODE"

	// ErrGoogleNoClientAssociated unable to find client from oauth scope.
	ErrGoogleNoClientAssociated = "NO_ASSOICATED_CLIENT"

	// ErrGoogleMissingClientRedirect client has not configured redirect uri.
	ErrGoogleMissingClientRedirect = "NO_REDIRECT_URI"

	// ErrGoogleMissingAuthEndpoint missing config.
	ErrGoogleMissingAuthEndpoint = "BAD_AUTH_ENDPOINT"

	// ErrGoogleInvalidGoogleResponse response when something unknown happens.
	ErrGoogleInvalidGoogleResponse = "BAD_GOOGLE_RESPONSE"

	// ErrBadImageType error string for bad image upload types.
	ErrBadImageType = "BAD_IMAGE_TYPE"

	// ErrBadImageUuid response for invalid image uuids.
	ErrBadImageUuid = "BAD_UUID_GENERATED"

	// ErrBadBearerToken response when user is required for auth but is not present.
	ErrBadBearerToken = "ERR_BAD_BEARER"

	// ErrBadS3Response response when s3 communication fails.
	ErrBadS3Response = "BAD_S3_RESPONSE"
)
