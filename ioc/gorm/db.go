/*
This package provides error-free methods over "*gorm.DB" object.

DbTemplate

You can use "NewDbTemplate()" to convert a "*gorm.DB" to use the error-free methods.

    dbTmpl := NewDbTemplate(gormDb)
    // Panic if something gets wrong
    dbTmpl.Create(yourObj)

SoftTx

You can implements this interface to use "DbTemplate.SoftTransaction()".

"SoftTxFunc" is funcational style of "SoftTx".
*/
package gorm

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"
)

// The result of transaction
type TxResult int8

// Used to commit transaction
const TxCommit TxResult = 1
// Used to rollback transaction
const TxRollback TxResult = 2

// The callback type used by "DbTemplate.SoftTransaction(SoftTx)"
type SoftTx interface {
	Tx(tx *gorm.DB) TxResult
}

// The functional type of "SoftTx"
type SoftTxFunc func(*gorm.DB) TxResult

// Delegate the function "SoftTx.Tx(*gorm.DB)" to the fucntion self.
func (f SoftTxFunc) Tx(tx *gorm.DB) TxResult {
	return f(tx)
}

// Constructs a "SoftTx" with "context.Context" and "*sql.TxOptions"
func NewSoftTxWithContext(
	softTxImpl SoftTx,
	ctx context.Context, opts *sql.TxOptions,
) SoftTx {
	newTx := &softTxWithContext{
		context: ctx, txOptions: opts,
	}
	newTx.SoftTx = softTxImpl
	return newTx
}

type softTxFunc struct {
	implFunc func(*gorm.DB) TxResult
}
func (self *softTxFunc) Tx(tx *gorm.DB) TxResult {
	return self.implFunc(tx)
}

type softTxWithContext struct {
	context context.Context
	txOptions *sql.TxOptions
	SoftTx
}

// Every function provided by this template would be PANIC
// with viable error after corresponding operation is executed.
type DbTemplate interface {
	// Gets underlying object of *gorm.DB
	GetDb() *gorm.DB

	// As corresponding function of *gorm.DB
	Create(value interface{})
	// As corresponding function of *gorm.DB
	Delete(value interface{}, where ...interface{})
	// As corresponding function of *gorm.DB
	// Return values:
	// 	True - at least one row of data is found
	// 	False - nothing found
	First(out interface{}, where ...interface{}) bool
	// As corresponding function of *gorm.DB
	FirstOrCreate(out interface{}, where ...interface{})
	// As corresponding function of *gorm.DB
	FirstOrInit(out interface{}, where ...interface{})
	// As corresponding function of *gorm.DB
	Find(out interface{}, where ...interface{})
	// As corresponding function of *gorm.DB
	Related(value interface{}, foreignKeys ...string)
	// As corresponding function of *gorm.DB
	// Return values:
	// 	True - at least one row of data is found
	// 	False - nothing found
	Last(out interface{}, where ...interface{}) bool
	// As corresponding function of *gorm.DB
	Row() *sql.Row
	// As corresponding function of *gorm.DB
	Rows() *sql.Rows
	// As corresponding function of *gorm.DB
	Save(value interface{})
	// As corresponding function of *gorm.DB
	Scan(dest interface{})
	// As corresponding function of *gorm.DB
	ScanRows(rows *sql.Rows, result interface{})
	// As corresponding function of *gorm.DB
	Transaction(func(tx *gorm.DB))
	// As corresponding function of *gorm.DB
	Pluck(column string, value interface{})
	// Use returned value of "SoftTx.Tx(*gorm.DB)" to decide whether or not commit a transaction.
	// If something gets wrong, the transaction would be rollback either.
	SoftTransaction(SoftTx)
	// As corresponding function of *gorm.DB
	// Return values:
	// 	True - at least one row of data is found
	// 	False - nothing found
	Take(out interface{}, where ...interface{}) bool
	// As corresponding function of *gorm.DB
	Update(attrs ...interface{})
	// As corresponding function of *gorm.DB
	UpdateColumn(attrs ...interface{})
	// As corresponding function of *gorm.DB
	UpdateColumns(values interface{})
	// As corresponding function of *gorm.DB
	Updates(values interface{}, ignoreProtectedAttrs ...bool)
}

// Constructs a wrapper over a *gorm.DB object
func NewDbTemplate(targetDb *gorm.DB) DbTemplate {
	return &templateImpl{ db: targetDb }
}

type templateImpl struct {
	db *gorm.DB
}

