package models

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type GrupoIsla struct {
	Id           int    `orm:"column(id);pk;auto"`
	Geometria    string `orm:"column(geometria)"`
	IdGrupoPadre int    `orm:"column(id_grupo_padre);null"`
}

func (t *GrupoIsla) TableName() string {
	return "grupo_isla"
}

func init() {
	orm.RegisterModel(new(GrupoIsla))
}

// AddGrupoIsla insert a new GrupoIsla into database and returns
// last inserted Id on success.
func AddGrupoIsla(m *GrupoIsla) (id int64, err error) {
	o := orm.NewOrm()
	valid := validation.Validation{}
	valid.Required(m.Geometria, "Geometria")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
		}
	} else {
		log.Println("Insert New Register")
		id, err = o.Insert(m)
	}
	return
}

// GetGrupoIslaById retrieves GrupoIsla by Id. Returns error if
// Id doesn't exist
func GetGrupoIslaById(id int) (v *GrupoIsla, err error) {
	o := orm.NewOrm()
	v = &GrupoIsla{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllGrupoIsla retrieves all GrupoIsla matches certain condition. Returns empty list if
// no records exist
func GetAllGrupoIsla(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(GrupoIsla))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []GrupoIsla
	qs = qs.OrderBy("id")
	if _, err := qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateGrupoIsla updates GrupoIsla by Id and returns error if
// the record to be updated doesn't exist
func UpdateGrupoIslaById(m *GrupoIsla) (err error) {
	o := orm.NewOrm()
	v := GrupoIsla{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteGrupoIsla deletes GrupoIsla by Id and returns error if
// the record to be deleted doesn't exist
func DeleteGrupoIsla(id int) (err error) {
	o := orm.NewOrm()
	v := GrupoIsla{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&GrupoIsla{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
