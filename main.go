package main

import (
	"github.com/comail/colog"
)

func main() {
	colog.Register()
	//	ファイル取得して、ファイル書き出し処理.

	/*
		resp, _ := http.Get("http://dawnlight.ovh/viptext/#ui-tabs-2")
		filename := "vip.txt"
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		println(string(body))
		writeBody := []byte(body)
		err := ioutil.WriteFile(filename, writeBody, 0664)
		if err != nil {
			fmt.Println(err)
		}
	*/
}