func (self *templateImpl) GetDb() *gorm.DB {
	return self.db
}
func (self *templateImpl) Create(value interface{}) {
	self.db = self.db.Create(value)
	panicIfAnyError("Create", self.db)
}
func (self *templateImpl) Delete(value interface{}, where ...interface{}) {
	self.db = self.db.Delete(value, where...)
	panicIfAnyError("Delete", self.db)
}
func (self *templateImpl) First(out interface{}, where ...interface{}) bool {
	self.db = self.db.First(out, where...)
	panicIfAnyError("First", self.db)
	return !self.db.RecordNotFound()
}
func (self *templateImpl) FirstOrCreate(out interface{}, where ...interface{}) {
	self.db = self.db.FirstOrCreate(out, where...)
	panicIfAnyError("FirstOrCreate", self.db)
}
func (self *templateImpl) FirstOrInit(out interface{}, where ...interface{}) {
	self.db = self.db.FirstOrInit(out, where...)
	panicIfAnyError("FirstOrInit", self.db)
}
func (self *templateImpl) Find(out interface{}, where ...interface{}) {
	self.db = self.db.Find(out, where...)
	panicIfAnyError("Find", self.db)
}
func (self *templateImpl) Related(value interface{}, foreignKeys ...string) {
	self.db = self.db.Related(value, foreignKeys...)
	panicIfAnyError("Related", self.db)
}
func (self *templateImpl) Last(out interface{}, where ...interface{}) bool {
	self.db = self.db.Last(out, where...)
	panicIfAnyError("Last", self.db)
	return !self.db.RecordNotFound()
}
func (self *templateImpl) Pluck(column string, value interface{}) {
	self.db = self.db.Pluck(column, value)
	panicIfAnyError("Pluck", self.db)
}
func (self *templateImpl) Row() *sql.Row {
	return self.db.Row()
}
func (self *templateImpl) Rows() *sql.Rows {
	rows, err := self.db.Rows()

	if err != nil {
		panicError("Rows", err)
	}

	return rows
}
func (self *templateImpl) Save(value interface{}) {
	self.db = self.db.Save(value)
	panicIfAnyError("Save", self.db)
}
func (self *templateImpl) Scan(dest interface{}) {
	self.db = self.db.Scan(dest)
	panicIfAnyError("Scan", self.db)
}
func (self *templateImpl) ScanRows(rows *sql.Rows, result interface{}) {
	if err := self.db.ScanRows(rows, result); err != nil {
		panicError("ScanRows", err)
	}
}
func (self *templateImpl) Transaction(txCallback func(tx *gorm.DB)) {
	freeErrorTxCallback := func(tx *gorm.DB) (err error) {
		defer func() {
			if exception := decoratePanic("Transaction", recover()); exception != nil {
				panic(exception)
			}
		}()

		txCallback(tx)
		return nil
	}

	panicError("Transaction", self.db.Transaction(freeErrorTxCallback))
}
func (self *templateImpl) SoftTransaction(softTx SoftTx) {
	var tx *gorm.DB

	if txWithContext, ok := softTx.(*softTxWithContext); ok {
		tx = self.db.BeginTx(txWithContext.context, txWithContext.txOptions)
	} else {
		tx = self.db.Begin()
	}

	defer func() {
		p := recover()
		if p != nil {
			tx.Rollback()
		}

		if exception := decoratePanic("SoftTransaction", p); exception != nil {
			panic(exception)
		}
	}()

	if softTx.Tx(tx) == TxRollback {
		panicIfAnyError("SoftTransaction", tx.Rollback())
		return
	}

	panicIfAnyError("SoftTransaction", tx.Commit())
}
func (self *templateImpl) Take(out interface{}, where ...interface{}) bool {
	self.db = self.db.Take(out, where...)
	panicIfAnyError("Take", self.db)
	return !self.db.RecordNotFound()
}
func (self *templateImpl) Update(attrs ...interface{}) {
	self.db = self.db.Update(attrs...)
	panicIfAnyError("Update", self.db)
}
func (self *templateImpl) UpdateColumn(attrs ...interface{}) {
	self.db = self.db.UpdateColumn(attrs...)
	panicIfAnyError("UpdateColumn", self.db)
}
func (self *templateImpl) UpdateColumns(values interface{}) {
	self.db = self.db.UpdateColumns(values)
	panicIfAnyError("UpdateColumns", self.db)
}
func (self *templateImpl) Updates(values interface{}, ignoreProtectedAttrs ...bool) {
	self.db = self.db.Updates(values, ignoreProtectedAttrs...)
	panicIfAnyError("Updates", self.db)
}
