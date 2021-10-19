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

//	iframe内のURLを指定しないとスクレイピングがうまく動かないので適宜修正必要
//	引数にURL or 外部テキストファイルにURLを置いておく事で管理のしやすさ、
//	バイナリで鯖に投げた後のメンテナンス性が向上する
//	[TODO]URL管理の外部化
const url = "http://dawnlight.ovh/test/read.cgi/viptext/1597046459"

func main() {

	//	ファイル取得して、ファイル書き出し処理。
	//	変数:urlからnet/http経由で「viptext.sjis.html」を書き出す
	//	この時点ではまだsjisでファイルエンコードがかかっているため直接的には
	//	ファイルが文字化けして読めない
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

	//	ShiftJISのデコーダーを噛ませたReaderを作成する。
	//	この場所でエンコードをShift_JISに変更するためデコーダを生成する
	reader := transform.NewReader(sjisFile, japanese.ShiftJIS.NewDecoder())

	//	書き込み先ファイルを用意
	//	エンコードしたファイルをこの場所でファイル名を変え保存している。
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
	//	いくつかパーサーをテストした結果、昔触った
	//	bluemondayが一番キレイにファイルを処理してくれているので
	//	現状はここでパーサーを通すために先程閉じたファイルをここで再度開きなおしている。
	//	メモリ管理上たぶん此処は要改善だと思われる。
	file, err := os.OpenFile("./viptext.uft8.html", os.O_RDWR, 0664) // For read access.
	fmt.Println("file Opens!!!")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	//	bluemondayにはプログラム作成者がポリシーを作ることができる。
	//	それを活用して、通すHTMLタグを決め、ポリシーを通してサニタイズしている。
	policy := bluemonday.NewPolicy()
	policy.AllowElements("dd")
	doc := policy.SanitizeReader(file)
	fmt.Println(doc)
	// 書き込み先ファイルを用意
	//	勿論サニタイズ後の作業ファイルは別名保存している。
	//	全てのファイルにおいて同じことが言えるが、別に同じファイル名でもよかった気がする

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
			//	過去の産物
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
