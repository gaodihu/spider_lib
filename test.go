package main

import (
	"fmt"
	//"os"
	"os/exec"
	//"regexp"
	"strconv"
	"strings"

	//"github.com/PuerkitoBio/goquery"
)

func main() {
	cmd := exec.Command("D:/data/n1k0-casperjs-e3a77d0/bin/casperjs", "D:/data/my2/a.js")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	n := string(out)
	n = strings.Trim(n, " \n\r")
	fmt.Println(out)
	fmt.Println(n)
	if p, ok := strconv.Atoi(n); ok == nil {
		fmt.Println(p)
	} else {
		fmt.Println("error")
	}
	fmt.Println(string(out))
}
