package utilities

import (
	"bambamload/constant"
	"bambamload/enum"
	"bambamload/logger"
	"bambamload/middleware"
	"bambamload/models"
	"bambamload/service/email"
	"bambamload/service/postgresrepository"
	"bambamload/service/redisService"
	"bambamload/service/uploadService"
	"bambamload/types"
	"bambamload/utils"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type ServiceUtilities struct {
	RedisService       redisService.RedisService
	PostgresRepository *postgresrepository.PostgresRepository
	EmailService       email.Email
	UploadService      *uploadService.UploadService
}

func NewServiceUtilities(redisService redisService.RedisService, postgresRepository *postgresrepository.PostgresRepository, emailService email.Email, uploadService *uploadService.UploadService) *ServiceUtilities {
	return &ServiceUtilities{
		RedisService:       redisService,
		PostgresRepository: postgresRepository,
		EmailService:       emailService,
		UploadService:      uploadService,
	}
}

func (su ServiceUtilities) Login(req models.LoginRequest) (any, error) {
	ctx := context.Background()
	reqEmail := strings.ToLower(strings.TrimSpace(req.Email))
	attemptKey := fmt.Sprintf("%s:%s", constant.LoginAttempts, reqEmail)
	lockoutKey := fmt.Sprintf("%s:%s", constant.LoginLockouts, reqEmail)

	// Check if the user is currently locked out
	if _, err := su.RedisService.GetRedisClient().Get(ctx, lockoutKey).Result(); err == nil {
		logger.Logger.Errorf("lockout error: %v", err)
		return nil, errors.New("account locked due to too many failed login attempts. Try again in 30 minutes")
	}

	user, err := su.PostgresRepository.GetUser(reqEmail, constant.Email)
	if err != nil {
		logger.Logger.Errorf("[Login] postgres get error: %v", err)
		return nil, errors.New("login failed, try again later")
	}

	if !utils.ComparePassword(req.Password, user.Password) {
		// Increment failed attempts
		attempts, _ := su.RedisService.GetRedisClient().Incr(ctx, attemptKey).Result()

		// Set expiry
		if attempts == 1 {
			su.RedisService.GetRedisClient().Expire(ctx, attemptKey, 30*time.Minute)
		}

		// Lock account if attempts >= 5
		if attempts >= 5 {
			su.RedisService.GetRedisClient().Set(ctx, lockoutKey, "locked", 30*time.Minute)
			su.RedisService.GetRedisClient().Del(ctx, attemptKey) // reset attempts
			return nil, errors.New("account locked due to too many failed login attempts. Try again in 30 hours")
		}

		return nil, fmt.Errorf("invalid password. %v of 5 attempts used", attempts)
	}

	// Successful login — reset failed attempts
	su.RedisService.GetRedisClient().Del(ctx, attemptKey)

	// Generate session token
	token, err := middleware.GenerateSessionToken()
	if err != nil {
		logger.Logger.Errorf("[LoginAdmin] generate session token failed: %v", err)
		return nil, errors.New("login failed, try again later")
	}

	atExp := time.Now().Add(time.Hour * 12)
	err = su.RedisService.SetSession(types.RedisSessionInfo{
		Token:  token,
		Expiry: atExp,
		Owner:  user.Role,
		ID:     user.ID,
	})
	if err != nil {
		logger.Logger.Errorf("[LoginAdmin] redis set session failed: %v", err)
		return nil, errors.New("login failed, try again later")
	}

	return models.SignInRes{
		AccessToken:       token,
		AccessTokenExpiry: atExp,
		User:              *user,
	}, nil
}

func (su ServiceUtilities) VerifyOtp(action, otp, email string, value any) error {

	email = strings.ToLower(strings.TrimSpace(email))

	switch action {

	case constant.LoginOtp:
		if os.Getenv(constant.AppEnv) == constant.Development && otp == constant.DefaultOtp {

		} else {
			key := fmt.Sprintf("%s:%s", constant.LoginOtp, email)
			var existingOtp string
			err := su.RedisService.GetValue(key, &existingOtp)
			if err != nil {
				if errors.Is(err, redis.Nil) { //nolint:typecheck
					logger.Logger.Errorf("[VerifyOtp]expired otp: %v", err)
					return errors.New("otp is expired")
				}
				logger.Logger.Errorf("[VerifyOtp]redis get error: %v", err)
				return err
			}

			if existingOtp != otp {
				return errors.New("invalid otp verification code")
			}
		}

		err := su.PostgresRepository.UpdateUser(email, constant.Email, map[string]interface{}{
			"last_login_time": time.Now().UTC(),
		})
		if err != nil {
			logger.Logger.Errorf("[VerifyOtp]postgres update error: %v", err)
			return err
		}

		return nil

	case constant.RegisterOtp:
		if os.Getenv(constant.AppEnv) == constant.Development && otp == constant.DefaultOtp {

		} else {
			key := fmt.Sprintf("%s:%s", constant.RegisterOtp, email)
			var existingOtp string
			err := su.RedisService.GetValue(key, &existingOtp)
			if err != nil {
				if errors.Is(err, redis.Nil) { //nolint:typecheck
					logger.Logger.Errorf("[VerifyOtp]expired otp: %v", err)
					return errors.New("otp is expired")
				}
				logger.Logger.Errorf("[VerifyOtp]redis get error: %v", err)
				return err
			}

			if existingOtp != otp {
				return errors.New("invalid otp verification code")
			}
		}

		user, err := su.PostgresRepository.GetUser(email, constant.Email)
		if err != nil {
			logger.Logger.Errorf("[VerifyOtp]postgres get error: %v", err)
			return errors.New("verification failed, try again later")
		}

		updateMap := make(map[string]interface{})
		if user.Role == enum.Supplier {
			updateMap["status"] = constant.Pending
		} else {
			updateMap["status"] = constant.Verified
		}

		err = su.PostgresRepository.UpdateUser(email, constant.Email, updateMap)
		if err != nil {
			logger.Logger.Errorf("[VerifyOtp]postgres update error: %v", err)
			return err
		}

		return nil

	case constant.ForgotPassword:

		newPassword, ok := value.(string)
		if !ok {
			return errors.New("invalid password")
		}
		hashedPassword, err := utils.HashPassword(newPassword)
		if err != nil {
			logger.Logger.Errorf("[VerifyOtp]hashing password failed: %v", err)
			return err
		}

		if os.Getenv(constant.AppEnv) == constant.Development && otp == constant.DefaultOtp {

		} else {
			key := fmt.Sprintf("%s:%s", constant.ForgotPassword, email)
			var existingOtp string
			err = su.RedisService.GetValue(key, &existingOtp)
			if err != nil {
				if errors.Is(err, redis.Nil) { //nolint:typecheck
					logger.Logger.Errorf("[VerifyOtp]expired otp: %v", err)
					return errors.New("otp is expired")
				}
				logger.Logger.Errorf("[VerifyOtp]redis get error: %v", err)
				return err
			}

			if existingOtp != otp {
				return errors.New("invalid otp verification code")
			}
		}

		err = su.PostgresRepository.UpdateUser(email, constant.Email, map[string]interface{}{
			"password": hashedPassword,
		})
		if err != nil {
			logger.Logger.Errorf("[VerifyOtp]postgres update error: %v", err)
			return err
		}

		return nil

	default:
		return errors.New("invalid action")
	}
}

func (su ServiceUtilities) SendOtp(action string, email string) error {
	otp, err := utils.GenerateOTP(6)
	if err != nil {
		logger.Logger.Errorf("[SendOtp]otp generation failed: %v", err)
		return err
	}

	email = strings.ToLower(strings.TrimSpace(email))

	switch action {

	case constant.LoginOtp:
		key := fmt.Sprintf("%s:%s", constant.LoginOtp, email)

		err = su.RedisService.SetValue(key, otp, 60)
		if err != nil {
			logger.Logger.Errorf("[LoginOtp]redis set failed: %v", err)
			return err
		}

		message := fmt.Sprintf(
			"Hello,\n\nYour login verification code is: %s\n\nThis code will expire in %d seconds.\nIf you didn’t request this, please ignore this message.",
			otp,
			30,
		)
		//_, _ = sa.SmsService.SendMessage(email, message)

		_ = su.EmailService.Send(email, "Login Verification Code", message)
		return nil

	case constant.RegisterOtp:
		key := fmt.Sprintf("%s:%s", constant.RegisterOtp, email)

		err = su.RedisService.SetValue(key, otp, 60)
		if err != nil {
			logger.Logger.Errorf("[SendOtp]redis set failed: %v", err)
			return err
		}

		message := fmt.Sprintf(
			"Hello,\n\nYour registration verification code is: %s\n\nThis code will expire in %d seconds.\nIf you didn’t request this, please ignore this message.",
			otp,
			30,
		)
		//_, _ = sa.SmsService.SendMessage(phoneNumber, message)
		_ = su.EmailService.Send(email, "Registration Verification Code", message)
		return nil

	case constant.ForgotPassword:
		key := fmt.Sprintf("%s:%s", constant.ForgotPassword, email)

		err = su.RedisService.SetValue(key, otp, 60)
		if err != nil {
			logger.Logger.Errorf("[ForgotPassword]redis set failed: %v", err)
			return err
		}

		message := fmt.Sprintf(
			"Hello,\n\nYour password reset code is: %s\n\nThis code will expire in %d seconds.\nIf you didn’t request this, please ignore this message.",
			otp,
			30,
		)
		//_, _ = sa.SmsService.SendMessage(phoneNumber, message)
		_ = su.EmailService.Send(email, "Forgot Password Verification Code", message)

		return nil
	default:
		return errors.New("invalid action")
	}
}

func (su ServiceUtilities) Logout(token string, user *models.User) error {
	updateMap := make(map[string]interface{})
	updateMap["last_login_time"] = time.Now().UTC()

	err := su.PostgresRepository.UpdateUser(user.ID, constant.ID, updateMap)
	if err != nil {
		logger.Logger.Errorf("[LogoutUser] postgres update error: %v", err)
		return err
	}

	return su.RedisService.DeleteSession(token)
}
