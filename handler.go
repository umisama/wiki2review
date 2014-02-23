package main

import (
	"github.com/stretchr/goweb"
	"github.com/stretchr/goweb/context"
	"io/ioutil"
)

func ConvertHandler(ctx context.Context) (err error){
	req := ctx.HttpRequest()

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}
	defer req.Body.Close()

	cvtr := NewConverter(string(buf))
	if cvtr == nil {
		println("a")
		return
	}

	err = cvtr.DoConvert()
	if err != nil {
		println("b")
		return
	}

	return goweb.Respond.With(ctx, 200, []byte(cvtr.GetResult()))
}
