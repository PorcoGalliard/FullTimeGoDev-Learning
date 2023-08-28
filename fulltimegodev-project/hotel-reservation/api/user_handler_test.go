package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fulltimegodev/hotel-reservation-nana/db"
	"github.com/fulltimegodev/hotel-reservation-nana/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http/httptest"
	"testing"
)

const dbnametest = "hotel-reservation-test"
const testmongouri = "mongodb://localhost:27017"

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testmongouri))
	if err != nil {
		log.Fatal(err)
	}

	return &testdb{
		UserStore: db.NewMongoUserStore(client, dbnametest),
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()

	userHandler := NewUserHandler(tdb)

	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParam{
		FirstName: "Apollo Norm",
		LastName:  "AAA-0003",
		Email:     "apollonorm@uncf.org",
		Password:  "advancedarmamentartillery",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(req)
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)

	if user.FirstName != params.FirstName {
		t.Errorf("expected first name %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected last name %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}
	if len(user.ID) == 0 {
		t.Errorf("expecting a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("Encrypted password must not shown in json")
	}
}

func TestGetUserById(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb)
	app.Get("/user/:id", userHandler.HandleGetUser)

	expectedUser := types.User{
		ID:                primitive.NewObjectID(),
		FirstName:         "Apollo Norm",
		LastName:          "AAA0003",
		Email:             "apollonrom@uncf.org",
		EncryptedPassword: "advancedartilleryarmament",
	}

	insertedUser, err := tdb.InsertUser(context.TODO(), &expectedUser)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("/user/%s", insertedUser.ID.Hex()), nil)
	req.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(req)

	var retrievedUser types.User
	json.NewDecoder(resp.Body).Decode(&retrievedUser)

	if retrievedUser.ID.Hex() != insertedUser.ID.Hex() {
		t.Error("The ID is not match")
	}

	if retrievedUser.FirstName != insertedUser.FirstName {
		t.Errorf("Expecting first name %s but got %s", retrievedUser.FirstName, insertedUser.FirstName)
	}

	if retrievedUser.LastName != insertedUser.LastName {
		t.Errorf("Expecting last name %s but got %s", retrievedUser.LastName, insertedUser.LastName)
	}

	if retrievedUser.Email != insertedUser.Email {
		t.Errorf("Expecting email %s but got %s", retrievedUser.Email, retrievedUser.Email)
	}
}

func TestGetUsers(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb)
	app.Get("/user", userHandler.HandleGetUsers)

	expectedUser := []types.User{
		{
			ID:                primitive.NewObjectID(),
			FirstName:         "Apollo",
			LastName:          "Norm",
			Email:             "apollonorm@uncf.org",
			EncryptedPassword: "advancedartilleryarmament",
		},
		{
			ID:                primitive.NewObjectID(),
			FirstName:         "Antares",
			LastName:          "AAA-0005",
			Email:             "antares@uncf.org",
			EncryptedPassword: "advancedartilleryarmament",
		},
	}

	for _, user := range expectedUser {
		_, err := tdb.InsertUser(context.TODO(), &user)
		if err != nil {
			t.Fatal(err)
		}
	}

	req := httptest.NewRequest("GET", "/user", nil)
	req.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(req)

	var retrievedUsers []types.User
	json.NewDecoder(resp.Body).Decode(&retrievedUsers)

	if len(retrievedUsers) != len(expectedUser) {
		t.Errorf("Expecting number of users is %d but only got %d", len(expectedUser), len(retrievedUsers))
	}

	for i, expected := range expectedUser {
		if retrievedUsers[i].ID.Hex() != expected.ID.Hex() {
			t.Errorf("expecting ID %s but got %s", expected.ID.Hex(), retrievedUsers[i].ID.Hex())
		}

		if retrievedUsers[i].FirstName != expected.FirstName {
			t.Errorf("expecting first name %s but got %s", expected.FirstName, retrievedUsers[i].FirstName)
		}

		if retrievedUsers[i].LastName != expected.LastName {
			t.Errorf("expecting first name %s but got %s", expected.LastName, retrievedUsers[i].LastName)
		}

		if retrievedUsers[i].Email != expected.Email {
			t.Errorf("expecting first name %s but got %s", expected.Email, retrievedUsers[i].Email)
		}
	}
}