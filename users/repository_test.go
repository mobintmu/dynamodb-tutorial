package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAddItem(t *testing.T) {

	userID := "5e87ea75-aa99-4517-b4ac-a464fd5b3122"
	fmt.Println(userID)
	tgID := "12121212"

	user := User{
		Uid:            userID,
		TgId:           tgID,
		Referrer:       uuid.New().String(),
		Name:           "Mobin",
		TgFirstName:    "Mobin",
		TgLastName:     "SH",
		TgUsername:     "m_sh",
		ProfilePicture: "https://picture.com",
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	db := New()
	db.AddItem(user)

}

func TestGetItemByUserID(t *testing.T) {
	userID := "5e87ea75-aa99-4517-b4ac-a464fd5b3122"
	db := New()
	user, err := db.GetByUserID(userID)
	if err != nil {
		t.Fatal(err)
	}
	if user.Uid != userID {
		t.Fatal("User ID did not match")
	}
}

func TestGetItemByTgID(t *testing.T) {
	tgID := "12121212"
	db := New()
	users, err := db.GetByTgID(tgID)
	if err != nil {
		t.Fatal(err)
	}

	if len(users) == 0 {
		t.Fatal("No user found")
	}
}
