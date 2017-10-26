package utility

import (
	"reflect"
	"fmt"
	"os"
	"github.com/axgle/mahonia"
	"bufio"
)


// To check if the list contains the elem
func Contains(list interface{}, elem interface{}) bool {
	value := reflect.ValueOf(list)
	if value.Kind() != reflect.Slice {
		fmt.Fprintf(os.Stderr, "Input type is not an array or slice type: %v, kind:%s", value, value.Kind())
		return false
	}

	for i:=0; i<value.Len();i++ {
		if value.Index(i).Interface() == elem.(string) {
			return true
		}
	}

	return false
}

// To get map keys
func Keys(i interface{}) (keys []string) {
	v := reflect.ValueOf(i)

	if v.Kind() != reflect.Map {
		fmt.Fprintf(os.Stderr, "Input type is not a map type: %v", v)
		return nil
	}

	for _,key := range v.MapKeys() {
		keys = append(keys, key.Interface().(string))
	}

	return keys
}

// Write one line to file.
func WriteToFile(path string, line string) error{
	file, err:= os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0777)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARN: Open file %s failed, %s\n", path, err)
		return err
	}
	defer file.Close()

	encoder := mahonia.NewEncoder("gbk")
	writer := bufio.NewWriter(encoder.NewWriter(file))
	writer.WriteString(line + "\n")
	writer.Flush()
	//io.Copy(writer, strings.NewReader(line))
	return nil
}


func IsFileExist(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}