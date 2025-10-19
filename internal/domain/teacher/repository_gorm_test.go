package teacher

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestTeacherModel is a simplified model for SQLite testing
type TestTeacherModel struct {
	ID         string `gorm:"column:id;primaryKey"`
	FirstName  string `gorm:"column:first_name;type:varchar(100);not null"`
	LastName   string `gorm:"column:last_name;type:varchar(100);not null"`
	Email      string `gorm:"column:email;type:varchar(255);uniqueIndex;not null"`
	Department string `gorm:"column:department;type:varchar(100);not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (TestTeacherModel) TableName() string {
	return "teachers"
}

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "failed to create test database")

	// Use TestTeacherModel for SQLite compatibility
	err = db.AutoMigrate(&TestTeacherModel{})
	require.NoError(t, err, "failed to migrate test database")

	return db
}

// TestGormRepository_Create tests the Create method
func TestGormRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormRepository(db)
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		teacher := &Teacher{
			UUID:       "550e8400-e29b-41d4-a716-446655440100",
			FirstName:  "Anna",
			LastName:   "Schmidt",
			Email:      "anna@test.com",
			Department: "Computer Science",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		err := repo.Create(ctx, teacher)
		assert.NoError(t, err)
		assert.NotEmpty(t, teacher.UUID)
	})

	t.Run("duplicate email error", func(t *testing.T) {
		// Create first teacher
		teacher1 := &Teacher{
			UUID:       "550e8400-e29b-41d4-a716-446655440101",
			FirstName:  "Thomas",
			LastName:   "Müller",
			Email:      "duplicate@test.com",
			Department: "Mathematics",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		err := repo.Create(ctx, teacher1)
		require.NoError(t, err)

		// Try to create second teacher with same email
		teacher2 := &Teacher{
			UUID:       "550e8400-e29b-41d4-a716-446655440102",
			FirstName:  "Maria",
			LastName:   "Weber",
			Email:      "duplicate@test.com",
			Department: "Physics",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		err = repo.Create(ctx, teacher2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

// TestGormRepository_GetByUUID tests the GetByUUID method
func TestGormRepository_GetByUUID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormRepository(db)
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		// Create a teacher first
		teacher := &Teacher{
			UUID:       "550e8400-e29b-41d4-a716-446655440200",
			FirstName:  "Lisa",
			LastName:   "Weber",
			Email:      "lisa@test.com",
			Department: "Chemistry",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		err := repo.Create(ctx, teacher)
		require.NoError(t, err)

		// Retrieve the teacher
		retrieved, err := repo.GetByUUID(ctx, teacher.UUID)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, teacher.UUID, retrieved.UUID)
		assert.Equal(t, teacher.FirstName, retrieved.FirstName)
		assert.Equal(t, teacher.LastName, retrieved.LastName)
		assert.Equal(t, teacher.Email, retrieved.Email)
		assert.Equal(t, teacher.Department, retrieved.Department)
	})

	t.Run("teacher not found", func(t *testing.T) {
		retrieved, err := repo.GetByUUID(ctx, "nonexistent-uuid")
		assert.Error(t, err)
		assert.Nil(t, retrieved)
		assert.Contains(t, err.Error(), "not found")
	})
}

// TestGormRepository_Update tests the Update method
func TestGormRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormRepository(db)
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		// Create a teacher first
		teacher := &Teacher{
			UUID:       "550e8400-e29b-41d4-a716-446655440300",
			FirstName:  "John",
			LastName:   "Doe",
			Email:      "john@test.com",
			Department: "Biology",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		err := repo.Create(ctx, teacher)
		require.NoError(t, err)

		// Update the teacher
		teacher.FirstName = "Jane"
		teacher.Email = "jane@test.com"
		teacher.Department = "Physics"
		err = repo.Update(ctx, teacher)
		assert.NoError(t, err)

		// Verify the update
		retrieved, err := repo.GetByUUID(ctx, teacher.UUID)
		assert.NoError(t, err)
		assert.Equal(t, "Jane", retrieved.FirstName)
		assert.Equal(t, "jane@test.com", retrieved.Email)
		assert.Equal(t, "Physics", retrieved.Department)
		assert.Equal(t, "Doe", retrieved.LastName) // Unchanged
	})

	t.Run("update nonexistent teacher", func(t *testing.T) {
		teacher := &Teacher{
			UUID:       "nonexistent-uuid",
			FirstName:  "Ghost",
			LastName:   "User",
			Email:      "ghost@test.com",
			Department: "Unknown",
		}
		err := repo.Update(ctx, teacher)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("update with duplicate email", func(t *testing.T) {
		// Create two teachers
		teacher1 := &Teacher{
			UUID:       "550e8400-e29b-41d4-a716-446655440301",
			FirstName:  "Teacher",
			LastName:   "One",
			Email:      "teacher1@test.com",
			Department: "Math",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		err := repo.Create(ctx, teacher1)
		require.NoError(t, err)

		teacher2 := &Teacher{
			UUID:       "550e8400-e29b-41d4-a716-446655440302",
			FirstName:  "Teacher",
			LastName:   "Two",
			Email:      "teacher2@test.com",
			Department: "Science",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		err = repo.Create(ctx, teacher2)
		require.NoError(t, err)

		// Try to update teacher2 with teacher1's email
		teacher2.Email = "teacher1@test.com"
		err = repo.Update(ctx, teacher2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

// TestGormRepository_Delete tests the Delete method
func TestGormRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormRepository(db)
	ctx := context.Background()

	t.Run("successful delete", func(t *testing.T) {
		// Create a teacher first
		teacher := &Teacher{
			UUID:       "550e8400-e29b-41d4-a716-446655440400",
			FirstName:  "Delete",
			LastName:   "Me",
			Email:      "delete@test.com",
			Department: "Temporary",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		err := repo.Create(ctx, teacher)
		require.NoError(t, err)

		// Delete the teacher
		err = repo.Delete(ctx, teacher.UUID)
		assert.NoError(t, err)

		// Verify deletion
		retrieved, err := repo.GetByUUID(ctx, teacher.UUID)
		assert.Error(t, err)
		assert.Nil(t, retrieved)
	})

	t.Run("delete nonexistent teacher", func(t *testing.T) {
		err := repo.Delete(ctx, "nonexistent-uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

// TestGormRepository_List tests the List method
func TestGormRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormRepository(db)
	ctx := context.Background()

	t.Run("list multiple teachers", func(t *testing.T) {
		// Create multiple teachers
		teachers := []*Teacher{
			{
				UUID:       "550e8400-e29b-41d4-a716-446655440500",
				FirstName:  "Alice",
				LastName:   "Smith",
				Email:      "alice@test.com",
				Department: "Mathematics",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			{
				UUID:       "550e8400-e29b-41d4-a716-446655440501",
				FirstName:  "Bob",
				LastName:   "Johnson",
				Email:      "bob@test.com",
				Department: "Physics",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			{
				UUID:       "550e8400-e29b-41d4-a716-446655440502",
				FirstName:  "Charlie",
				LastName:   "Brown",
				Email:      "charlie@test.com",
				Department: "Chemistry",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		}

		for _, teacher := range teachers {
			err := repo.Create(ctx, teacher)
			require.NoError(t, err)
		}

		// List all teachers
		retrieved, err := repo.List(ctx)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(retrieved), 3) // At least the 3 we created
	})

	t.Run("list empty", func(t *testing.T) {
		// Use a fresh database
		freshDB := setupTestDB(t)
		freshRepo := NewGormRepository(freshDB)

		retrieved, err := freshRepo.List(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Len(t, retrieved, 0)
	})
}

// TestGormRepository_DepartmentUpdate tests department field updates
func TestGormRepository_DepartmentUpdate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormRepository(db)
	ctx := context.Background()

	t.Run("update department", func(t *testing.T) {
		teacher := &Teacher{
			UUID:       "550e8400-e29b-41d4-a716-446655440600",
			FirstName:  "Department",
			LastName:   "Changer",
			Email:      "dept.change@test.com",
			Department: "Computer Science",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		err := repo.Create(ctx, teacher)
		require.NoError(t, err)

		// Update department
		teacher.Department = "Data Science"
		err = repo.Update(ctx, teacher)
		assert.NoError(t, err)

		// Verify department change
		retrieved, err := repo.GetByUUID(ctx, teacher.UUID)
		assert.NoError(t, err)
		assert.Equal(t, "Data Science", retrieved.Department)
	})
}
