package main

import (
    "os"
    "bytes"
    "io"
    "log"
    "strings"
)

func init() {
    log.SetPrefix(os.Args[0] + ": ")
    log.SetFlags(0)
}

type DecompressionCheck func (f *os.File) bool
type Decompressor func (f *os.File, target string)

type CompressionCheck func (filename string) bool
type Compressor func (filename string, target string)

type ArchiveType struct {
    Name string

    CheckDecompress DecompressionCheck
    Decompress Decompressor

    CheckCompression CompressionCheck
    Compress Compressor
}

func check(e error) {
    if e != nil {
        log.Fatal(e)
    }
}

var ArchiveIndex []ArchiveType

func main() {
    if len(os.Args) < 3 {
        log.Fatal("usage: pressure source target")
    }
    f, err := os.Open(os.Args[1])
    check(err)

    var compressor Compressor
    for _, archiveType := range ArchiveIndex {
        if archiveType.CheckDecompress(f) {
            archiveType.Decompress(f, os.Args[2] + "/")
            return
        } else if archiveType.CheckCompression(os.Args[2]) {
            compressor = archiveType.Compress
        }
    }
    if compressor != nil {
        compressor(os.Args[1], os.Args[2])
    } else {
        log.Fatal("No decompressor available for file " + os.Args[1])
    }
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

// Generates a function that checks file extension
// We take a slice of strings because formats like 
// zip have many different names (jar, epub, docx etc)
func checkExt(extensions []string) CompressionCheck {
    return func (filename string) bool {
        for _, ext := range extensions {
            // strings.HasSuffix is used instead of filepath.Ext because
            // filepath.Ext will only capture the final dot, so formats
            // like tar.gz won't be able to be detected
            if strings.HasSuffix(filename, ext) {
                return true
            }
        }
        return false
    }
}

func NoDecompressionCheck(f *os.File) bool {
    return false
}

func NoCompressionCheck(f string) bool {
    return false
}

func NoCompression(f string, target string) {}
func NoDecompression(f *os.File, target string) {}
