package main

import (
    "os"
    "bytes"
    "io"
    "fmt"
)

type DecompressionCheck func (f *os.File) bool
type Decompressor func (f *os.File, target string)

type ArchiveType struct {
    Name string
    CheckDecompress DecompressionCheck
    Decompress Decompressor
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}



var ArchiveIndex []ArchiveType

func main() {
    if len(os.Args) < 3 {
        fmt.Println("usage: pressure source target")
        return
    }
    f, err := os.Open(os.Args[1])
    check(err)

    for _, archiveType := range ArchiveIndex {
        if archiveType.CheckDecompress(f) {
            archiveType.Decompress(f, os.Args[2] + "/")
            return
        }
    }
    fmt.Println("No decompressor available")
}

// Generates a function that checks magic values at a certain offset
func checkMagic(offset int64, sequence []byte) DecompressionCheck {
    return func (f *os.File) bool {
        fileMagic := make([]byte, len(sequence))
        _, err := f.ReadAt(fileMagic, offset)
        if err != nil && err != io.EOF {
             panic(err)
        }
        return bytes.Equal(fileMagic, sequence)
    }
}


