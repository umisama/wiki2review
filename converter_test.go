package main

import (
	"testing"
)

type TestCase struct {
	Src string // source
	Res string // result string
	Err bool   // returns error?
}

var TestCasesConvertSection = []TestCase{
	TestCase{"**あいうえお", "==あいうえお\n", false}, // case from "コンバータ未完了タスク"
	TestCase{"* 見出し見出し", "= 見出し見出し\n", false},
	TestCase{"** Header2", "== Header2\n", false},
	TestCase{"*** Header3", "=== Header3\n", false},
	TestCase{"***スペース無し", "===スペース無し\n", false},
}

var TestCasesConvertSize = []TestCase{
	TestCase{`&size(18){''hogehoge''};`, `@<b>{hogehoge}`, false}, // base from 1400.eduwiki
	TestCase{`&size(18){''hogehoge''}`, `@<b>{hogehoge}`, false},
	TestCase{`&size(18){hogehoge}`, `@<b>{hogehoge}`, false},
}

var TestCasesConvertColor = []TestCase{
	TestCase{`&color(#0000ff){hogehoge}`, `@<color:#0000ff>{hogehoge}`, false}, // base from 1000.eduwiki
	TestCase{`&color(#0000ff){''hogehoge''}`, `@<color:#0000ff>{hogehoge}`, false},
	TestCase{`&color(#0000ff){''hogehoge''};`, `@<color:#0000ff>{hogehoge}`, false},
	TestCase{`&color(#0000ff){@<b>{hello}};`, `@<color:#0000ff>{@<b>{hello}}`, false},
	TestCase{`&color(Red){hogehoge}`, `@<color:Red>{hogehoge}`, false},
	TestCase{`&color(Red){''hogehoge''}`, `@<color:Red>{hogehoge}`, false},
	TestCase{`&color(Red){''hogehoge''};`, `@<color:Red>{hogehoge}`, false},
	TestCase{`&color(Red){@<b>{hello}};`, `@<color:Red>{@<b>{hello}}`, false},
}

var TestCasesConvertAmazon = []TestCase{
	TestCase{`#amazon(http://xxxx/dp/123456)`, `@<href>{amazon://dp/123456}`, false}, // case from "コンバータ未完了タスク"
	TestCase{`#amazon(http://www.amazon.co.jp/%E4%B8%80%E6%9E%9A%E7%B5%B5%E3%83%BB%E5%A0%B4%E9%9D%A2%E7%B5%B5%E3%81%A7%E5%89%B5%E3%82%8B%E5%86%8D%E7%8F%BE%E6%A7%8B%E6%88%90%E6%B3%95%E3%81%AE%E9%81%93%E5%BE%B3%E6%8E%88%E6%A5%AD-%E5%B0%8F%E5%AD%A6%E6%A0%A1%E7%B7%A8-%E6%96%B0%E3%81%97%E3%81%84%E9%81%93%E5%BE%B3%E6%8E%88%E6%A5%AD%E3%81%A5%E3%81%8F%E3%82%8A%E3%81%B8%E3%81%AE%E6%8F%90%E5%94%B1-20-%E5%85%AB%E6%9C%A8%E4%B8%8B/dp/4188695129)`, `@<href>{amazon://dp/4188695129}`, false}, // from 1200.eduwiki
}

var TestCasesConvertYoutube = []TestCase{
	TestCase{`#youtube(http://www.youtube.com/v/hogehoge)`, `@<href>{youtube://video/hogehoge}`, false}, // case from "コンバータ未完了タスク"
	TestCase{``, ``, false},
}

var TestCasesConvertImg = []TestCase{
	TestCase{`#img(144,zengo.jpg,560)`, `@<href>{image://attachment/144/zengo.jpg?width=560}`, false}, // case from "コンバータ未完了タスク"
	TestCase{`#img(508,月と星p.97.jpg,560)`, `@<href>{image://attachment/508/月と星p.97.jpg?width=560}`, false},	// case from 500.eduwiki
	TestCase{`#img(508,月と星p.97.jpg)`, `@<href>{image://attachment/508/月と星p.97.jpg}`, false},	// case from 500.eduwiki
}

var TestCasesConvertLink = []TestCase{
	TestCase{`http://google.com`, `@<href>{http://google.com}`, false}, // case from "コンバータ未完了タスク"
	TestCase{`#link(http://google.com)`, `@<href>{http://google.com}`, false}, // case from "コンバータ未完了タスク"
	TestCase{`#link(3251,小さな勇気みつけ.docx)`, `@<href>{attachment://id/3251}`, false},
}

