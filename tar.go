// +build !notar tar

package main

import (
    "archive/tar"
    "io"
    "path/filepath"
    "os"
)

func init() {
    ArchiveIndex = append(ArchiveIndex, ArchiveType{
        "tar",
        checkMagic(
            0x101,
            []byte{0x75,0x73,0x74,0x61,0x72,0x00,0x30,0x30},
        ),
        untar,
    })
}

func untar(f *os.File, target string) {
    tr := tar.NewReader(f)

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
            dst, err := os.Create(fileTarget)
            check(err)
            _, err = io.Copy(dst, tr)
            check(err)
            dst.Close()
        }
    }
}
