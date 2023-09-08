package errors

func strp(s string) *string {
	return &s
}

var PtBtMessages = MessageList{
	EntityNotFound: strp(
		"não foi possível encontrar a entidade",
	),
	InvalidObjectID: strp(
		"o object id do request não tem um encoding hexadecimal válido",
	),
	UserUnauthorized: strp(
		"usuário não existe ou a senha não é a correta",
	),
	InvalidJwtTokenFormat: strp(
		"o token de autenticação é inválido, faça login novamente",
	),
	JwtTokenExpired: strp(
		"o token de autenticação está expirado, faça login novamente",
	),
	InvalidAuthSignature: strp(
		"a assinatura de autenticação apresentada não é válida, faça login novamente",
	),
	InvalidAuthSignatureEncodingMethod: strp(
		"a assinatura de autenticação apresentada não é válida, faça login novamente",
	),
	MalformedOrTooBigBody: strp(
		"o corpo da requisição está incompleto ou é muito grande, verifique sua conexão com a internet",
	),
	InternalServerError: strp(
		"alguma coisa deu errado enquando processavamos sua requisição, tente novamente mais tarde",
	),
	InvalidRequestEntity: strp(
		"a entidade no corpo da requisição é inválida",
	),
	EntityAlreadyExists: strp(
		"não foi possível criar pois a entidade já existe",
	),
	ServerOperationTookTooLong: strp(
		"a operação demorou muito para ser concluída e foi cancelada, tente novamente mais tarde",
	),
	SurnameSearchInvalid: strp(
		"pelo menos um sobrenome precisa ser fornecido",
	),
	InvalidEmail: strp(
		"o email apresentado é inválido",
	),
	NoSuchRoute: strp(
		"essa rota não existe na api",
	),
	MethodNotAllowed: strp(
		"requisições com esse método não são aceitas nessa rota",
	),
	RouteRequiresAdminAuth: strp(
		"a rota precisa de autorização administrativa",
	),
	InvalidAuthStrategy: strp(
		"a estratégia do cabeçalho de autenticação é inválida",
	),
	InvalidFMTQueryParam: strp(
		"o parametro query-string 'fmt' pode ser apenas 'old' ou 'new'",
	),
	InvalidNIDQueryParam: strp(
		"o parametro query-string 'nid' precisa ser um inteiro sem sinal (>0) válido",
	),
}
