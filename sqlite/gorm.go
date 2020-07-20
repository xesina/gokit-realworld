package sqlite

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	realworld "github.com/xesina/go-kit-realworld-example-app"
	"time"
)

type Model struct {
	ID        int64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type Storage struct {
	DB *gorm.DB
}

func NewStorage(filename string) (*Storage, error) {
	db, err := gorm.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	db.DB().SetMaxIdleConns(3)
	db.LogMode(true)
	return &Storage{DB: db}, nil
}

func (s *Storage) Migrate() {
	s.DB.AutoMigrate(
		&User{},
		&Follow{},
		&Article{},
		&Comment{},
		&Tag{},
	)
}

func (s *Storage) NewUserRepository() realworld.UserRepo {
	return &userRepository{
		db: s.DB,
	}
}

func (s *Storage) NewArticleRepository() realworld.ArticleRepo {
	return &articleRepository{
		db: s.DB,
	}
}
