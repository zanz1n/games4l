package utils

import (
	"github.com/games4l/backend/libs/utils/httpcodes"
)

type MessageList struct {
	// When a entity of any type was not found
	EntityNotFound *string `json:"entity_not_found"`
	// When a ObjectID cannot be parsed because does not contains a valid hexadecimal sequence
	InvalidObjectID *string `json:"invalid_object_id"`
	// When the user cannot be not found or the password does not match in a auth proccess
	UserUnauthorized *string `json:"user_unauthorized"`
	// When the jwt is not valid or has an old format
	InvalidJwtTokenFormat *string `json:"invalid_jwt_token_format"`
	// When the jwt token is expired
	JwtTokenExpired *string `json:"jwt_token_expired"`
	// "Unauthorization" when the provided signature does not match the intended one
	InvalidAuthSignature *string `json:"invalid_auth_signature"`
	// When the provided signature is encoded with a unexpected method
	InvalidAuthSignatureEncodingMethod *string `json:"invalid_auth_signature_encoding_method"`
	// When something goes wrong reading or manipulating the body blob/buffer
	MalformedOrTooBigBody *string `json:"malformed_or_too_big_body"`
	// When something unexpected happens in some server routine
	InternalServerError *string `json:"internal_server_error"`
	// When the request body is invalidated after parsed
	InvalidRequestEntity *string `json:"invalid_request_entity"`
	// When the user tries to create something that (or it's unique index) already exists
	EntityAlreadyExists *string `json:"failed_to_create_conflict"`
	// When some server operation took longer than expected
	ServerOperationTookTooLong *string `json:"server_operation_took_too_long"`
	// Query string param errors
	SurnameSearchInvalid *string `json:"surname_search_invalid"`
	// When a provided email address is invalid
	InvalidEmail *string `json:"invalid_email"`
	// When the requested route does not exist in the server
	NoSuchRoute *string `json:"no_such_route"`
	// When the request method is not allowed
	MethodNotAllowed *string `json:"method_not_allowed"`
	// When the route requires admin privilleges
	RouteRequiresAdminAuth *string `json:"route_requires_admin_auth"`
	// When the provided auth header strategy is not supported
	InvalidAuthStrategy *string `json:"invalid_auth_strategy"`
	// Query string param errors
	InvalidFMTQueryParam *string `json:"invalid_fmt_query_param"`
	// Query string param errors
	InvalidNIDQueryParam *string `json:"invalid_nid_query_param"`
}

type ErrorList struct {
	EntityNotFound                     StatusCodeErr
	InvalidObjectID                    StatusCodeErr
	UserUnauthorized                   StatusCodeErr
	InvalidJwtTokenFormat              StatusCodeErr
	JwtTokenExpired                    StatusCodeErr
	InvalidAuthSignature               StatusCodeErr
	InvalidAuthSignatureEncodingMethod StatusCodeErr
	MalformedOrTooBigBody              StatusCodeErr
	InternalServerError                StatusCodeErr
	InvalidRequestEntity               StatusCodeErr
	EntityAlreadyExists                StatusCodeErr
	ServerOperationTookTooLong         StatusCodeErr
	SurnameSearchInvalid               StatusCodeErr
	InvalidEmail                       StatusCodeErr
	NoSuchRoute                        StatusCodeErr
	MethodNotAllowed                   StatusCodeErr
	RouteRequiresAdminAuth             StatusCodeErr
	InvalidAuthStrategy                StatusCodeErr
	InvalidFMTQueryParam               StatusCodeErr
	InvalidNIDQueryParam               StatusCodeErr
}

