package service

import (
	"context"
	"fmt"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/seyuta/ecommerce-microservices-example/constant"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var insecuredPaths = [...]string{
	"/pb.AuthApi/Login",
	"/pb.AuthApi/Register",
	"/pb.ProductApi/GetProductByID",
	"/pb.ProductApi/CreateProduct",
	"/pb.ProductApi/ListProduct",
}

func insecuredPath(path string) bool {
	for _, p := range insecuredPaths {
		if p == path {
			return true
		}
	}
	return false
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	md, _ := metadata.FromIncomingContext(ctx)
	ua, _ := md["user-agent"]
	peer, _ := peer.FromContext(ctx)
	token, tokenFound := md["authorization"]
	var tkn string
	if len(token) <= 0 {
		tkn = "n/a"
	} else {
		tkn = token[0]
		tkn = fmt.Sprintf("%s...%s", tkn[0:17], tkn[len(tkn)-10:])
	}
	insecured := insecuredPath(info.FullMethod)
	var secured = "Y"
	if insecured {
		secured = "N"
	}
	log.Debugf("rpc sec: %s, mtd: %s, tkn: %v, ip: %s, ua: %v",
		secured, info.FullMethod, tkn, peer.Addr.String(), ua)
	// allowed method without authorization
	if insecured {
		return handler(ctx, req)
	}
	// bearer token not found
	if !tokenFound {
		log.Errorf("unauthenticated: no authorization metadata found")
		return nil, status.Errorf(codes.Unauthenticated, "no authorization metadata found")
	}

	tokenStr := strings.Replace(strings.Join(token, ""), "Bearer ", "", -1)
	jwtToken, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Signing method invalid")
		} else if method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("Signing method invalid")
		}

		return jwt.SigningMethodHS256, nil
	})

	claims, tokenFound := jwtToken.Claims.(jwt.MapClaims)
	ctx = context.WithValue(ctx, constant.CtxKeyUserID, claims[constant.CtxKeyUserID])
	ctx = context.WithValue(ctx, constant.CtxKeyDeviceID, claims[constant.CtxKeyDeviceID])
	ctx = context.WithValue(ctx, constant.CtxKeyUsername, claims[constant.CtxKeyUsername])
	ctx = context.WithValue(ctx, constant.CtxKeyUserEmail, claims[constant.CtxKeyUserEmail])
	ctx = context.WithValue(ctx, constant.CtxKeyUserMobile, claims[constant.CtxKeyUserMobile])

	return handler(ctx, req)
}

// JwtClaims is wrapper for claiming JWT during login process
type JwtClaims struct {
	jwt.StandardClaims
	UserID   string `json:"userid"`
	DeviceID string `json:"dvcid"`
	UserNo   int    `json:"userno"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}
