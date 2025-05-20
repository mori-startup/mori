package models

import "time"

// defines User data type
type User struct {
	ID                   string     `json:"id"`
	Email                string     `json:"login"`
	FirstName            string     `json:"firstName"`
	LastName             string     `json:"lastName"`
	Password             string     `json:"password,omitempty"`
	Nickname             string     `json:"nickname"`
	About                string     `json:"about"`
	DateOfBirth          string     `json:"dateOfBirth"`
	ImagePath            string     `json:"avatar"`
	VerificationToken    string     `json:"verificationToken"`
	Verified             bool       `json:"verified"`
	Status               string     `json:"status"`         // "private" / "public"
	CurrentUser          bool       `json:"currentUser"`    //  returns true for current, false otherwise
	Follower             bool       `json:"follower"`       //  if this user is following another user
	Following            bool       `json:"following"`      //  if curr user is following this one
	FollowRequestPending bool       `json:"requestPending"` // true if requested to follow
	ResetToken           string     `json:"resetToken"`
	ResetTokenExpires    *time.Time `json:"resetTokenExpires,omitempty"`
}

// Repository represent all possible actions available to deal with Users
// all db packages (in case of different db) should implement these methods
type UserRepository interface {
	Add(User) error                           // save new user in db
	EmailNotTaken(email string) (bool, error) // returns true if not taken
	FindUserByEmail(email string) (User, error)
	VerifyEmail(token string) error

	// Follow/friend logic
	GetAllAndFollowing(userID string) ([]User, error)
	GetFollowers(userId string) ([]User, error)
	GetFollowing(userId string) ([]User, error)
	SaveFollower(userId, followerId string) error
	DeleteFollower(userId, followerId string) error
	IsFollowing(userID, currentUserID string) (bool, error)

	// Profile/Status logic
	GetDataMin(userID string) (User, error)
	ProfileStatus(userID string) (string, error)
	GetProfileMax(userID string) (User, error)
	GetProfileMin(userID string) (User, error)
	GetStatus(userID string) (string, error)
	SetStatus(User) error

	// Reset token logic
	SetResetToken(userID, token string, expires time.Time) error
	FindUserByResetToken(token string) (User, error)
	UpdatePasswordAndClearToken(userID, newHashedPwd string) error

	UpdateNickname(userID, newNickname string) error // update user nickname
	UpdateAvatar(userID, avatarPath string) error    // update user avatar
	DeleteUser(userID string) error
}
