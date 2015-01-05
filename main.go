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

	result := make(chan checkResult)
	throttle := make(chan int, 60)

	words := strings.Split(string(fbytes), "\n")
	go func() {
		for {
			outr := <-result
			if outr.available {
				fmt.Printf("[Ohh Yeah] %s %s\n", outr.word, outr.domain)
			} else {
				fmt.Printf("\t\t\t %s %s %s\n", outr.word, outr.domain, outr.summary)
			}
		}
	}()

	for _, w := range words {
		if strings.TrimSpace(w) == "" {
			continue
		}

		py := pinyin.Convert(w)
		pydowncase := strings.ToLower(py)
		domain := pydowncase + ".com"

		throttle <- 1

		go domainAvailable(w, domain, result, throttle)
	}

}

type checkResult struct {
	word      string
	domain    string
	available bool
	summary   string
}

func domainAvailable(word string, domain string, in chan checkResult, throttle chan int) {
	cmd := exec.Command("whois", domain)
	var available bool
	var summary string

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
	summary = summary + firstLineOf(outputstring, "Expiration Date")

	in <- checkResult{word, domain, available, summary}
	<-throttle
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
