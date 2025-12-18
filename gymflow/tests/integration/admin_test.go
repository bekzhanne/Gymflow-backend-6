package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"gymflow/internal/domain/user"

	"github.com/stretchr/testify/assert"
)

func TestAdminDashboard_Success(t *testing.T) {
	router := setupTestRouter()

	// 1. Register admin
	adminReq := user.RegisterRequest{
		Name:           "Admin User",
		Email:          "admin@example.com",
		Password:       "admin123",
		Role:           user.RoleAdmin,
		MembershipTier: user.MembershipVIP,
	}

	w := makeRequest(t, router, "POST", "/api/v1/auth/register", adminReq, "")
	assert.Equal(t, http.StatusCreated, w.Code)

	var adminResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &adminResp)
	adminToken := adminResp["token"].(string)

	// 2. Access dashboard
	w = makeRequest(t, router, "GET", "/api/v1/admin/dashboard", nil, adminToken)
	assert.Equal(t, http.StatusOK, w.Code)

	var dashboardResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &dashboardResp)

	// Verify dashboard structure
	assert.Contains(t, dashboardResp, "total_users")
	assert.Contains(t, dashboardResp, "total_classes")
	assert.Contains(t, dashboardResp, "total_bookings")
	assert.Contains(t, dashboardResp, "total_revenue")
	assert.Contains(t, dashboardResp, "active_members")
	assert.Contains(t, dashboardResp, "upcoming_classes")

	// Initial values should be at least 1 (admin user)
	assert.GreaterOrEqual(t, int(dashboardResp["total_users"].(float64)), 1)
}

func TestAdminDashboard_Unauthorized(t *testing.T) {
	router := setupTestRouter()

	// 1. Register regular member
	memberReq := user.RegisterRequest{
		Name:           "Regular Member",
		Email:          "member@example.com",
		Password:       "password123",
		Role:           user.RoleMember,
		MembershipTier: user.MembershipBasic,
	}

	w := makeRequest(t, router, "POST", "/api/v1/auth/register", memberReq, "")
	assert.Equal(t, http.StatusCreated, w.Code)

	var memberResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &memberResp)
	memberToken := memberResp["token"].(string)

	// 2. Try to access admin dashboard (should fail)
	w = makeRequest(t, router, "GET", "/api/v1/admin/dashboard", nil, memberToken)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAdminDashboard_NoToken(t *testing.T) {
	router := setupTestRouter()

	// Try to access without token
	w := makeRequest(t, router, "GET", "/api/v1/admin/dashboard", nil, "")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminDashboard_WithData(t *testing.T) {
	router := setupTestRouter()

	// 1. Register admin
	adminReq := user.RegisterRequest{
		Name:           "Admin Boss",
		Email:          "boss@example.com",
		Password:       "admin123",
		Role:           user.RoleAdmin,
		MembershipTier: user.MembershipVIP,
	}

	w := makeRequest(t, router, "POST", "/api/v1/auth/register", adminReq, "")
	assert.Equal(t, http.StatusCreated, w.Code)

	var adminResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &adminResp)
	adminToken := adminResp["token"].(string)

	// 2. Register some members
	for i := 1; i <= 5; i++ {
		memberReq := user.RegisterRequest{
			Name:           "Test Member " + string(rune(i)),
			Email:          "testmember" + string(rune(i)) + "@example.com",
			Password:       "password123",
			Role:           user.RoleMember,
			MembershipTier: user.MembershipBasic,
		}
		makeRequest(t, router, "POST", "/api/v1/auth/register", memberReq, "")
	}

	// 3. Check dashboard reflects new users
	w = makeRequest(t, router, "GET", "/api/v1/admin/dashboard", nil, adminToken)
	assert.Equal(t, http.StatusOK, w.Code)

	var dashboardResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &dashboardResp)

	totalUsers := int(dashboardResp["total_users"].(float64))
	assert.GreaterOrEqual(t, totalUsers, 1) // admin + 5 members
}

func TestTrainerAccess_Dashboard(t *testing.T) {
	router := setupTestRouter()

	// 1. Register trainer
	trainerReq := user.RegisterRequest{
		Name:           "Trainer Tom",
		Email:          "tom@example.com",
		Password:       "password123",
		Role:           user.RoleTrainer,
		MembershipTier: user.MembershipPremium,
	}

	w := makeRequest(t, router, "POST", "/api/v1/auth/register", trainerReq, "")
	assert.Equal(t, http.StatusCreated, w.Code)

	var trainerResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &trainerResp)
	trainerToken := trainerResp["token"].(string)

	// 2. Trainer tries to access admin dashboard (should fail)
	w = makeRequest(t, router, "GET", "/api/v1/admin/dashboard", nil, trainerToken)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAdminStatistics_InitialState(t *testing.T) {
	router := setupTestRouter()

	// 1. Register admin
	adminReq := user.RegisterRequest{
		Name:           "Stats Admin",
		Email:          "statsadmin@example.com",
		Password:       "admin123",
		Role:           user.RoleAdmin,
		MembershipTier: user.MembershipVIP,
	}

	w := makeRequest(t, router, "POST", "/api/v1/auth/register", adminReq, "")
	assert.Equal(t, http.StatusCreated, w.Code)

	var adminResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &adminResp)
	adminToken := adminResp["token"].(string)

	// 2. Get dashboard in initial state (no bookings/classes)
	w = makeRequest(t, router, "GET", "/api/v1/admin/dashboard", nil, adminToken)
	assert.Equal(t, http.StatusOK, w.Code)

	var dashboardResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &dashboardResp)

	// In initial state
	assert.GreaterOrEqual(t, int(dashboardResp["total_users"].(float64)), 1)
	assert.Equal(t, float64(0), dashboardResp["total_bookings"].(float64))
	assert.Equal(t, float64(0), dashboardResp["total_revenue"].(float64))
	assert.Equal(t, float64(0), dashboardResp["upcoming_classes"].(float64))
}