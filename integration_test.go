package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func Test_IntegrationTest(t *testing.T) {
	cases := testcase_generator(t)
	t.Log("found", len(cases), "test cases.")

	for k, v := range cases {
		c := Converter{Src: v.Src}
		c.DoConvert()

		if c.GetResult() != v.Res {
			t.Error("fail on", k)
			continue
		}
	}
	return
}

func testcase_generator(t *testing.T)map[string]TestCase {
	f, err := os.Open("./integration_test")
	if err != nil {
		t.Error(err)
		return nil
	}

	finfos, err := f.Readdir(-1)
	if err != nil {
		t.Error(err)
		return nil
	}

	ret := make(map[string]TestCase)
	for _, v := range finfos {
		if path.Ext(v.Name())[1:] == "eduwiki" {
			// source file was found.
			f, err := os.Open("./integration_test/" + v.Name())
			if err != nil {
				t.Error(err)
				return nil
			}
			src, _ := ioutil.ReadAll(f)

			// result file.
			f, err = os.Open("./integration_test/" + v.Name()+".re")
			if err != nil {
				t.Error(err)
				return nil
			}
			result, _ := ioutil.ReadAll(f)

			ret[v.Name()] = TestCase {
				Src : string(src),
				Res : string(result),
				Err : false,
			}
		}
	}

	return ret
}
