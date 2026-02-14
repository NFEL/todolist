package storage_postgres

import (
	"context"
	"fmt"
	"graph-interview/internal/api/handlers/dto"
	"graph-interview/internal/domain"
	"graph-interview/internal/repository/storage"
	"math"

	"gorm.io/gorm"
)

type userImp struct {
	db *gorm.DB
}

func NewUserRepo(db *storage.DB) *userImp {
	return &userImp{
		db: db.DB,
	}
}

func (i *userImp) Create(ctx context.Context, user *domain.User) (uint, error) {
	err := gorm.G[domain.User](i.db).Create(ctx, user)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}
func (i *userImp) GetByID(ctx context.Context, ID uint) (domain.User, error) {
	return gorm.G[domain.User](i.db).Where("id = ?", ID).Take(ctx)
}

func (i *userImp) GetByField(ctx context.Context, field string, value any) (domain.User, error) {
	return gorm.G[domain.User](i.db).Where(fmt.Sprintf("%s = ?", field), value).Take(ctx)
}

func (i *userImp) List(ctx context.Context, limit, offset int) ([]domain.User, error) {
	r := make([]domain.User, int(math.Abs(float64(limit-offset))))
	err := gorm.G[domain.User](i.db).Select("*").Limit(limit).Offset(offset).Scan(ctx, &r)
	if err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (i *userImp) ListByFilter(ctx context.Context, filter dto.UserListFilter, limit, offset int) ([]domain.User, error) {
	q := gorm.G[domain.User](i.db).Select("*").Limit(limit).Offset(offset)
	return q.Find(ctx)
}

func (i *userImp) UpdateByID(ctx context.Context, user *domain.User, fields []string) error {
	_, err := gorm.G[domain.User](i.db).Where("id = ?", user.ID).Select(fields[0], fields[1:]).Updates(ctx, *user)
	return err
}
