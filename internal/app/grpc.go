package app

import (
	"context"
	"currency-converter/internal/model"
	"currency-converter/internal/service"

	"currency-converter/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CurrencyServer struct {
	proto.UnimplementedCurrencyServiceServer
	svc service.Service
}

func NewCurrencyServer(svc service.Service) *CurrencyServer {
	return &CurrencyServer{svc: svc}
}

func (s *CurrencyServer) CreateCurrency(ctx context.Context, req *proto.CreateCurrencyRequest) (*proto.Currency, error) {
	if req.Currency == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Currency object is required")
	}
	
	switch {
	case req.Currency.Code == "":
		return nil, status.Errorf(codes.InvalidArgument, "Currency code cannot be empty")
	case req.Currency.Rate <= 0:
		return nil, status.Errorf(codes.InvalidArgument, "Exchange rate must be greater than zero")
	case req.Currency.Name == "":
		return nil, status.Errorf(codes.InvalidArgument, "Currency name cannot be empty")
	case req.Currency.Symbol == "":
		return nil, status.Errorf(codes.InvalidArgument, "Currency symbol cannot be empty")
	}

	cur := &model.Currency{
		Code:   req.Currency.Code,
		Rate:   req.Currency.Rate,
		Name:   req.Currency.Name,
		Symbol: req.Currency.Symbol,
	}

	created, err := s.svc.CreateCurrency(cur)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create currency: %v", err)
	}

	return &proto.Currency{
		Code:   created.Code,
		Rate:   created.Rate,
		Name:   created.Name,
		Symbol: created.Symbol,
	}, nil
}

func (s *CurrencyServer) ListCurrencies(ctx context.Context, _ *emptypb.Empty) (*proto.ListCurrenciesResponse, error) {
	data, err := s.svc.ListCurrencies()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to retrieve currency list: %v", err)
	}
	
	result := make([]*proto.Currency, 0, len(data))
	for _, v := range data {
		result = append(result, &proto.Currency{
			Code:   v.Code,
			Rate:   v.Rate,
			Name:   v.Name,
			Symbol: v.Symbol,
		})
	}
	return &proto.ListCurrenciesResponse{Currencies: result}, nil
}

func (s *CurrencyServer) GetCurrency(ctx context.Context, req *proto.Currency) (*proto.Currency, error) {
	if req.Code == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Currency code is required")
	}
	
	data, err := s.svc.GetCurrency(req.Code)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Currency '%s' not found in the system", req.Code)
	}
	
	return &proto.Currency{
		Code:   data.Code,
		Rate:   data.Rate,
		Name:   data.Name,
		Symbol: data.Symbol,
	}, nil
}

func (s *CurrencyServer) UpdateCurrency(ctx context.Context, req *proto.Currency) (*proto.Currency, error) {
	switch {
	case req == nil:
		return nil, status.Errorf(codes.InvalidArgument, "Request body cannot be empty")
	case req.Code == "":
		return nil, status.Errorf(codes.InvalidArgument, "Currency code is required for update")
	case req.Rate <= 0:
		return nil, status.Errorf(codes.InvalidArgument, "Exchange rate must be a positive value")
	case req.Name == "":
		return nil, status.Errorf(codes.InvalidArgument, "Currency name cannot be empty")
	case req.Symbol == "":
		return nil, status.Errorf(codes.InvalidArgument, "Currency symbol cannot be empty")
	}
	
	cur := &model.Currency{
		Code:   req.Code,
		Rate:   req.Rate,
		Name:   req.Name,
		Symbol: req.Symbol,
	}
	
	updated, err := s.svc.UpdateCurrency(cur)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update currency: %v", err)
	}
	
	return &proto.Currency{
		Code:   updated.Code,
		Rate:   updated.Rate,
		Name:   updated.Name,
		Symbol: updated.Symbol,
	}, nil
}

func (s *CurrencyServer) DeleteCurrency(ctx context.Context, req *proto.Currency) (*emptypb.Empty, error) {
	if req.Code == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Currency code is required for deletion")
	}
	
	if err := s.svc.DeleteCurrency(req.Code); err != nil {
		return nil, status.Errorf(codes.NotFound, "Currency '%s' not found - cannot delete", req.Code)
	}
	
	return &emptypb.Empty{}, nil
}

// *********************************Conversions*****************************************

type ConversionServer struct {
	proto.UnimplementedConversionServiceServer
	svc service.Service
}

func NewConversionServer(svc service.Service) *ConversionServer {
	return &ConversionServer{svc: svc}
}

func (s *ConversionServer) ListConversions(ctx context.Context, _ *emptypb.Empty) (*proto.ListConversionsResponse, error) {
	data, err := s.svc.ListConversions()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to retrieve conversion history: %v", err)
	}

	result := make([]*proto.Conversion, 0, len(data))
	for _, v := range data {
		result = append(result, &proto.Conversion{
			Amount: v.Amount,
			From: &proto.Currency{
				Code:   v.From.Code,
				Rate:   v.From.Rate,
				Name:   v.From.Name,
				Symbol: v.From.Symbol,
			},
			To: &proto.Currency{
				Code:   v.To.Code,
				Rate:   v.To.Rate,
				Name:   v.To.Name,
				Symbol: v.To.Symbol,
			},
			Result: v.Result,
		})
	}
	return &proto.ListConversionsResponse{Conversions: result}, nil
}

func (s *ConversionServer) CreateConversion(ctx context.Context, req *proto.CreateConversionRequest) (*proto.Conversion, error) {
	if req.Amount <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Conversion amount must be greater than zero")
	}
	if req.From == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Source currency code is required")
	}
	if req.To == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Target currency code is required")
	}

	conv, err := s.svc.CreateConversion(req.Amount, req.From, req.To)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Conversion failed: %v", err)
	}

	return &proto.Conversion{
		Amount: req.Amount,
		From: &proto.Currency{
			Code:   conv.From.Code,
			Rate:   conv.From.Rate,
			Name:   conv.From.Name,
			Symbol: conv.From.Symbol,
		},
		To: &proto.Currency{
			Code:   conv.To.Code,
			Rate:   conv.To.Rate,
			Name:   conv.To.Name,
			Symbol: conv.To.Symbol,
		},
		Result: conv.Result,
	}, nil
}