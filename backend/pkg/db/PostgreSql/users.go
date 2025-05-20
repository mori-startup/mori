package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"mori/pkg/models"
)

type UserRepository struct {
	DB *sql.DB
}

// Add inserts a new user into the database and sends a verification email.
func (repo *UserRepository) Add(user models.User) error {
	// 1) Generate a random verification token
	token, err := generateRandomToken()
	if err != nil {
		return fmt.Errorf("could not generate token: %w", err)
	}

	user.VerificationToken = token
	user.Verified = false

	// 2) Insert user with the verification_token and verified=false
	query := `
		INSERT INTO users (
			user_id, email, first_name, last_name, nickname, about, 
			password, birthday, image, verification_token, verified
		)
		VALUES(
			$1, $2, $3, $4, NULLIF($5, ''), $6,
			$7, $8, $9, $10, $11
		);
	`
	_, errDB := repo.DB.Exec(
		query,
		user.ID,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Nickname,
		user.About,
		user.Password,
		user.DateOfBirth,
		user.ImagePath,
		user.VerificationToken,
		user.Verified,
	)
	if errDB != nil {
		return fmt.Errorf("could not insert new user: %w", errDB)
	}

	// 3) Send verification email
	if errEmail := repo.sendVerificationEmail(user); errEmail != nil {
		// Optionally remove the user if email fails, or leave as is
		return fmt.Errorf("error sending verification email: %w", errEmail)
	}

	return nil
}

