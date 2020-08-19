package gorm

import (
	"fmt"
	"database/sql"
	"context"

	"github.com/jinzhu/gorm"
)

func ExampleDbTemplate() {
	db := exampleDb()

	defer db.Close()
	setupExampleSchema(db)

	/**
	 * Uses the template to perform operations to database
	 */
	tmpl := NewDbTemplate(db)
	tmpl.Create(&Company{ Id: 981, Name: "Green Inc.", EstablishTime: 1248954 })

	foundCompany := &Company{ Id: 981 }
	tmpl.Take(foundCompany)
	// :~)

	fmt.Printf("Company: %s", foundCompany.Name)

	// Output:
	// Company: Green Inc.
}

func ExampleDbTemplate_softTransaction() {
	db := exampleDb()

	defer db.Close()
	setupExampleSchema(db)

	/**
	 * Executes database operations in transaction
	 */
	NewDbTemplate(db).SoftTransaction(
		SoftTxFunc(func(db *gorm.DB) TxResult {
			dbTmpl := NewDbTemplate(db)

			dbTmpl.Create(&Company{ Id: 4091, Name: "Yangmei Inc.", EstablishTime: 6644698 })
			dbTmpl.Create(&Company{ Id: 4092, Name: "Grapefruit Inc.", EstablishTime: 7086887 })

			// Commits current transaction
			return TxCommit
		}),
	)
	// :~)

	foundCompany := &Company{ Id: 4092 }
	NewDbTemplate(db).Take(foundCompany)

	fmt.Printf("Company: %s", foundCompany.Name)

	// Output:
	// Company: Grapefruit Inc.
}

func ExampleNewSoftTxWithContext() {
	db := exampleDb()

	defer db.Close()
	setupExampleSchema(db)

	/**
	 * Executes database operations in transaction
	 */
	customizedTx := NewSoftTxWithContext(
		SoftTxFunc(func(db *gorm.DB) TxResult {
			dbTmpl := NewDbTemplate(db)

			dbTmpl.Create(&Company{ Id: 1871, Name: "Yangmei Inc.", EstablishTime: 6644698 })
			dbTmpl.Create(&Company{ Id: 1872, Name: "Honeycrisp apple Inc.", EstablishTime: 7086887 })

			// Commits current transaction
			return TxCommit
		}),
		context.Background(),
		&sql.TxOptions { Isolation: sql.LevelReadUncommitted },
	)
	NewDbTemplate(db).SoftTransaction(customizedTx)
	// :~)

	foundCompany := &Company{ Id: 1872 }
	NewDbTemplate(db).Take(foundCompany)

	fmt.Printf("Company: %s", foundCompany.Name)

	// Output:
	// Company: Honeycrisp apple Inc.
}
