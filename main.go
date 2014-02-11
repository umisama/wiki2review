package main

import (
	"os"
	"fmt"
	"io/ioutil"
)

func Test( src string ) string {
	return "A"
}

func outputUsage() {
	fmt.Printf("usage : %s [input-file]", os.Args[0])
	return
}

func GetStringsFromFile( path string ) string {
	f, err := os.Open( path )
	if err != nil {
		fmt.Printf( err.Error() )
		return ""
	}

	buf, err := ioutil.ReadAll( f )
	if err != nil {
		fmt.Printf( err.Error() )
		return ""
	}

	return string(buf)
}

func main() {
	if len(os.Args) <= 1 {
		outputUsage()
		return
	}

	src := GetStringsFromFile(os.Args[1])

	cvtr := NewConverter( src )
	if cvtr == nil {
		return
	}

	cvtr.DoConvert()
	fmt.Print( cvtr.GetResult() )
}
