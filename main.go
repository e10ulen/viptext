package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"bufio"
	"fmt"
	"io"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"gopkg.in/alecthomas/kingpin.v2"
)

//	iframe内のURLを指定しないとスクレイピングがうまく動かないので適宜修正必要
//	引数にURL or 外部テキストファイルにURLを置いておく事で管理のしやすさ、
//	バイナリで鯖に投げた後のメンテナンス性が向上する
//	[TODO]URL管理の外部化
const (
	//	url          = "http://dawnlight.ovh/test/read.cgi/viptext/1597046459"
	filesjis     = "./viptext.sjis.html"
	fileutf8     = "./viptext.utf8.html"
	filesanitize = "./viptext.sanitize.html"
)

//	kingpin.v2で引数を受け付けるための変数宣言
var (
	Cmd = kingpin.CommandLine
	//	起動時は--urlsで起動しないと動かない
	//	この時、URLを文字列として要求しているため、
	//	README.mdの実行処理の通りに実行する
	url = Cmd.Flag("urls", "http urls").String()
)

//	kingpinの初期処理
func init() {
	Cmd.Name = "ViP de TextSite"
	Cmd.Help = "VTS portal"
	Cmd.Version("0.0.1")
}

func main() {
	//	kingpin実行処理
	_, err := Cmd.Parse(os.Args[1:])
	if err != nil {
		Cmd.FatalUsage(fmt.Sprintf("\x1b[33m%v\x1b[0m", err))
	}
	if err := run(); err != nil {
		log.Fatalf("!!%v", err)
	}

	utf8toSANITIZE(filesanitize, fileutf8)
	//	HTML化処理を行う
	//	現在実装途中で動かすとサニタイズしたファイル内容が吹っ飛ぶ
	htmlParse()

}

//	kingpin起動処理
func run() error {
	url := *url
	fmt.Println(url)
	fileDown(url, fileutf8, filesjis)

	return nil
}

//	ファイル取得して、ファイル書き出し処理。
//	変数:urlからnet/http経由で「viptext.sjis.html」を書き出す
//	この時点ではまだsjisでファイルエンコードがかかっているため直接的には
//	ファイルが文字化けして読めない

func fileDown(url string, fileutf8 string, filesjis string) {
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
	policy.AllowElements("dd", "a")
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

/*
	ファイルを開いて、replaceで<dd>と</dd>を削除する。
	あと、レス番号、>>1の行を削除（名無しの先行者）の行
	目標物：<dd> !サンプルサイト\これはサンプルです  http://example.com/example.html テストテキスト</dd>
	これをこうしてしたい
	成果物：<a href="http://example.com/example.html" target="_blank">サンプルサイト｜これはサンプルです</a><br />
*/
func htmlParse() {
	//	Open
	file, err := os.OpenFile("viptext.sanitize.html", os.O_RDWR, 0664)
	if err != nil {
		log.Println("file Read error")
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		log.Println(err)
	}

	data := make([]byte, info.Size())
	count, err := file.Read(data)
	if err != nil {
		log.Println(err)
	}
	scanner := bufio.NewScanner(file)
	//	Write
	doc, err := os.Create("viptext.html")
	if err != nil {
		log.Fatal(err)
	}
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		doc.WriteString(strings.ReplaceAll(string(scanner.Text()), "<dd> !", "<a name='"))
		doc.WriteString(strings.ReplaceAll(string(scanner.Text()), "\\", "|"))
		doc.WriteString(strings.ReplaceAll(string(scanner.Text()), "\n", "' href='"))
		count += 1
	}

	//	URL正規表現
	//	https?://[\w/:%#\$&\?\(\)~\.=\+\-]+
	//	http://www.google.co.jp/search?hl=ja&q=%U&lr=
	//	or
	//	 regexp.MustCompile("http(.*)://([a-z]+)/([a-z]+)/([a-zA-Z0-9]+)")
	//	http://example.com/example.html
	//	通るのであれば上の正規表現で精査を掛けた方がいい
	/*
		//	 書き込み
		teesanitize := io.TeeReader(doc, file)
		s := bufio.NewScanner(teesanitize)
		for s.Scan() {
		}
		if err := s.Err(); err != nil {
			log.Fatal(err)
		}
		log.Println("Sanitize done")
	*/
}
