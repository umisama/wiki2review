package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrParsingError = errors.New("parsing error")
)

type Converter struct {
	Src    string
	IsDone bool
	buf    string
}

func NewConverter(src string) *Converter {
	return &Converter{
		Src:    src,
		IsDone: false,
		buf:    "",
	}
}

func (c *Converter) GetResult() string {
	if c.IsDone {
		return c.buf
	}

	return ""
}

func (c *Converter) DoConvert() error {
	// call converter functions
	// please custom for partial converting
	funclist_outer := [](func() error){
		c.convertSection,
		c.convertYoutube,
		c.convertAmazon,
		c.convertTildeNl,
	}

	funclist_inner := [](func() error){
	//	c.convertTripleComma,
		c.convertSize,
		c.convertLink,
		c.convertColor,
		c.convertUnderline,
		c.convertDelete,
		c.convertImg,
		c.convertBoldString,
	}

	if c.IsDone {
		return nil
	}

	c.buf = c.Src
	// convert inner attr
	for _, fn := range funclist_inner {
		err := fn()
		if err != nil {
			return err
		}
	}

	// convert outer attr
	for _, fn := range funclist_outer {
		err := fn()
		if err != nil {
			return err
		}
	}

	// Complete
	c.IsDone = true

	return nil
}

// convertSection() converts wiki style section(*section) to ReVIEW style(#section)
func (c *Converter) convertSection() (err error) {
	fnConv := func(src string) (ret string) {
		for i := 0; i < len(src); i++ {
			ret += "=" // replace * to #
		}
		return
	}

	result := ""
	for _, line := range strings.Split(c.buf, "\r\n") { // CR+LF? LF only?
		conv_line, ierr := simpleReplacer(`^\*+`, line, fnConv)
		if ierr != nil {
			err = ierr
			return
		}

		result += conv_line + "\n"
	}

	c.buf = result

	return
}

// convertBoldString() converts wiki style bold string("bold") to ReVIEW style(@<b>{"bold"})
func (c *Converter) convertSize() (err error) {
	fnConv := func(src string) string {
		comment := regexp.MustCompile(`\{(.*?)}(;*)`).FindString(src)
		return "@<b>{" + strings.Trim(comment, "\"'{};") + "}"
	}
	c.buf, err = simpleReplacer(`&size(.*?){(.*?)}(;*)`, c.buf, fnConv)
	return
}

// convertBoldString() converts wiki style bold string("bold") to ReVIEW style(@<b>{"bold"})
func (c *Converter) convertBoldString() (err error) {
	fnConv := func(src string) string {
		return "@<b>{" + strings.Trim(src, "'") + "}"
	}

	c.buf, err = simpleReplacer(`(\"|'').*?(\"|'')`, c.buf, fnConv)
	return
}

// convertUnderline() converts wiki style bold string(%%%underline%%%) to ReVIEW style(@<u>{"underline"})
func (c *Converter) convertUnderline() (err error) {
	fnConv := func(src string) string {
		return "@<u>{" + strings.Trim(src, "%") + "}"
	}

	c.buf, err = simpleReplacer(`%%%.*?%%%`, c.buf, fnConv)
	return
}

// convertDelete() converts wiki style bold string(%%%underline%%%) to ReVIEW style(@<u>{"underline"})
func (c *Converter) convertDelete() (err error) {
	fnConv := func(src string) string {
		return "@<del>{" + strings.Trim(src, "%") + "}"
	}

	c.buf, err = simpleReplacer(`%%.*?%%`, c.buf, fnConv)
	return
}

// convertYoutube() converts wiki style youtube link(#youtube()) to ReVIEW style(@<href> with youtube:// url scheme)
func (c *Converter) convertYoutube() (err error) {
	fnConv := func(src string) string {
		keyword := "watch?v="

		if strings.Index(src, keyword) == -1 {
			keyword = "youtube.com/v/"
		}

		// get video-id stringrange
		start := strings.Index(src, keyword) + len(keyword)
		end := len(src) - 1

		return "@<href>{youtube://video/" + src[start:end] + "}"
	}

	c.buf, err = simpleReplacer(`#youtube\((.*?)\)`, c.buf, fnConv)
	return
}

// convertAmazon() converts wiki style amazon link(#amazon()) to ReVIEW style(@<link> with amazon:// url scheme)
func (c *Converter) convertAmazon() (err error) {
	fnConv := func(src string) string {
		keyword := "/dp/"

		// get video-id stringrange
		start := strings.Index(src, keyword) + len(keyword)
		end := len(src) - 1

		return "@<href>{amazon://dp/" + src[start:end] + "}"
	}

	c.buf, err = simpleReplacer(`#amazon\((.*?)\)`, c.buf, fnConv)
	return
}

