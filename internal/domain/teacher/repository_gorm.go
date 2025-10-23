package teacher

import (
	"context"
	"time"

	"github.com/JustDoItBetter/FITS-backend/internal/common/errors"
	"github.com/JustDoItBetter/FITS-backend/internal/common/pagination"
	"gorm.io/gorm"
)

// TeacherModel represents the GORM model for teachers table
type TeacherModel struct {
	ID         string     `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	FirstName  string     `gorm:"column:first_name;type:varchar(100);not null"`
	LastName   string     `gorm:"column:last_name;type:varchar(100);not null"`
	Email      string     `gorm:"column:email;type:varchar(255);uniqueIndex;not null"`
	Department string     `gorm:"column:department;type:varchar(100);not null"`
	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt  *time.Time `gorm:"column:deleted_at;index"` // Soft delete support
}

// TableName specifies the table name for GORM
func (TeacherModel) TableName() string {
	return "teachers"
}

// ToTeacher converts TeacherModel to Teacher domain entity
func (m *TeacherModel) ToTeacher() *Teacher {
	return &Teacher{
		UUID:       m.ID,
		FirstName:  m.FirstName,
		LastName:   m.LastName,
		Email:      m.Email,
		Department: m.Department,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

// FromTeacher converts Teacher domain entity to TeacherModel
func FromTeacher(t *Teacher) *TeacherModel {
	return &TeacherModel{
		ID:         t.UUID,
		FirstName:  t.FirstName,
		LastName:   t.LastName,
		Email:      t.Email,
		Department: t.Department,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}

// GormRepository implements Repository interface using GORM
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new GORM-based teacher repository
func NewGormRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

// WithDB returns a new repository instance using the provided database connection
// This enables the repository to participate in transactions
func (r *GormRepository) WithDB(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

// Create adds a new teacher to the database
func (r *GormRepository) Create(ctx context.Context, teacher *Teacher) error {
	model := FromTeacher(teacher)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if errors.IsUniqueViolation(err) {
			return errors.Conflict("teacher with this email already exists")
		}
		return errors.Internal("failed to create teacher: " + err.Error())
	}

	// Update teacher with DB-generated fields
	*teacher = *model.ToTeacher()
	return nil
}

// GetByUUID retrieves a teacher by UUID
func (r *GormRepository) GetByUUID(ctx context.Context, uuid string) (*Teacher, error) {
	var model TeacherModel

	if err := r.db.WithContext(ctx).Where("id = ?", uuid).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("teacher not found")
		}
		return nil, errors.Internal("failed to get teacher: " + err.Error())
	}

	return model.ToTeacher(), nil
}

// Update updates an existing teacher
func (r *GormRepository) Update(ctx context.Context, teacher *Teacher) error {
	model := FromTeacher(teacher)

	result := r.db.WithContext(ctx).
		Model(&TeacherModel{}).
		Where("id = ?", teacher.UUID).
		Updates(map[string]interface{}{
			"first_name": model.FirstName,
			"last_name":  model.LastName,
			"email":      model.Email,
			"department": model.Department,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		if errors.IsUniqueViolation(result.Error) {
			return errors.Conflict("teacher with this email already exists")
		}
		return errors.Internal("failed to update teacher: " + result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return errors.NotFound("teacher not found")
	}

	return nil
}

// Delete removes a teacher from the database
func (r *GormRepository) Delete(ctx context.Context, uuid string) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", uuid).
		Delete(&TeacherModel{})

	if result.Error != nil {
		return errors.Internal("failed to delete teacher: " + result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return errors.NotFound("teacher not found")
	}

	return nil
}

// List retrieves all teachers (deprecated: use ListPaginated for better performance)
func (r *GormRepository) List(ctx context.Context) ([]*Teacher, error) {
	var models []TeacherModel

	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, errors.Internal("failed to list teachers: " + err.Error())
	}

	teachers := make([]*Teacher, len(models))
	for i, model := range models {
		teachers[i] = model.ToTeacher()
	}

	return teachers, nil
}

// ListPaginated retrieves teachers with pagination using efficient OFFSET/LIMIT query
// Performs two queries: COUNT for total, SELECT with LIMIT/OFFSET for data
func (r *GormRepository) ListPaginated(ctx context.Context, params pagination.Params) ([]*Teacher, int64, error) {
	var models []TeacherModel
	var totalCount int64

	// Count total records first for pagination metadata
	if err := r.db.WithContext(ctx).Model(&TeacherModel{}).Count(&totalCount).Error; err != nil {
		return nil, 0, errors.Internal("failed to count teachers: " + err.Error())
	}

	// Fetch paginated results with OFFSET and LIMIT
	if err := r.db.WithContext(ctx).
		Offset(params.Offset()).
		Limit(params.Limit).
		Order("created_at DESC"). // Most recent teachers first for better UX
		Find(&models).Error; err != nil {
		return nil, 0, errors.Internal("failed to list teachers: " + err.Error())
	}

	teachers := make([]*Teacher, len(models))
	for i, model := range models {
		teachers[i] = model.ToTeacher()
	}

	return teachers, totalCount, nil
}
