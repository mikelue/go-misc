package gorm

import (
	"flag"
	"testing"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	. "github.com/onsi/ginkgo"
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
	var err error
	db, err = gorm.Open("sqlite3", "file:testdb?mode=memory&cache=shared")
	if err != nil {
		panic(err)
	}

	db.LogMode(*gorm_logMode)
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