var TestCasesConvertTildaLn = []TestCase{
	TestCase{"test~\n", "test\n", false},
}

var TestCasesConvertTripleComma = []TestCase{
	TestCase{`'''`, `'`, false},
	TestCase{`'''''`, `'''`, false},
	TestCase{`&color(Red){''hogehoge''};`, `&color(Red){''hogehoge''};`, false},
}

var TestCasesBase = []TestCase{
	TestCase{``, ``, false},
	TestCase{``, ``, false},
}

func Test_converter_convertSection(t *testing.T) {
	for k, v := range TestCasesConvertSection {
		c := Converter{IsDone: false, buf: v.Src}
		err := c.convertSection()
		if (err != nil) != v.Err {
			t.Error("fail on", k, "with", err)
			return
		}

		if c.buf != v.Res {
			t.Error("fail on", k)
			t.Log(c.buf, v.Res)
			return
		}
	}
}

func Test_converter_convertSize(t *testing.T) {
	for k, v := range TestCasesConvertSize {
		c := Converter{IsDone: false, buf: v.Src}
		err := c.convertSize()
		if (err != nil) != v.Err {
			t.Error("fail on", k, "with", err)
			return
		}

		if c.buf != v.Res {
			t.Error("fail on", k)
			t.Log(c.buf, v.Res)
		}
	}
}

func Test_converter_convertColor(t *testing.T) {
	for k, v := range TestCasesConvertColor {
		c := Converter{IsDone: false, buf: v.Src}
		err := c.convertColor()
		if (err != nil) != v.Err {
			t.Error("fail on", k, "with", err)
			return
		}

		if c.buf != v.Res {
			t.Error("fail on", k)
			t.Log(c.buf, v.Res)
		}
	}
}

func Test_converter_convertAmazon(t *testing.T) {
	for k, v := range TestCasesConvertAmazon {
		c := Converter{IsDone: false, buf: v.Src}
		err := c.convertAmazon()
		if (err != nil) != v.Err {
			t.Error("fail on", k, "with", err)
			return
		}

		if c.buf != v.Res {
			t.Error("fail on", k)
			t.Log(c.buf, v.Res)
		}
	}
}

func Test_converter_convertYoutube(t *testing.T) {
	for k, v := range TestCasesConvertYoutube {
		c := Converter{IsDone: false, buf: v.Src}
		err := c.convertYoutube()
		if (err != nil) != v.Err {
			t.Error("fail on", k, "with", err)
			return
		}

		if c.buf != v.Res {
			t.Error("fail on", k)
			t.Log(c.buf, v.Res)
		}
	}
}

func Test_converter_convertetImg(t *testing.T) {
	for k, v := range TestCasesConvertImg{
		c := Converter{IsDone: false, buf: v.Src}
		err := c.convertImg()
		if (err != nil) != v.Err {
			t.Error("fail on", k, "with", err)
			return
		}

		if c.buf != v.Res {
			t.Error("fail on", k)
			t.Log(c.buf, v.Res)
		}
	}
}

func Test_converter_convertetLink(t *testing.T) {
	for k, v := range TestCasesConvertLink {
		c := Converter{IsDone: false, buf: v.Src}
		err := c.convertLink()
		if (err != nil) != v.Err {
			t.Error("fail on", k, "with", err)
			return
		}

		if c.buf != v.Res {
			t.Error("fail on", k)
			t.Log(c.buf, v.Res)
		}
	}
}

func Test_converter_convertetTildaLn(t *testing.T) {
	for k, v := range TestCasesConvertTildaLn {
		c := Converter{IsDone: false, buf: v.Src}
		err := c.convertTildeNl()
		if (err != nil) != v.Err {
			t.Error("fail on", k, "with", err)
			return
		}

		if c.buf != v.Res {
			t.Error("fail on", k)
			t.Log(c.buf, v.Res)
		}
	}
}

func Test_converter_convertetTripleComma (t *testing.T) {
	for k, v := range TestCasesConvertTripleComma {
		c := Converter{IsDone: false, buf: v.Src}
		err := c.convertTripleComma()
		if (err != nil) != v.Err {
			t.Error("fail on", k, "with", err)
			return
		}

		if c.buf != v.Res {
			t.Error("fail on", k)
			t.Log(c.buf, v.Res)
		}
	}
}
