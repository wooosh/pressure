// +build !notar tar

package main

import (
    "archive/tar"
    "compress/gzip"
    "io"
    "path/filepath"
    "os"
    "strings"
    "log"
)

func init() {
    ArchiveIndex = append(ArchiveIndex, ArchiveType{
        "tar",
        checkMagic(
            0x101,
            []byte{0x75,0x73,0x74,0x61,0x72,0x00,0x30,0x30},
        ),
        untarWrap,
        checkExt([]string{"tar"}),
        tarWrap,
    })
    ArchiveIndex = append(ArchiveIndex, ArchiveType{
        "tar.gz",
        checkMagic(
            0,
            []byte{0x1F, 0x8B},
        ),
        untargz,
        NoCompressionCheck,
        NoCompression,
   })
}

func tarWrap(f string, target string) {
    tarFile(f, target, false)
}

func tarFile(f string, target string, gzip bool) {
    archive, err := os.Create(target)
    check(err)
    defer archive.Close()

    tw := tar.NewWriter(archive)

    err = filepath.Walk(f, func(path string, info os.FileInfo, err error) error {
        if err != nil  {
            log.Println(err)
            return nil
        }
        if !info.Mode().IsRegular() {
            return nil
        }

        file, err := os.Open(path)
        check(err)
        defer file.Close()

        header, err := tar.FileInfoHeader(info, info.Name())
        check(err)

        header.Name = strings.TrimPrefix(strings.Replace(path, f, "", -1), string(filepath.Separator))

        if header.Name == "" {
            return nil
        }

        err = tw.WriteHeader(header)
        check(err)

        _, err = io.Copy(tw, file)
        check(err)

        return nil
    })
    check(err)
    check(tw.Close())
}

func untargz(f *os.File, target string) {
    gz, err := gzip.NewReader(f)
    check(err)
    untar(gz, target)
    gz.Close()
}

func untarWrap(f *os.File, target string) {
        untar(f, target)
}

func untar(f io.Reader, target string) {
    tr := tar.NewReader(f)
    
    err := os.MkdirAll(target, os.ModePerm)
    check(err)

    for {
        header, err := tr.Next()

        // Return once there are no more files left to decompress
        if err == io.EOF {
            return
        }
        check(err)
        // Create the path for the next file to decompress
        fileTarget := filepath.Join(target + header.Name)

        if header.Typeflag == tar.TypeDir { // Directories
            err := os.MkdirAll(fileTarget, os.ModePerm)
            check(err)
        } else if header.Typeflag == tar.TypeReg { // Normal Files
            os.MkdirAll(filepath.Dir(fileTarget), os.ModePerm)
            dst, err := os.Create(fileTarget)
            check(err)
            _, err = io.Copy(dst, tr)
            check(err)
            dst.Close()
        }
    }
}
