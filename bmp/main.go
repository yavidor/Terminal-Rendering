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

type CompressionType int

const (
	BI_RGB            CompressionType = 0
	BI_RLE8           CompressionType = 1
	BI_RLE4           CompressionType = 2
	BI_BITFIELDS      CompressionType = 3
	BI_JPEG           CompressionType = 4
	BI_PNG            CompressionType = 5
	BI_ALPHABITFIELDS CompressionType = 6
	BI_CMYK           CompressionType = 11
	BI_CMYKRLE8       CompressionType = 12
	BI_CMYKRLE4       CompressionType = 13
)

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Structure interface {
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

// Overall 14 bytes
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

type InfoHeader struct {
	Size                *chunk[uint32]
	Width               *chunk[int32]
	Height              *chunk[int16]
	Planes              *chunk[uint16]
	BitCount            *chunk[uint16]
	Compression         *chunk[uint32]
	ImageSize           *chunk[uint32]
	XPixelsPerMeter     *chunk[int32]
	YPixelsPerMeter     *chunk[int32]
	ColorCount          *chunk[uint32]
	ImportantColorCount *chunk[uint32]
}

type BMP struct {
	Header *FileHeader
}

func readStructure(reader *bufio.Reader, st Structure) {
	v := reflect.ValueOf(st).Elem()
	for i := range v.NumField() {
		field := v.Field(i)
		field.MethodByName("Read").Call([]reflect.Value{reflect.ValueOf(reader)})
	}
}

func printStructure(st Structure) {
	v := reflect.ValueOf(st).Elem()
	for i := range v.NumField() {
		field := v.Field(i)
		if field.Kind() == reflect.Int {
			fmt.Printf("%X\n", field.Elem().FieldByName("Data").Int())
		} else {

			fmt.Printf("%s - 0x%X\n", v.Type().Field(i).Name, field.Elem().FieldByName("Data").Uint())
		}
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
	printStructure(header)
}
