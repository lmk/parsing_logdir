package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/suapapa/go_hangul/encoding/cp949"
)

func saveMsg(rfn string, wfn string) {
	//fmt.Println(rfn, wfn)

	rfi, err := os.Open(rfn)
	if err != nil {
		fmt.Println("File Open Error: " + rfn)
		panic(err)
	}
	defer rfi.Close()

	wfi, err := os.Create(wfn)
	if err != nil {
		fmt.Println("File Create Error: " + wfn)
		panic(err)
	}
	defer wfi.Close()

	// euc-kr
	reader, _ := cp949.NewReader(rfi)
	writer, _ := cp949.NewWriter(wfi)

	lineNo := 1
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		if line == "" { // empty line
			continue
		}

		splitStr := strings.Split(line[1:len(line)-1], "][")
		if len(splitStr) != 10 { // parsing fail
			continue
		}

		str := strings.TrimSpace(splitStr[9])
		if str == "" { // empty msg
			continue
		}
		_, err := writer.Write([]byte(str))
		if err != nil {
			fmt.Println("File Write Error: " + str + ":" + strconv.Itoa(lineNo))
			panic(err)
		}
		writer.Write([]byte("\r\n"))
	}

	if err = scanner.Err(); err != nil {
		panic(err)
	}
}

func main() {
	// usage cpu 8
	runtime.GOMAXPROCS(8)

	path := "D:/Desktop/gwrecv_20190122"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("ReadDir Error: " + path)
		panic(err)
	}

	var wait sync.WaitGroup

	for _, f := range files {
		go func(f os.FileInfo) {
			start := time.Now()
			wait.Add(1)
			saveMsg(path+"/"+f.Name(), path+"/"+f.Name()+".txt")
			fmt.Printf("%s\t%v\r\n", f.Name(), time.Since(start))
			defer wait.Done()
		}(f)
	}

	wait.Wait()
}
