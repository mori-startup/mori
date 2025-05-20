package utils

import (
	"errors"
	"unicode"

	"github.com/dchest/captcha"

	"mori/pkg/models"
)

// ValidateNewUser checks captcha, user fields, and password rules.
func ValidateNewUser(user models.User, captchaID, captchaValue string) error {
	// 1) Validate Captcha
	if err := validateCaptcha(captchaID, captchaValue); err != nil {
		return err
	}

	// 2) Check empty fields
	if err := validateFirstName(user.FirstName); err != nil {
		return err
	}
	if err := validateLastName(user.LastName); err != nil {
		return err
	}
	if err := validateBirth(user.DateOfBirth); err != nil {
		return err
	}
	if err := validatePassword(user.Password); err != nil {
		return err
	}
	if err := validateEmail(user.Email); err != nil {
		return err
	}

	return nil
}

// ---------------------------------------------------------------------------------
// Helper validation methods
// ---------------------------------------------------------------------------------

func validateCaptcha(captchaID, captchaValue string) error {
	if captchaID == "" || captchaValue == "" {
		return errors.New("invalid captcha")
	}
	if !captcha.VerifyString(captchaID, captchaValue) {
		return errors.New("invalid captcha")
	}
	return nil
}

func validateFirstName(name string) error {
	if fieldEmpty(name) {
		return errors.New("first name is required")
	}
	return nil
}

func validateLastName(name string) error {
	if fieldEmpty(name) {
		return errors.New("last name is required")
	}
	return nil
}

func validateBirth(birthDate string) error {
	if fieldEmpty(birthDate) {
		return errors.New("birthday is required")
	}
	return nil
}

func validateEmail(email string) error {
	if fieldEmpty(email) {
		return errors.New("email is required")
	}
	// Optionally check format if you want
	return nil
}

func validatePassword(password string) error {
	if fieldEmpty(password) {
		return errors.New("password is required")
	}
	// 1) At least 10 characters
	if len(password) < 10 {
		return errors.New("password must be at least 10 characters")
	}
	// 2) Must contain at least one uppercase, one digit, one special char
	var hasUpper, hasDigit, hasSpecial bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case isSpecialChar(ch):
			hasSpecial = true
		}
	}
	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasDigit {
		return errors.New("password must contain at least one digit")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

// If you prefer the broad definition of "special" as anything not letter or digit:
func isSpecialChar(ch rune) bool {
	return !(unicode.IsLetter(ch) || unicode.IsDigit(ch))
}

func fieldEmpty(value string) bool {
	return len(value) == 0
}
