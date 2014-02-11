package main

import (
	"regexp"
	"strings"
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
	funclist := [](func() error){
		c.convertSection,
		c.convertBoldString,
		c.convertUnderline,
		c.convertDelete,
		c.convertYoutube,
		c.convertAmazon,
		c.convertLink,
		c.convertColor,
	}

	if c.IsDone {
		return nil
	}

	// do convert
	c.buf = c.Src
	for _, fn := range funclist {
		err := fn()
		if err != nil {
			print(err.Error())
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
			ret += "#" // replace * to #
		}
		return
	}

	result := ""
	for _, line := range strings.Split(c.buf, "\n") {
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
func (c *Converter) convertBoldString() (err error) {
	fnConv := func(src string) string {
		return "@<b>{" + strings.Trim(src, "\"") + "}"
	}

	c.buf, err = simpleReplacer(`\".*\"`, c.buf, fnConv)
	return
}

// convertUnderline() converts wiki style bold string(%%%underline%%%) to ReVIEW style(@<u>{"underline"})
func (c *Converter) convertUnderline() (err error) {
	fnConv := func(src string) string {
		return "@<u>{" + strings.Trim(src, "%") + "}"
	}

	c.buf, err = simpleReplacer(`%%%.*%%%`, c.buf, fnConv)
	return
}

// convertDelete() converts wiki style bold string(%%%underline%%%) to ReVIEW style(@<u>{"underline"})
func (c *Converter) convertDelete() (err error) {
	fnConv := func(src string) string {
		return "@<del>{" + strings.Trim(src, "%") + "}"
	}

	c.buf, err = simpleReplacer(`%%.*%%`, c.buf, fnConv)
	return
}

// convertYoutube() converts wiki style youtube link(#youtube()) to ReVIEW style(@<href> with youtube:// url scheme)
func (c *Converter) convertYoutube() (err error) {
	fnConv := func(src string) string {
		keyword := "watch?v="

		// get video-id stringrange
		start := strings.Index(src, keyword) + len(keyword)
		end := len(src) - 1

		return "@<href>{youtube://video/" + src[start:end] + "}"
	}

	c.buf, err = simpleReplacer(`#youtube(.*)`, c.buf, fnConv)
	return
}

// convertAmazon() converts wiki style amazon link(#amazon()) to ReVIEW style(@<link> with amazon:// url scheme)
func (c *Converter) convertAmazon() (err error) {
	fnConv := func(src string) string {
		keyword := "http://"

		// get video-id stringrange
		start := strings.Index(src, keyword) + len(keyword)
		end := len(src) - 1

		return "@<link>{amazon://" + src[start:end] + "}"
	}

	c.buf, err = simpleReplacer(`#amazon(.*)`, c.buf, fnConv)
	return
}

// convertLink() converts wiki style link(#link() to ReVIEW style(@<link>)
func (c *Converter) convertLink() (err error) {
	fnConv := func(src string) string {
		return "@<href>{" + strings.Trim(strings.Trim(src, "#link("), ")") + "}"
	}

	c.buf, err = simpleReplacer(`#link(.*)`, c.buf, fnConv)
	return
}

// convertColor() converts wiki style link(#color(#**){***}) to custom ReVIEW style
func (c *Converter) convertColor() (err error) {
	fnConv := func(src string) string {
		color := regexp.MustCompile(`#([0-9a-fA-F]{6})`).FindString(src)
		comment := regexp.MustCompile(`\{''.*''};`).FindString(src)
		return "@<color:" + color + ">{" + strings.Trim(strings.Trim(comment, "{''"), "''};") + "}"
	}

	c.buf, err = simpleReplacer(`#color\(#(.*)\){(.*)}`, c.buf, fnConv)
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
