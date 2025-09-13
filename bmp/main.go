package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"unsafe"
)

// const FILE_NAME = "./images/greenland_grid_velo.bmp"
const FILE_NAME = "./images/video-001.bmp"

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
	Init()
}

func BytesToInteger[T Integer](bytes []byte, size int) T {
	var result T
	size = min(len(bytes), size)
	for i := range size {
		result |= T(bytes[i]) << (i * 8)
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

func (fh *FileHeader) Init() {
	fh.Signature = InitChunk[uint16]()
	fh.Size = InitChunk[uint32]()
	fh.ReservedI = InitChunk[uint16]()
	fh.ReservedII = InitChunk[uint16]()
	fh.Offset = InitChunk[uint32]()
}

type InfoHeader struct {
	Size                *chunk[uint32]
	Width               *chunk[uint32]
	Height              *chunk[uint32]
	Planes              *chunk[uint16]
	BitCount            *chunk[uint16]
	Compression         *chunk[uint32]
	ImageSize           *chunk[uint32]
	XPixelsPerMeter     *chunk[uint32]
	YPixelsPerMeter     *chunk[uint32]
	ColorCount          *chunk[uint32]
	ImportantColorCount *chunk[uint32]
}

func (ih *InfoHeader) Init() {
	ih.Size = InitChunk[uint32]()
	ih.Width = InitChunk[uint32]()
	ih.Height = InitChunk[uint32]()
	ih.Planes = InitChunk[uint16]()
	ih.BitCount = InitChunk[uint16]()
	ih.Compression = InitChunk[uint32]()
	ih.ImageSize = InitChunk[uint32]()
	ih.XPixelsPerMeter = InitChunk[uint32]()
	ih.YPixelsPerMeter = InitChunk[uint32]()
	ih.ColorCount = InitChunk[uint32]()
	ih.ImportantColorCount = InitChunk[uint32]()
}

type BMP struct {
	fileHeader *Structure
	infoHeader *Structure
}

func readStructure(reader *bufio.Reader, st Structure) {
	v := reflect.ValueOf(st).Elem()
	for i := range v.NumField() {
		field := v.Field(i)
		field.MethodByName("Read").Call([]reflect.Value{reflect.ValueOf(reader)})
	}
}

func printStructure(st Structure) error {
	v := reflect.ValueOf(st).Elem()
	fmt.Println(v.Type().Name())
	for i := range v.NumField() {
		fieldValue := v.Field(i).Elem()
		fieldType := v.Type().Field(i)
		if fieldValue.FieldByName("Data").CanUint() {
			fmt.Printf("\t%s - %s: 0x%X | %d\n", fieldType.Name, fieldValue.Field(1).Type(), fieldValue.FieldByName("Data").Uint(), fieldValue.FieldByName("Data").Uint())
		} else {
			return fmt.Errorf("field %s of struct %s has invalid type %s", fieldType.Name, v.Type().Name(), fieldValue.Field(1).Type())
		}
	}
	return nil
}

func main() {
	data, err := os.Open(FILE_NAME)
	reader := bufio.NewReader(data)
	if err != nil {
		log.Fatal(err)
	}
	fileHeader := &FileHeader{}
	fileHeader.Init()
	readStructure(reader, fileHeader)
	err = printStructure(fileHeader)
	if err != nil {
		panic(err)
	}
	infoHeader := &InfoHeader{}
	infoHeader.Init()
	readStructure(reader, infoHeader)
	err = printStructure(infoHeader)
	if err != nil {
		panic(err)
	}
}
