package handler

import (
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/internal/utils"
	"github.com/games4l/internal/utils/apigwt"
	"github.com/games4l/pkg/errors"
)

func (s *Server) HandleGetMany(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	limit, _ := strconv.Atoi(req.PathParameters["l"])

	qs, err := s.h.GetMany(limit)
	if err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      qs.StatusCode,
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
		Body:            utils.MarshalJSON(qs),
	}, nil
}

func (s *Server) HandleGetOne(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	id, err := strconv.Atoi(req.PathParameters["id"])
	if err != nil {
		return nil, errors.ErrInvalidIntegerIdPathParam
	}

	q, err := s.h.GetOne(id)
	if err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      q.StatusCode,
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
		Body:            utils.MarshalJSON(q),
	}, nil
}

func (s *Server) HandlePostOne(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	err := s.ap.AuthenticateAdminHeader(req.Headers["authorization"], utils.S2B(req.Body))
	if err != nil {
		return nil, err
	}

	mp, err := apigwt.ParseMultipartForm(&req, 1<<24) // 16 MB
	if err != nil {
		return nil, err
	}

	q, err := s.h.PostOne(mp)
	if err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      q.StatusCode,
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
		Body:            utils.MarshalJSON(q),
	}, nil
}

func (s *Server) HandleUpdateOne(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	err := s.ap.AuthenticateAdminHeader(req.Headers["authorization"], utils.S2B(req.Body))
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(req.PathParameters["id"])
	if err != nil {
		return nil, errors.ErrInvalidIntegerIdPathParam
	}

	q, err := s.h.UpdateOne(id, utils.S2B(req.Body))
	if err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      q.StatusCode,
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
		Body:            utils.MarshalJSON(q),
	}, nil
}
