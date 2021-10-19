package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"bufio"
	"fmt"
	"io"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const url = "http://dawnlight.ovh/test/read.cgi/viptext/1597046459"

func main() {

	//	ファイル取得して、ファイル書き出し処理.
	resp, _ := http.Get(url)
	filename := "viptext.sjis.html"
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
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

	//	ここから先でパーサーを通す事
	//	File Open
	file, err := os.OpenFile("./viptext.uft8.html", os.O_RDWR, 0664) // For read access.
	fmt.Println("file Opens!!!")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	policy := bluemonday.NewPolicy()

	policy.AllowElements("dd")
	doc := policy.SanitizeReader(file)
	fmt.Println(doc)
	// 書き込み先ファイルを用意
	sanitizeFile, err := os.Create("./viptext.sinitize.html")
	if err != nil {
		log.Fatal(err)
	}
	defer sanitizeFile.Close()

	// 書き込み
	teesanitize := io.TeeReader(doc, sanitizeFile)
	docs := bufio.NewScanner(teesanitize)
	for docs.Scan() {
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	log.Println("Sanitize done")
	/*
	   //	goquery
	   	doc, err := goquery.NewDocumentFromReader(file)

	   	if err != nil {
	   		log.Fatal(err)
	   	}
	   	doc.Find("table > dl > dd").Each(func(i int, s *goquery.Selection) {
	   		band := s.Find("a").Text()
	   		urls, _ := s.Attr("href")
	   		fmt.Println("URLs: %s", urls)
	   		fmt.Println("SiteTitle: %s", band)

	   	})
	   //	html2text
	   		text, err := html2text.FromString(string(file), html2text.Options{PrettyTables: true})
	   		if err != nil {
	   			panic(err)
	   		}
	   		fmt.Println(text)
	*/
}
