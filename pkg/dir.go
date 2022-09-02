package pkg

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func ValidateDir(createDirIfNotExist bool, paths ...string) error {
	for _, path := range paths {
		if !IsFileOrFolderExist(path) && createDirIfNotExist {
			if err := os.MkdirAll(path, 0700); err != nil {
				return err
			}
		}
	}

	return nil
}

func IsFileOrFolderExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}

	return false, err
}

func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	return os.Chmod(dst, srcinfo.Mode())
}

func CopyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}

	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}

	return nil
}

func DeleteDir(folderPath string) error {
	if err := os.RemoveAll(folderPath); err != nil {
		return err
	}

	return nil
}

func GetAbsFileUrl(filePath string) (string, error) {
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("file://%v", filePath), nil
}

func GetListFolderFromDir(dirpath string) ([]string, error) {
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return nil, err
	}

	var res = []string{}
	for _, f := range files {
		if f.IsDir() {
			res = append(res, f.Name())
		}
	}

	return res, nil
}
