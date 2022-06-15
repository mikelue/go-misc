package gorm

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var dbTemplate DbTemplate

var _ = Describe("DbTemplate", func() {
	BeforeEach(func() {
		dbTemplate = NewDbTemplate(db)
		prepareData()
	})
	AfterEach(func() {
		cleanData()
	})

	Context("Normal case", normalCase)
	Context("Panic case", panicCase)
})

var normalCase = func () {
	It("Create(interface{})", func() {
		newCompany := &Company {
			Id: 107, Name: "Leopard Seals Inc.",
			EstablishTime: 39871,
		}
		newCompany.SetAddress("No. 234, Penny County.")

		/**
		 * Creates data
		 */
		dbTemplate.Create(newCompany)
		Expect(dbTemplate.GetDb().RowsAffected).
			To(BeEquivalentTo(1))
		// :~)

		/**
		 * Reload data from database
		 */
		addedCompany := &Company {}
		NewDbTemplate(db.Where(newCompany)).
			First(addedCompany)
		Expect(addedCompany.Name).To(BeIdenticalTo(newCompany.Name))
		// :~)
	})

	It("Delete(interface{}, ...inteface{})", func() {
		dbTemplate.Delete(&Company{ Id: 2 })

		Expect(dbTemplate.GetDb().RowsAffected).
			To(BeEquivalentTo(1))
	})

	It("First(interface{}, ...interface{})", func() {
		Expect(
			dbTemplate.First(&Company{ Id: 2 }),
		).To(BeTrue())

		Expect(
			dbTemplate.First(&Company{ Id: 27 }),
		).To(BeFalse())
	})

	It("FirstOrCreate(interface{}, ...interface{})", func() {
		var sampleId int32 = 113

		/**
		 * Creates data
		 */
		company := &Company{ Id: sampleId, Name: "Auto-Created Inc.", EstablishTime: 23123 }
		dbTemplate.FirstOrCreate(company)
		// :~)

		/**
		 * Reload data from database
		 */
		addedCompany := &Company { Id: sampleId }
		NewDbTemplate(db.Where(addedCompany)).
			First(addedCompany)
		Expect(addedCompany.Name).To(BeIdenticalTo(company.Name))
		// :~)
	})

	It("FirstOrInit(interface{}, ...interface{})", func() {
		dbTemplate = NewDbTemplate(
			db.Attrs(&Company{ Id: 34, Name: "Working..." }),
		)

		initCompany := &Company{ Id: 34 }
		dbTemplate.FirstOrInit(initCompany)

		Expect(initCompany.Name).To(BeEquivalentTo("Working..." ))
	})

	It("Last(interface{}, ...interface{})", func() {
		Expect(
			dbTemplate.Last(&Company{ Id: 2 }),
		).To(BeTrue())

		Expect(
			dbTemplate.Last(&Company{ Id: 27 }),
		).To(BeFalse())
	})

	It("Find(interface{}, ...interface{})", func() {
		companies := make([]*Company, 0, 4)

		dbTemplate.Find(&companies, &Company{ Color: 3 })
		Expect(companies).To(HaveLen(2))

		dbTemplate.Find(&companies, &Company{ Color: 8 })
		Expect(companies).To(HaveLen(0))
	})

	It("Pluck(string, interface{})", func() {
		testedResult := []string{}

		NewDbTemplate(
			db.Model(&Company{}),
		).Pluck("cp_name", &testedResult)

		Expect(testedResult).To(HaveLen(4))
	})

	It("Related(interface{}, ...string)", func() {
		employees := make([]*Employee, 0, 2)

		NewDbTemplate(db.Model(&Company{ Id: 1 })).
			Related(&employees)

		Expect(employees).To(HaveLen(2))
	})

	It("Row() *sql.Row", func() {
		row := NewDbTemplate(
			db.
				Table("ts_company").
				Select([]string{ "cp_id", "cp_name" }).
				Where("cp_id = 1"),
		).Row()

		var id int
		var name string

		Expect(row.Scan(&id, &name)).To(Succeed())
		Expect(id).To(BeEquivalentTo(1))
		Expect(name).To(BeEquivalentTo("Blue Inc."))
	})

	It("Rows() *sql.Rows", func() {
		rows := NewDbTemplate(
			db.
				Table("ts_company").
				Select([]string{ "cp_id", "cp_name" }).
				Where("cp_color = 3"),
		).Rows()
		defer rows.Close()

		var id int
		var name string

		/**
		 * Asserts the matched rows(2 records)
		 */
		counter := 0
		for rows.Next() {
			Expect(rows.Scan(&id, &name)).To(Succeed())
			counter++
		}
		Expect(counter).To(BeEquivalentTo(2))
		// :~)
	})

	It("Save(interface{})", func() {
		sampleName := "Blue2 Inc."
		dbTemplate.Save(&Company{ Id: 1, Name: sampleName})

		updatedData := &Company{ Id: 1 }
		dbTemplate.Take(updatedData)

		Expect(updatedData.Name).To(BeEquivalentTo(sampleName))
	})

	It("Scan(interface{})", func() {
		var loadedCompany Company
		NewDbTemplate(
			db.
				Table("ts_company").
				Select([]string{ "cp_id", "cp_name" }).
				Where("cp_id = 2"),
		).
			Scan(&loadedCompany)

		Expect(loadedCompany.Name).To(BeEquivalentTo("Red Inc."))
	})

	It("ScanRows(*sql.Rows, interface{})", func() {
		rows := NewDbTemplate(
			db.
				Table("ts_company").
				Select([]string{ "cp_id", "cp_name" }).
				Where("cp_color = 3"),
		).Rows()
		defer rows.Close()

		/**
		 * Asserts the matched rows(2 records)
		 */
		counter := 0
		loadedCompany := &Company{}
		for rows.Next() {
			dbTemplate.ScanRows(rows, loadedCompany)

			if counter == 0 {
				Expect(loadedCompany.Id).To(BeEquivalentTo(3))
			} else {
				Expect(loadedCompany.Id).To(BeEquivalentTo(4))
			}
			counter++
		}
		Expect(counter).To(BeEquivalentTo(2))
		// :~)
	})

	It("Transaction(func(*gorm.DB))", func() {
		dbTemplate.Transaction(func(tx *gorm.DB) {
			txTemplate := NewDbTemplate(tx)

			txTemplate.Create(&Company{ Id: 181, Name: "Nice 1" })
			txTemplate.Create(&Company{ Id: 182, Name: "Nice 2" })
		})

		savedCompany := &Company{ Id: 182 }
		dbTemplate.Take(savedCompany)
		Expect(savedCompany.Name).To(BeEquivalentTo("Nice 2"))
	})

	It("SoftTransaction(func(*gorm.DB)) with Context", func() {
		findCompany := &Company{ Id: 3 }
		var softTxFunc SoftTxFunc = func(tx *gorm.DB) TxResult {
			txTemplate := NewDbTemplate(tx)

			txTemplate.Take(findCompany)

			return TxCommit
		}
		testedTx := NewSoftTxWithContext(
			softTxFunc, context.Background(),
			&sql.TxOptions { Isolation: sql.LevelReadUncommitted, ReadOnly: true },
		)

		dbTemplate.SoftTransaction(testedTx)

		Expect(findCompany.Name).To(BeEquivalentTo("Green Inc."))
	})

	It("SoftTransaction(func(*gorm.DB) TxResult) - Commit", func() {
		var testedTx SoftTxFunc = func(tx *gorm.DB) TxResult {
			txTemplate := NewDbTemplate(tx)

			txTemplate.Create(&Company{ Id: 933, Name: "葉葉聲聲 1" })
			txTemplate.Create(&Company{ Id: 934, Name: "葉葉聲聲 2" })

			return TxCommit
		}

		dbTemplate.SoftTransaction(testedTx)

		savedCompany := &Company{ Id: 934 }

		Expect(dbTemplate.First(savedCompany)).To(BeTrue())
	})

	It("SoftTransaction(func(*gorm.DB) TxResult) - Rollback", func() {
		var testedTx SoftTxFunc = func(tx *gorm.DB) TxResult {
			txTemplate := NewDbTemplate(tx)

			txTemplate.Create(&Company{ Id: 1703, Name: "西樓夜 1" })
			txTemplate.Create(&Company{ Id: 1704, Name: "西樓夜 2" })

			return TxRollback
		}

		dbTemplate.SoftTransaction(testedTx)

		savedCompany := &Company{ Id: 1703 }
		Expect(dbTemplate.First(savedCompany)).To(BeFalse())
	})

	It("Take(interface{}, ...interface{})", func() {
		takenCompany := &Company{ Id: 1 }

		Expect(dbTemplate.Take(takenCompany)).To(BeTrue())

		takenCompany.Id = 21
		Expect(dbTemplate.Take(takenCompany)).To(BeFalse())
	})

	It("Update(..interface{})", func() {
		var sampleCompanyId int32 = 3

		NewDbTemplate(db.Model(&Company{ Id: sampleCompanyId })).Update("Color", 10)

		updatedCompany := &Company{ Id: sampleCompanyId }
		Expect(dbTemplate.Take(updatedCompany)).To(BeTrue())
		Expect(updatedCompany.Color).To(BeEquivalentTo(10))
	})

	It("UpdateColumn(..interface{})", func() {
		var sampleCompanyId int32 = 3

		NewDbTemplate(db.Model(&Company{ Id: sampleCompanyId })).UpdateColumn("Color", 13)

		updatedCompany := &Company{ Id: sampleCompanyId }
		Expect(dbTemplate.Take(updatedCompany)).To(BeTrue())
		Expect(updatedCompany.Color).To(BeEquivalentTo(13))
	})

	It("UpdateColumns(interface{})", func() {
		var sampleCompanyId int32 = 3

		NewDbTemplate(db.Model(&Company{ Id: sampleCompanyId })).UpdateColumns(Company{ Color: 15 })

		updatedCompany := &Company{ Id: sampleCompanyId }
		Expect(dbTemplate.Take(updatedCompany)).To(BeTrue())
		Expect(updatedCompany.Color).To(BeEquivalentTo(15))
	})

	It("Updates(interface{}, ...bool)", func() {
		var sampleCompanyId int32 = 3

		NewDbTemplate(db.Model(&Company{ Id: sampleCompanyId })).Updates(Company{ Color: 17 })

		updatedCompany := &Company{ Id: sampleCompanyId }
		Expect(dbTemplate.Take(updatedCompany)).To(BeTrue())
		Expect(updatedCompany.Color).To(BeEquivalentTo(17))
	})
}
var panicCase = func() {
	It("Create(interface{})", func() {
		dbTemplate = NewDbTemplate(db.Set("gorm:insert_option", "<ERROR>"))

		Expect(
			func() { dbTemplate.Create(&Company{}) },
		).To(PanicWith(MatchRegexp("Create.*db_test\\.go")))
	})

	It("Delete(interface{}, ...inteface{})", func() {
		dbTemplate = NewDbTemplate(db.Set("gorm:delete_option", "<ERROR>"))

		Expect(
			func() { dbTemplate.Delete(&Company{ Id: 3 }) },
		).To(PanicWith(MatchRegexp("Delete.*db_test\\.go")))
	})

	It("First(interface{}, ...interface{})", func() {
		dbTemplate = NewDbTemplate(db.Set("gorm:query_option", "<ERROR>"))

		Expect(
			func() { dbTemplate.First(&Company{ Id: 3 }) },
		).To(PanicWith(MatchRegexp("First.*db_test\\.go")))
	})

	It("FirstOrCreate(interface{}, ...interface{})", func() {
		dbTemplate = NewDbTemplate(db.Set("gorm:insert_option", "<ERROR>"))

		company := &Company{ Id: 32, Name: "Auto-Created Inc.", EstablishTime: 23123 }
		Expect(
			func() { dbTemplate.FirstOrCreate(company) },
		).To(PanicWith(MatchRegexp("FirstOrCreate.*db_test\\.go")))
	})

	It("FirstOrInit(interface{}, ...interface{})", func() {
		dbTemplate = NewDbTemplate(
			db.Set("gorm:query_option", "<ERROR>").
				Attrs(&Company{ Id: 34, Name: "Working..." }),
		)

		initCompany := &Company{ Id: 34 }
		Expect(
			func() { dbTemplate.FirstOrInit(initCompany) },
		).To(PanicWith(MatchRegexp("FirstOrInit.*db_test\\.go")))
	})

	It("Find(interface{}, ...interface{})", func() {
		dbTemplate = NewDbTemplate(
			db.Set("gorm:query_option", "<ERROR>"),
		)

		companies := make([]*Company, 0, 4)
		Expect(
			func() { dbTemplate.Find(&companies, &Company{ Color: 3 }) },
		).To(PanicWith(MatchRegexp("Find.*db_test\\.go")))
	})

	It("Pluck(string, interface{})", func() {
		testedResult := []string{}

		panicFunc := func() {
			NewDbTemplate(
				db.Set("gorm:query_option", "<ERROR>").
					Model(&Company{}),
			).Pluck("cp_name", &testedResult)
		}

		Expect(panicFunc).To(PanicWith(MatchRegexp("Pluck.*db_test\\.go")))
	})

	It("Related(interface{}, ...string)", func() {
		dbTemplate = NewDbTemplate(
			db.Set("gorm:query_option", "<ERROR>").
				Model(&Company{ Id: 1 }),
		)

		employees := make([]*Employee, 0, 2)
		Expect(
			func() { dbTemplate.Related(&employees) },
		).To(PanicWith(MatchRegexp("Related.*db_test\\.go")))
	})

	It("Last(interface{}, ...interface{})", func() {
		dbTemplate = NewDbTemplate(db.Set("gorm:query_option", "<ERROR>"))

		Expect(
			func() { dbTemplate.Last(&Company{ Id: 3 }) },
		).To(PanicWith(MatchRegexp("Last.*db_test\\.go")))
	})

	It("Rows() *sql.Rows", func() {
		rawSqlDb := db.
			Table("ts_company").
			Select([]string{ "no_column", "cp_id", "cp_name" }).
			Where("cp_color = 3")

		Expect(
			func() { NewDbTemplate(rawSqlDb).Rows() },
		).To(PanicWith(MatchRegexp("Rows.*db_test\\.go")))
	})

	It("Save(interface{})", func() {
		dbTemplate = NewDbTemplate(db.Set("gorm:update_option", "<ERROR>"))

		Expect(
			func() { dbTemplate.Save(&Company{ Id: 1, Name: "Blue3 Inc." }) },
		).To(PanicWith(MatchRegexp("Save.*db_test\\.go")))
	})

	It("Scan(interface{})", func() {
		dbTemplate = NewDbTemplate(
			db.
				Table("ts_company2").
				Select([]string{ "cp_id", "cp_name" }).
				Where("cp_id = 2"),
		)

		tempCompany := Company{}
		Expect(
			func() { dbTemplate.Scan(&tempCompany) },
		).To(PanicWith(MatchRegexp("Scan.*db_test\\.go")))
	})

	It("ScanRows(*sql.Rows, interface{})", func() {
		rows := NewDbTemplate(
			db.
				Table("ts_company").
				Select([]string{ "cp_id", "cp_name" }).
				Where("cp_color = 3"),
		).Rows()
		defer rows.Close()

		/**
		 * Asserts the matched rows(2 records)
		 */
		counter := 0
		loadedCompany := &struct {
			ErrName int `gorm:"column:cp_name"`
		} {}
		for rows.Next() {
			Expect(
				func() { dbTemplate.ScanRows(rows, loadedCompany) },
			).
				To(PanicWith(MatchRegexp("ScanRows.*db_test\\.go")))

			counter++
		}
		Expect(counter).To(BeEquivalentTo(2))
		// :~)
	})

	It("Transaction(func(*gorm.DB))", func() {
		panicFunc := func() {
			dbTemplate.Transaction(func(tx *gorm.DB) {
				txTemplate := NewDbTemplate(tx)

				txTemplate.Create(&Company{ Id: 391, Name: "Great 1" })
				txTemplate.Create(&Company{ Id: 391, Name: "Great 2" })
			})
		}

		/**
		 * Asserts the panic information
		 */
		Expect(panicFunc).To(PanicWith(MatchRegexp("Transaction.*db_test\\.go")))
		// :~)

		/**
		 * Asserts the rollback result
		 */
		nonExistingCompany := &Company{ Id: 391 }
		Expect(dbTemplate.First(nonExistingCompany)).To(BeFalse())
		// :~)
	})

	It("SoftTransaction(func(*gorm.DB) TxResult)", func() {
		var testedTx SoftTxFunc = func(tx *gorm.DB) TxResult {
			txTemplate := NewDbTemplate(tx)

			txTemplate.Create(&Company{ Id: 501, Name: "花前坐 1" })
			txTemplate.Create(&Company{ Id: 501, Name: "花前坐 2" })

			return TxCommit
		}

		/**
		 * Asserts the panic information
		 */
		Expect(
			func() { dbTemplate.SoftTransaction(testedTx) },
		).
			To(PanicWith(MatchRegexp("SoftTransaction.*db_test\\.go")))
		// :~)

		/**
		 * Asserts the rollback result
		 */
		nonExistingCompany := &Company{ Id: 501 }
		Expect(dbTemplate.First(nonExistingCompany)).To(BeFalse())
		// :~)
	})

	It("Take(interface{}, ...interface{})", func() {
		dbTemplate = NewDbTemplate(db.Set("gorm:query_option", "<ERROR>"))

		takenCompany := &Company{ Id: 1 }
		Expect(
			func() { dbTemplate.Take(takenCompany) },
		).
			To(PanicWith(MatchRegexp("Take.*db_test\\.go")))
	})

	It("Update(..interface{})", func() {
		panicFunc := func() {
			NewDbTemplate(
				db.Set("gorm:update_option", "<ERROR>").
					Model(&Company{ Id: 2 }),
			).
				Update("Color", 10)
		}

		Expect(panicFunc).To(PanicWith(MatchRegexp("Update.*db_test\\.go")))
	})

	It("UpdateColumn(..interface{})", func() {
		panicFunc := func() {
			NewDbTemplate(
				db.Set("gorm:update_option", "<ERROR>").
					Model(&Company{ Id: 2 }),
			).
				UpdateColumn("Color", 8)
		}

		Expect(panicFunc).To(PanicWith(MatchRegexp("UpdateColumn.*db_test\\.go")))
	})

	It("UpdateColumns(interface{})", func() {
		panicFunc := func() {
			NewDbTemplate(
				db.Set("gorm:update_option", "<ERROR>").
					Model(&Company{ Id: 2 }),
			).
				UpdateColumns(Company{ Color: 8 })
		}

		Expect(panicFunc).To(PanicWith(MatchRegexp("UpdateColumns.*db_test\\.go")))
	})

	It("Updates(interface{}, ...bool)", func() {
		panicFunc := func() {
			NewDbTemplate(
				db.Set("gorm:update_option", "<ERROR>").
					Model(&Company{ Id: 2 }),
			).
				Updates(Company{ Color: 8 })
		}

		Expect(panicFunc).To(PanicWith(MatchRegexp("Update.*db_test\\.go")))
	})
}
