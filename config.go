package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func getDefaultConfigPath() (file, directory string, err error) {
	configDirPath, err := os.UserConfigDir()
	if err != nil {
		return
	}

	directory = filepath.Join(configDirPath, "dxhd")
	file = filepath.Join(directory, "dxhd.sh")
	return
}

func isPathToConfigValid(path string) (isValid bool, err error) {
	stat, err := os.Stat(path)

	if err != nil {
		return
	}

	if !stat.Mode().IsRegular() {
		err = errors.New(fmt.Sprintf("%s is not a regular file", path))
		return
	}

	isValid = true

	return
}

func createDefaultConfig() (err error) {
	var (
		file, directory string
	)

	file, directory, err = getDefaultConfigPath()
	if err != nil {
		return
	}

	_, err = os.Stat(directory)

	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(directory, 0744)
			if err != nil {
				return
			}
		} else {
			return
		}
	}

	_, err = os.Stat(file)

	if err != nil {
		if os.IsNotExist(err) {
			f, e := os.Create(file)
			if e != nil {
				err = e
				return
			}
			f.Write([]byte("#!/bin/sh\n"))
			err = f.Close()
			if err != nil {
				return
			}
		}
	}
	return
}
