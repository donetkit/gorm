package gorm

import (
	"fmt"
	"gorm.io/gorm/clause"
	"reflect"
	"strings"
)

const (
	OrderByAsc  = "ASC"
	OrderByDesc = "DESC"
)

// Limit specify the number of records to be retrieved
func (db *DB) Page(pageIndex, pageSize int) (tx *DB) {
	if pageIndex < 1 {
		pageIndex = 1
	}
	tx = db.getInstance()
	tx.Limit(pageSize).Offset((pageIndex - 1) * pageSize)
	return
}

// Limit specify the number of records to be retrieved pageSize 50
func (db *DB) PageLimit(pageIndex, pageSize int) (tx *DB) {
	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize > 50 {
		pageSize = 50
	}
	tx = db.getInstance()
	tx.Limit(pageSize).Offset((pageIndex - 1) * pageSize)
	return
}

// Order specify order when retrieve records from database
//     db.Order("name DESC")
//     db.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: true})
func (db *DB) OrderBy(orderName, orderType string) (tx *DB) {
	tx = db.getInstance()
	if orderName == "" {
		return
	}
	orderTypes := "ASC"
	if strings.ToUpper(orderType) == "DESC" {
		orderTypes = "DESC"
	}
	tx.Statement.AddClause(clause.OrderBy{
		Columns: []clause.OrderByColumn{{
			Column: clause.Column{Name: fmt.Sprintf("%s %s", orderName, orderTypes), Raw: true},
		}},
	})
	return
}

// Order specify order when retrieve records from database
//     db.Order("name DESC")
//     db.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: true})
func (db *DB) OrderByAsc(orderName string) (tx *DB) {
	if orderName == "" {
		return
	}
	tx = db.getInstance()
	tx.Statement.AddClause(clause.OrderBy{
		Columns: []clause.OrderByColumn{{
			Column: clause.Column{Name: fmt.Sprintf("%s %s", orderName, "ASC"), Raw: true},
		}},
	})
	return
}

// Order specify order when retrieve records from database
//     db.Order("name DESC")
//     db.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: true})
func (db *DB) OrderByDesc(orderName string) (tx *DB) {
	if orderName == "" {
		return
	}
	tx = db.getInstance()
	tx.Statement.AddClause(clause.OrderBy{
		Columns: []clause.OrderByColumn{{
			Column: clause.Column{Name: fmt.Sprintf("%s %s", orderName, "DESC"), Raw: true},
		}},
	})
	return
}

// Order specify order when retrieve records from database
//     db.Order("name DESC")
//     db.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: true})
func (db *DB) OrderByStruct(v interface{}, orderName, orderType string) (tx *DB) {
	tx = db.getInstance()
	if orderName == "" || orderType == "" {
		return
	}
	orderName, orderType = orderByString(structToTagArray(v), orderName, orderType)
	if orderName == "" {
		return
	}
	orderTypes := "ASC"
	if strings.ToUpper(orderType) == "DESC" {
		orderTypes = "DESC"
	}
	tx.Statement.AddClause(clause.OrderBy{
		Columns: []clause.OrderByColumn{{
			Column: clause.Column{Name: fmt.Sprintf("%s %s", orderName, orderTypes), Raw: true},
		}},
	})
	return
}

func orderByString(field []string, sortName, sortOrder string) (string, string) {
	if len(field) <= 0 || sortName == "" || sortOrder == "" {
		return "", ""
	}
	for _, val := range field {
		if val == sortName {
			var sortOrderData = "DESC"
			if strings.ToLower(sortOrder) == "ASC" {
				sortOrderData = "ASC"
			}
			return sortName, sortOrderData
		}
	}
	return "", ""
}

// Struct Tag
//func structToTag(v interface{}) string {
//	json := ""
//	s := reflect.TypeOf(v).Elem() //通过反射获取type定义
//	for i := 0; i < s.NumField(); i++ {
//		var tag = getStructTagGorm(s.Field(i))
//		if tag != "-" {
//			json += getStructTagJson(s.Field(i))
//			if i < s.NumField()-1 {
//				json += ", "
//			}
//		}
//	}
//	return json
//}

func structToTagArray(v interface{}) []string {
	json := make([]string, 0)
	s := reflect.TypeOf(v).Elem() //通过反射获取type定义
	for i := 0; i < s.NumField(); i++ {
		var tag = getStructTagGorm(s.Field(i))
		if tag != "-" {
			json = append(json, getStructTagJson(s.Field(i)))
		}
	}
	return json
}

func getStructTagJson(f reflect.StructField) string {
	return f.Tag.Get("json")
}

func getStructTagGorm(f reflect.StructField) string {
	return f.Tag.Get("gorm")
}
