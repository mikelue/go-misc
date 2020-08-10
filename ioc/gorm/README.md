**Table of Content**

* [Features](#features)
* [Error-Free operator](#efo)
	* [DbTemplate](#dbtemplate)
	* [AssociationTemplate](#associationtemplate)
* [Test environment](#test-environment)
	* [Flags](#test-flags)

This library depends on [gorm](https://gorm.io/),
which provides some convenient features.

# Features <a name="features"></a>

## Error-Free operator <a name="efo"></a>

`DbTemplate` and `AssociationTemplate` provides some error-free operators,<br>
which would **[panic](https://golang.org/ref/spec#Handling_panics)** the caller with [DbException](./errors.go).

### [DbTemplate](./db.go) <a name="dbtemplate"></a>

Usage example:
```go
// Panic if something gets wrong
NewDbTemplate(db).Create(yourObject)
```

### [AssociationTemplate](./association.go) <a name="associationtemplate"></a>

Usage example:
```go
// Panic if something gets wrong
NewAssociationTemplate(
    db.Model(parentObject).Association("SubObjects"),
).Find(&subObjects)
```

# Test environment <a name="test-environment"></a>

This library uses [gingkgo](https://github.com/onsi/ginkgo) as underlying test framework.

The database used by tests is [Sqlite](https://www.sqlite.org/index.html).

## Flags <a name="test-flags"></a>

Usage example:
```sh
# Run all tests
ginkgo -- -gorm.log-mode=true
# Run matched(by RegExp) tests
ginkgo --focus 'Create' -- -gorm.log-mode=true
```

`gorm.log-mode` - see [gorm/DB.LogMode](https://godoc.org/github.com/jinzhu/gorm#DB.LogMode)
* `true`(`1`) - Verbose logging by *gorm*
* `false0`(`0`)(default) - For logging by *gorm*
