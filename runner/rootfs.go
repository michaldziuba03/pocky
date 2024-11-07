package runner

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const url = "https://dl-cdn.alpinelinux.org/alpine/v3.20/releases/x86_64/alpine-minirootfs-3.20.3-x86_64.tar.gz"
const dest = "/var/lib/pocky"

func Download() {
	fmt.Println("Preparing destination dir...")
	err := os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = os.MkdirAll(dest+"/alpine", os.ModePerm)
	if err != nil {
		log.Fatal(err)
		return
	}

	file, err := os.Create(dest + "/alpine.tar.gz")
	if err != nil {
		log.Fatal(err)
		return
	}

	defer file.Close()

	fmt.Println("Downloading...")
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, _ = io.Copy(file, res.Body)
	fmt.Println("Download finished.")

	cmd := exec.Command("tar", "-xzf", dest+"/alpine.tar.gz", "-C", dest+"/alpine")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Finished.")
}
