package sqlite

import (
	"github.com/jinzhu/gorm"
	realworld "github.com/xesina/gokit-realworld"
)

type Article struct {
	Model
	Slug        string `gorm:"unique_index;not null"`
	Title       string `gorm:"not null"`
	Description string
	Body        string
	Author      User
	AuthorID    int64
	Comments    []Comment
	Favorites   []User `gorm:"many2many:favorites;"`
	Tags        []Tag  `gorm:"many2many:article_tags;association_autocreate:false"`
}

type Comment struct {
	Model
	Article   Article
	ArticleID int64
	User      User
	UserID    int64
	Body      string
}

type Tag struct {
	Model
	Tag      string    `gorm:"unique_index"`
	Articles []Article `gorm:"many2many:article_tags;"`
}

type articleRepository struct {
	db *gorm.DB
}

func (s articleRepository) Get(slug string) (*realworld.Article, error) {
	var m Article

	err := s.db.Where(&Article{Slug: slug}).Preload("Favorites").Preload("Tags").Preload("Author").Find(&m).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, realworld.ErrArticleNotFound
		}

		return nil, err
	}

	return s.domainArticle(&m), err
}

func (s articleRepository) List(offset, limit int) ([]*realworld.Article, int, error) {
	var articles []Article

	err := s.db.Preload("Favorites").
		Preload("Tags").
		Preload("Author").
		Offset(offset).
		Limit(limit).
		Order("created_at desc").Find(&articles).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []*realworld.Article{}, 0, nil
		}
		return nil, 0, err
	}

	return s.domainArticles(articles), len(articles), nil
}

func (s articleRepository) ListByTag(tag string, offset, limit int) ([]*realworld.Article, int, error) {
	var (
		t        Tag
		articles []Article
	)

	err := s.db.Where(&Tag{Tag: tag}).First(&t).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []*realworld.Article{}, 0, nil
		}
		return nil, 0, err
	}

	err = s.db.Model(&t).
		Preload("Favorites").
		Preload("Tags").
		Preload("Author").
		Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Association("Articles").
		Find(&articles).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []*realworld.Article{}, 0, nil
		}
		return nil, 0, err
	}

	return s.domainArticles(articles), len(articles), nil
}

func (s articleRepository) ListByAuthorID(id int64, offset, limit int) ([]*realworld.Article, int, error) {
	var articles []Article

	err := s.db.Where(Article{AuthorID: id}).
		Preload("Favorites").
		Preload("Tags").
		Preload("Author").
		Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Find(&articles).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []*realworld.Article{}, 0, nil
		}
		return nil, 0, err
	}

	return s.domainArticles(articles), len(articles), nil
}

func (s articleRepository) ListByFavoriterID(id int64, offset, limit int) ([]*realworld.Article, int, error) {
	var articles []Article

	err := s.db.Model(&User{Model: Model{ID: id}}).
		Preload("Favorites").
		Preload("Tags").
		Preload("Author").
		Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Association("Favorites").
		Find(&articles).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []*realworld.Article{}, 0, nil
		}
		return nil, 0, err
	}

	return s.domainArticles(articles), len(articles), nil
}

func (s articleRepository) Feed(req realworld.FeedRequest) ([]*realworld.Article, int, error) {
	var (
		u        User
		articles []Article
	)

	err := s.db.First(&u, req.UserID).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []*realworld.Article{}, 0, nil
		}
		return nil, 0, err
	}

	var followings []Follow

	err = s.db.Model(&u).
		Preload("Following").
		Preload("Follower").
		Association("Followings").
		Find(&followings).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []*realworld.Article{}, 0, nil
		}
		return nil, 0, err
	}

	if len(followings) == 0 {
		return []*realworld.Article{}, 0, nil
	}

	ids := make([]int64, len(followings))
	for i, f := range followings {
		ids[i] = f.FollowingID
	}

	err = s.db.Where("author_id in (?)", ids).
		Preload("Favorites").
		Preload("Tags").
		Preload("Author").
		Offset(req.Offset).
		Limit(req.Limit).
		Order("created_at desc").
		Find(&articles).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []*realworld.Article{}, 0, nil
		}
		return nil, 0, err
	}

	return s.domainArticles(articles), len(articles), nil
}

