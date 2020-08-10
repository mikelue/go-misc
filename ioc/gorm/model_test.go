package gorm

import (
	"strings"
)

type Company struct {
	Id int32 `gorm:"column:cp_id;type:integer;primary_key"`
	Name string `gorm:"column:cp_name;type:varchar(64);unique;not null"`
	Color int8 `gorm:"column:cp_color;type:tinyint;not null;default: 1"`
	Phone *string `gorm:"column:cp_phone;type:varchar(64)"`
	Address *string `gorm:"column:cp_address;type:varchar(64)"`
	EstablishTime int32 `gorm:"column:cp_establish_time;type:integer;not null"`
	Employees []*Employee `gorm:"foreignkey:CompanyId"`
	PreviousEmployees []*Employee `gorm:"foreignkey:PreviousCompanyId"`
}

func (Company) TableName() string { return "ts_company" }
func (s *Company) SetAddress(v string) {
	v = strings.TrimSpace(v)
	if v == "" {
		s.Address = nil
	} else {
		s.Address = &v
	}
}

type Employee struct {
	Id int32 `gorm:"column:em_id;type:integer;primary_key"`
	Name string `gorm:"column:em_name;type:varchar(64);unique;not null"`
	Birthday int32 `gorm:"column:em_birthday;type:integer;not null"`

	CompanyId int32 `gorm:"column:em_cp_id;type:integer;not null"`
	PreviousCompanyId *int32 `gorm:"column:em_previous_cp_id;type:integer"`
	WorkFor Company `gorm:"foreignkey:CompanyId"`
}

func (Employee) TableName() string { return "ts_employee" }

func prepareData() {
	db.Exec(`
	INSERT INTO ts_company(cp_id, cp_name, cp_color, cp_phone, cp_establish_time)
	VALUES
		(1, 'Blue Inc.', 1, '9801-9821', 19081),
		(2, 'Red Inc.', 2, '1823-1906', 2619),
		(3, 'Green Inc.', 3, '0295-2942', 24359),
		(4, 'Green/2 Inc.', 3, '1395-7742', 7634);
	`)
	if db.Error != nil { panic(db.Error) }

	db.Exec(`
	INSERT INTO ts_employee(em_id, em_name, em_birthday, em_cp_id, em_previous_cp_id)
	VALUES
		(1, 'Hayes', 3750, 1, null),
		(2, 'Gilbert', 4652, 1, 3),
		(3, 'Robinson', 4320, 2, 4),
		(4, 'Lewis', 3676, 2, null);
	`)
	if db.Error != nil { panic(db.Error) }
}
func cleanData() {
	db.Exec(`DELETE FROM ts_employee`)
	if db.Error != nil { panic(db.Error) }
	db.Exec(`DELETE FROM ts_company`)
	if db.Error != nil { panic(db.Error) }
}
