package lifectrl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

func newTestEnv(t *testing.T) string {
	path, err := ioutil.TempDir(".", "tmp")
	if err != nil {
		t.Fatal(err)
	}

	return path
}

func removeTestEnv(t *testing.T, path string) {
	err := os.RemoveAll(path)
	if err != nil {
		t.Fatal(err)
	}
}

func equalStage(s1, s2 Stage) error {
	if s1.Title != s2.Title ||
		!s1.From.Equal(s2.From) ||
		!s1.To.Equal(s2.To) ||
		s1.Desc != s2.Desc {
		e := fmt.Sprintf("Expect %v was %v", s1, s2)
		return errors.New(e)
	}

	return nil
}

func Test_ParseFile_OK(t *testing.T) {
	fContent := "Title: Test\nFrom: 1/2010\nTo: 2/2010\nTitle\n====="
	reader := bytes.NewReader([]byte(fContent))

	stage, err := ParseFile(reader)
	if err != nil {
		t.Fatal(err)
	}

	expect := Stage{
		Title: "Test",
		From:  time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC),
		To:    time.Date(2010, time.February, 1, 0, 0, 0, 0, time.UTC),
		Desc:  "Title\n=====",
	}

	if err = equalStage(expect, stage); err != nil {
		t.Fatal(err)
	}

}

func Test_ParseFile_Fail(t *testing.T) {
	fContent := "From: 1/2010\nTo: 2/2010\nTitle\n====="
	reader := bytes.NewReader([]byte(fContent))

	_, err := ParseFile(reader)
	if err.Error() != "Wrong Title line" {
		t.Fatal("Expect error was", err)
	}

	fContent = "Title: Test\nTo: 2/2010\nTitle\n====="
	reader = bytes.NewReader([]byte(fContent))

	_, err = ParseFile(reader)
	if err.Error() != "Wrong From line" {
		t.Fatal("Expect error was", err)
	}

	fContent = "Title: Test\nFrom: 2/2010\nTitle\n====="
	reader = bytes.NewReader([]byte(fContent))

	_, err = ParseFile(reader)
	if err.Error() != "Wrong To line" {
		t.Fatal("Expect error was", err)
	}

	fContent = "Title: Test\nFrom: 2/2010\nTo: 2/2010\n"
	reader = bytes.NewReader([]byte(fContent))

	_, err = ParseFile(reader)
	if err.Error() != "Missing description" {
		t.Fatal("Expect error was", err)
	}
}

func Test_NewStageFile_OK(t *testing.T) {
	tmpPath := newTestEnv(t)
	defer removeTestEnv(t, tmpPath)

	inPath := path.Join(tmpPath, "test.md")
	input := []byte("Title:Funky bacon!\nFrom:1/2010\nTo:1/2011\nwhy?")
	err := ioutil.WriteFile(inPath, input, 0755)
	if err != nil {
		t.Fatal(err)
	}

	outPath := path.Join(tmpPath, "output.json")
	err = NewStageFile(tmpPath, outPath)
	if err != nil {
		t.Fatal(err)
	}

	in, err := ioutil.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}

	result := StageFile{}
	err = json.Unmarshal(in, &result)
	if err != nil {
		t.Fatal(err)
	}

	expect := Stage{
		Title: "Funky bacon!",
		From:  time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
		To:    time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC),
		Desc:  "why?",
	}

	err = equalStage(expect, result.Stages[0])
	if err != nil {
		t.Fatal(err)
	}
}
