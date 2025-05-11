package db

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/songphuc19102004/social/internal/store"
)

type faker struct {
	username, firstName, lastName string
	age                           int
}

var usernames = []string{
	"alice",
	"bob",
	"charlie",
	"diana",
	"evan",
	"fiona",
	"gabriel",
	"hannah",
	"isaac",
	"jade",
}

var firstNames = []string{
	"alex",
	"ben",
	"casey",
	"dana",
	"eric",
	"freya",
	"greg",
	"helen",
	"ian",
	"julia",
}

var lastNames = []string{
	"smith",
	"jones",
	"garcia",
	"chen",
	"patel",
	"kim",
	"lee",
	"brown",
	"miller",
	"wilson",
}

func Seed(store store.Storage) error {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, &user); err != nil {
			log.Printf("Error seeding user: %v with err: %v", user, err)
			return err
		}
	}

	return nil
}

// posts and comments later, im too lazy

func generateUsers(quantity int) []store.User {
	users := make([]store.User, quantity)

	for i := range quantity {
		faker := randomData()
		users[i] = store.User{
			Username:  faker.username,
			FirstName: faker.firstName,
			LastName:  faker.lastName,
			Age:       faker.age,
		}
	}

	return users
}

func randomData() *faker {
	randomSalt := rand.IntN(1000) + 100
	randomIndex := rand.IntN(11)
	randomUsername := usernames[randomIndex]
	age := rand.IntN(30-18+1) + 18

	return &faker{
		username:  fmt.Sprintf("%s%d", randomUsername, randomSalt),
		firstName: firstNames[randomIndex],
		lastName:  lastNames[randomIndex],
		age:       age,
	}
}
