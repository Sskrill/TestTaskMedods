package service

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Sskrill/TestTaskMedods/internal/domain"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
	"math/rand"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

const (
	smtpHost = "smtp.mail.ru"
	smtpPort = "587"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByParams(ctx context.Context, email, password string) (domain.User, error)
	GetEmailById(ctx context.Context, id int64) (string, error)
}

type TokenRepository interface {
	Get(ctx context.Context, token string) (domain.Tokens, error)
	Create(ctx context.Context, token domain.Tokens) error
	DeleteToken(ctx context.Context, userId int64) error
	GetCoupleOfTokens(ctx context.Context, guid string) (aT, rT string, err error)
}
type Hasher interface {
	Hash(string) (string, error)
}
type UserService struct {
	hasher    Hasher
	tokenRepo TokenRepository
	userRepo  UserRepository
	secretKey []byte
}

func NewUserService(hasher Hasher, tokenRepository TokenRepository, userRepository UserRepository, key []byte) *UserService {
	return &UserService{tokenRepo: tokenRepository, hasher: hasher, userRepo: userRepository, secretKey: key}
}
func (uS *UserService) ParseToken(ctx context.Context, token string) (int64, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return uS.secretKey, nil
	})
	if err != nil {

		return 0, err
	}

	if !t.Valid {

		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {

		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["jti"].(string)
	if !ok {

		return 0, errors.New("invalid subject")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {

		return 0, errors.New("invalid subject")
	}

	return int64(id), nil
}
func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
func (uS *UserService) generateTokens(ctx context.Context, userId int64, ipAddr string) (string, string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   ipAddr,
		Id:        strconv.Itoa(int(userId)),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
	})

	accessToken, err := t.SignedString(uS.secretKey)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}
	err = uS.tokenRepo.DeleteToken(ctx, userId)
	if err != nil {
		return "", "", err
	}
	namespace := uuid.NameSpaceDNS
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(userId))
	id := uuid.NewSHA1(namespace, []byte(fmt.Sprintf("%d", userId)))
	if err := uS.tokenRepo.Create(ctx, domain.Tokens{
		UserId:    userId,
		RToken:    refreshToken,
		AToken:    accessToken,
		Uuid:      id.String(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
		IpAddr:    ipAddr,
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
func (uS *UserService) RefreshTokens(ctx context.Context, refreshToken, ipAddr string) (string, string, error) {
	refreshToken = strings.ReplaceAll(refreshToken, "'", "")
	rToken, err := uS.tokenRepo.Get(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}
	if rToken.IpAddr != ipAddr {
		email, err := uS.userRepo.GetEmailById(ctx, rToken.UserId)
		if err != nil {

			return "", "", err
		}
		err = sendEmailOfWarning(email)
		if err != nil {
			log.Println("Error to send warning email message:", err.Error())
		}
	}
	if rToken.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", domain.ErrRefreshTokenExpired
	}

	return uS.generateTokens(ctx, rToken.UserId, rToken.IpAddr)
}
func (uS *UserService) SignIn(ctx context.Context, inp domain.SignInInput, ipAddr string) (string, string, error) {
	password, err := uS.hasher.Hash(inp.Password)
	if err != nil {
		return "", "", err
	}

	user, err := uS.userRepo.GetByParams(ctx, inp.Email, password)
	if err != nil {
		return "", "", err
	}

	return uS.generateTokens(ctx, user.Id, ipAddr)
}
func (uS *UserService) SignUp(ctx context.Context, param domain.SignUpInput) error {

	password, err := uS.hasher.Hash(param.Password)
	if err != nil {

		return err
	}
	user := domain.User{Name: param.Name, Password: password, Email: param.Email}
	return uS.userRepo.Create(ctx, user)
}
func sendEmailOfWarning(toEmail string) error {
	cfg := ConfigParamEmail{}
	if err := godotenv.Load(); err != nil {
		return err
	}
	if err := envconfig.Process("eml", &cfg); err != nil {
		return err
	}

	to := []string{toEmail}

	message := []byte("Subject: Security Warning\r\n" +
		"\r\n" +
		"Warning: different IP detected.")

	auth := smtp.PlainAuth("", cfg.Email, cfg.Password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, cfg.Email, to, message)
	if err != nil {
		return err
	}
	log.Println("Send email successfully")
	return nil
}
func (uS *UserService) GetTokensByGUID(ctx context.Context, guid string) (aT, rT string, err error) {

	aT, rT, err = uS.tokenRepo.GetCoupleOfTokens(ctx, guid)
	if err != nil {
		return "", "", err
	}
	return aT, rT, nil
}
