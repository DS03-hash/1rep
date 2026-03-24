package userService

import "gorm.io/gorm"

type Repository interface {
	Create(u *User) error
	List() ([]User, error)
	GetByID(id uint) (*User, error)
	Update(u *User) error
	DeleteByID(id uint) (rowsAffected int64, err error)
}

type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates user repository backed by GORM.
func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// Create saves new user in database.
func (r *GormRepository) Create(u *User) error {
	return r.db.Create(u).Error
}

// List returns all non-deleted users.
func (r *GormRepository) List() ([]User, error) {
	var users []User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetByID finds user by id.
func (r *GormRepository) GetByID(id uint) (*User, error) {
	var u User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// Update persists user changes.
func (r *GormRepository) Update(u *User) error {
	return r.db.Save(u).Error
}

// DeleteByID soft-deletes user by id.
func (r *GormRepository) DeleteByID(id uint) (int64, error) {
	res := r.db.Delete(&User{}, id)
	return res.RowsAffected, res.Error
}
