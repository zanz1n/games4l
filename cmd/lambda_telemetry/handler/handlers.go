package handler

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/internal/utils"
)

func (s *Server) HandleGetByName(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	err := s.ap.AuthenticateAdminHeader(req.Headers["authorization"], utils.S2B(req.Body))
	if err != nil {
		return nil, err
	}

	t, err := s.h.GetByName(req.PathParameters["name"])
	if err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      t.StatusCode,
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
		Body:            utils.MarshalJSON(t),
	}, nil
}

func (s *Server) HandleGetById(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	err := s.ap.AuthenticateAdminHeader(req.Headers["authorization"], utils.S2B(req.Body))

	t, err := s.h.GetById(req.PathParameters["id"], err == nil)
	if err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      t.StatusCode,
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
		Body:            utils.MarshalJSON(t),
	}, nil
}

func (s *Server) HandlePostOne(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	t, err := s.h.PostOne(utils.S2B(req.Body))
	if err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      t.StatusCode,
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
		Body:            utils.MarshalJSON(t),
	}, nil
}
