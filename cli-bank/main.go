package main

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

// Account data model
type Account struct {
	gorm.Model
	Username     string
	Balance      float64
	PasswordHash string
}

var db *gorm.DB
var loggenIn Account

func hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func checkHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func register() {
	fmt.Println("Username:")

	// take username
	var username string
	fmt.Scanln(&username)
	fmt.Println("")

	fmt.Println("Password: ")

	// take password
	var password string
	fmt.Scanln(&password)

	var exists Account
	if err := db.Where("username = ?", username).First(&exists).Error; err == nil {
		fmt.Println("User already exists")
		return
	}

	hash, hashErr := hash(password)
	if hashErr != nil {
		fmt.Println("Problem with hashing password")
		return
	}

	account := Account{
		Balance:      0.0,
		PasswordHash: hash,
		Username:     username,
	}

	db.NewRecord(account)
	db.Create(&account)
	return
}

func login() {
	fmt.Println("Username:")

	// take username
	var username string
	fmt.Scanln(&username)
	fmt.Println("")

	fmt.Println("Password: ")

	// take password
	var password string
	fmt.Scanln(&password)

	var account Account
	if err := db.Where("username = ?", username).First(&account).Error; err != nil {
		fmt.Println("No user has been found")
		return
	}

	if !checkHash(password, account.PasswordHash) {
		fmt.Println("Wrong credentials")
		return
	}

	// set the global variable
	loggenIn = account
	return
}

func main() {
	var err error
	// load environment variables
	err = godotenv.Load()
	if err != nil {
		panic(err)
	}
	dbConfig := os.Getenv("DB_CONFIG")

	fmt.Println("Connecting to database...")

	// connect tot database
	db, err = gorm.Open("mysql", dbConfig)

	// handle errors
	if err != nil {
		panic(err)
	}
	fmt.Println("Database connection established.")

	// input loop
	for {
		fmt.Println("Choose action")
		fmt.Println("1. Login")
		fmt.Println("2. Register")

		var input string
		fmt.Scanln(&input)

		number, err := strconv.Atoi(input)
		if err != nil {
			panic(err)
		}

		if number > 2 {
			fmt.Println("Please choose between 1, 2.")
			continue
		}
	}

	defer db.Close()
}
