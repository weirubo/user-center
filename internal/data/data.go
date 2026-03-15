package data

import (
	"context"
	"time"

	"user-center/internal/biz"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewDB,
	NewCache,
	NewUserRepo,
)

type Data struct {
	db    *DB
	cache *Cache
}

func NewData(db *DB, cache *Cache) (func(), error) {
	return func() {
		db.Close()
		cache.Close()
	}, nil
}

func (d *Data) GetUserByEmail(email string) (*biz.User, error) {
	user, err := d.db.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return biz.UserFromEntity(user), nil
}

func (d *Data) GetUserByPhone(phone string) (*biz.User, error) {
	user, err := d.db.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}
	return biz.UserFromEntity(user), nil
}

func (d *Data) CreateUser(user *biz.User) error {
	return d.db.CreateUser(biz.UserToEntity(user))
}

func (d *Data) GetUserByID(id int64) (*biz.User, error) {
	user, err := d.db.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return biz.UserFromEntity(user), nil
}

func (d *Data) SetCache(key string, value interface{}) error {
	return d.cache.Set(key, value)
}

func (d *Data) GetCache(key string) (string, error) {
	return d.cache.Get(key)
}

type UserRepo struct {
	db    *DB
	cache *Cache
}

func NewUserRepo(db *DB, cache *Cache) biz.UserRepo {
	return &UserRepo{db: db, cache: cache}
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*biz.User, error) {
	user, err := r.db.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return biz.UserFromEntity(user), nil
}

func (r *UserRepo) GetUserByPhone(ctx context.Context, phone string) (*biz.User, error) {
	user, err := r.db.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}
	return biz.UserFromEntity(user), nil
}

func (r *UserRepo) CreateUser(ctx context.Context, user *biz.User) error {
	return r.db.CreateUser(biz.UserToEntity(user))
}

func (r *UserRepo) GetUserByID(ctx context.Context, id int64) (*biz.User, error) {
	user, err := r.db.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return biz.UserFromEntity(user), nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, user *biz.User) error {
	return r.db.UpdateUser(biz.UserToEntity(user))
}

func (r *UserRepo) DeleteUser(ctx context.Context, id int64) error {
	return r.db.DeleteUser(id)
}

func (r *UserRepo) SetCache(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	return r.cache.SetWithExpire(key, value, expire)
}

func (r *UserRepo) GetCache(ctx context.Context, key string) (string, error) {
	return r.cache.Get(key)
}

func (r *UserRepo) DeleteCache(ctx context.Context, key string) error {
	return r.cache.Delete(key)
}
