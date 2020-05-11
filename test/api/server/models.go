package server

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"time"
)

type MyModel struct {
	ID        uint       `gorm:"primary_key",json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-";sql:"index"`
}

type User struct {
	MyModel
	Surname      string `json:"surname"`
	Name         string `json:"name"`
	MiddleName   string `json:"middle_name"`
	INN          int    `json:"inn"`
	Sex          string `json:"sex"`
	Login        string `json:"-";gorm:"unique_index"`
	Password     string `json:"-"`
	Email        string `json:"email"`
	CompanyRefer int    `json:"company_id"`
}

type Event struct {
	MyModel
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DateStart   time.Time `json:"date_start"`
	DateFin     time.Time `json:"date_fin"`
	PlaceName   string    `json:"place_name"`
	PlaceLat    float64   `json:"place_lat"`
	PlaceLon    float64   `json:"place_lon"`
	Manager     User      `gorm:"foreignkey:ManagerId;",json:"-"`
	ManagerId   int       `json:"manager"`
	Members     []User    `gorm:"many2many:events_members;",json:"members"`
	Companies   []Company `gorm:"many2many:events_companies;",json:"companies"`
}

type Company struct {
	MyModel
	Name  string `json:"name"`
	Users []User `gorm:"foreignkey:CompanyRefer"`
}

var db *gorm.DB

func DbStart() {
	var err error
	db, err = gorm.Open("postgres", "host="+os.Getenv("POSTGRES_HOST")+" port="+os.Getenv("POSTGRES_PORT")+" user="+os.Getenv("POSTGRES_USER")+" dbname="+os.Getenv("POSTGRES_DB")+" password="+os.Getenv("POSTGRES_PASSWORD")+" sslmode=disable")
	//db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&User{},
		&Event{},
		&Company{})
	fmt.Println("Migrations succeeded")
}
