package service

import (
	"context"
	"encoding/hex"
	"regexp"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/seyuta/ecommerce-microservices-example/constant"
	"github.com/seyuta/ecommerce-microservices-example/s-auth/pkg/model"
	"github.com/seyuta/ecommerce-microservices-example/s-auth/pkg/pb"
	"github.com/seyuta/ecommerce-microservices-example/s-auth/repository"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthSvc Service
type AuthSvc struct {
	urepo repository.UserRepository
	log   *logrus.Logger
}

// NewAuthSvc Instantiate new auth Service
func NewAuthSvc(urepo repository.UserRepository, log *logrus.Logger) *AuthSvc {
	return &AuthSvc{
		urepo: urepo,
		log:   log,
	}
}

// Register ...
func (s *AuthSvc) Register(ctx context.Context, udto *pb.RegisterDto) (*pb.UserAuthResDto, error) {
	var (
		reqUname = strings.ToLower(udto.GetUsername())
		reqEmail = strings.ToLower(udto.GetEmail())
	)

	hasPwd, err := hashPassword(udto.GetPassword())
	if err != nil {
		s.log.Errorln(err)
		return nil, status.Errorf(codes.DataLoss, "password is not hashed")
	}

	// blocker username
	checkUser, _ := s.urepo.FindByUsername(ctx, reqUname)
	if checkUser.Username == udto.Username {
		s.log.Errorf("username exist: %s, err: %v", reqUname, err)
		return nil, status.Errorf(codes.NotFound, "username exist: %s", reqUname)
	}

	// blocker email
	checkEmail, _ := s.urepo.FindByEmail(ctx, reqEmail)
	if checkEmail.Email == udto.Email {
		s.log.Errorf("email exist: %s, err: %v", reqEmail, err)
		return nil, status.Errorf(codes.NotFound, "email exist: %s", reqEmail)
	}

	u := &model.UserAuth{
		Username: reqUname,
		Email:    reqEmail,
		Phone:    udto.GetPhone(),
		Password: hasPwd,
		Status:   model.UserStatusActive,
	}

	uauth, err := s.urepo.Create(ctx, u)
	if err != nil {
		s.log.Errorln(err)
		return nil, status.Errorf(codes.DataLoss, "user is not created")
	}

	user := &pb.UserAuthResDto{
		Username: uauth.Username,
		Email:    uauth.Email,
		Phone:    uauth.Phone,
	}

	s.log.Infof("new user(%s) created", user.Username)
	return user, nil
}

// Login ...
func (s *AuthSvc) Login(ctx context.Context, login *pb.LoginDto) (*pb.AccessTokenDto, error) {
	uname := strings.ToLower(login.GetUsername())
	dvcid := login.GetDeviceId()
	s.log.Debugf("Login(Username : %s, DeviceId : %s)", uname, dvcid)

	if uname == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty username")
	}
	if !validUsername(uname) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid username")
	}
	if dvcid == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty deviceId")
	}

	// find username
	user, err := s.urepo.FindByUsername(ctx, uname)
	if err != nil {
		s.log.Errorf("not found username: %s, err: %v", uname, err)
		return nil, status.Errorf(codes.NotFound, "not found username: %s", uname)
	}
	// password is not match
	passOk := checkPasswordHash(login.GetPassword(), user.Password)
	if !passOk {
		s.log.Errorln("invalid password")
		return nil, status.Errorf(codes.InvalidArgument, "invalid password")
	}
	// user is not active
	if user.Status != model.UserStatusActive {
		s.log.Errorln("user is not active")
		return nil, status.Errorf(codes.FailedPrecondition, "user is not active")
	}

	// find existing token and return
	curToken, err := s.urepo.FindTokenByUserAndDevice(ctx, uname, dvcid)
	if err == nil {

		expIn := curToken.ExpiredDt.Unix() - time.Now().Unix()
		s.log.Debugf("current token will expired in %v seconds", expIn)
		if expIn >= 420 { //420 secs = 7 minutes
			accesToken := &pb.AccessTokenDto{
				AccessToken: curToken.Token,
				ExpiresIn:   expIn,
				Scope:       "all",
				TokenType:   "Bearer",
				Username:    user.Username,
				Email:       user.Email,
				Phone:       user.Phone,
			}
			return accesToken, nil
		} else {
			td, _ := s.urepo.DeleteTokenByUserAndDevice(ctx, uname, dvcid)
			if td {
				s.log.Debugf("expired token removed for %s on %s", uname, dvcid)
			}
		}

	}

	var issuer = viper.GetString(constant.CfgJwtIssuer)
	var issTime = time.Now().Unix()
	var expTime = issTime + viper.GetInt64(constant.CfgJwtExpiresSeconds)
	var signkey = []byte(viper.GetString(constant.CfgJwtISignKey))

	uid, _ := hex.DecodeString(user.ID.Hex())
	claims := &JwtClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.New().String(),
			Audience:  uname,
			Issuer:    issuer,
			IssuedAt:  issTime,
			ExpiresAt: expTime,
		},
		UserID:   string(uid),
		DeviceID: dvcid,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString(signkey)
	if err != nil {
		s.log.Errorf("unable to sign token: %v", err)
		return nil, status.Errorf(codes.Internal, "unable to sign token: %v", err)
	}

	logDt := time.Now()
	expDt := time.Unix(expTime, 0)
	userToken := &model.UserToken{
		UserID:    user.ID,
		Username:  user.Username,
		DeviceID:  dvcid,
		LoginDt:   &logDt,
		ExpiredDt: &expDt,
		Token:     signedToken,
	}

	utoken, err := s.urepo.CreateToken(ctx, userToken)
	if err != nil {
		s.log.Errorf("unable to save token: %v", err)
		return nil, status.Errorf(codes.Internal, "unable to save token: %v", err)
	}

	accesToken := &pb.AccessTokenDto{
		AccessToken: utoken.Token,
		ExpiresIn:   utoken.ExpiredDt.Unix() - time.Now().Unix(),
		Scope:       "all",
		TokenType:   "Bearer",
		Username:    user.Username,
		Email:       user.Email,
		Phone:       user.Phone,
	}

	return accesToken, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// username only any word character (letter, number, underscore)
func validUsername(uname string) bool {
	matched, _ := regexp.MatchString(`^\w*$`, uname)
	return matched
}
