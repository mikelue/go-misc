package gorm

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func ExampleDbTemplate() {
	var err error
	db, err = gorm.Open("sqlite3", "file:exampledb?mode=memory&cache=shared")

	if err != nil {
		panic(err)
	}

	tmpl := NewDbTemplate(db)
	// tmpl.Create(value)
	// tmpl.Take(out)

	_ = tmpl
}
