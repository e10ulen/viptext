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
const filesjis = "./viptext.sjis.html"
const fileutf8 = "./viptext.utf8.htnl"
const filesanitize = "./viptext.sanitize.html"

func main() {

	fileDown(fileutf8, filesjis)

	utf8toSANITIZE(filesanitize, fileutf8)

}

//	ファイル取得して、ファイル書き出し処理。
//	変数:urlからnet/http経由で「viptext.sjis.html」を書き出す
//	この時点ではまだsjisでファイルエンコードがかかっているため直接的には
//	ファイルが文字化けして読めない

func fileDown(fileutf8 string, filesjis string) {
	resp, _ := http.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	writeBody := []byte(body)
	err := ioutil.WriteFile(filesjis, writeBody, 0664)
	if err != nil {
		fmt.Println(err)
	}
	sjisFile, err := os.Open(filesjis)
	if err != nil {
		log.Fatal(err)
	}
	defer sjisFile.Close()

	encodingSJIStoUTF8(fileutf8, sjisFile)
}

//	ShiftJISのデコーダーを噛ませたReaderを作成する。
//	この場所でエンコードをShift_JISに変更するためデコーダを生成する
//	書き込み先ファイルを用意
//	エンコードしたファイルをこの場所でファイル名を変え保存している。

func encodingSJIStoUTF8(fileutf8 string, sjisFile *os.File) {
	reader := transform.NewReader(sjisFile, japanese.ShiftJIS.NewDecoder())

	utf8File, err := os.Create(fileutf8)
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

//	ここから先でパーサーを通す事
//	いくつかパーサーをテストした結果、昔触った
//	bluemondayが一番キレイにファイルを処理してくれているので
//	現状はここでパーサーを通すために先程閉じたファイルをここで再度開きなおしている。
//	メモリ管理上たぶん此処は要改善だと思われる。

func utf8toSANITIZE(filesanitize string, fileutf8 string) {
	file, err := os.OpenFile(fileutf8, os.O_RDWR, 0664) // For read access.
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

	sanitizeFile, err := os.Create(filesanitize)
	if err != nil {
		log.Fatal(err)
	}
	defer sanitizeFile.Close()

	// 書き込み
	teesanitize := io.TeeReader(doc, sanitizeFile)
	s := bufio.NewScanner(teesanitize)
	for s.Scan() {
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	log.Println("Sanitize done")
}
