package main

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/robfig/cron/v3"
	// "github.com/go-co-op/gocron"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

// A Response struct to map the Entire Response
type User1 struct {
	Id       int    `json:"user_id" gorm:"size:32;  "`
	Name     string `json:"name" gorm:"size:100"`
	Username string `json:"username" gorm:"size:50"`
	Email    string `json:"email" gorm:"size:50"`
	Address  JSONB  `json:"address" gorm:"type:json"`
	Phone    string `json:"phone" gorm:"size:50"`
	Website  string `json:"website" gorm:"size:50"`
	Company  JSONB  `json:"company" gorm:"type:json"`
}

var responseData []byte

func main() {
	cron := cron.New()

	cron.AddFunc("@every 0m59s", Schedule)

	cron.Start()

	http.ListenAndServe(":8003", nil)

	// time.Sleep(time.Minute * 5)
	// s.Every(5).Seconds().Do(Schedule)

	// <-s.Start()

	// s.StartBlocking()
}

func Schedule() {
	response, err := http.Get("https://jsonplaceholder.typicode.com/users")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, _ = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("...")

	var responseObject []User1

	json.Unmarshal(responseData, &responseObject)

	fmt.Println("Vikas")
	fmt.Println(responseObject)
	fmt.Println("Vikas1")

	db, err := gorm.Open(mysql.Open("root:@/company?parseTime=true"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&User1{})
	if err != nil {
		panic(err)
	}

	result := db.Create(responseObject)
	if result.Error != nil {
		panic(result.Error)
	}
}
