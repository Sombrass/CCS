package shadowbins

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Launcher struct {
	conf *tomlConfig
}

func (l *Launcher) convertBin(bin string) (newBin string, err error) {
	newBin, found := l.conf.Hyper.Binaries[bin]
	if !found {
		err = errors.New("Unknown binary " + bin)
	}
	return newBin, err
}

func (l *Launcher) convertArgs(shadowId string, args []string) ([]string, error) {
	var err error
	shadowConf, found := l.conf.Shadows[shadowId]
	if !found {
		err = errors.New("Unknown shadow identifier: " + shadowId)
	}

	reverseMap := make(map[string]string)
	for _, m := range shadowConf.DirMap {
		reverseMap[m["shadow_path"]] = m["hyper_path"]
	}
	var newArgs []string
	for _, arg := range args {
		newArg := arg
		for sp, hp := range reverseMap {
			if strings.HasPrefix(arg, sp) {
				newArg = filepath.Clean(strings.Replace(arg, sp, hp, 1))
				break
			}
		}
		newArgs = append(newArgs, newArg)
	}
	return newArgs, err
}

func (l *Launcher) Launch(c Command, reply *int) error {
	bin, err := l.convertBin(c.Bin)
	if err != nil {
		log.Fatal(err)
	}

	args, err := l.convertArgs(c.Identifier, c.Args)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Executing ", bin)
	cmd := exec.Command(bin, args...)
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(c)
	return err
	return nil
}

func ServeRPC() {
	conf := readConfig()
	launcher := new(Launcher)
	launcher.conf = &conf
	rpc.Register(launcher)
	addr := "127.0.0.1:" + strconv.Itoa(conf.Hyper.Port)
	log.Println("Starting RPC Server on", addr)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(c)
	}
}
