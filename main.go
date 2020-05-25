package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
)

func main() {
	if runtime.GOOS != "linux" {
		log.Fatal("dxhd is only supported on linux")
		os.Exit(1)
	}

	customConfigPath := flag.String("c", "", "reads the config from custom path")

	flag.Parse()

	var (
		configStat     os.FileInfo
		configFilePath string
		err            error
	)

	// we default to "", no need to make sure it's not nil
	if *customConfigPath != "" {
		configStat, err = os.Stat(*customConfigPath)

		if err != nil {
			log.Fatalf("can't read from %s file (%s)", *customConfigPath, err.Error())
			os.Exit(1)
		}

		if !configStat.Mode().IsRegular() {
			log.Fatalf("%s is not a regular file", configFilePath)
			os.Exit(1)
		}

		configFilePath = *customConfigPath
	} else {
		configDirPath, err := os.UserConfigDir()
		if err != nil {
			log.Fatalf("couldn't get config directory (%s)", err.Error())
			os.Exit(1)
		}

		configDirPath = filepath.Join(configDirPath, "dxhd")
		configFilePath = filepath.Join(configDirPath, "dxhd.sh")

		configStat, err = os.Stat(configDirPath)

		if err != nil {
			if os.IsNotExist(err) {
				err = os.Mkdir(configDirPath, 0744)
				if err != nil {
					log.Fatalf("couldn't create %s directory (%s)", configDirPath, err.Error())
					os.Exit(1)
				}
				configStat, err = os.Stat(configDirPath)
				if err != nil {
					log.Fatalf("error occurred - %s", err.Error())
					os.Exit(1)
				}
			} else {
				log.Fatalf("error occurred - %s", err.Error())
				os.Exit(1)
			}
		}

		if !configStat.Mode().IsDir() {
			log.Fatalf("%s is not a directory", configDirPath)
			os.Exit(1)
		}

		configStat, err = os.Stat(configFilePath)

		if err != nil {
			if os.IsNotExist(err) {
				file, err := os.Create(configFilePath)
				if err != nil {
					log.Fatalf("couldn't create %s file (%s)", configFilePath, err.Error())
					os.Exit(1)
				}
				// write to the file, and exit
				file.Write([]byte("#!/bin/sh\n"))
				err = file.Close()
				if err != nil {
					log.Fatalf("can't close newly created file %s (%s)", configFilePath, err.Error())
					os.Exit(1)
				}
				os.Exit(0)
			} else {
				log.Fatalf("error occurred - %s", err.Error())
				os.Exit(1)
			}
		}

		if !configStat.Mode().IsRegular() {
			log.Fatalf("%s is not a regular file", configFilePath)
			os.Exit(1)
		}
	}

	var (
		data  []filedata
		shell string
	)
	shell, err = parse(configFilePath, &data)
	if err != nil {
		log.Fatalf("failed to parse file %s (%s)", configFilePath, err.Error())
		os.Exit(0)
	}

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatalf("can not open connection to Xorg (%s)", err.Error())
		os.Exit(1)
	}

	keybind.Initialize(X)

	for _, d := range data {
		err = listenKeybinding(X, shell, d.binding.String(), d.action.String())
		if err != nil {
			log.Printf("error occurred whilst trying to register keybinding %s (%s)", d.binding.String(), err.Error())
		}
	}

	xevent.Main(X)
}
