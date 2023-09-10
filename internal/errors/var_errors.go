package errors

import "github.com/games4l/internal/httpcodes"

var (
	ErrEntityNotFound = New(
		"the entity could not be found",
	)
	ErrInvalidObjectID = New(
		"the requested object id does not have a valid hex encoding",
	)
	ErrUserUnauthorized = New(
		"user does not exist or the password do not match",
	)
	ErrInvalidJwtTokenFormat = New(
		"invalid auth token metadata format, please login again",
	)
	ErrJwtTokenExpired = New(
		"the token is expired, please login again",
	)
	ErrInvalidAuthSignature = New(
		"the provided signature is not valid",
	)
	ErrInvalidAuthSignatureEncodingMethod = New(
		"the signature encoding method is not valid",
	)
	ErrMalformedOrTooBigBody = New(
		"the request body is partial, malformed or too big",
	)
	ErrInternalServerError = New(
		"something went wrong on our part while processing your request",
	)
	ErrInvalidRequestEntity = New(
		"the request body entity is invalid or could not be parsed",
	)
	ErrEntityAlreadyExists = New(
		"failed to create because this entity already exists",
	)
	ErrServerOperationTookTooLong = New(
		"the operation took too long and was canceled in the process",
	)
	ErrSurnameSearchInvalid = New(
		"at least one surname must be provided",
	)
	ErrInvalidEmail = New(
		"the provided email address is not valid",
	)
	ErrNoSuchRoute = New(
		"no such route",
	)
	ErrMethodNotAllowed = New(
		"method not allowed",
	)
	ErrRouteRequiresAdminAuth = New(
		"this route requires admin authorization",
	)
	ErrInvalidAuthStrategy = New(
		"invalid authorization strategy",
	)
	ErrInvalidIntegerIdPathParam = New(
		"the path param 'id' must be a valid integer",
	)
)

var mpe = map[error]StatusError{
	ErrEntityNotFound: &statusErrorImpl{
		code:     40401,
		httpCode: httpcodes.StatusNotFound,
		message:  "Não foi possível encontrar a entidade",
	},

	ErrInvalidObjectID: &statusErrorImpl{
		code:     40001,
		httpCode: httpcodes.StatusBadRequest,
		message:  "O object id do request não tem um encoding hexadecimal válido",
	},

	ErrUserUnauthorized: &statusErrorImpl{
		code:     40101,
		httpCode: httpcodes.StatusUnauthorized,
		message:  "Usuário não existe ou a senha não é a correta",
	},

	ErrInvalidJwtTokenFormat: &statusErrorImpl{
		code:     40102,
		httpCode: httpcodes.StatusUnauthorized,
		message:  "O token de autenticação é inválido, faça login novamente",
	},

	ErrJwtTokenExpired: &statusErrorImpl{
		code:     40103,
		httpCode: httpcodes.StatusUnauthorized,
		message:  "O token de autenticação está expirado, faça login novamente",
	},

	ErrInvalidAuthSignature: &statusErrorImpl{
		code:     40104,
		httpCode: httpcodes.StatusUnauthorized,
		message:  "A assinatura de autenticação apresentada não é válida, faça login novamente",
	},

	ErrInvalidAuthSignatureEncodingMethod: &statusErrorImpl{
		code:     40105,
		httpCode: httpcodes.StatusUnauthorized,
		message:  "A assinatura de autenticação apresentada não é válida, faça login novamente",
	},

	ErrMalformedOrTooBigBody: &statusErrorImpl{
		code:     40002,
		httpCode: httpcodes.StatusBadRequest,
		message:  "O corpo da requisição está incompleto ou é muito grande, verifique sua conexão com a internet",
	},

	ErrInternalServerError: &statusErrorImpl{
		code:     50000,
		httpCode: httpcodes.StatusInternalServerError,
		message:  "Algo deu errado enquando processavamos sua requisição, tente novamente mais tarde",
	},

	ErrInvalidRequestEntity: &statusErrorImpl{
		code:     40003,
		httpCode: httpcodes.StatusBadRequest,
		message:  "A entidade no corpo da requisição é inválida",
	},

	ErrEntityAlreadyExists: &statusErrorImpl{
		code:     40901,
		httpCode: httpcodes.StatusConflict,
		message:  "Não foi possível criar pois a entidade já existe",
	},

	ErrServerOperationTookTooLong: &statusErrorImpl{
		code:     50401,
		httpCode: httpcodes.StatusGatewayTimeout,
		message:  "A operação demorou muito para ser concluída e foi cancelada, tente novamente mais tarde",
	},

	ErrSurnameSearchInvalid: &statusErrorImpl{
		code:     40004,
		httpCode: httpcodes.StatusBadRequest,
		message:  "Pelo menos um sobrenome precisa ser fornecido",
	},

	ErrInvalidEmail: &statusErrorImpl{
		code:     40005,
		httpCode: httpcodes.StatusBadRequest,
		message:  "O email apresentado é inválido",
	},

	ErrNoSuchRoute: &statusErrorImpl{
		code:     40400,
		httpCode: httpcodes.StatusNotFound,
		message:  "Essa rota não existe na api",
	},

	ErrMethodNotAllowed: &statusErrorImpl{
		code:     40500,
		httpCode: httpcodes.StatusMethodNotAllowed,
		message:  "Requisições com esse método não são aceitas nessa rota",
	},

	ErrRouteRequiresAdminAuth: &statusErrorImpl{
		code:     40301,
		httpCode: httpcodes.StatusForbidden,
		message:  "A rota precisa de autorização administrativa",
	},

	ErrInvalidAuthStrategy: &statusErrorImpl{
		code:     40106,
		httpCode: httpcodes.StatusUnauthorized,
		message:  "A estratégia do cabeçalho de autenticação é inválida",
	},

	ErrInvalidIntegerIdPathParam: &statusErrorImpl{
		code:     40006,
		httpCode: httpcodes.StatusBadRequest,
		message:  "O parâmetro de diretório 'id' precisa ser um inteiro válido",
	},
}
