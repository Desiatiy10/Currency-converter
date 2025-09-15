package grpc_server

import (
	"context"
	"currency-converter/internal/model"
	"currency-converter/proto"
	"currency-converter/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CurrencyServer struct {
	proto.UnimplementedCurrencyServiceServer
}

func (s *CurrencyServer) CreateCurrency(ctx context.Context, req *proto.CreateCurrencyRequest) (*proto.Currency, error) {
	switch {
	case req.Currency.Code == "":
		return nil, status.Errorf(codes.InvalidArgument, "enter the currency code.")
	case req.Currency.Rate <= 0:
		return nil, status.Errorf(codes.InvalidArgument, "the rate should be > 0.")
	case req.Currency.Name == "":
		return nil, status.Errorf(codes.InvalidArgument, "enter the currency name.")
	case req.Currency.Symbol == "":
		return nil, status.Errorf(codes.InvalidArgument, "enter the currency symbol.")
	}

	cur := &model.Currency{
		Code:   req.Currency.Code,
		Rate:   req.Currency.Rate,
		Name:   req.Currency.Name,
		Symbol: req.Currency.Symbol,
	}

	repository.Store(cur)

	return &proto.Currency{
		Code:   cur.Code,
		Rate:   cur.Rate,
		Name:   cur.Name,
		Symbol: cur.Symbol,
	}, nil
}

func (s *CurrencyServer) ListCurrencies(ctx context.Context, _ *emptypb.Empty) (*proto.ListCurrenciesResponse, error) {
	data := repository.GetCurrencies()
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
		return nil, status.Errorf(codes.InvalidArgument, "enter the currency code.")
	}

	data := repository.GetCurrencies()

	cur, ok := data[req.Code]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "the currency %s not found.", cur.Code)
	}

	return &proto.Currency{
		Code:   cur.Code,
		Rate:   cur.Rate,
		Name:   cur.Name,
		Symbol: cur.Symbol,
	}, nil
}

func (s *CurrencyServer) UpdateCurrency(ctx context.Context, req *proto.Currency) (*proto.Currency, error) {
	switch {
	case req == nil:
		return nil, status.Errorf(codes.InvalidArgument, "empty request body.")
	case req.Code == "":
		return nil, status.Errorf(codes.InvalidArgument, "enter the currency code.")
	case req.Rate <= 0:
		return nil, status.Errorf(codes.InvalidArgument, "the rate should be > 0.")
	case req.Name == "":
		return nil, status.Errorf(codes.InvalidArgument, "enter the currency name.")
	case req.Symbol == "":
		return nil, status.Errorf(codes.InvalidArgument, "enter the currency symbol.")
	}

	curMap := repository.GetCurrencies()
	old, ok := curMap[req.Code]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "the currency %s not found.", req.Code)
	}

	old.Rate = req.Rate
	old.Name = req.Name
	old.Symbol = req.Symbol

	if err := repository.UpdateCurInMap(old); err != nil {
		return nil, status.Errorf(codes.Internal, "update error : %v", err)
	}
	return req, nil
}

func (s *CurrencyServer) DeleteCurrency(ctx context.Context, req *proto.Currency) (*emptypb.Empty, error) {
	if err := repository.DeleteCurFromMap(req.Code); err != nil {
		return nil, status.Errorf(codes.NotFound, "the currency %s not found", req.Code)
	}
	return &emptypb.Empty{}, nil
}

// *********************************Конверсии*****************************************

type ConversionServer struct {
	proto.UnimplementedConversionServiceServer
}

func (s *ConversionServer) ListConversions(ctx context.Context, _ *emptypb.Empty) (*proto.ListConversionsResponse, error) {
	data := repository.GetConversions()
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
		return nil, status.Errorf(codes.InvalidArgument, "the rate should be > 0")
	} else if req.From == "" || req.To == "" {
		return nil, status.Errorf(codes.InvalidArgument, "need to enter the currency code of the source and target")
	}

	curs := repository.GetCurrencies()
	from, ok1 := curs[req.From]
	to, ok2 := curs[req.To]
	if !ok1 || !ok2 {
		return nil, status.Error(codes.NotFound, "the source or target currency was not found..")
	} else if from.Rate <= 0 || to.Rate <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "the rate should be > 0")
	}

	result := req.Amount * (to.Rate / from.Rate)
	conversion := model.NewConversion(req.Amount,
		from, to, result)
	repository.Store(conversion)

	return &proto.Conversion{
		Amount: req.Amount,
		From: &proto.Currency{
			Code:   from.Code,
			Rate:   from.Rate,
			Name:   from.Name,
			Symbol: from.Symbol,
		},
		To: &proto.Currency{
			Code:   to.Code,
			Rate:   to.Rate,
			Name:   to.Name,
			Symbol: to.Symbol,
		},
		Result: result,
	}, nil
}
