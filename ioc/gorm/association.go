package gorm

import (
	"github.com/jinzhu/gorm"
)

// Every function provided by this template would be PANIC
// with viable error after corresponding operation is executed.
type AssociationTemplate interface {
	// As corresponding function of *gorm.Association
    Append(values ...interface{})
	// As corresponding function of *gorm.Association
    Clear()
	// As corresponding function of *gorm.Association
    Count() int
	// As corresponding function of *gorm.Association
    Delete(values ...interface{})
	// As corresponding function of *gorm.Association
    Find(value interface{})
	// As corresponding function of *gorm.Association
    Replace(values ...interface{})
}

func NewAssociationTemplate(association *gorm.Association) AssociationTemplate {
	return &associationTemplateImpl{ association: association }
}

type associationTemplateImpl struct {
	association *gorm.Association
}

func (self *associationTemplateImpl) Append(values ...interface{}) {
	self.association = self.association.Append(values...)
	panicIfAnyErrorOfAssociation("Append", self.association)
}
func (self *associationTemplateImpl) Clear() {
	self.association = self.association.Clear()
	panicIfAnyErrorOfAssociation("Clear", self.association)
}
func (self *associationTemplateImpl) Count() int {
	result := self.association.Count()
	panicIfAnyErrorOfAssociation("Count", self.association)
	return result
}
func (self *associationTemplateImpl) Delete(values ...interface{}) {
	self.association = self.association.Delete(values...)
	panicIfAnyErrorOfAssociation("Delete", self.association)
}
func (self *associationTemplateImpl) Find(value interface{}) {
	self.association = self.association.Find(value)
	panicIfAnyErrorOfAssociation("Find", self.association)
}
func (self *associationTemplateImpl) Replace(values ...interface{}) {
	self.association = self.association.Replace(values...)
	panicIfAnyErrorOfAssociation("Replace", self.association)
}