// convertLink() converts wiki style link(#link() to ReVIEW style(@<link>)
func (c *Converter) convertLink() (err error) {
	err = c.convertLinkWithHashLink()
	if err != nil {
		return
	}

	err = c.convertLinkWithLinkString()
	return
}

func (c *Converter) convertLinkWithHashLink() (err error) {
	fnConv := func(src string) string {
		paramstr := strings.TrimFunc(regexp.MustCompile(`\((.*?)\)`).FindString(src), trimmer)
		params := strings.Split(paramstr, ",")

		if len(params) == 2 {
			// attachment
			return fmt.Sprintf("@<href>{attachment://id/%s}", params[0])
		}

		// simple link
		return fmt.Sprintf("@<href>{%s}", params[0])
	}

	c.buf, err = simpleReplacer(`#link\((.*?)\)`, c.buf, fnConv)
	return
}

func (c *Converter) convertLinkWithLinkString() (err error) {
	regexp_url := `https?://[\w/:%#\$&\?\(\)~\.=\+\-]+`

	fnConv := func(src string) string {
		if regexp.MustCompile("{" + regexp_url + "}").MatchString(src) {
			return src
		}
		if strings.Index(src, "amazon.") != -1 || strings.Index(src, "youtube.") != -1 {
			return src
		}

		return fmt.Sprintf("@<href>{%s}", src)
	}

	c.buf, err = simpleReplacer("({*)"+regexp_url+"(}*)", c.buf, fnConv)
	return
}

// convertColor() converts wiki style link(#color(#**){***}) to custom ReVIEW style
func (c *Converter) convertColor() (err error) {
	err = c.convertColorByHEX()
	if err != nil {
		return
	}

	err = c.convertColorByName()
	return
}

func (c *Converter) convertColorByHEX() (err error) {
	fnConv := func(src string) string {
		color := regexp.MustCompile(`#([0-9a-fA-F]{6})`).FindString(src)
		comment := regexp.MustCompile(`\{(.*?)(}+)`).FindString(src)

		nested := strings.Count(comment, "}")		// support for nested elements(ex:testcase 3)

		return "@<color:" + color + ">{" + strings.Trim(comment, "{}'") + strings.Repeat("}", nested)
	}

	c.buf, err = simpleReplacer(`&color\(#(.*?)\){(.*?)}((}|;)*)`, c.buf, fnConv)
	return
}

func (c *Converter) convertColorByName() (err error) {
	colors := `(Fuchsia|fuchsia|Lime|lime|Teal|teal|Navy|navy|red|Red|blue|Blue|green|Green)`
	fnConv := func(src string) string {
		color := regexp.MustCompile(colors).FindString(src)
		comment := regexp.MustCompile(`{(.*?)(}+)`).FindString(src)

		nested := strings.Count(comment, "}")		// support for nested elements(ex:testcase 3)

		return "@<color:" + color + ">{" + strings.Trim(comment, "{}'") + strings.Repeat("}", nested)
	}

	c.buf, err = simpleReplacer(`&color\(`+colors+`\){(.*?)}((}|;)*)`, c.buf, fnConv)
	return
}

func (c *Converter) convertImg() (err error) {
	fnConv := func(src string) string {
		paramstr := strings.TrimFunc(regexp.MustCompile(`\((.*?)\)`).FindString(src), trimmer)
		params := strings.Split(paramstr, ",")

		if len(params) >= 3 {
			return fmt.Sprintf("@<href>{image://attachment/%s/%s?width=%s}", params[0], params[1], params[2])
		} else if len(params) >= 2 {
			return fmt.Sprintf("@<href>{image://attachment/%s/%s}", params[0], params[1])
		}

		panic("error")
	}

	c.buf, err = simpleReplacer(`#img\((.*?)\)`, c.buf, fnConv)
	return
}

func (c *Converter) convertTildeNl() (err error) {
	fnConv := func(src string) string {
		return "\n"
	}

	c.buf, err = simpleReplacer("~\n", c.buf, fnConv)
	return
}

func (c *Converter) convertTripleComma() (err error) {
	fnConv := func(src string) string {
		return "'"
	}

	c.buf, err = simpleReplacer(`'''`, c.buf, fnConv)
	return
}

// simpleReplacer() reprecents usful regexp.ReplaceAllStringFunc util.(private use)
func simpleReplacer(reg string, src string, callback func(string) string) (ret string, err error) {
	r, err := regexp.Compile(reg)
	if err != nil {
		return
	}

	ret = r.ReplaceAllStringFunc(src, callback)
	return
}

func trimmer(t rune) bool {
	switch t {
	case '{', '}', '\'', '(', ')':
		return true
	}

	return false
}
