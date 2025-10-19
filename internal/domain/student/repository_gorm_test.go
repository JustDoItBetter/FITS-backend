package student

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestStudentModel is a simplified model for SQLite testing
type TestStudentModel struct {
	ID        string  `gorm:"column:id;primaryKey"`
	FirstName string  `gorm:"column:first_name;type:varchar(100);not null"`
	LastName  string  `gorm:"column:last_name;type:varchar(100);not null"`
	Email     string  `gorm:"column:email;type:varchar(255);uniqueIndex;not null"`
	TeacherID *string `gorm:"column:teacher_id;type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (TestStudentModel) TableName() string {
	return "students"
}

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "failed to create test database")

	// Use TestStudentModel for SQLite compatibility
	err = db.AutoMigrate(&TestStudentModel{})
	require.NoError(t, err, "failed to migrate test database")

	return db
}

// TestGormRepository_Create tests the Create method
func TestGormRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormRepository(db)
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		student := &Student{
			UUID:      "550e8400-e29b-41d4-a716-446655440100",
			FirstName: "Max",
			LastName:  "Mustermann",
			Email:     "max@test.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, student)
		assert.NoError(t, err)
		assert.NotEmpty(t, student.UUID)
	})

	t.Run("duplicate email error", func(t *testing.T) {
		// Create first student
		student1 := &Student{
			UUID:      "550e8400-e29b-41d4-a716-446655440101",
			FirstName: "Anna",
			LastName:  "Schmidt",
			Email:     "duplicate@test.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, student1)
		require.NoError(t, err)

		// Try to create second student with same email
		student2 := &Student{
			UUID:      "550e8400-e29b-41d4-a716-446655440102",
			FirstName: "Thomas",
			LastName:  "Müller",
			Email:     "duplicate@test.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = repo.Create(ctx, student2)
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
		// Create a student first
		student := &Student{
			UUID:      "550e8400-e29b-41d4-a716-446655440200",
			FirstName: "Lisa",
			LastName:  "Weber",
			Email:     "lisa@test.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, student)
		require.NoError(t, err)

		// Retrieve the student
		retrieved, err := repo.GetByUUID(ctx, student.UUID)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, student.UUID, retrieved.UUID)
		assert.Equal(t, student.FirstName, retrieved.FirstName)
		assert.Equal(t, student.LastName, retrieved.LastName)
		assert.Equal(t, student.Email, retrieved.Email)
	})

	t.Run("student not found", func(t *testing.T) {
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
		// Create a student first
		student := &Student{
			UUID:      "550e8400-e29b-41d4-a716-446655440300",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@test.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, student)
		require.NoError(t, err)

		// Update the student
		student.FirstName = "Jane"
		student.Email = "jane@test.com"
		err = repo.Update(ctx, student)
		assert.NoError(t, err)

		// Verify the update
		retrieved, err := repo.GetByUUID(ctx, student.UUID)
		assert.NoError(t, err)
		assert.Equal(t, "Jane", retrieved.FirstName)
		assert.Equal(t, "jane@test.com", retrieved.Email)
		assert.Equal(t, "Doe", retrieved.LastName) // Unchanged
	})

	t.Run("update nonexistent student", func(t *testing.T) {
		student := &Student{
			UUID:      "nonexistent-uuid",
			FirstName: "Ghost",
			LastName:  "User",
			Email:     "ghost@test.com",
		}
		err := repo.Update(ctx, student)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("update with duplicate email", func(t *testing.T) {
		// Create two students
		student1 := &Student{
			UUID:      "550e8400-e29b-41d4-a716-446655440301",
			FirstName: "User",
			LastName:  "One",
			Email:     "user1@test.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, student1)
		require.NoError(t, err)

		student2 := &Student{
			UUID:      "550e8400-e29b-41d4-a716-446655440302",
			FirstName: "User",
			LastName:  "Two",
			Email:     "user2@test.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = repo.Create(ctx, student2)
		require.NoError(t, err)

		// Try to update student2 with student1's email
		student2.Email = "user1@test.com"
		err = repo.Update(ctx, student2)
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
		// Create a student first
		student := &Student{
			UUID:      "550e8400-e29b-41d4-a716-446655440400",
			FirstName: "Delete",
			LastName:  "Me",
			Email:     "delete@test.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, student)
		require.NoError(t, err)

		// Delete the student
		err = repo.Delete(ctx, student.UUID)
		assert.NoError(t, err)

		// Verify deletion
		retrieved, err := repo.GetByUUID(ctx, student.UUID)
		assert.Error(t, err)
		assert.Nil(t, retrieved)
	})

	t.Run("delete nonexistent student", func(t *testing.T) {
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

	t.Run("list multiple students", func(t *testing.T) {
		// Create multiple students
		students := []*Student{
			{
				UUID:      "550e8400-e29b-41d4-a716-446655440500",
				FirstName: "Alice",
				LastName:  "Smith",
				Email:     "alice@test.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				UUID:      "550e8400-e29b-41d4-a716-446655440501",
				FirstName: "Bob",
				LastName:  "Johnson",
				Email:     "bob@test.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				UUID:      "550e8400-e29b-41d4-a716-446655440502",
				FirstName: "Charlie",
				LastName:  "Brown",
				Email:     "charlie@test.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		for _, s := range students {
			err := repo.Create(ctx, s)
			require.NoError(t, err)
		}

		// List all students
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

// TestGormRepository_WithTeacher tests student creation with teacher reference
func TestGormRepository_WithTeacher(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormRepository(db)
	ctx := context.Background()

	t.Run("create student with teacher ID", func(t *testing.T) {
		teacherID := "550e8400-e29b-41d4-a716-446655440001"
		student := &Student{
			UUID:      "550e8400-e29b-41d4-a716-446655440600",
			FirstName: "Student",
			LastName:  "WithTeacher",
			Email:     "student.teacher@test.com",
			TeacherID: &teacherID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, student)
		assert.NoError(t, err)

		// Verify teacher ID is stored
		retrieved, err := repo.GetByUUID(ctx, student.UUID)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved.TeacherID)
		assert.Equal(t, teacherID, *retrieved.TeacherID)
	})

	t.Run("update student to remove teacher", func(t *testing.T) {
		teacherID := "550e8400-e29b-41d4-a716-446655440002"
		student := &Student{
			UUID:      "550e8400-e29b-41d4-a716-446655440601",
			FirstName: "Student",
			LastName:  "RemoveTeacher",
			Email:     "remove.teacher@test.com",
			TeacherID: &teacherID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, student)
		require.NoError(t, err)

		// Remove teacher
		student.TeacherID = nil
		err = repo.Update(ctx, student)
		assert.NoError(t, err)

		// Verify teacher is removed
		retrieved, err := repo.GetByUUID(ctx, student.UUID)
		assert.NoError(t, err)
		assert.Nil(t, retrieved.TeacherID)
	})
}
