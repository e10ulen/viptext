# VIP で テキストサイトやろうずｗｗｗｗ
目下目標ではポータルサイト用管理ツール開発レポジトリです。
非公式非公認なのでのんびり、せっつかれること無く開発を行っています

# 現状の機能について

避難所にアクセス→HTML(Shift-JIS)を保存→HTML(Shift-JIS)を開き、HTML(UTF8)として保存
```html
<dd> !サンプルサイト\これはサンプルです  http://example.com/example.html テストテキスト</dd>
```
この部分だけ必要なので、抜き出して保存。
~~現在は不要な情報までもサニタイズで保存されているので再加工が必要~~

# これからの開発目標
更新したらあげるスレにて行われている
!サイトタイトル\更新内容
URL
この形式からサイトタイトルを取得。リストにappendしていき、一定の形式を持ってHTMLとして吐き出してポータルサイトとしての活用を行おうと思っています。
現状考えているのは、リストかスライスにサニタイズしたのを一行ずつ詰め込んでいき、判定で空白とか、  
要らない部分を除去、上下反転したのちにまたファイルを生成して、htmlとしてリメイクして保存。  

# 実行について
```bash
./viptext --urls [掲示板URL] 
```
現状は引数を実装した。
これによってプログラム内部的に参照していたURLがなくなったのでメンテナンス性向上に繋がったと思われる。  

```html
	ファイルを開いて、replaceかなにかで<dd>～～</dd>を<a>～～</a>に置換する。
	あと、レス番号、>>1の行を削除（名無しの先行者）の行
	目標物：<dd> !サンプルサイト\これはサンプルです  http://example.com/example.html テストテキスト</dd>
	成果物：<a href="http://example.com/example.html" target="_blank">サンプルサイト｜これはサンプルです</a><br />
```

## サニタイズをした
現状では実行時にいくつかのファイルが自動的に生成される。  
viptext.sjis.html  
viptext.utf8.html  
viptext.sanitize.html  
上から順番に、スクレイピングした生データ  
エンコードを変えただけの生データ  
サニタイズを行って必要な個所だけ残した加工一段階目のデータ  

<blockquote class="twitter-tweet" data-partner="tweetdeck"><p lang="ja" dir="ltr">正規表現で検出したやつを次々にスライスに突っ込む…？<br>1個目サイトタイトル<br>2個目記事タイトル<br>改行挟んで<br>3個目URL<br>あとは切り捨てで処理できればいいんだけどね</p>&mdash; 乗っ取り食らった依藤 (@e10ulen) <a href="https://twitter.com/e10ulen/status/1453004195943026699?ref_src=twsrc%5Etfw">October 26, 2021</a></blockquote>
<script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>
