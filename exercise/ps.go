package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func main() {
	var cmdlinename string
	f, err := os.Open("/proc")
	if err != nil {
		log.Fatal(err)

	}
	infos, _ := f.Readdir(-1)
	for _, info := range infos {
		if !info.IsDir() {
			continue
		}
		infoname, err := strconv.Atoi(info.Name())
		if err != nil {
			continue
		}

		cmdlinename = "/proc" + "/" + info.Name() + "/cmdline"
		buf, _ := ioutil.ReadFile(cmdlinename)
		fmt.Println(infoname, string(buf))

	}

	f.Close()
}