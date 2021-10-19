package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"bufio"
	"fmt"
	"io"

	"github.com/comail/colog"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func main() {
	colog.Register()
	//	ファイル取得して、ファイル書き出し処理.
	resp, _ := http.Get("http://dawnlight.ovh/viptext/#ui-tabs-2")
	filename := "viptext.sjis.html"
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	println(string(body))
	writeBody := []byte(body)
	err := ioutil.WriteFile(filename, writeBody, 0664)
	if err != nil {
		fmt.Println(err)
	}
	sjisFile, err := os.Open("./viptext.sjis.html")
	if err != nil {
		log.Fatal(err)
	}
	defer sjisFile.Close()

	// ShiftJISのデコーダーを噛ませたReaderを作成する
	reader := transform.NewReader(sjisFile, japanese.ShiftJIS.NewDecoder())

	// 書き込み先ファイルを用意
	utf8File, err := os.Create("./viptext.uft8.html")
	if err != nil {
		log.Fatal(err)
	}
	defer utf8File.Close()

	// 書き込み
	tee := io.TeeReader(reader, utf8File)
	s := bufio.NewScanner(tee)
	for s.Scan() {
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	log.Println("done")

}
