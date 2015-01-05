package main

import (
	"flag"
	"fmt"
	"github.com/axgle/pinyin"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var file = flag.String("file", "", "每行一个词语的列表文本文件")

func main() {
	flag.Parse()

	fmt.Println(*file)
	f, err := os.Open(*file)
	if err != nil {
		fmt.Println(err)
		return
	}

	fbytes, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		return
	}

	words := strings.Split(string(fbytes), "\n")

	for _, w := range words {
		py := pinyin.Convert(w)
		pydowncase := strings.ToLower(py)
		domain := pydowncase + ".com"

		yes, summary := domainAvailable(domain)
		if yes {
			fmt.Printf("[Ohh Yeah] %s %s\n", w, domain)
		} else {
			fmt.Printf("\t\t\t %s %s %s\n", w, domain, summary)
		}
	}
}

func domainAvailable(domain string) (available bool, summary string) {
	cmd := exec.Command("whois", domain)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	outputstring := string(output)
	if strings.Contains(outputstring, "No match for \"") {
		available = true
		return
	}

	summary = firstLineOf(outputstring, "Registrant Name") + " => "
	summary = summary + firstLineOf(outputstring, "Registrar Registration Expiration Date")

	return
}

func firstLineOf(content string, keyword string) (line string) {
	lines := strings.Split(content, "\n")
	for _, l := range lines {
		if strings.Contains(l, keyword) {
			line = l
			return
		}
	}
	return
}
