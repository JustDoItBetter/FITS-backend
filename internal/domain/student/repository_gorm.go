package student

import (
	"context"
	"time"

	"github.com/JustDoItBetter/FITS-backend/internal/common/errors"
	"github.com/JustDoItBetter/FITS-backend/internal/common/pagination"
	"gorm.io/gorm"
)

// StudentModel represents the GORM model for students table
type StudentModel struct {
	ID        string     `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	FirstName string     `gorm:"column:first_name;type:varchar(100);not null"`
	LastName  string     `gorm:"column:last_name;type:varchar(100);not null"`
	Email     string     `gorm:"column:email;type:varchar(255);uniqueIndex;not null"`
	TeacherID *string    `gorm:"column:teacher_id;type:uuid"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index"` // Soft delete support
}

// TableName specifies the table name for GORM
func (StudentModel) TableName() string {
	return "students"
}

// ToStudent converts StudentModel to Student domain entity
func (m *StudentModel) ToStudent() *Student {
	return &Student{
		UUID:      m.ID,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		Email:     m.Email,
		TeacherID: m.TeacherID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// FromStudent converts Student domain entity to StudentModel
func FromStudent(s *Student) *StudentModel {
	return &StudentModel{
		ID:        s.UUID,
		FirstName: s.FirstName,
		LastName:  s.LastName,
		Email:     s.Email,
		TeacherID: s.TeacherID,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

// GormRepository implements Repository interface using GORM
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new GORM-based student repository
func NewGormRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

// WithDB returns a new repository instance using the provided database connection
// This enables the repository to participate in transactions
func (r *GormRepository) WithDB(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

// Create adds a new student to the database
func (r *GormRepository) Create(ctx context.Context, student *Student) error {
	model := FromStudent(student)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if errors.IsUniqueViolation(err) {
			return errors.Conflict("student with this email already exists")
		}
		return errors.Internal("failed to create student: " + err.Error())
	}

	// Update student with DB-generated fields
	*student = *model.ToStudent()
	return nil
}

// GetByUUID retrieves a student by UUID
func (r *GormRepository) GetByUUID(ctx context.Context, uuid string) (*Student, error) {
	var model StudentModel

	if err := r.db.WithContext(ctx).Where("id = ?", uuid).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("student not found")
		}
		return nil, errors.Internal("failed to get student: " + err.Error())
	}

	return model.ToStudent(), nil
}

// Update updates an existing student
func (r *GormRepository) Update(ctx context.Context, student *Student) error {
	model := FromStudent(student)

	result := r.db.WithContext(ctx).
		Model(&StudentModel{}).
		Where("id = ?", student.UUID).
		Updates(map[string]interface{}{
			"first_name": model.FirstName,
			"last_name":  model.LastName,
			"email":      model.Email,
			"teacher_id": model.TeacherID,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		if errors.IsUniqueViolation(result.Error) {
			return errors.Conflict("student with this email already exists")
		}
		return errors.Internal("failed to update student: " + result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return errors.NotFound("student not found")
	}

	return nil
}

// Delete removes a student from the database
func (r *GormRepository) Delete(ctx context.Context, uuid string) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", uuid).
		Delete(&StudentModel{})

	if result.Error != nil {
		return errors.Internal("failed to delete student: " + result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return errors.NotFound("student not found")
	}

	return nil
}

// List retrieves all students (deprecated: use ListPaginated for better performance)
func (r *GormRepository) List(ctx context.Context) ([]*Student, error) {
	var models []StudentModel

	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, errors.Internal("failed to list students: " + err.Error())
	}

	students := make([]*Student, len(models))
	for i, model := range models {
		students[i] = model.ToStudent()
	}

	return students, nil
}

// ListPaginated retrieves students with pagination using efficient OFFSET/LIMIT query
// Performs two queries: COUNT for total, SELECT with LIMIT/OFFSET for data
func (r *GormRepository) ListPaginated(ctx context.Context, params pagination.Params) ([]*Student, int64, error) {
	var models []StudentModel
	var totalCount int64

	// Count total records first for pagination metadata
	if err := r.db.WithContext(ctx).Model(&StudentModel{}).Count(&totalCount).Error; err != nil {
		return nil, 0, errors.Internal("failed to count students: " + err.Error())
	}

	// Fetch paginated results with OFFSET and LIMIT
	if err := r.db.WithContext(ctx).
		Offset(params.Offset()).
		Limit(params.Limit).
		Order("created_at DESC"). // Most recent students first for better UX
		Find(&models).Error; err != nil {
		return nil, 0, errors.Internal("failed to list students: " + err.Error())
	}

	students := make([]*Student, len(models))
	for i, model := range models {
		students[i] = model.ToStudent()
	}

	return students, totalCount, nil
}
