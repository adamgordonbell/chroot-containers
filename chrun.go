package main

import (
	"context"
	"errors"
	"fmt"
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
		image := os.Args[2]
		tar := fmt.Sprintf("./assets/%s.tar.gz", image)

		if _, err := os.Stat(tar); errors.Is(err, os.ErrNotExist) {
			panic(err)
		}

		cmd := ""
		if len(os.Args) > 3 {
			cmd = os.Args[3]
		} else {
			buf, err := os.ReadFile(fmt.Sprintf("./assets/%s-cmd", image))
			if err != nil {
				panic(err)
			}
			cmd = string(buf)
		}

		dir := createTempDir(tar)
		defer os.RemoveAll(dir)
		must(unTar(tar, dir))
		chroot(dir, cmd)
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
	//Hold onto old root
	oldrootHandle, err := os.Open("/")
	if err != nil {
		panic(err)
	}
	defer oldrootHandle.Close()

	//Create command
	fmt.Printf("Running %s in %s\n", call, root)
	cmd := exec.Command(call)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//New Root time
	must(syscall.Chdir(root))
	must(syscall.Chroot(root))

	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	//Go back to old root
	//So that we can clean up the temp dir
	must(syscall.Fchdir(int(oldrootHandle.Fd())))
	must(syscall.Chroot("."))

}

func createTempDir(name string) string {
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

	prefix := nonAlphanumericRegex.ReplaceAllString(name, "_")
	dir, err := os.MkdirTemp("", prefix)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("created %s\n", dir)
	return dir
}

func unTar(source string, dst string) error {
	// fmt.Printf("Extracting %s %s\n", source, dst)
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
