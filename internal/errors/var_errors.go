package errors

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
	ErrInvalidFMTQueryParam = New(
		"the 'fmt' query param mus be 'old' or 'new', no other string is accepted",
	)
	ErrInvalidNIDQueryParam = New(
		"the 'nid' query param is required and it must be a valid unsigned integer",
	)
)

var mpe = map[error]StatusError{
	ErrEntityNotFound: &statusErrorImpl{
		code:     -0,
		httpCode: -0,
		message:  "Não foi possível encontrar a entidade",
	},

	ErrInvalidObjectID: &statusErrorImpl{
		code:     -0,
		httpCode: -0,
		message:  "O object id do request não tem um encoding hexadecimal válido",
	},

	ErrUserUnauthorized: &statusErrorImpl{
		code:     -0,
		httpCode: -0,
		message:  "Usuário não existe ou a senha não é a correta",
	},

	ErrInvalidJwtTokenFormat: &statusErrorImpl{
		code:     -0,
		httpCode: -0,
		message:  "O token de autenticação é inválido, faça login novamente",
	},

	ErrJwtTokenExpired: &statusErrorImpl{
		code:     -0,
		httpCode: -0,
		message:  "O token de autenticação está expirado, faça login novamente",
	},

	ErrInvalidAuthSignature: &statusErrorImpl{
		code:     -0,
		httpCode: -0,
		message:  "A assinatura de autenticação apresentada não é válida, faça login novamente",
	},

	ErrInvalidAuthSignatureEncodingMethod: &statusErrorImpl{
		code:     -0,
		httpCode: -0,
		message:  "A assinatura de autenticação apresentada não é válida, faça login novamente",
	},

	ErrMalformedOrTooBigBody: &statusErrorImpl{
		code:     -0,
		httpCode: -0,
		message:  "O corpo da requisição está incompleto ou é muito grande, verifique sua conexão com a internet",
	},

	ErrInternalServerError: &statusErrorImpl{
		code:     -0,
		httpCode: -0,
		message:  "Algo deu errado enquando processavamos sua requisição, tente novamente mais tarde",
	},

	ErrInvalidRequestEntity: &statusErrorImpl{
		code:     -0,
		httpCode: -0,
		message:  "A entidade no corpo da requisição é inválida",
	},

	ErrEntityAlreadyExists: &statusErrorImpl{
		code:     -0,
		httpCode: -0,
		message:  "Não foi possível criar pois a entidade já existe",
	},

	ErrServerOperationTookTooLong: &statusErrorImpl{
		code:     -0,
		httpCode: -0,
		message:  "A operação demorou muito para ser concluída e foi cancelada, tente novamente mais tarde",
	},

	ErrSurnameSearchInvalid: &statusErrorImpl{
		code:     -0,
		httpCode: -0,
		message:  "Pelo menos um sobrenome precisa ser fornecido",
	},

	ErrInvalidEmail: &statusErrorImpl{
		code:     -0,
		httpCode: -0,
		message:  "O email apresentado é inválido",
	},

	ErrNoSuchRoute: &statusErrorImpl{
		code:     -0,
		httpCode: -0,
		message:  "Essa rota não existe na api",
	},

	ErrMethodNotAllowed: &statusErrorImpl{
		code:     -0,
		httpCode: -0,
		message:  "Requisições com esse método não são aceitas nessa rota",
	},

	ErrRouteRequiresAdminAuth: &statusErrorImpl{
		code:     -0,
		httpCode: -0,
		message:  "A rota precisa de autorização administrativa",
	},

	ErrInvalidAuthStrategy: &statusErrorImpl{
		code:     -0,
		httpCode: -0,
		message:  "A estratégia do cabeçalho de autenticação é inválida",
	},

	ErrInvalidFMTQueryParam: &statusErrorImpl{
		code:     -0,
		httpCode: -0,
		message:  "O parâmetro query-string 'fmt' pode ser apenas 'old' ou 'new'",
	},

	ErrInvalidNIDQueryParam: &statusErrorImpl{
		code:     -0,
		httpCode: -0,
		message:  "O parâmetro query-string 'nid' precisa ser um inteiro sem sinal (>0) válido",
	},
}
