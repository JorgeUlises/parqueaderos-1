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

type Vehiculo struct {
	Id            int          `orm:"column(id);pk;auto"`
	Placa         string       `orm:"column(placa)"`
	IdNfc         int16        `orm:"column(id_nfc)"`
	IdPropietario *Propietario `orm:"column(id_propietario);rel(fk)"`
}

func (t *Vehiculo) TableName() string {
	return "vehiculo"
}

func init() {
	orm.RegisterModel(new(Vehiculo))
}

// AddVehiculo insert a new Vehiculo into database and returns
// last inserted Id on success.
func AddVehiculo(m *Vehiculo) (id int64, err error) {
	o := orm.NewOrm()
	valid := validation.Validation{}
	valid.Required(m.Placa, "Placa")
	valid.MaxSize(m.Placa, 6, "PlacaMax")
	valid.Required(m.IdNfc, "IdNfc")
	valid.AlphaNumeric(m.Placa, "PlacaAlphaNum")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
		}
	} else {
		log.Println("Insert New Register")
		m.IdPropietario, _ = GetPropietarioById(m.IdPropietario.Id)
		id, err = o.Insert(m)
	}
	return
}

// GetVehiculoById retrieves Vehiculo by Id. Returns error if
// Id doesn't exist
func GetVehiculoById(id int) (v *Vehiculo, err error) {
	o := orm.NewOrm()
	v = &Vehiculo{Id: id}
	if err = o.Read(v); err == nil {
		v.IdPropietario, _ = GetPropietarioById(v.IdPropietario.Id)
		return v, nil
	}
	return nil, err
}

// GetVehiculoByIdNfc retrieves Vehiculo by IdNfc. Returns error if
// Id doesn't exist
func GetVehiculoByIdNfc(id int16) (v *Vehiculo, err error) {
	o := orm.NewOrm()
	v = &Vehiculo{IdNfc: id}
	if err = o.Read(v, "IdNfc"); err == nil {
		v.IdPropietario, _ = GetPropietarioById(v.IdPropietario.Id)
		return v, nil
	}
	return nil, err
}

// GetAllVehiculo retrieves all Vehiculo matches certain condition. Returns empty list if
// no records exist
func GetAllVehiculo(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Vehiculo))
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

	var l []Vehiculo
	qs = qs.OrderBy("id")
	if _, err := qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				v.IdPropietario, _ = GetPropietarioById(v.IdPropietario.Id)
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

// UpdateVehiculo updates Vehiculo by Id and returns error if
// the record to be updated doesn't exist
func UpdateVehiculoById(m *Vehiculo) (err error) {
	o := orm.NewOrm()
	v := Vehiculo{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteVehiculo deletes Vehiculo by Id and returns error if
// the record to be deleted doesn't exist
func DeleteVehiculo(id int) (err error) {
	o := orm.NewOrm()
	v := Vehiculo{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Vehiculo{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
