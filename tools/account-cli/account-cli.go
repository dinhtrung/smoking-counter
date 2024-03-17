package main

import (
	"flag"
	"github.com/dinhtrung/smoking-counter/internal/app/api-gateway/services/impl"
	"github.com/dinhtrung/smoking-counter/pkg/fiber/shared"
	"log"
	"time"
)

var configFile, username, email, password, authorities, operator string

func init() {
	flag.StringVar(&configFile, "o", "user.db", "BuntDB database path")
	flag.StringVar(&username, "username", "admin", "Username to create")
	flag.StringVar(&email, "email", "admin@localhost", "Email to create")
	flag.StringVar(&password, "password", "admin", "Password to create")
	flag.StringVar(&authorities, "authorities", "ROLE_USER,ROLE_ADMIN", "Authorities of user")
	flag.StringVar(&operator, "action", "list", "list / delete / create / update / truncate ")
}
func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.Parse()

	userRepo := impl.NewUserRepositoryBuntDB(configFile)
	userSvc := impl.NewUserServiceBuntDB(userRepo)

	switch operator {
	case "list":
		log.Print("user lists:")
		entities, err := userRepo.FindAll()
		if err != nil {
			log.Fatal(err)
		}
		for _, entity := range entities {
			log.Printf("%+v", entity)
		}

	case "truncate":
		userRepo.DeleteAll()
		log.Printf("successfully delete all users")
	case "delete":
		user, err := userRepo.FindByLogin(username)
		if err != nil {
			log.Fatal(err)
		}
		userRepo.DeleteById(user.Id)
		userRepo.DeleteById(user.Login)
	case "save", "create", "update":
		hash, err := userSvc.HashPassword(password)
		if err != nil {
			log.Fatal(err)
		}

		entity := shared.ManagedUserDTO{
			UserDTO: shared.UserDTO{
				Login:       username,
				Email:       email,
				LangKey:     "en",
				Activated:   true,
				Authorities: []string{"ROLE_USER", "ROLE_ADMIN"},
			},
			CreatedBy:        "system",
			CreatedDate:      time.Now().Local().Format("2006-01-02"),
			Password:         hash,
			LastModifiedBy:   "system",
			LastModifiedDate: time.Now().Local().Format("2006-01-02"),
		}

		err = userRepo.Save(&entity)

		if err != nil {
			log.Fatal(err)
		}
		log.Printf("successfully create user with ID: %+v", entity)
	}
}
