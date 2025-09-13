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
		return nil, status.Errorf(codes.InvalidArgument, "требутся указать код валюты.")
	case req.Currency.Rate <= 0:
		return nil, status.Errorf(codes.InvalidArgument, "курс должен быть > 0.")
	case req.Currency.Name == "":
		return nil, status.Errorf(codes.InvalidArgument, "название не должно быть пустым.")
	case req.Currency.Symbol == "":
		return nil, status.Errorf(codes.InvalidArgument, "укажите символ.")
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

func (s *CurrencyServer) ListCurrencies(ctx context.Context, _ *proto.ListCurrenciesRequest) (*proto.ListCurrenciesResponse, error) {
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

func (s *CurrencyServer) GetCurrency(ctx context.Context, req *proto.GetCurrencyRequest) (*proto.Currency, error) {
	if req.Code == "" {
		return nil, status.Errorf(codes.InvalidArgument, "требутся указать код валюты.")
	}

	data := repository.GetCurrencies()

	cur, ok := data[req.Code]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "валюта %s не найдена.", cur.Code)
	}

	return &proto.Currency{
		Code:   cur.Code,
		Rate:   cur.Rate,
		Name:   cur.Name,
		Symbol: cur.Symbol,
	}, nil
}

func (s *CurrencyServer) UpdateCurrency(ctx context.Context, req *proto.UpdateCurrencyRequest) (*proto.Currency, error) {
	switch {
	case req.Currency == nil:
		return nil, status.Errorf(codes.InvalidArgument, "пустое тело запроса для обновления.")
	case req.Currency.Code == "":
		return nil, status.Errorf(codes.InvalidArgument, "требутся указать код валюты.")
	case req.Currency.Rate <= 0:
		return nil, status.Errorf(codes.InvalidArgument, "курс должен быть > 0.")
	case req.Currency.Name == "":
		return nil, status.Errorf(codes.InvalidArgument, "название не должно быть пустым.")
	case req.Currency.Symbol == "":
		return nil, status.Errorf(codes.InvalidArgument, "укажите символ.")
	}

	curMap := repository.GetCurrencies()
	old, ok := curMap[req.Currency.Code]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "валюта %s не найдена.", req.Currency.Code)
	}

	old.Rate = req.Currency.Rate
	old.Name = req.Currency.Name
	old.Symbol = req.Currency.Symbol

	if err := repository.UpdateCurInMap(old); err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка обновления: %v", err)
	}
	return req.Currency, nil
}

func (s *CurrencyServer) DeleteCurrency(ctx context.Context, req *proto.DeleteCurrencyRequest) (*emptypb.Empty, error) {
	if err := repository.DeleteCurFromMap(req.Code); err != nil {
		return nil, status.Errorf(codes.NotFound, "валюта %s не найдена", req.Code)
	}
	return &emptypb.Empty{}, nil
}

// *********************************Конверсии***************************************** 

type ConversionServer struct {
	proto.UnimplementedConversionServiceServer
}

func (s *ConversionServer) ListConversions(ctx context.Context, _ *proto.ListConversionsRequest) (*proto.ListConversionsResponse, error) {
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
		return nil, status.Errorf(codes.InvalidArgument, "сумма должна быть > 0")
	} else if req.From == "" || req.To == "" {
		return nil, status.Errorf(codes.InvalidArgument, "требуется указать код валюты From и To")
	}

	curs := repository.GetCurrencies()
	from, ok1 := curs[req.From]
	to, ok2 := curs[req.To]
	if !ok1 || !ok2 {
		return nil, status.Error(codes.NotFound, "не найдена исходная или целевая валюта.")
	} else if from.Rate <= 0 || to.Rate <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "курс валют должен быть > 0")
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