func (s articleRepository) Create(a realworld.Article) (*realworld.Article, error) {
	var found Article
	err := s.db.Where("slug = ?", a.Slug).First(&found).Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	}

	// already exists
	if found.ID > 0 {
		return nil, realworld.ErrArticleAlreadyExists
	}

	m := s.articleModel(&a)

	tx := s.db.Begin()

	if err := tx.Create(m).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// need this block to not add the duplicate tags and getting sql unique errors
	for _, t := range m.Tags {
		err := tx.Where(&Tag{Tag: t.Tag}).First(&t).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			tx.Rollback()
			return nil, err
		}

		if err := tx.Model(m).Association("Tags").Append(t).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err = tx.Where(m.ID).
		Preload("Favorites").
		Preload("Tags").
		Preload("Author").
		Find(m).Error

	if err != nil {
		tx.Rollback()
		if gorm.IsRecordNotFoundError(err) {
			return nil, realworld.ErrArticleNotFound
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return s.domainArticle(m), nil
}

func (s articleRepository) Update(slug string, a realworld.Article) (*realworld.Article, error) {
	var found Article
	err := s.db.Where("slug = ?", slug).First(&found).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, realworld.ErrArticleNotFound
		}
	}

	m := s.articleModel(&a)
	m.ID = found.ID

	tx := s.db.Begin()

	if err := tx.Model(&found).Update(m).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tags := make([]Tag, 0)

	for _, t := range m.Tags {
		tag := Tag{Tag: t.Tag}

		err := tx.Where(&tag).First(&tag).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			tx.Rollback()
			return nil, err
		}

		tags = append(tags, tag)
	}

	if err := tx.Model(m).Association("Tags").Replace(tags).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Where(m.ID).
		Preload("Favorites").
		Preload("Tags").
		Preload("Author").
		Find(m).Error

	if err != nil {
		tx.Rollback()
		if gorm.IsRecordNotFoundError(err) {
			return nil, realworld.ErrArticleNotFound
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return s.domainArticle(m), nil
}

func (s articleRepository) Delete(a realworld.Article) error {
	return s.db.Where("slug = ?", a.Slug).Delete(Article{}).Error
}

func (s articleRepository) AddFavorite(a realworld.Article, u realworld.User) (*realworld.Article, error) {
	var m Article
	err := s.db.Where("slug = ?", a.Slug).First(&m).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, realworld.ErrArticleNotFound
		}
	}

	err = s.db.Model(m).Association("Favorites").Append(userModel(&u)).Error

	err = s.db.Where(m.ID).
		Preload("Favorites").
		Preload("Tags").
		Preload("Author").
		Find(&m).Error

	return s.domainArticle(&m), nil
}

func (s articleRepository) RemoveFavorite(a realworld.Article, u realworld.User) (*realworld.Article, error) {
	var m Article
	err := s.db.Where("slug = ?", a.Slug).First(&m).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, realworld.ErrArticleNotFound
		}
	}

	err = s.db.Model(m).Association("Favorites").Delete(userModel(&u)).Error

	err = s.db.Where(m.ID).
		Preload("Favorites").
		Preload("Tags").
		Preload("Author").
		Find(&m).Error

	return s.domainArticle(&m), nil
}

func (s articleRepository) AddComment(c realworld.Comment) (*realworld.Comment, error) {
	var article Article

	err := s.db.Where("slug = ?", c.Article.Slug).First(&article).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, realworld.ErrArticleNotFound
		}
	}

	cm := s.commentModel(&c)

	err = s.db.Model(&article).Association("Comments").Append(cm).Error

	if err != nil {
		return nil, err
	}

	err = s.db.Where(cm.ID).Preload("User").First(cm).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, realworld.ErrResourceNotFound
		}
		return nil, err
	}

	return s.domainComment(cm), nil
}

