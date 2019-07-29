package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func main() {
	defer func() {
		for {
			time.Sleep(1 * time.Hour)
		}
	}()

	// check golang executable
	fmt.Print("[1/5]check 'go' command ... ")
	_, err := exec.LookPath("go")
	if err != nil {
		fmt.Println("ERROR: can not find 'go' command in PATH")
		return
	}
	fmt.Println("ok")

	// check git
	fmt.Print("[2/5]check 'git' command ... ")
	_, err = exec.LookPath("git")
	if err != nil {
		fmt.Println("ERROR: can not find 'git' command in PATH")
		return
	}
	fmt.Println("ok")

	// mkdir
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		return
	}
	dir := cwd + string(os.PathSeparator) + "dev"
	if _, err = os.Stat(dir); err != nil {
		fmt.Print("[3/5]mkdir dev ...")
		err = os.Mkdir(dir, 0777)
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			return
		}
		fmt.Println("ok")
	} else {
		fmt.Println("[3/5]directory created")
	}

	// download build
	if _, err = os.Stat(dir + "/build/src"); err == nil {
		fmt.Print("[4/5]pulling build codes ... ")
		cmd := exec.Command("git", "pull")
		cmd.Dir = dir + string(os.PathSeparator) + "build"
		err = cmd.Run()
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			return
		}

		fmt.Println("done")
	} else {
		fmt.Print("[4/5]downloading build codes ... ")

		cmd := exec.Command("git", "clone", "https://github.com/TeaWeb/build.git", "--depth", "1")
		cmd.Dir = dir
		err = cmd.Run()
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			return
		}

		fmt.Println("done")
	}

	// download dependents
	fmt.Println("[5/5]downloading dependents ...")
	data, err := ioutil.ReadFile(dir + "/build/src/main/init.sh")
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		return
	}

	match := regexp.MustCompile(`(?U)go_get\s+"(.+)"`).FindAllStringSubmatch(string(data), -1)
	for _, m := range match {
		codePackage := m[1]
		stderr := bytes.NewBuffer([]byte{})

		if _, err = os.Stat(dir + "/build/src/" + codePackage); err == nil {
			continue
		}

		fmt.Println("===")
		fmt.Println("go get \"" + codePackage + "\" ...")
		cmd := exec.Command("go", "get", "-u", "-v", codePackage)
		cmd.Dir = dir + "/build"
		cmd.Stderr = stderr

		cmd.Env = append(os.Environ(), []string{
			"GOPATH=" + dir + "/build/",
			"GO111MODULE=off",
		}...)
		err = cmd.Run()
		if err != nil {
			if stderr.Len() > 0 {
				errString := string(stderr.Bytes())
				if strings.Contains(errString, "version control system") ||
					strings.Contains(errString, "build constraints exclude all Go files") ||
					strings.Contains(errString, "no Go files") {
					continue
				}
				fmt.Println(errString)
			} else {
				fmt.Println("ERROR:", err.Error())
			}
			return
		}
	}
	fmt.Println("... finished")

	// success
	fmt.Println("")
	fmt.Println("==============================")
	fmt.Println("All codes were downloaded to your local disk.")
	fmt.Println("Now, you can create a project from 'dev/build' directory!")

	// start
	{
		fmt.Println("==============================")
		fmt.Println("starting TeaWeb server ...")
		fmt.Println("go run build/src/github.com/TeaWeb/code/main/main.go")
		fmt.Println("==============================")
		stderr := bytes.NewBuffer([]byte{})
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			return
		}
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			return
		}
		cmd := exec.Command("go", "run", dir+"/build/src/github.com/TeaWeb/code/main/main.go")
		cmd.Env = append(os.Environ(), []string{
			"GOPATH=" + dir + "/build/",
			"GO111MODULE=off",
		}...)
		pipe, err := cmd.StderrPipe()
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			return
		}
		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := pipe.Read(buf)
				if n > 0 {
					fmt.Print(string(buf[:n]))
				}
				if err != nil {
					break
				}
			}
		}()
		cmd.Run()
		if stderr.Len() > 0 {
			fmt.Println("ERROR:", string(stderr.Bytes()))
		}
	}
}
