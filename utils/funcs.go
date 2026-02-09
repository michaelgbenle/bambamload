package utils

import (
	"bambamload/constant"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ValidMSISDN checks if a phone number is a valid MSISDN format
func ValidMSISDN(number string) bool {

	number = strings.TrimSpace(number)

	//if !strings.HasPrefix(number, "+") {
	//	return false
	//}

	totalDigits := len(number)
	if totalDigits < 8 || totalDigits > 13 {
		return false
	}

	numRegex := regexp.MustCompile(`^[0-9]+$`)
	if !numRegex.MatchString(number[1:]) {
		return false
	}

	return true
}

// GetWATTime Get West African Time
func GetWATTime() time.Time {
	location, locationErr := time.LoadLocation(constant.AfricaLagos)
	if locationErr != nil {
		return time.Now()
	}
	return time.Now().In(location)
}

func GenerateReference(prefix string) string {

	randomBytes := make([]byte, 4)
	_, _ = rand.Read(randomBytes)

	if prefix == "" {
		return fmt.Sprintf("%s%s", GetWATTime().Format("20060102150405"), hex.EncodeToString(randomBytes))
	}

	return fmt.Sprintf("%s-%s%s", strings.ToUpper(prefix), GetWATTime().Format("20060102150405"), hex.EncodeToString(randomBytes))
}

func StandardiseMSISDN(msisdn string) string {
	msisdn = strings.TrimSpace(msisdn)

	switch {
	case strings.HasPrefix(msisdn, constant.NigeriaMSISDNPrefixPlus):
		// +2348012345678 → 2348012345678
		return msisdn[1:]

	case strings.HasPrefix(msisdn, constant.NigeriaMSISDNPrefix):
		// Already correct → 2348012345678
		return msisdn

	case strings.HasPrefix(msisdn, constant.Zero):
		// 08012345678 → 2348012345678
		return "234" + msisdn[1:]

	default:
		// Handle cases like 8012345678 → 2348012345678
		if len(msisdn) == 10 || len(msisdn) == 11 {
			return "234" + strings.TrimPrefix(msisdn, "0")
		}
		return msisdn
	}
}

// HashPassword hashes password and returns the hashed value
func HashPassword(p string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "unable to encrypt password", err
	}
	return string(hash), nil
}

// ComparePassword decrypts user password and compares the hash with the supplied password
func ComparePassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func PercentageChange(current, previous float64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100 // went from 0 to something
	}
	return ((current - previous) / previous) * 100
}

func DateStringToTime(format, date string) (time.Time, error) {
	timeObj, err := time.Parse(format, date)
	if err != nil {
		return time.Time{}, err
	}
	return timeObj, nil
}

// GenerateOTP generates OTP of length n
func GenerateOTP(n int) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("OTP length must be greater than 0")
	}

	otp := make([]byte, n)

	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(constant.Digits))))
		if err != nil {
			return "", fmt.Errorf("failed to generate OTP: %v", err)
		}
		otp[i] = constant.Digits[num.Int64()]
	}

	return string(otp), nil
}

func GetFileExtension(filename string) string {
	return filepath.Ext(filename)
}

func SplitFileName(filename string) (name string, ext string) {
	ext = filepath.Ext(filename)

	return strings.TrimSuffix(filename, ext), ext
}
