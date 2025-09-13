package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"unsafe"
)

const FILE_NAME = "./images/greenland_grid_velo.bmp"

type Stucture interface {
}

func BytesToInteger[T Integer](bytes []byte, size int) T {
	var result T
	size = min(len(bytes), size)
	for i := range size {
		result |= T(bytes[i]) << ((size - 1 - i) * 8)
	}
	return result
}

type chunk[T Integer] struct {
	Size int
	Data T
}

func (c *chunk[T]) Read(b *bufio.Reader) error {
	rawData := make([]byte, c.Size)
	for i := range c.Size {
		nextByte, err := b.ReadByte()
		if err != nil {
			return err
		}
		rawData[i] = nextByte
	}
	c.Data = BytesToInteger[T](rawData, c.Size)
	return nil
}

func InitChunk[T Integer]() *chunk[T] {
	var zero T
	return &chunk[T]{
		Size: int(unsafe.Sizeof(zero)),
		Data: 0,
	}
}

type FileHeader struct {
	Signature  *chunk[uint16]
	Size       *chunk[uint32]
	ReservedI  *chunk[uint16]
	ReservedII *chunk[uint16]
	Offset     *chunk[uint32]
}

func InitHeader() *FileHeader {
	return &FileHeader{
		Signature:  InitChunk[uint16](),
		Size:       InitChunk[uint32](),
		ReservedI:  InitChunk[uint16](),
		ReservedII: InitChunk[uint16](),
		Offset:     InitChunk[uint32](),
	}
}

type BMP struct {
}

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func readStructure(reader *bufio.Reader, st Stucture) {
	v := reflect.ValueOf(st).Elem()
	for i := range v.NumField() {
		field := v.Field(i)
		field.MethodByName("Read").Call([]reflect.Value{reflect.ValueOf(reader)})
	}
}

func main() {
	data, err := os.Open(FILE_NAME)
	reader := bufio.NewReader(data)
	if err != nil {
		log.Fatal(err)
	}
	header := InitHeader()
	readStructure(reader, header)
	fmt.Printf("%x\n%x\n", header.Signature.Data, header.Size.Data)
}
