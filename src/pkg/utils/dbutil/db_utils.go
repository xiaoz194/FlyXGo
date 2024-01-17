package dbutil

import (
	"encoding/json"
	"reflect"
)

// AssignMapToStruct 把map的值赋值给model的结构体 并写入数据库
func AssignMapToStruct(m map[string]interface{}, model interface{}) (interface{}, error) {
	jsonStr, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonStr, &model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

// CreateRecord 通用的插入消息记录方法 应用场景:插入部分数据 data 具体需要插入的数据; record 通用的记录结构
func CreateRecord(data map[string]interface{}, record interface{}) error {
	// 获取传入参数的类型
	recordType := reflect.TypeOf(record)
	// 根据传入参数的类型创建一个新的对象
	newRecord := reflect.New(recordType.Elem()).Interface()
	// 将传入的数据映射到新的对象中
	_, err := AssignMapToStruct(data, newRecord)
	if err != nil {
		return err
	}
	// 将新的对象转换为传入参数的类型
	recordValue := reflect.ValueOf(record).Elem()
	newValue := reflect.ValueOf(newRecord).Elem()
	recordValue.Set(newValue)
	// 将新记录插入数据库
	if err := DB.Create(record).Error; err != nil {
		return err
	}
	return nil
}