func (el *ErrorList) Apply(ml MessageList) {
	if ml.EntityNotFound != nil {
		el.EntityNotFound.SetMsg(*ml.EntityNotFound)
	}
	if ml.InvalidObjectID != nil {
		el.InvalidObjectID.SetMsg(*ml.InvalidObjectID)
	}
	if ml.UserUnauthorized != nil {
		el.UserUnauthorized.SetMsg(*ml.UserUnauthorized)
	}
	if ml.InvalidJwtTokenFormat != nil {
		el.InvalidJwtTokenFormat.SetMsg(*ml.InvalidJwtTokenFormat)
	}
	if ml.JwtTokenExpired != nil {
		el.JwtTokenExpired.SetMsg(*ml.JwtTokenExpired)
	}
	if ml.InvalidAuthSignature != nil {
		el.InvalidAuthSignature.SetMsg(*ml.InvalidAuthSignature)
	}
	if ml.InvalidAuthSignatureEncodingMethod != nil {
		el.InvalidAuthSignatureEncodingMethod.SetMsg(*ml.InvalidAuthSignatureEncodingMethod)
	}
	if ml.MalformedOrTooBigBody != nil {
		el.MalformedOrTooBigBody.SetMsg(*ml.MalformedOrTooBigBody)
	}
	if ml.InternalServerError != nil {
		el.InternalServerError.SetMsg(*ml.InternalServerError)
	}
	if ml.InvalidRequestEntity != nil {
		el.InvalidRequestEntity.SetMsg(*ml.InvalidRequestEntity)
	}
	if ml.EntityAlreadyExists != nil {
		el.EntityAlreadyExists.SetMsg(*ml.EntityAlreadyExists)
	}
	if ml.ServerOperationTookTooLong != nil {
		el.ServerOperationTookTooLong.SetMsg(*ml.ServerOperationTookTooLong)
	}
	if ml.SurnameSearchInvalid != nil {
		el.SurnameSearchInvalid.SetMsg(*ml.SurnameSearchInvalid)
	}
	if ml.InvalidEmail != nil {
		el.InvalidEmail.SetMsg(*ml.InvalidEmail)
	}
	if ml.NoSuchRoute != nil {
		el.NoSuchRoute.SetMsg(*ml.NoSuchRoute)
	}
	if ml.MethodNotAllowed != nil {
		el.MethodNotAllowed.SetMsg(*ml.MethodNotAllowed)
	}
	if ml.RouteRequiresAdminAuth != nil {
		el.RouteRequiresAdminAuth.SetMsg(*ml.RouteRequiresAdminAuth)
	}
	if ml.InvalidAuthStrategy != nil {
		el.InvalidAuthStrategy.SetMsg(*ml.InvalidAuthStrategy)
	}
	if ml.InvalidFMTQueryParam != nil {
		el.InvalidFMTQueryParam.SetMsg(*ml.InvalidFMTQueryParam)
	}
	if ml.InvalidNIDQueryParam != nil {
		el.InvalidNIDQueryParam.SetMsg(*ml.InvalidNIDQueryParam)
	}
}

var DefaultErrorList = ErrorList{
	EntityNotFound: NewStatusCodeErr(
		"the entity could not be found",
		httpcodes.StatusNotFound,
	),
	InvalidObjectID: NewStatusCodeErr(
		"the requested object id does not have a valid hex encoding",
		httpcodes.StatusBadRequest,
	),
	UserUnauthorized: NewStatusCodeErr(
		"user does not exist or the password do not match",
		httpcodes.StatusUnauthorized,
	),
	InvalidJwtTokenFormat: NewStatusCodeErr(
		"invalid auth token metadata format, please login again",
		httpcodes.StatusUnauthorized,
	),
	JwtTokenExpired: NewStatusCodeErr(
		"the token is expired, please login again",
		httpcodes.StatusUnauthorized,
	),
	InvalidAuthSignature: NewStatusCodeErr(
		"the provided signature is not valid",
		httpcodes.StatusUnauthorized,
	),
	InvalidAuthSignatureEncodingMethod: NewStatusCodeErr(
		"the signature encoding method is not valid",
		httpcodes.StatusUnauthorized,
	),
	MalformedOrTooBigBody: NewStatusCodeErr(
		"the request body is partial, malformed or too big",
		httpcodes.StatusBadRequest,
	),
	InternalServerError: NewStatusCodeErr(
		"something went wrong on our part while processing your request",
		httpcodes.StatusInternalServerError,
	),
	InvalidRequestEntity: NewStatusCodeErr(
		"the request body entity is invalid or could not be parsed",
		httpcodes.StatusBadRequest,
	),
	EntityAlreadyExists: NewStatusCodeErr(
		"failed to create because this entity already exists",
		httpcodes.StatusConflict,
	),
	ServerOperationTookTooLong: NewStatusCodeErr(
		"the operation took too long and was canceled in the process",
		httpcodes.StatusInternalServerError,
	),
	SurnameSearchInvalid: NewStatusCodeErr(
		"at least one surname must be provided",
		httpcodes.StatusBadRequest,
	),
	InvalidEmail: NewStatusCodeErr(
		"the provided email address is not valid",
		httpcodes.StatusBadRequest,
	),
	NoSuchRoute: NewStatusCodeErr(
		"no such route",
		httpcodes.StatusNotFound,
	),
	MethodNotAllowed: NewStatusCodeErr(
		"method not allowed",
		httpcodes.StatusMethodNotAllowed,
	),
	RouteRequiresAdminAuth: NewStatusCodeErr(
		"this route requires admin authorization",
		httpcodes.StatusUnauthorized,
	),
	InvalidAuthStrategy: NewStatusCodeErr(
		"invalid authorization strategy",
		httpcodes.StatusUnauthorized,
	),
	InvalidFMTQueryParam: NewStatusCodeErr(
		"the 'fmt' query param mus be 'old' or 'new', no other string is accepted",
		httpcodes.StatusBadRequest,
	),
	InvalidNIDQueryParam: NewStatusCodeErr(
		"the nid query param is required and it must be a valid unsigned integer",
		httpcodes.StatusBadRequest,
	),
}
