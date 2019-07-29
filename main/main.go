package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	// check golang executable
	fmt.Println("check 'go' command ...")
	_, err := exec.LookPath("go")
	if err != nil {
		fmt.Println("can not find 'go' command in PATH")
		return
	}

	// check git
	fmt.Println("check 'git' command ...")
	_, err = exec.LookPath("git")
	if err != nil {
		fmt.Println("can not find 'git' command in PATH")
		return
	}

	// mkdir
	dir := "dev"
	if _, err = os.Stat(dir); err != nil {
		fmt.Println("mkdir dev ...")
		err = os.Mkdir(dir, 0777)
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			return
		}
	}

	// download build
	if _, err = os.Stat(dir + "/build/src"); err == nil {
		fmt.Println("pulling build codes ...")
		cmd := exec.Command("git", "pull")
		cmd.Dir = dir
		err = cmd.Run()
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			return
		}

		fmt.Println("done")
	} else {
		fmt.Println("downloading build codes ...")

		cmd := exec.Command("git", "clone", "https://github.com/TeaWeb/build.git")
		cmd.Dir = dir
		err = cmd.Run()
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			return
		}

		fmt.Println("done")
	}

	// download dependents
	fmt.Println("downloading dependents ...")
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
			"GOPATH=/Users/liuxiangchao/Documents/Projects/pp/apps/TeaWebGit/code-setup/dev/build/",
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

	// success
	fmt.Println("")
	fmt.Println("==============================")
	fmt.Println("All codes were downloaded to your local disk.")
	fmt.Println("Now, you can create a project from 'dev/build' directory!")
}
