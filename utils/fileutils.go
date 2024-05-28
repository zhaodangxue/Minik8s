package utils

import (
	"io"
	"os"
	"path/filepath"
)

func CopyDir(src string, dst string) error {
    srcInfo, err := os.Stat(src)
    if err != nil {
        return err
    }

    err = os.MkdirAll(dst, srcInfo.Mode())
    if err != nil {
        return err
    }

    entries, err := os.ReadDir(src)
    if err != nil {
        return err
    }

    for _, entry := range entries {
        srcPath := filepath.Join(src, entry.Name())
        dstPath := filepath.Join(dst, entry.Name())

        if entry.IsDir() {
            err = CopyDir(srcPath, dstPath)
            if err != nil {
                return err
            }
        } else {
            err = CopyFile(srcPath, dstPath)
            if err != nil {
                return err
            }
        }
    }

    return nil
}

func CopyFile(src string, dst string) error {
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer dstFile.Close()

    _, err = io.Copy(dstFile, srcFile)
    if err != nil {
        return err
    }

    srcInfo, err := os.Stat(src)
    if err != nil {
        return err
    }

    err = os.Chmod(dst, srcInfo.Mode())
    if err != nil {
        return err
    }

    return nil
}