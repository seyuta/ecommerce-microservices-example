package service

import (
	"context"
	"encoding/hex"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/seyuta/ecommerce-microservices-example/s-catalog/pkg/model"
	"github.com/seyuta/ecommerce-microservices-example/s-catalog/pkg/pb"
	"github.com/seyuta/ecommerce-microservices-example/s-catalog/repository"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProductService
type ProductSvc struct {
	prepo repository.ProductRepository
	log   *logrus.Logger
}

// Instantiate new user Service
func NewProductSvc(prepo repository.ProductRepository, log *logrus.Logger) *ProductSvc {
	return &ProductSvc{
		prepo: prepo,
		log:   log,
	}
}

// CreateProduct ...
func (s *ProductSvc) Create(ctx context.Context, p *pb.ProductDto) (*pb.ProductDto, error) {
	var (
		pname = p.Name
		price = p.Price
		stock = p.Stock
	)

	if pname == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty product name")
	}
	if price <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty product price")
	}
	if stock <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty product stock")
	}

	var pdto *pb.ProductDto
	rm := &model.Product{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Status:      p.Status,
	}

	pres, err := s.prepo.Create(ctx, rm)
	if err != nil {
		s.log.Error(err)
		return pdto, err
	}

	pid, _ := hex.DecodeString(pres.ID.Hex())
	pdto = &pb.ProductDto{
		Id:          string(pid),
		Name:        pres.Name,
		Description: pres.Description,
		Price:       pres.Price,
		Stock:       pres.Stock,
		Status:      pres.Status,
	}

	s.log.Debugf("created Product(ID: %s, Name: %s)", pdto.Id, pdto.Name)

	return pdto, nil

}

// GetProductByID will return product detail
func (s *ProductSvc) GetProductByID(ctx context.Context, id *wrappers.StringValue) (*pb.ProductDto, error) {

	pdetl, err := s.prepo.FindByID(id.GetValue())
	if err != nil {
		s.log.Error(err)
		return nil, status.Errorf(codes.NotFound, "Product not found.")
	}

	pid, _ := hex.DecodeString(pdetl.ID.Hex())
	response := &pb.ProductDto{
		Id:          string(pid),
		Name:        pdetl.Name,
		Description: pdetl.Description,
		Price:       pdetl.Price,
		Stock:       pdetl.Stock,
		Status:      pdetl.Status,
	}

	return response, nil

}

//ListProducts will return all product
func (s *ProductSvc) ListProduct(ctx context.Context, e *empty.Empty) (*pb.ProductDtoList, error) {
	var (
		plist *pb.ProductDtoList
		pdtos []*pb.ProductDto
		pdto  *pb.ProductDto
	)

	product, err := s.prepo.FindAll()
	if err != nil {
		s.log.Error(err)
		return plist, err
	}

	for _, p := range product {
		pid, _ := hex.DecodeString(p.ID.Hex())
		pdto = &pb.ProductDto{
			Id:          string(pid),
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
			Status:      p.Status,
		}
		pdtos = append(pdtos, pdto)
	}

	plist = new(pb.ProductDtoList)
	plist.List = pdtos

	return plist, nil
}
