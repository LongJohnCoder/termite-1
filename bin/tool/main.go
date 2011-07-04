package main

import (
	"github.com/hanwen/go-fuse/termite"
	"log"
	"net"
	"os"
	"path/filepath"
	"rpc"
)

const _SOCKET = ".termite-socket"

func main() {
	path := os.Getenv("PATH")

	args := os.Args
	_, base := filepath.Split(args[0])
	binary := ""
	for _, c := range filepath.SplitList(path) {
		_, dirBase := filepath.Split(c)
		if dirBase == "termite" {
			continue
		}

		try := filepath.Join(c, base)
		_, err := os.Stat(try)
		if err == nil {
			binary = try
		}
	}
	if binary == "" {
		log.Fatal("could not find", base)
	}
	args[0] = binary
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Getwd", err)
	}
	req := termite.WorkRequest{
	Binary: binary,
	Argv: args,
	Env: os.Environ(),
	Dir: wd,
	}

	socket := os.Getenv("TERMITE_SOCKET")
	if socket == "" {
		socketPath := wd
		for socketPath != "/" {
			cand := filepath.Join(socketPath, _SOCKET)
			fi, _ := os.Lstat(cand)
			if fi != nil && fi.IsSocket() {
				socket = cand
				break
			}
			socketPath = filepath.Clean(filepath.Join(socketPath, ".."))
		}
	}

	conn, err := net.Dial("unix", socket)
	if err != nil {
		log.Fatal("Dial:", err)
	}

	client := rpc.NewClient(conn)

	rep := termite.WorkReply{}
	err = client.Call("LocalMaster.Run", &req, &rep)
	if err != nil {
		log.Fatal("LocalMaster.Run: ", err)
	}

	os.Stdout.Write([]byte(rep.Stdout))
	os.Stderr.Write([]byte(rep.Stderr))
	// TODO -something with signals.
	os.Exit(rep.Exit.ExitStatus())
}



