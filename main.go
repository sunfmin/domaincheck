package main

import (
	"flag"
	"fmt"
	"github.com/axgle/pinyin"
	"github.com/sunfmin/fanout"
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

	inputs := []interface{}{}
	for _, word := range words {
		inputs = append(inputs, word)
	}

	results, err2 := fanout.ParallelRun(60, func(input interface{}) (interface{}, error) {
		word := input.(string)
		if strings.TrimSpace(word) == "" {
			return nil, nil
		}

		py := pinyin.Convert(word)
		pydowncase := strings.ToLower(py)
		domain := pydowncase + ".com"
		outr, err := domainAvailable(word, domain)

		if outr.available {
			fmt.Printf("[Ohh Yeah] %s %s\n", outr.word, outr.domain)
		} else {
			fmt.Printf("\t\t\t %s %s %s\n", outr.word, outr.domain, outr.summary)
		}

		if err != nil {
			fmt.Println("Error: ", err)
		}

		return outr, nil
	}, inputs)

	fmt.Println("Finished ", len(results), ", Error:", err2)

}

type checkResult struct {
	word      string
	domain    string
	available bool
	summary   string
}

func domainAvailable(word string, domain string) (ch checkResult, err error) {
	var summary string
	var output []byte

	ch.word = word
	ch.domain = domain

	cmd := exec.Command("whois", domain)
	output, err = cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}

	outputstring := string(output)
	if strings.Contains(outputstring, "No match for \"") {
		ch.available = true
		return
	}

	summary = firstLineOf(outputstring, "Registrant Name") + " => "
	summary = summary + firstLineOf(outputstring, "Expiration Date")
	ch.summary = summary
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
