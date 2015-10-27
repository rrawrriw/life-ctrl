package lifectrl

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/knieriem/markdown"
)

type (
	Stage struct {
		Title string
		From  time.Time
		To    time.Time
		Desc  string
	}
)

func CleanLine(l []byte) []byte {
	return bytes.TrimSpace(l)
}

func ReadParam(l []byte, field string) ([]byte, error) {
	cutset := []byte(fmt.Sprintf("%v:", field))
	s := bytes.SplitAfter(l, cutset)

	if len(s) != 2 {
		m := fmt.Sprintf("Wrong %v line", field)
		err := errors.New(m)
		return []byte{}, err
	}

	return CleanLine(s[1]), nil
}

func NewDate(l []byte) (time.Time, error) {
	d := bytes.Split(l, []byte("/"))

	if len(d) != 2 {
		return time.Time{}, errors.New("Wrong date format")
	}

	tmp, err := strconv.Atoi(string(d[0]))
	if err != nil {
		return time.Time{}, err
	}
	month := time.Month(tmp)
	year, err := strconv.Atoi(string(d[1]))
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC), nil
}

func ParseFile(in io.Reader) (Stage, error) {
	stage := Stage{}
	reader := bufio.NewReader(in)

	l, err := reader.ReadBytes('\n')
	if err != nil {
		return Stage{}, err
	}
	tmp, err := ReadParam(l, "Title")
	if err != nil {
		return Stage{}, err
	}
	stage.Title = string(tmp)

	l, err = reader.ReadBytes('\n')
	if err != nil {
		return Stage{}, err
	}
	tmp, err = ReadParam(l, "From")
	if err != nil {
		return Stage{}, err
	}
	from, err := NewDate(tmp)
	if err != nil {
		return Stage{}, err
	}
	stage.From = from

	l, err = reader.ReadBytes('\n')
	if err != nil {
		return Stage{}, err
	}
	tmp, err = ReadParam(l, "To")
	if err != nil {
		return Stage{}, err
	}
	to, err := NewDate(tmp)
	if err != nil {
		return Stage{}, err
	}
	stage.To = to

	p := markdown.NewParser(&markdown.Extensions{})
	html := bytes.NewBuffer([]byte{})
	p.Markdown(reader, markdown.ToHTML(html))

	if html.Len() < 2 {
		return Stage{}, errors.New("Missing description")
	}

	stage.Desc = html.String()

	return stage, nil
}
