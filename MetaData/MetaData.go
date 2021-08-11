package main

import "bytes"
import "bufio"
import "encoding/json"
import "flag"
import "fmt"
import "net/http"
import "os"
import "strings"

type MetaData interface{}

var pretty bool

func init() {
	flag.BoolVar(&pretty, "prettyprint", false, "pretty print")
}

func readData(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	rdr := bufio.NewReader(resp.Body)
	data := make([]string, 0)
	for {
		line, err := rdr.ReadString('\n')
		data = append(data, strings.TrimRight(line, "\n"))
		if err != nil {
			break
		}
	}
	return data
}

func getMetaData(url string) MetaData {
	md := make(map[string]interface{})
	data := readData(url)
	for _, line := range data {
		switch {
		default:
			md[line] = readData(url + line)[0]
		case line == "metrics/":
			break
		case strings.HasSuffix(line, "/"):
			md[line[:len(line)-1]] = getMetaData(url + line)
		case strings.HasSuffix(url, "public-keys/"):
			keyId := strings.SplitN(line, "=", 2)[0]
			md[line] = readData(url + keyId + "/openssh-key")[0]
		}
	}
	return md
}

func main() {
	flag.Parse()
	url := "http://169.254.169.254/latest/meta-data/"
	md := getMetaData(url)
	b, err := json.Marshal(md)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	if pretty {
		var p bytes.Buffer
		err = json.Indent(&p, b, "", "  ")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		b, err = p.ReadBytes(0)
	}
	os.Stdout.Write(b)
	fmt.Println()
}
