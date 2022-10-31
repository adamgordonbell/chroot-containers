package main

import (
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("/hello")
	must(syscall.Chroot("./assets/hello-world_fs"))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
