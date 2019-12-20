package earPicking

import (
	"database/sql"
	"errors"
	"reflect"
	"strconv"
)

type (

	// queryRes 查询结果
	queryRes struct {
		data []map[string]string // 存储数据库查询数据
		err  error               // 异常

	}

	// tableInfo 表信息结构体
	tableInfo struct {
		tabName       string // 表名
		primaryKey    string // 主键标识
		limit         int    // 分页起始页
		offset        int    // 分页结束页
		sqlTemp       string // sql临时存储
		whereTemp     string // sql条件临时存储

		orderByTemp   string // 排序存储
		groupByTemp   string // 分组聚合
	}
)

var typeRegistry = make(map[string]reflect.Type)

// formatRes 处理查询响应数据
// @param rows 行是查询的结果
func (qr *queryRes) formatRes(rows *sql.Rows) error {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([][]byte, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	recordList := make([]map[string]string, 0)
	for rows.Next() {

		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		checkErr(err)
		record := make(map[string]string)
		for i, col := range values {

			if col != nil {
				record[columns[i]] = string(col)
			}
		}
		recordList = append(recordList, record)

	}

	qr.data = recordList
	// jsonStr, err := util.ToJsonStr(recordList)
	// fmt.Println(jsonStr)
	// checkErr(err)
	rows.Close()

	return nil

}



// Unique ...
func (qr *queryRes) Unique(in interface{}) error {
	if len(qr.data) > 0 {
		return qr.mapping(qr.data[0], reflect.ValueOf(in))
	}
	return nil
}

// List 集合对象映射 (公共方法)
func (qr *queryRes) List(in interface{}) error {

	if qr.err != nil {
		return qr.err
	}

	length := len(qr.data)


	if length > 0 {
		v := reflect.ValueOf(in)
		qr.list(v)
	}

	return nil
}

func (qr *queryRes) unqiue(v reflect.Value)error{
	if len(qr.data) > 0 {
		return qr.mapping(qr.data[0], v)
	}
	return nil
}

// List 集合对象映射（私有方法）
func (qr *queryRes) list(v reflect.Value) error{

	// 1.reflect.ValueOf->获取接口保管的具体值(实例化)
	// 2.Indirect->获取该值的指针值
	sliceValue := reflect.Indirect(v)

	// 如果该对象类型不是切片类型，则返回类型错误
	if sliceValue.Kind() != reflect.Slice {
		return ErrorNeedPointerToSlice
	}

	// 获取切片集合中的类型
	sliceElementType := sliceValue.Type().Elem()

	for _, results := range qr.data {
		// 根据切片集合中的类型，创建新的实体
		newValue := reflect.New(sliceElementType)
		// 映射数据
		err := qr.mapping(results, newValue)
		if err != nil {
			return err
		}

		//fmt.Printf("sliceValue ： %+v\n", sliceValue)

		// 获取映射后的对象实例的指针值
		rTmep := reflect.Indirect(reflect.ValueOf(newValue.Interface()))

		//fmt.Printf("rTmep ： %+v\n", rTmep)
		// 将映射后的对象实例的指针值复制到源对象指针中
		sliceValue.Set(reflect.Append(sliceValue, rTmep))

		//fmt.Printf("sliceValue111 ： %+v\n", sliceValue)
	}


	return nil
}



// 处理结构体于查询数据之间的映射关系
func (qr *queryRes) mapping(m map[string]string, v reflect.Value) error {

	t := v.Type()
	val := v.Elem()
	typ := t.Elem()

	// 检查结构体的数据类型
	if !val.IsValid() {
		return errors.New("数据类型不正确")
	}

	for i := 0; i < val.NumField(); i++ {
		value := val.Field(i)
		kind := value.Kind()
		// 获取字段注解
		tag := typ.Field(i).Tag.Get("col")

		if len(tag) > 0 {
			meta, ok := m[tag]
			if !ok {

				continue
			}
			// 判断字段是否有读写权限
			if !value.CanSet() {
				return errors.New("结构体字段没有读写权限")
			}

			if len(meta) == 0 {
				continue
			}

			if kind == reflect.String {
				value.SetString(meta)
			} else if kind == reflect.Float32 {
				f32, err := strconv.ParseFloat(meta, 32)

				if err != nil {
					return err
				}
				value.SetFloat(f32)
			} else if kind == reflect.Float64 {
				f64, err := strconv.ParseFloat(meta, 64)

				if err != nil {
					return err
				}
				value.SetFloat(f64)
			} else if kind == reflect.Int64 {
				i64, err := strconv.ParseInt(meta, 10, 64)
				if err != nil {
					return err
				}
				value.SetInt(i64)
			} else if kind == reflect.Int {
				i, err := strconv.Atoi(meta)
				if err != nil {
					return err
				}

				value.SetInt(int64(i))
			} else if kind == reflect.Bool {
				b, err := strconv.ParseBool(meta)
				if err != nil {
					return err
				}
				value.SetBool(b)
			} else {
				return errors.New("没有对应的数据类型")
			}
		}
	}
	return nil
}
