package shadowbins

import (
	"log"
	"net/rpc"
	"os"
	"path/filepath"
	"strconv"
)

type Command struct {
	Bin        string
	Args       []string
	Identifier string
}

func buildCommadFromArgs() Command {
	c := Command{}
	c.Bin = os.Args[1]
	args := os.Args[2:len(os.Args)]
	for _, arg := range args {
		c.Args = append(c.Args, resolveFilePathInArg(arg))
	}
	return c
}

func resolveFilePathInArg(arg string) string {
	path, err := filepath.Abs(arg)
	if err != nil {
		log.Fatal("Failed to resolve path", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// this argument is not a file path
		return arg
	}
	return path
}

func CallRemoteCommand() {
	conf := readConfig()
	addr := "127.0.0.1:" + strconv.Itoa(conf.Hyper.Port)

	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Dialing:", err)
	}

	c := buildCommadFromArgs()
	c.Identifier = conf.ShadowIdentifier

	var reply int
	err = client.Call("Launcher.Launch", c, &reply)
	if err != nil {
		log.Fatal("Launch error:", err)
	}
}
