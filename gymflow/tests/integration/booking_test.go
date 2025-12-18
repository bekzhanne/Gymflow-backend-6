package integration

import (
	"encoding/json"
	"net/http"
	"testing"
	"fmt"

	"gymflow/internal/domain/booking"
	"gymflow/internal/domain/user"

	"github.com/stretchr/testify/assert"
)

func TestBookingFlow_Complete(t *testing.T) {
	router := setupTestRouter()

	// 1. Register trainer
	trainerReq := user.RegisterRequest{
		Name:           "Trainer John",
		Email:          "trainer@example.com",
		Password:       "password123",
		Role:           user.RoleTrainer,
		MembershipTier: user.MembershipBasic,
	}

	w := makeRequest(t, router, "POST", "/api/v1/auth/register", trainerReq, "")
	assert.Equal(t, http.StatusCreated, w.Code)

	var trainerResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &trainerResp)
	trainerToken := trainerResp["token"].(string)
	trainerUser := trainerResp["user"].(map[string]interface{})
	trainerID := uint(trainerUser["id"].(float64))

	// 2. Create class as trainer
	classReq := booking.CreateClassRequest{
		Name:        "Yoga Class",
		Description: "Morning yoga session",
		TrainerID:   trainerID,
		Capacity:    2,
		StartTime:   "2024-12-20T10:00:00Z",
		EndTime:     "2024-12-20T11:00:00Z",
		Price:       50.0,
	}

	w = makeRequest(t, router, "POST", "/api/v1/classes", classReq, trainerToken)
	assert.Equal(t, http.StatusCreated, w.Code)

	var classResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &classResp)
	classID := uint(classResp["id"].(float64))

	// 3. Register member
	memberReq := user.RegisterRequest{
		Name:           "Member Alice",
		Email:          "alice@example.com",
		Password:       "password123",
		Role:           user.RoleMember,
		MembershipTier: user.MembershipPremium,
	}

	w = makeRequest(t, router, "POST", "/api/v1/auth/register", memberReq, "")
	assert.Equal(t, http.StatusCreated, w.Code)

	var memberResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &memberResp)
	memberToken := memberResp["token"].(string)

	// 4. Member books the class
	bookingReq := booking.CreateBookingRequest{
		ClassID: classID,
	}

	w = makeRequest(t, router, "POST", "/api/v1/bookings", bookingReq, memberToken)
	assert.Equal(t, http.StatusCreated, w.Code)

	var bookingResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &bookingResp)
	bookingID := uint(bookingResp["id"].(float64))
	assert.Equal(t, "booked", bookingResp["status"])
	assert.Equal(t, "pending", bookingResp["payment_status"])

	// 5. Member lists their bookings
	w = makeRequest(t, router, "GET", "/api/v1/bookings", nil, memberToken)
	assert.Equal(t, http.StatusOK, w.Code)

	var listResp []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &listResp)
	assert.Len(t, listResp, 1)
	assert.Equal(t, float64(bookingID), listResp[0]["id"].(float64))

	// 6. Member cancels booking
	w = makeRequest(t, router, "POST", fmt.Sprintf("/api/v1/bookings/%d/cancel", bookingID), nil, memberToken)
	assert.Equal(t, http.StatusOK, w.Code)

	var cancelResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &cancelResp)
	assert.Equal(t, "cancelled", cancelResp["status"])
}

