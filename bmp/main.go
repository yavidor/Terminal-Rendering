package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"unsafe"
)

func BytesToInteger[T Integer](bytes []byte, offset int) T {
	var result T
	size := min(len(bytes)-offset, int(unsafe.Sizeof(result)))
	for i := range size {
		result |= T(bytes[i+offset]) << ((size - 1 - i) * 8)
	}
	return result
}

type FileHeader struct {
	Signature  uint16
	Size       uint32
	ReservedI  uint16
	ReservedII uint16
	Offset     uint32
}

type BMP struct {
}

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func (fh *FileHeader) read(data []byte, offset int) int {
	v := reflect.ValueOf(fh).Elem()
	for i := range v.NumField() {
		field := v.Field(i)
		switch field.Kind() {
		case reflect.Int8:
			field.SetInt(int64(BytesToInteger[int8](data, offset)))
			offset += 1
		case reflect.Int16:
			field.SetInt(int64(BytesToInteger[int16](data, offset)))
			offset += 2
		case reflect.Int32:
			field.SetInt(int64(BytesToInteger[int32](data, offset)))
			offset += 3
		case reflect.Int64:
			field.SetInt(int64(BytesToInteger[int64](data, offset)))
			offset += 4
		case reflect.Uint8:
			field.SetUint(uint64(BytesToInteger[uint8](data, offset)))
			offset += 1
		case reflect.Uint16:
			field.SetUint(uint64(BytesToInteger[uint16](data, offset)))
			offset += 2
		case reflect.Uint32:
			field.SetUint(uint64(BytesToInteger[uint32](data, offset)))
			offset += 3
		case reflect.Uint64:
			field.SetUint(uint64(BytesToInteger[uint64](data, offset)))
			offset += 4
		}
	}
	return offset
}

const FILE_NAME = "./images/greenland_grid_velo.bmp"

func main() {
	data, err := os.ReadFile(FILE_NAME)
	if err != nil {
		log.Fatal(err)
	}
	header := FileHeader{}
	header.read(data, 0)
	fmt.Printf("%+v\n", header)
	fmt.Printf("%#v\n", header)
}