func (s articleRepository) DeleteComment(c realworld.Comment) error {
	cm := s.commentModel(&c)
	return s.db.Delete(cm).Error
}

func (s articleRepository) Comments(a realworld.Article) ([]*realworld.Comment, error) {
	var m Article
	err := s.db.Where(&Article{Slug: a.Slug}).Preload("Comments").Preload("Comments.User").First(&m).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, realworld.ErrArticleNotFound
		}

		return nil, err
	}

	return s.domainComments(m.Comments), nil
}

func (s articleRepository) Tags() ([]*realworld.Tag, error) {
	var tags []Tag
	if err := s.db.Find(&tags).Error; err != nil {
		return nil, err
	}

	return s.domainTags(tags), nil
}

func (s *articleRepository) articleModel(a *realworld.Article) *Article {
	return &Article{
		Model: Model{
			ID:        a.ID,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
			DeletedAt: nil,
		},
		Slug:        a.Slug,
		Title:       a.Title,
		Description: a.Description,
		Body:        a.Body,
		AuthorID:    a.Author.ID,
		Tags:        s.tags(a),
	}
}

func (s *articleRepository) domainArticle(m *Article) *realworld.Article {
	return &realworld.Article{
		ID:          m.ID,
		Slug:        m.Slug,
		Title:       m.Title,
		Description: m.Description,
		Body:        m.Body,
		Author:      realworld.User{ID: m.AuthorID},
		Favorites:   s.favoriteMap(m.Favorites),
		Tags:        s.tagMap(m.Tags),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (s *articleRepository) domainArticles(m []Article) []*realworld.Article {
	aa := make([]*realworld.Article, 0)

	for _, a := range m {
		aa = append(aa, s.domainArticle(&a))
	}

	return aa
}

func (s *articleRepository) favoriteMap(ff []User) realworld.Favorites {
	fm := make(realworld.Favorites)
	for _, f := range ff {
		fm[f.ID] = struct{}{}
	}
	return fm
}

func (s *articleRepository) tagMap(tt []Tag) realworld.Tags {
	tm := make(realworld.Tags)
	for _, t := range tt {
		tm[t.Tag] = realworld.Tag{
			ID:  t.ID,
			Tag: t.Tag,
		}
	}
	return tm
}

func (s *articleRepository) tags(a *realworld.Article) []Tag {
	var tm []Tag
	for _, t := range a.Tags {
		tm = append(tm, Tag{
			Model: Model{ID: t.ID},
			Tag:   t.Tag,
		})
	}

	return tm
}

func (s *articleRepository) commentModel(c *realworld.Comment) *Comment {
	return &Comment{
		Model:     Model{ID: c.ID},
		ArticleID: c.ArticleID,
		UserID:    c.UserID,
		Body:      c.Body,
	}
}

func (s *articleRepository) domainComment(c *Comment) *realworld.Comment {
	return &realworld.Comment{
		ID:        c.ID,
		ArticleID: c.ArticleID,
		UserID:    c.UserID,
		Body:      c.Body,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func (s *articleRepository) domainTags(tt []Tag) []*realworld.Tag {
	tags := make([]*realworld.Tag, 0)

	for _, t := range tt {
		tags = append(tags, &realworld.Tag{
			ID:  t.ID,
			Tag: t.Tag,
		})
	}

	return tags
}

func (s *articleRepository) domainComments(cc []Comment) []*realworld.Comment {
	var comments []*realworld.Comment
	for _, c := range cc {
		comments = append(comments, &realworld.Comment{
			ID:        c.ID,
			ArticleID: c.ArticleID,
			UserID:    c.UserID,
			Body:      c.Body,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})
	}

	return comments
}