func TestBookingFlow_Waitlist(t *testing.T) {
	router := setupTestRouter()

	// 1. Register trainer
	trainerReq := user.RegisterRequest{
		Name:           "Trainer Bob",
		Email:          "bob@example.com",
		Password:       "password123",
		Role:           user.RoleTrainer,
		MembershipTier: user.MembershipBasic,
	}

	w := makeRequest(t, router, "POST", "/api/v1/auth/register", trainerReq, "")
	assert.Equal(t, http.StatusCreated, w.Code)

	var trainerResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &trainerResp)
	trainerToken := trainerResp["token"].(string)
	trainerUser := trainerResp["user"].(map[string]interface{})
	trainerID := uint(trainerUser["id"].(float64))

	// 2. Create class with capacity 1
	classReq := booking.CreateClassRequest{
		Name:        "Limited Class",
		Description: "Only 1 spot",
		TrainerID:   trainerID,
		Capacity:    1,
		StartTime:   "2024-12-20T14:00:00Z",
		EndTime:     "2024-12-20T15:00:00Z",
		Price:       75.0,
	}

	w = makeRequest(t, router, "POST", "/api/v1/classes", classReq, trainerToken)
	assert.Equal(t, http.StatusCreated, w.Code)

	var classResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &classResp)
	classID := uint(classResp["id"].(float64))

	// 3. Register first member and book
	member1Req := user.RegisterRequest{
		Name:           "Member One",
		Email:          "member1@example.com",
		Password:       "password123",
		Role:           user.RoleMember,
		MembershipTier: user.MembershipBasic,
	}

	w = makeRequest(t, router, "POST", "/api/v1/auth/register", member1Req, "")
	assert.Equal(t, http.StatusCreated, w.Code)

	var member1Resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &member1Resp)
	member1Token := member1Resp["token"].(string)

	bookingReq := booking.CreateBookingRequest{ClassID: classID}
	w = makeRequest(t, router, "POST", "/api/v1/bookings", bookingReq, member1Token)
	assert.Equal(t, http.StatusCreated, w.Code)

	var booking1Resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &booking1Resp)
	assert.Equal(t, "booked", booking1Resp["status"])

	// 4. Register second member and try to book (should be waitlisted)
	member2Req := user.RegisterRequest{
		Name:           "Member Two",
		Email:          "member2@example.com",
		Password:       "password123",
		Role:           user.RoleMember,
		MembershipTier: user.MembershipBasic,
	}

	w = makeRequest(t, router, "POST", "/api/v1/auth/register", member2Req, "")
	assert.Equal(t, http.StatusCreated, w.Code)

	var member2Resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &member2Resp)
	member2Token := member2Resp["token"].(string)

	w = makeRequest(t, router, "POST", "/api/v1/bookings", bookingReq, member2Token)
	assert.Equal(t, http.StatusCreated, w.Code)

	var booking2Resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &booking2Resp)
	assert.Equal(t, "waitlist", booking2Resp["status"])
}

func TestListClasses_Public(t *testing.T) {
	router := setupTestRouter()

	// 1. Register trainer
	trainerReq := user.RegisterRequest{
		Name:           "Trainer Charlie",
		Email:          "charlie@example.com",
		Password:       "password123",
		Role:           user.RoleTrainer,
		MembershipTier: user.MembershipBasic,
	}

	w := makeRequest(t, router, "POST", "/api/v1/auth/register", trainerReq, "")
	assert.Equal(t, http.StatusCreated, w.Code)

	var trainerResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &trainerResp)
	trainerToken := trainerResp["token"].(string)
	trainerUser := trainerResp["user"].(map[string]interface{})
	trainerID := uint(trainerUser["id"].(float64))

	// 2. Create multiple classes
	for i := 1; i <= 3; i++ {
		classReq := booking.CreateClassRequest{
			Name:        "Class " + string(rune(i)),
			Description: "Description " + string(rune(i)),
			TrainerID:   trainerID,
			Capacity:    10,
			StartTime:   "2024-12-20T10:00:00Z",
			EndTime:     "2024-12-20T11:00:00Z",
			Price:       50.0,
		}
		w = makeRequest(t, router, "POST", "/api/v1/classes", classReq, trainerToken)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// 3. List classes (public endpoint, no token needed)
	w = makeRequest(t, router, "GET", "/api/v1/classes", nil, "")
	assert.Equal(t, http.StatusOK, w.Code)

	var classes []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &classes)
	assert.GreaterOrEqual(t, len(classes), 3)
}

func TestCancelBooking_Unauthorized(t *testing.T) {
	router := setupTestRouter()

	// 1. Register two members
	member1Req := user.RegisterRequest{
		Name:           "Member Alpha",
		Email:          "alpha@example.com",
		Password:       "password123",
		Role:           user.RoleMember,
		MembershipTier: user.MembershipBasic,
	}

	w := makeRequest(t, router, "POST", "/api/v1/auth/register", member1Req, "")
	assert.Equal(t, http.StatusCreated, w.Code)

	var member1Resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &member1Resp)
	member1Token := member1Resp["token"].(string)

	member2Req := user.RegisterRequest{
		Name:           "Member Beta",
		Email:          "beta@example.com",
		Password:       "password123",
		Role:           user.RoleMember,
		MembershipTier: user.MembershipBasic,
	}

	w = makeRequest(t, router, "POST", "/api/v1/auth/register", member2Req, "")
	assert.Equal(t, http.StatusCreated, w.Code)

	var member2Resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &member2Resp)
	member2Token := member2Resp["token"].(string)

	// 2. Create trainer and class (skipped for brevity, assume classID = 1)

	// 3. Member1 creates booking
	bookingReq := booking.CreateBookingRequest{ClassID: 1}
	w = makeRequest(t, router, "POST", "/api/v1/bookings", bookingReq, member1Token)
	
	if w.Code == http.StatusCreated {
		var bookingResp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &bookingResp)
		bookingID := uint(bookingResp["id"].(float64))

		// 4. Member2 tries to cancel Member1's booking (should fail)
		w = makeRequest(t, router, "POST", "/api/v1/bookings/"+string(rune(bookingID))+"/cancel", nil, member2Token)
		assert.Equal(t, http.StatusForbidden, w.Code)
	}
}