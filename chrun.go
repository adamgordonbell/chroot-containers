package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"syscall"

	"github.com/codeclysm/extract"
)

func main() {
	switch os.Args[1] {
	case "run":
		tar := os.Args[2]
		cmd := os.Args[3]
		dir := createTempDir(tar)
		must(unTar(tar, dir))
		chroot(dir, cmd)
		defer must(os.RemoveAll(dir))
	case "pull":
		image := os.Args[2]
		pullImage(image)
	default:
		panic("what?")
	}
}

func pullImage(image string) {
	cmd := exec.Command("./pull", image)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	must(cmd.Run())
}

func chroot(root string, call string) {
	fmt.Printf("Running %s in %s\n", call, root)
	cmd := exec.Command(call)
	must(syscall.Chroot(root))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	must(cmd.Run())
}

func createTempDir(name string) string {
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

	prefix := nonAlphanumericRegex.ReplaceAllString(name, "_")
	dir, err := ioutil.TempDir("", prefix)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("created %s\n", dir)
	}
	return dir
}

func unTar(source string, dst string) error {
	fmt.Printf("Extracting %s %s\n", source, dst)
	r, err := os.Open(source)
	if err != nil {
		return err
	}
	defer r.Close()

	ctx := context.Background()
	return extract.Archive(ctx, r, dst, nil)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
