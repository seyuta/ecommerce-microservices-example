package service

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/seyuta/ecommerce-microservices-example/s-transactions/pkg/model"
	"github.com/seyuta/ecommerce-microservices-example/s-transactions/pkg/pb"
	"github.com/seyuta/ecommerce-microservices-example/s-transactions/repository"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OrderService
type OrderSvc struct {
	orepo repository.OrderRepository
	log   *logrus.Logger
}

// Instantiate new user Service
func NewOrderSvc(orepo repository.OrderRepository, log *logrus.Logger) *OrderSvc {
	return &OrderSvc{
		orepo: orepo,
		log:   log,
	}
}

// CreateOrder ...
func (s *OrderSvc) Create(ctx context.Context, ord *pb.OrderReqDto) (*pb.OrderDto, error) {
	var (
		now            = time.Now()
		ordDto         *pb.OrderDto
		ordDetailSlice []*pb.OrderDetailDto
		ordModelSlice  []model.OrderDetail
	)

	/*
		Generate invoice number
		format => INV/currentTime/randomNumber
		INV = invoice
		randomNumber from 1000 to 9999
	*/
	rand.Seed(time.Now().UnixNano())
	noInv := fmt.Sprintf("INV/%v/%v", now.Format("060201150405"), rand.Intn(9999-1000)+1000)

	om := &model.Order{
		NoInvoice: noInv,
		Status:    false,
	}

	for _, o := range ord.Order {
		od := model.OrderDetail{
			ProductID: o.ProductId,
			Qty:       o.Qty,
		}
		ordModelSlice = append(ordModelSlice, od)
	}
	om.OrderDetail = ordModelSlice

	ores, err := s.orepo.Create(ctx, om)
	if err != nil {
		s.log.Error(err)
		return ordDto, err
	}

	pid, _ := hex.DecodeString(ores.ID.Hex())
	ordDto = &pb.OrderDto{
		Id:          string(pid),
		UserId:      ores.UserID,
		NoInv:       ores.NoInvoice,
		Status:      ores.Status,
		OrderDetail: ordDetailSlice,
	}

	for _, o := range ores.OrderDetail {
		od := pb.OrderDetailDto{
			ProductId: o.ProductID,
			Qty:       o.Qty,
		}
		ordDetailSlice = append(ordDetailSlice, &od)
	}
	ordDto.OrderDetail = ordDetailSlice

	s.log.Debugf("created Order(ID: %s, No Invoice: %s)", ordDto.Id, ordDto.NoInv)

	return ordDto, nil
}

// GetOrderByID will return product detail
func (s *OrderSvc) GetOrderByID(ctx context.Context, id *wrappers.StringValue) (*pb.OrderDto, error) {

	pdetl, err := s.orepo.FindByID(id.GetValue())
	if err != nil {
		s.log.Error(err)
		return nil, status.Errorf(codes.NotFound, "Order not found.")
	}

	pid, _ := hex.DecodeString(pdetl.ID.Hex())
	response := &pb.OrderDto{
		Id:          string(pid),
		UserId:      pdetl.UserID,
		NoInv:       pdetl.NoInvoice,
		Status:      pdetl.Status,
		OrderDetail: []*pb.OrderDetailDto{},
	}

	return response, nil

}

// ListOrderByUserID will return all product with date choosen
func (s *OrderSvc) ListOrderByUserID(ctx context.Context, userId *wrappers.StringValue) (*pb.OrderDtoList, error) {
	var (
		plist   *pb.OrderDtoList
		ordDtos []*pb.OrderDto
		ordDto  *pb.OrderDto
	)

	rooms, err := s.orepo.FindAll()
	if err != nil {
		s.log.Error(err)
		return plist, err
	}

	for _, p := range rooms {
		pid, _ := hex.DecodeString(p.ID.Hex())
		ordDto = &pb.OrderDto{
			Id:          string(pid),
			UserId:      "",
			NoInv:       "",
			Status:      p.Status,
			OrderDetail: []*pb.OrderDetailDto{},
		}
		ordDtos = append(ordDtos, ordDto)
	}

	plist = new(pb.OrderDtoList)
	plist.List = ordDtos

	return plist, nil
}