func (repo *UserRepository) sendVerificationEmail(user models.User) error {
	from := "mori.team.contact@gmail.com"
	password := "qeey kngz gmyn bzwi" // Put this in env variables for production
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	toList := []string{user.Email}
	subject := "Verify Your Email Address"

	// Build the verification link
	verifyLink := fmt.Sprintf("http://localhost:8080/verified?token=%s", user.VerificationToken)

	// Updated HTML + CSS
	body := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Verification</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f4f4f4;
            color: #333;
        }
        .container {
            max-width: 600px;
            margin: 20px auto;
            padding: 20px;
            background-color: #2c2c2c;
            border-radius: 10px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            color: #e4e4e4;
        }
        .header {
            text-align: center;
            padding: 10px 0;
        }
        .header img {
            width: 100px;
            height: auto;
            margin-bottom: 10px;
            border-radius: 180px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        }
        .header h1 {
            font-size: 24px;
            margin: 0;
            color: #9146bc;
        }
        .content {
            margin: 20px 0;
            font-size: 16px;
            line-height: 1.6;
        }
        .content p {
            color: #e4e4e4 !important;
            font-size: 16px !important;
        }
        .button {
            display: inline-block;
            padding: 12px 20px;
            margin: 20px 0;
            font-size: 16px;
            color: #fff !important;
            background-color: #9146bc;
            border-radius: 5px;
            text-decoration: none;
        }
        .footer {
            margin-top: 80px;
            text-align: center;
            font-size: 12px;
            color: #aaa;
        }
        .footer a {
            color: #9146bc;
            text-decoration: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Mori Team</h1>
        </div>
        <div class="content">
            <p>Hello %s,</p>
            <p>Thank you for signing up with Mori! Please verify your email address by clicking the button below.</p>
            <p>
                <a href="%s" class="button">Verify Email</a>
            </p>
            <p>If the button doesn't work, please copy and paste the following link into your browser.</p>
            <p>%s</p>
            <p>Thank you, and welcome to Mori.</p>
        </div>
        <div class="footer">
            <p>&copy; 2024 Mori Team. All Rights Reserved.</p>
            <p>
                <a href="#">Privacy Policy</a> | <a href="#">Contact Us</a>
            </p>
        </div>
    </div>
</body>
</html>`,
		user.FirstName, verifyLink, verifyLink)

	// Build MIME headers + body
	msg := strings.Join([]string{
		fmt.Sprintf("Subject: %s", subject),
		"MIME-version: 1.0;",
		`Content-Type: text/html; charset="UTF-8";`,
		"",
		body,
	}, "\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	errSend := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, toList, []byte(msg))
	if errSend != nil {
		return fmt.Errorf("failed to send email: %w", errSend)
	}

	return nil
}

func (repo *UserRepository) VerifyEmail(token string) error {
	// 1) If token is empty, return error
	if token == "" {
		return fmt.Errorf("missing token")
	}

	// 2) Check if the token actually exists in the DB
	var count int
	err := repo.DB.QueryRow(`
        SELECT COUNT(*)
        FROM users
        WHERE verification_token = $1
    `, token).Scan(&count)
	if err != nil {
		return fmt.Errorf("database error while checking token: %w", err)
	}

	if count == 0 {
		// Means no row found with this token
		return fmt.Errorf("invalid or expired token")
	}

	// 3) If the token exists, update to verified
	_, err = repo.DB.Exec(`
        UPDATE users
        SET verified = TRUE,
            verification_token = NULL
        WHERE verification_token = $1
    `, token)
	if err != nil {
		return fmt.Errorf("could not update user verification: %w", err)
	}

	// 4) If everything is good, return nil
	return nil
}

// generateRandomToken creates a 32-character hex string
func generateRandomToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (repo *UserRepository) SetResetToken(userID, token string, expires time.Time) error {
	query := `
		UPDATE users
		SET reset_token = $1,
		    reset_token_expires = $2
		WHERE user_id = $3
	`
	_, err := repo.DB.Exec(query, token, expires, userID)
	return err
}

func (repo *UserRepository) UpdatePasswordAndClearToken(userID, newHashedPwd string) error {
	query := `
		UPDATE users
		SET password = $1,
		    reset_token = NULL,
		    reset_token_expires = NULL
		WHERE user_id = $2
	`
	_, err := repo.DB.Exec(query, newHashedPwd, userID)
	return err
}

func (repo *UserRepository) FindUserByResetToken(token string) (models.User, error) {
	query := `
		SELECT user_id, email, reset_token_expires
		FROM users
		WHERE reset_token = $1
		LIMIT 1
	`
	var user models.User
	var expires time.Time

	err := repo.DB.QueryRow(query, token).Scan(
		&user.ID,
		&user.Email,
		&expires,
	)
	if err != nil {
		return user, err
	}
	user.ResetToken = token
	user.ResetTokenExpires = &expires
	return user, nil
}

// EmailNotTaken checks if an email is not already registered.
func (repo *UserRepository) EmailNotTaken(email string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM users 
		WHERE email = $1;
	`
	var count int
	err := repo.DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// FindUserByEmail retrieves a user ID and password by email.
func (repo *UserRepository) FindUserByEmail(email string) (models.User, error) {
	query := `
		SELECT user_id, password, verified 
		FROM users 
		WHERE email = $1;
	`
	var user models.User
	err := repo.DB.QueryRow(query, email).Scan(&user.ID, &user.Password, &user.Verified)
	return user, err
}

// GetAllAndFollowing retrieves a list of users with follower and following status.
func (repo *UserRepository) GetAllAndFollowing(userID string) ([]models.User, error) {
	query := `
		SELECT user_id, 
		       COALESCE(nickname, first_name || ' ' || last_name), 
		       (SELECT COUNT(*) FROM followers WHERE followers.user_id = $1 AND follower_id = users.user_id) AS follower, 
		       (SELECT COUNT(*) FROM followers WHERE followers.user_id = users.user_id AND follower_id = $1) AS following, 
		       image 
		FROM users 
		WHERE user_id != $1;
	`
	rows, err := repo.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var follower, following int
		if err := rows.Scan(&user.ID, &user.Nickname, &follower, &following, &user.ImagePath); err != nil {
			return nil, err
		}
		user.Follower = follower > 0
		user.Following = following > 0
		users = append(users, user)
	}
	return users, rows.Err()
}

// GetDataMin retrieves basic user information (ID, nickname, image).
func (repo *UserRepository) GetDataMin(userID string) (models.User, error) {
	query := `
		SELECT user_id, 
		       COALESCE(nickname, first_name || ' ' || last_name), 
		       image 
		FROM users 
		WHERE user_id = $1;
	`
	var user models.User
	err := repo.DB.QueryRow(query, userID).Scan(&user.ID, &user.Nickname, &user.ImagePath)
	return user, err
}

