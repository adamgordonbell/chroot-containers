package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"syscall"

	"github.com/codeclysm/extract"
	_ "github.com/codeclysm/extract"
)

func main() {
	dir := extract1("./assets/hello-world_fs.tar.gz")
	chroot(dir, "/hello")
	fmt.Printf("removing %s\n", dir)
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

func extract1(tar string) string {
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

	prefix := nonAlphanumericRegex.ReplaceAllString(tar, "_")
	dir, err := ioutil.TempDir("", prefix)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("created %s\n", dir)
	}
	must(Untar2(tar, dir))
	fmt.Printf("Extracted to %s\n", dir)
	return dir
}
func Untar2(source string, dst string) error {
	r, err := os.Open(source)
	if err != nil {
		return err
	}
	defer r.Close()

	// data, _ := ioutil.ReadFile(source)
	ctx := context.Background()
	// buffer := bytes.NewBuffer(data)
	return extract.Archive(ctx, r, dst, nil)
}

func Untar(source string, dst string) error {
	fmt.Printf("Extracting %s to %s\n", source, dst)

	r, err := os.Open(source)
	if err != nil {
		return err
	}
	defer r.Close()

	gzr, err := gzip.NewReader(r)
	if err != nil {
		log.Fatal("ExtractTarGz: NewReader failed: ", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

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
