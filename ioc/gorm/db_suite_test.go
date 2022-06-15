package gorm

import (
	"fmt"
	"flag"
	"testing"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var gorm_logMode *bool

func init() {
	gorm_logMode = flag.Bool("gorm.log-mode", false, "Whether or not to enable logging of GORM(default: false)")
}

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Db Suite")
}

var db *gorm.DB

var _ = BeforeSuite(func() {
	db = newDbBySqliteMemory("testDb")
	setupSchema()
})
var _ = AfterSuite(func() {
	if db != nil {
		db.Close()
	}

	db = nil
})

func setupSchema() {
	if errors := db.AutoMigrate(&Company{}, &Employee{}).
		GetErrors();
		len(errors) > 0 {
		panic(errors)
	}
}

func setupExampleSchema(db *gorm.DB) {
	if errors := db.AutoMigrate(&Company{}, &Employee{}).
		GetErrors();
		len(errors) > 0 {
		panic(errors)
	}
}

func newDbBySqliteMemory(name string) *gorm.DB {
	db, err := gorm.Open("sqlite3", fmt.Sprintf("file:%s?mode=memory&cache=shared", name))
	db.LogMode(*gorm_logMode)
	if err != nil {
		panic(err)
	}

	return db
}
func exampleDb() *gorm.DB {
	return newDbBySqliteMemory("exampleDb")
}
