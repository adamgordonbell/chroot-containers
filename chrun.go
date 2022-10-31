package main

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func main() {
	dir := extract("./assets/hello-world_fs.tar.gz")
	chroot(dir, "/hello")
	fmt.Printf("removing %s", dir)
	defer os.RemoveAll(dir)

}

func chroot(root string, call string) {
	cmd := exec.Command(call)
	must(syscall.Chroot(root))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

func extract(tar string) string {
	dir, err := ioutil.TempDir("", tar)
	if err != nil {
		log.Fatal(err)
	}
	must(Untar(tar, dir))
	fmt.Printf("Extracted to %s", dir)
	return dir
}

func Untar(source string, dst string) error {
	r, err := os.Open(source)
	if err != nil {
		return err
	}
	defer r.Close()

	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			panic(err)
		case header == nil:
			continue
		}

		target := filepath.Join(dst, header.Name)

		switch header.Typeflag {

		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					panic(err)
				}
			}

		case tar.TypeReg:
			f, err := os.Create(target)
			if err != nil {
				panic(err)
			}

			if _, err := io.Copy(f, tr); err != nil {
				panic(err)
			}

			f.Close()
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
