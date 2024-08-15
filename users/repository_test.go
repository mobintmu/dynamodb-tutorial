package repository

import (
	"fmt"
	"testing"
	"time"
)

func TestGetTableListIfNotExistCreateTable(t *testing.T) {
	db := New()
	tableName := "users"
	result := db.TableIsExist(tableName)
	if !result {
		db.CreateTable()
	}
}

func TestAddItem(t *testing.T) {

	userID := "5e87ea75-aa99-4517-b4ac-a464fd5b3122"
	tgID := "12121212"
	referrerID := "52aaa63e-e2b7-42f5-853d-ffa03a15e900"

	user := User{
		Uid:             userID,
		TgId:            tgID,
		Referrer:        referrerID,
		Name:            "Mobin",
		TgFirstName:     "Mobin TG",
		TgLastName:      "SH",
		TgUsername:      "m_sh",
		ProfilePicture:  "https://picture.com",
		CreatedAt:       time.Now().UTC().UnixNano(),
		UpdatedAt:       time.Now().UTC().UnixNano(),
		CounterReferrer: 0,
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

func TestUpdateItem(t *testing.T) {

	fmt.Println(time.Now().UTC().UnixNano())

	userID := "5e87ea75-aa99-4517-b4ac-a464fd5b3122"
	db := New()
	userFirst, err := db.GetByUserID(userID)
	if err != nil {
		t.Fatal(err)
	}

	db.IncrementUserReferrer(userID)
	userAgain, err := db.GetByUserID(userID)
	if err != nil {
		t.Fatal(err)
	}

	if userAgain.CounterReferrer != (userFirst.CounterReferrer + 1) {
		fmt.Println(userFirst.CounterReferrer)
		fmt.Println(userAgain.CounterReferrer)
		t.Fatal("Counter did not increment")
	}
}

func TestCount(t *testing.T) {
	db := New()
	count, err := db.CountAllData()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(count)
}

func TestGetByReferrer(t *testing.T) {
	referrer := "52aaa63e-e2b7-42f5-853d-ffa03a15e900"
	db := New()
	users, err := db.GetByReferrer(referrer)
	if err != nil {
		t.Fatal(err)
	}

	if len(users) == 0 {
		t.Fatal("No user found")
	}
}
