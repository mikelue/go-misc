package gorm

import (
	"fmt"
)

func ExampleDbException() {
	db := exampleDb()

	defer func() {
		if p := recover(); p != nil {
			dbException := p.(*DbException)

			fmt.Printf("DbException: Func[%s]. Content: %s",
				dbException.GetGormFuncName(),
				dbException.GetCause(),
			)
		}
	}()

	defer db.Close()
	setupExampleSchema(db)

	dbTmpl := NewDbTemplate(db)
	dbTmpl.Create(&Company{ Id: 6701, Name: "Green Inc." })
	dbTmpl.Create(&Company{ Id: 6702, Name: "Green Inc." })

	// Output:
	// DbException: Func[Create]. Content: UNIQUE constraint failed: ts_company.cp_name
}
