package postgresql

import (
	"fmt"
	"log"

	"github.com/bzhtux/sample_apps/postgresql/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDB() *gorm.DB {

	// wokeignore:rule=disable
	// dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", models.PGConfig.Host, models.PGConfig.Port, models.PGConfig.Username, models.PGConfig.DB, models.PGConfig.Password)
	pgc := new(config.PGConfig)
	pgc.NewConfig()

	var dsn = fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s", pgc.Host, pgc.Port, pgc.Username, pgc.DB, pgc.Password)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// log.Fatalf("%s", err)
		fmt.Println("*** Error connectinng to DB ...")
		log.Printf("%s", err)
	}

	return conn
}

func HealthCheck(dsn string) bool {
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return false
	} else {
		return true
	}
}

func AutoMigrate(db *gorm.DB, database interface{}) {

	db.AutoMigrate(database)

}
