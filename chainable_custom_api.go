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

// Scan scan value to a struct
func (db *DB) ScanOne(dest interface{}) (tx *DB) {
	tx = db.getInstance()
	tx.Limit(1).Scan(dest)
	return
}

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
func (db *DB) OrderByName(orderName string, desc bool) (tx *DB) {
	if orderName == "" {
		return
	}
	tx = db.getInstance()
	tx.Statement.AddClause(clause.OrderBy{
		Columns: []clause.OrderByColumn{clause.OrderByColumn{Column: clause.Column{Name: orderName}, Desc: desc}},
	})
	return
}

// Order specify order when retrieve records from database
//     db.Order("name DESC")
//     db.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: true})
func (db *DB) OrderByNameGbk(orderName string, desc bool) (tx *DB) {
	if orderName == "" {
		return
	}
	tx = db.getInstance()
	orderTypes := "ASC"
	if desc {
		orderTypes = "DESC"
	}
	tx.Statement.AddClause(clause.OrderBy{
		Columns: []clause.OrderByColumn{{
			Column: clause.Column{Name: fmt.Sprintf("CONVERT(%s USING gbk) %s", orderName, orderTypes), Raw: true},
		}},
	})
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
			Column: clause.Column{Name: fmt.Sprintf(" CONVERT(%s USING gbk) %s", orderName, orderTypes), Raw: true},
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
			Column: clause.Column{Name: fmt.Sprintf(" CONVERT(%s USING gbk) %s", orderName, "ASC"), Raw: true},
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
			Column: clause.Column{Name: fmt.Sprintf(" CONVERT(%s USING gbk) %s", orderName, "DESC"), Raw: true},
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
	orderName, orderType = orderByString(getStructFieldTagArray(v), orderName, orderType)
	if orderName == "" {
		return
	}
	orderTypes := "ASC"
	if strings.ToUpper(orderType) == "DESC" {
		orderTypes = "DESC"
	}
	tx.Statement.AddClause(clause.OrderBy{
		Columns: []clause.OrderByColumn{{
			Column: clause.Column{Name: fmt.Sprintf(" CONVERT(%s USING gbk) %s", orderName, orderTypes), Raw: true},
		}},
	})
	return
}

// Select specify fields that you want when querying, creating, updating
func (db *DB) SelectByStruct(v interface{}, args ...interface{}) (tx *DB) {
	tx = db.getInstance()
	tx.Select(structToTag(v))
	return
}

func orderByString(field []string, sortName, sortOrder string) (string, string) {
	if len(field) <= 0 || sortName == "" || sortOrder == "" {
		return "", ""
	}
	for _, val := range field {
		if val == sortName {
			var sortOrderData = "DESC"
			if strings.ToUpper(sortOrder) == "ASC" {
				sortOrderData = "ASC"
			}
			return sortName, sortOrderData
		}
	}
	return "", ""
}

// Delete delete value match given conditions, if the value has primary key, then will including the primary key as condition
func (db *DB) DeleteByNil() (tx *DB) {
	tx = db.getInstance()
	tx.Statement.Dest = ""
	tx.callbacks.Delete().Execute(tx)
	return
}

// Struct Tag
func structToTag(v interface{}) string {
	jsonArray := getStructFieldTagArray(v)
	return strings.Join(jsonArray, ", ")
}

func getStructFieldTagArray(v interface{}) []string {
	jsonArray := make([]string, 0)
	s := reflect.TypeOf(v).Elem() //??????????????????type??????
	for i := 0; i < s.NumField(); i++ {
		var tag = getStructFieldTag(s.Field(i), "gorm")
		if tag == "-" {
			continue
		}
		if tag != "" {
			switch tag {
			case "name":
				tag = fmt.Sprintf("`%s`", tag)
			case "describe":
				tag = fmt.Sprintf("`%s`", tag)
			case "status":
				tag = fmt.Sprintf("`%s`", tag)
			}
			jsonArray = append(jsonArray, tag)
			continue
		}
		//json,omitempty
		tag = getStructFieldTag(s.Field(i), "json")
		if tag != "" {
			tag = strings.ReplaceAll(tag, ",omitempty", "")
			switch tag {
			case "name":
				tag = fmt.Sprintf("`%s`", tag)
			case "describe":
				tag = fmt.Sprintf("`%s`", tag)
			case "status":
				tag = fmt.Sprintf("`%s`", tag)
			}
			jsonArray = append(jsonArray, tag)
		}
	}
	return jsonArray
}

func getStructFieldTag(f reflect.StructField, tag string) string {
	return f.Tag.Get(tag)
}