// IsFollowing checks if the current user is following another user.
func (repo *UserRepository) IsFollowing(userID, currentUserID string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM followers 
		WHERE user_id = $1 AND follower_id = $2;
	`
	var count int
	err := repo.DB.QueryRow(query, userID, currentUserID).Scan(&count)
	return count > 0, err
}

// ProfileStatus retrieves the status of a user profile.
func (repo *UserRepository) ProfileStatus(userID string) (string, error) {
	query := `
		SELECT status 
		FROM users 
		WHERE user_id = $1;
	`
	var status string
	err := repo.DB.QueryRow(query, userID).Scan(&status)
	return status, err
}

// GetProfileMax retrieves full user information.
func (repo *UserRepository) GetProfileMax(userID string) (models.User, error) {
	query := `
		SELECT COALESCE(nickname, first_name || ' ' || last_name), 
		       first_name, 
		       last_name, 
		       image, 
		       email, 
		       TO_CHAR(birthday, 'DD.MM.YYYY'), 
		       about 
		FROM users 
		WHERE user_id = $1;
	`
	var user models.User
	err := repo.DB.QueryRow(query, userID).Scan(&user.Nickname, &user.FirstName, &user.LastName, &user.ImagePath, &user.Email, &user.DateOfBirth, &user.About)
	user.ID = userID
	return user, err
}

// GetProfileMin retrieves minimal user profile information.
func (repo *UserRepository) GetProfileMin(userID string) (models.User, error) {
	query := `
		SELECT COALESCE(nickname, first_name || ' ' || last_name), 
		       image 
		FROM users 
		WHERE user_id = $1;
	`
	var user models.User
	err := repo.DB.QueryRow(query, userID).Scan(&user.Nickname, &user.ImagePath)
	user.ID = userID
	return user, err
}

// GetFollowers retrieves the followers of a user.
func (repo *UserRepository) GetFollowers(userID string) ([]models.User, error) {
	query := `
		SELECT user_id, 
		       COALESCE(nickname, first_name || ' ' || last_name),
			   image
		FROM users 
		WHERE (SELECT COUNT(*) FROM followers WHERE followers.user_id = $1 AND follower_id = users.user_id) = 1;
	`
	rows, err := repo.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Nickname, &user.ImagePath); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

// GetFollowing retrieves the users the current user is following.
func (repo *UserRepository) GetFollowing(userID string) ([]models.User, error) {
	query := `
		SELECT user_id, 
		       COALESCE(nickname, first_name || ' ' || last_name),
			   image
		FROM users 
		WHERE (SELECT COUNT(*) FROM followers WHERE followers.follower_id = $1 AND user_id = users.user_id) = 1;
	`
	rows, err := repo.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Nickname, &user.ImagePath); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

// GetStatus retrieves the current status of a user.
func (repo *UserRepository) GetStatus(userID string) (string, error) {
	query := `
		SELECT status 
		FROM users 
		WHERE user_id = $1;
	`
	var status string
	err := repo.DB.QueryRow(query, userID).Scan(&status)
	return status, err
}

// SetStatus updates the status of a user.
func (repo *UserRepository) SetStatus(user models.User) error {
	query := `
		UPDATE users 
		SET status = $1 
		WHERE user_id = $2;
	`
	_, err := repo.DB.Exec(query, user.Status, user.ID)
	return err
}

// SaveFollower adds a follower to a user.
func (repo *UserRepository) SaveFollower(userID, followerID string) error {
	query := `
		INSERT INTO followers(user_id, follower_id) 
		VALUES ($1, $2);
	`
	_, err := repo.DB.Exec(query, userID, followerID)
	return err
}

// DeleteFollower removes a follower from a user.
func (repo *UserRepository) DeleteFollower(userID, followerID string) error {
	query := `
		DELETE FROM followers 
		WHERE user_id = $1 AND follower_id = $2;
	`
	_, err := repo.DB.Exec(query, userID, followerID)
	return err
}

// UpdateNickname modifies the user's nickname.
func (repo *UserRepository) UpdateNickname(userID, newNickname string) error {
	query := `
		UPDATE users 
		SET nickname = $1 
		WHERE user_id = $2;
	`
	_, err := repo.DB.Exec(query, newNickname, userID)
	return err
}

// UpdateAvatar modifies the user's avatar.
func (repo *UserRepository) UpdateAvatar(userID, avatarPath string) error {
	query := `
		UPDATE users 
		SET image = $1 
		WHERE user_id = $2;
	`
	_, err := repo.DB.Exec(query, avatarPath, userID)
	return err
}

func (repo *UserRepository) DeleteUser(userID string) error {
	// Requête SQL pour supprimer l'utilisateur
	query := `DELETE FROM users WHERE user_id = $1`

	// Exécute la requête en passant userID comme paramètre
	_, err := repo.DB.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("could not delete account: %w", err)
	}

	return nil
}
