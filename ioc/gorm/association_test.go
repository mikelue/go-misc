package gorm

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AssociationTemplate", func() {
	BeforeEach(func() {
		dbTemplate = NewDbTemplate(db)
		prepareData()
	})
	AfterEach(func() {
		cleanData()
	})

	Context("Normal Case", normalCaseOfAssociation)
	Context("Panic Case", panicCaseOfAssociation)
})

var normalCaseOfAssociation = func() {
	It("Append(...interface{})", func() {
		association := NewAssociationTemplate(
			db.Model(&Company{ Id: 4 }).Association("Employees"),
		)

		association.Append(&Employee{ Id: 401, Name: "Seth", Birthday: 3573 })

		/**
		 * Asserts the saved employees
		 */
		savedEmployee := &Employee{ Id: 401 }
		dbTemplate.Take(savedEmployee)
		Expect(savedEmployee.CompanyId).To(BeEquivalentTo(4))
		// :~)
	})
    It("Clear()", func() {
		association := NewAssociationTemplate(
			db.Model(&Company{ Id: 3 }).Association("PreviousEmployees"),
		)

		association.Clear()

		/**
		 * Asserts the removed FK
		 */
		employee := &Employee{ CompanyId: 2 }
		dbTemplate.Take(employee)
		Expect(employee.PreviousCompanyId).To(BeNil())
		// :~)
	})
    It("Count() int", func() {
		association := NewAssociationTemplate(
			db.Model(&Company{ Id: 2 }).Association("Employees"),
		)

		Expect(association.Count()).To(BeEquivalentTo(2))
	})
    It("Delete(...interface{})", func() {
		association := NewAssociationTemplate(
			db.Model(&Company{ Id: 4 }).Association("PreviousEmployees"),
		)

		association.Delete(&Employee{ Id: 3 })

		/**
		 * Asserts the removed FK
		 */
		employee := &Employee{ CompanyId: 3 }
		dbTemplate.Take(employee)
		Expect(employee.PreviousCompanyId).To(BeNil())
		// :~)
	})
    It("Find(interface{})", func() {
		association := NewAssociationTemplate(
			db.Model(&Company{ Id: 1 }).Association("Employees"),
		)

		findEmployees := make([]*Employee, 0, 2)
		association.Find(&findEmployees)

		Expect(findEmployees).To(HaveLen(2))
	})
    It("Replace(...interface{})", func() {
		association := NewAssociationTemplate(
			db.Model(&Company{ Id: 3 }).Association("PreviousEmployees"),
		)

		changedEmployee := &Employee{ Id: 1 }
		dbTemplate.Take(changedEmployee)
		association.Replace([]*Employee{ changedEmployee })

		/**
		 * Asserts the removed FK
		 */
		employee := &Employee{ Id: 2 }
		dbTemplate.Take(employee)
		Expect(employee.PreviousCompanyId).To(BeNil())
		// :~)

		/**
		 * Asserts the assigned FK
		 */
		employee = &Employee{ Id: 1 }
		dbTemplate.Take(employee)
		Expect(*employee.PreviousCompanyId).To(BeEquivalentTo(3))
		// :~)
	})
}
var panicCaseOfAssociation = func() {
	It("Append(...interface{})", func() {
		association := NewAssociationTemplate(
			db.Set("gorm:insert_option", "<ERROR>").
				Model(&Company{ Id: 4 }).Association("Employees"),
		)

		Expect(func() {
			association.Append(&Employee{ Id: 921, Name: "Ian", Birthday: 10279 })
		}).To(PanicWith(MatchRegexp("Append.*association_test\\.go")))
	})
    It("Clear()", func() {
		association := NewAssociationTemplate(
			db.Set("gorm:update_option", "<ERROR>").
				Model(&Company{ Id: 3 }).Association("PreviousEmployees"),
		)

		Expect(func() {
			association.Clear()
		}).To(PanicWith(MatchRegexp("Clear.*association_test\\.go")))
	})
    It("Count() int", func() {
		association := NewAssociationTemplate(
			db.Set("gorm:query_option", "<ERROR>").
				Model(&Company{ Id: 2 }).Association("Employees"),
		)

		Expect(func() {
			association.Count()
		}).To(PanicWith(MatchRegexp("Count.*association_test\\.go")))
	})
    It("Delete(...interface{})", func() {
		association := NewAssociationTemplate(
			db.Set("gorm:update_option", "<ERROR>").
				Model(&Company{ Id: 4 }).Association("PreviousEmployees"),
		)

		Expect(func() {
			association.Delete(&Employee{ Id: 3 })
		}).To(PanicWith(MatchRegexp("Delete.*association_test\\.go")))
	})
    It("Find(interface{})", func() {
		association := NewAssociationTemplate(
			db.Set("gorm:query_option", "<ERROR>").
				Model(&Company{ Id: 1 }).Association("Employees"),
		)

		findEmployees := make([]*Employee, 0, 2)

		Expect(func() {
			association.Find(&findEmployees)
		}).To(PanicWith(MatchRegexp("Find.*association_test\\.go")))
	})
    It("Replace(...interface{})", func() {
		association := NewAssociationTemplate(
			db.Set("gorm:update_option", "<ERROR>").
				Model(&Company{ Id: 3 }).Association("PreviousEmployees"),
		)

		changedEmployee := &Employee{ Id: 1 }
		dbTemplate.Take(changedEmployee)

		Expect(func() {
			association.Replace([]*Employee{ changedEmployee })
		}).To(PanicWith(MatchRegexp("Replace.*association_test\\.go")))
	})
}
