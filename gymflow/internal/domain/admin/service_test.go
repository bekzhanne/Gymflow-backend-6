package admin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func TestGetDashboard_Success(t *testing.T) {
	db := setupTestDB()
	service := NewService(db)

	// Create service
	dashboard, err := service.GetDashboard()

	// Should not error even with empty database
	assert.NoError(t, err)
	assert.NotNil(t, dashboard)
	
	// Initial values should be zero
	assert.GreaterOrEqual(t, dashboard.TotalUsers, int64(0))
	assert.GreaterOrEqual(t, dashboard.TotalClasses, int64(0))
	assert.GreaterOrEqual(t, dashboard.TotalBookings, int64(0))
	assert.GreaterOrEqual(t, dashboard.TotalRevenue, 0.0)
}

func TestGetDashboard_EmptyDatabase(t *testing.T) {
	db := setupTestDB()
	service := NewService(db)

	dashboard, err := service.GetDashboard()

	assert.NoError(t, err)
	assert.NotNil(t, dashboard)
	assert.Equal(t, int64(0), dashboard.TotalUsers)
	assert.Equal(t, int64(0), dashboard.TotalClasses)
	assert.Equal(t, int64(0), dashboard.TotalBookings)
	assert.Equal(t, 0.0, dashboard.TotalRevenue)
}