package gorm

import (
	"fmt"
	"gorm.io/gorm/clause"
	"reflect"
	"strings"
)

// ScanCount scan count value to a int64
func (db *DB) ScanCount() int64 {
	tx := db.getInstance()
	var total int64
	tx.Count(&total)
	return total
}

// ScanOne Scan scan value to a struct
func (db *DB) ScanOne(dest interface{}) (tx *DB) {
	tx = db.getInstance()
	tx.Limit(1).Scan(dest)
	return
}

// Page Limit specify the number of records to be retrieved
func (db *DB) Page(pageIndex, pageSize int) (tx *DB) {
	if pageIndex < 1 {
		pageIndex = 1
	}
	tx = db.getInstance()
	tx.Limit(pageSize).Offset((pageIndex - 1) * pageSize)
	return
}

// PageLimit Limit specify the number of records to be retrieved pageSize 50
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

// OrderByAsc Order specify order when retrieve records from database
//
//	db.Order("name DESC")
//	db.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: true})
func (db *DB) OrderByAsc(orderName string) (tx *DB) {
	if orderName == "" {
		return
	}
	tx = db.getInstance()
	tx.Statement.AddClause(clause.OrderBy{
		Columns: []clause.OrderByColumn{{Column: clause.Column{Name: orderName}, Desc: false}},
	})
	return
}

// OrderByAscGBK Order specify order when retrieve records from database
//
//	db.Order("name DESC")
//	db.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: true})
func (db *DB) OrderByAscGBK(orderName string) (tx *DB) {
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

// OrderByDesc Order specify order when retrieve records from database
//
//	db.Order("name DESC")
//	db.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: true})
func (db *DB) OrderByDesc(orderName string) (tx *DB) {
	if orderName == "" {
		return
	}
	tx = db.getInstance()
	tx.Statement.AddClause(clause.OrderBy{
		Columns: []clause.OrderByColumn{{Column: clause.Column{Name: orderName}, Desc: true}},
	})
	return
}

// OrderByDescGBK Order specify order when retrieve records from database
//
//	db.Order("name DESC")
//	db.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: true})
func (db *DB) OrderByDescGBK(orderName string) (tx *DB) {
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

// OrderByStructColumn Order specify order when retrieve records from database
//
// orderNames["OrderByName"]=OrderByColumn (true Desc false ASC)
// orderNames["name"] = true (Desc)
// orderNames["sort"] = false (ASC)
//
//	db.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: true})
func (db *DB) OrderByStructColumn(v interface{}, orderColumns []*OrderColumn) (tx *DB) {
	tx = db.getInstance()
	if len(orderColumns) <= 0 {
		return
	}
	field := getStructFieldTagArray(v)
	for _, item := range orderColumns {
		columnName := orderByString(field, item.Name)
		if len(columnName) <= 0 {
			continue
		}
		if item.GBK {
			orderTypes := "ASC"
			if item.Desc {
				orderTypes = "DESC"
			}
			tx.Statement.AddClause(clause.OrderBy{
				Columns: []clause.OrderByColumn{{
					Column: clause.Column{Name: fmt.Sprintf(" CONVERT(%s USING gbk) %s", columnName, orderTypes), Raw: true},
				}},
			})
		} else {
			tx.Statement.AddClause(clause.OrderBy{
				Columns: []clause.OrderByColumn{{Column: clause.Column{Name: columnName}, Desc: item.Desc}},
			})
		}

	}
	return
}

// OrderByStruct Order specify order when retrieve records from database
//
//	db.Order("name DESC")
//	db.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: true})
func (db *DB) OrderByStruct(v interface{}, orderName string, desc ...bool) (tx *DB) {
	tx = db.getInstance()
	if orderName == "" {
		return
	}
	orderName = orderByString(getStructFieldTagArray(v), orderName)
	if orderName == "" {
		return
	}
	orderDesc := false
	if len(desc) > 0 {
		orderDesc = desc[0]
	}

	tx.Statement.AddClause(clause.OrderBy{
		Columns: []clause.OrderByColumn{{Column: clause.Column{Name: orderName}, Desc: orderDesc}},
	})
	return
}

// OrderByStructGBK Order specify order when retrieve records from database
//
//	db.Order("name DESC")
//	db.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: true})
func (db *DB) OrderByStructGBK(v interface{}, orderName string, desc ...bool) (tx *DB) {
	tx = db.getInstance()
	if orderName == "" {
		return
	}
	orderName = orderByString(getStructFieldTagArray(v), orderName)
	if orderName == "" {
		return
	}
	orderTypes := "ASC"
	if len(desc) > 0 {
		orderTypes = "DESC"
	}
	tx.Statement.AddClause(clause.OrderBy{
		Columns: []clause.OrderByColumn{{
			Column: clause.Column{Name: fmt.Sprintf(" CONVERT(%s USING gbk) %s", orderName, orderTypes), Raw: true},
		}},
	})
	return
}

// SelectByStruct Select specify fields that you want when querying, creating, updating
func (db *DB) SelectByStruct(v interface{}, args ...interface{}) (tx *DB) {
	tx = db.getInstance()
	tx.Select(structToTag(v), args...)
	return
}

func orderByString(field []string, sortName string) string {
	if len(field) <= 0 || sortName == "" {
		return ""
	}
	for _, val := range field {
		if val == sortName {
			return sortName
		}
	}
	return ""
}

// DeleteByNil delete value match given conditions, if the value has primary key, then will including the primary key as condition
func (db *DB) DeleteByNil() (tx *DB) {
	tx = db.getInstance()
	tx.Statement.Dest = ""
	tx.callbacks.Delete().Execute(tx)
	return
}

// Clone get new Instance
func (db *DB) Clone() (tx *DB) {
	tx = &DB{Config: db.Config, Error: db.Error}
	// with clone statement
	tx.Statement = db.Statement.clone()
	tx.Statement.DB = tx
	return tx
}

// Struct Tag
func structToTag(v interface{}) string {
	jsonArray := getStructFieldTagArray(v)
	return strings.Join(jsonArray, ", ")
}

func getStructFieldTagArray(v interface{}) []string {
	jsonArray := make([]string, 0)
	s := reflect.TypeOf(v).Elem() //通过反射获取type定义
	for i := 0; i < s.NumField(); i++ {
		var tag = getStructFieldTag(s.Field(i), "gorm")
		if tag == "-" {
			continue
		}
		if tag != "" {
			tag = replaceKeyWord(tag)
			jsonArray = append(jsonArray, tag)
			continue
		}
		//json,omitempty
		tag = getStructFieldTag(s.Field(i), "json")
		if len(tag) > 0 {
			tag = strings.ReplaceAll(tag, ",omitempty", "")
			tag = replaceKeyWord(tag)
			jsonArray = append(jsonArray, tag)
		}
	}
	return jsonArray
}

func getStructFieldTag(f reflect.StructField, tag string) string {
	return f.Tag.Get(tag)
}

func replaceKeyWord(tag string) string {
	switch tag {
	case "name":
		tag = fmt.Sprintf("`%s`", tag)
	case "describe":
		tag = fmt.Sprintf("`%s`", tag)
	case "status":
		tag = fmt.Sprintf("`%s`", tag)
	}
	return tag
}

type OrderColumn struct {
	Name string `json:"name"` // 字段名称
	Desc bool   `json:"desc"` // 排序类型 true Desc false ASC default ASC
	GBK  bool   `json:"gbk"`  // 是否中文GBK排序
}
