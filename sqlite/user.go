package sqlite

import (
	"database/sql"
	"errors"
	"github.com/jinzhu/gorm"
	realworld "github.com/xesina/gokit-realworld"
)

type User struct {
	Model
	Username   string `gorm:"unique_index;not null"`
	Email      string `gorm:"unique_index;not null"`
	Password   string `gorm:"not null"`
	Bio        sql.NullString
	Image      sql.NullString
	Followers  []Follow  `gorm:"foreignkey:FollowingID"`
	Followings []Follow  `gorm:"foreignkey:FollowerID"`
	Favorites  []Article `gorm:"many2many:favorites;"`
}

type Follow struct {
	Follower    User
	FollowerID  int64 `gorm:"primary_key" sql:"type:int not null"`
	Following   User
	FollowingID int64 `gorm:"primary_key" sql:"type:int not null"`
}

type userRepository struct {
	db *gorm.DB
}

func (s *userRepository) Create(u realworld.User) (*realworld.User, error) {
	m, err := s.GetByUsername(u.Username)
	if err != nil && !errors.Is(err, realworld.ErrUserNotFound) {
		return nil, err
	}

	if m != nil {
		return nil, realworld.ErrUserAlreadyExists
	}

	user := userModel(&u)
	err = s.db.Create(user).Error
	return s.domainUser(user), err
}

func (s *userRepository) Update(u realworld.User) (*realworld.User, error) {
	old, err := s.GetByUsername(u.Username)
	if err != nil {
		return nil, err
	}

	if u.Password == "" {
		u.Password = old.Password
	}

	model := userModel(&u)
	err = s.db.Model(model).Update(model).Error
	return &u, err
}

func (s *userRepository) Get(e string) (*realworld.User, error) {
	var m User
	if err := s.db.Where(&User{Email: e}).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, realworld.ErrUserNotFound
		}
		return nil, err
	}
	return s.domainUser(&m), nil
}

func (s *userRepository) GetByID(id int64) (*realworld.User, error) {
	m, err := s.getByID(id)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, realworld.ErrUserNotFound
		}

		return nil, err
	}
	return s.domainUser(m), nil
}

func (s *userRepository) getByID(id int64) (*User, error) {
	var m User
	if err := s.db.First(&m, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, realworld.ErrUserNotFound
		}
		return nil, err
	}
	return &m, nil
}

func (s *userRepository) GetByUsername(username string) (u *realworld.User, err error) {
	m, err := s.getByUsername(username)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, realworld.ErrUserNotFound
		}
		return nil, err
	}
	return s.domainUser(m), nil
}

func (s *userRepository) getByUsername(username string) (u *User, err error) {
	var m User
	if err := s.db.Where(&User{Username: username}).Preload("Followers").First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *userRepository) AddFollower(followerID, followeeID int64) (*realworld.User, error) {
	// TODO: should we do this check in here or in service layer?
	_, err := s.getByID(followerID)
	if err != nil {
		return nil, err
	}

	followee, err := s.getByID(followeeID)
	if err != nil {
		return nil, err
	}

	err = s.db.Model(followee).
		Association("Followers").
		Append(
			&Follow{FollowerID: followerID, FollowingID: followeeID},
		).Error

	if err != nil {
		return nil, err
	}

	f, err := s.getByUsername(followee.Username)
	if err != nil {
		return nil, err
	}

	return s.domainUser(f), nil
}

func (s *userRepository) RemoveFollower(followerID, followeeID int64) (*realworld.User, error) {
	_, err := s.getByID(followerID)
	if err != nil {
		return nil, err
	}

	followee, err := s.getByID(followeeID)
	if err != nil {
		return nil, err
	}

	f := Follow{
		FollowerID:  followerID,
		FollowingID: followeeID,
	}

	if err := s.db.Model(followee).Association("Followers").Find(&f).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return s.domainUser(followee), nil
		}
		return nil, err
	}

	if err := s.db.Delete(f).Error; err != nil {
		return nil, err
	}

	return s.domainUser(followee), nil
}

func userModel(u *realworld.User) *User {
	return &User{
		Model: Model{
			ID: u.ID,
		},
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
		Bio: sql.NullString{
			String: u.Bio.Value,
			Valid:  u.Bio.Valid,
		},
		Image: sql.NullString{
			String: u.Image.Value,
			Valid:  u.Image.Valid,
		},
		Followers:  nil,
		Followings: nil,
		Favorites:  nil,
	}
}

func (s *userRepository) domainUser(u *User) *realworld.User {
	return &realworld.User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
		Email:    u.Email,
		Bio: realworld.Bio{
			Value: u.Bio.String,
			Valid: u.Bio.Valid,
		},
		Image: realworld.Image{
			Value: u.Image.String,
			Valid: u.Image.Valid,
		},
		Followers:  s.followersMap(u.Followers),
		Followings: s.followingMap(u.Followings),
	}
}

func (s *userRepository) followersMap(ff []Follow) realworld.Follows {
	fm := make(realworld.Follows)
	for _, f := range ff {
		fm[f.FollowerID] = struct{}{}
	}
	return fm
}

func (s *userRepository) followingMap(ff []Follow) realworld.Follows {
	fm := make(realworld.Follows)
	for _, f := range ff {
		fm[f.FollowingID] = struct{}{}
	}
	return fm
}
