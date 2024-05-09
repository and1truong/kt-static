package book_crawler

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"os"
	"text/template"
	
	"go.temporal.io/sdk/activity"
)

type (
	templateArgs struct {
		BookName        string
		ChapterNumber   int
		TranslationName string
		LanguageName    string
		Slug            string
		Content         string
	}
)

//go:embed resources/templates/chapter.tpl
var templateFile string

type writer func(name string, data []byte, perm os.FileMode) error

type ResultWriter struct {
	writer writer
}

func (w *ResultWriter) getWriter() writer {
	if w.writer != nil {
		return w.writer
	}
	
	return os.WriteFile
}

func (w *ResultWriter) ActivityHandler(ctx context.Context, book BookInfo, chapter ChapterInfo) error {
	logger := activity.GetLogger(ctx)
	t, err := template.New("chapter").Parse(templateFile)
	
	if err != nil {
		logger.Error("failed to parse template", "error", err)
		return err
	}
	
	args := templateArgs{
		BookName:        book.BookName,
		ChapterNumber:   chapter.Number,
		TranslationName: book.Tran,
		LanguageName:    "TODO",
		Slug:            fmt.Sprintf("/%s/%d/%s", book.BookCode, chapter.Number, book.Tran),
		Content:         chapter.String(),
	}
	
	buf := bytes.NewBufferString("")
	err = t.ExecuteTemplate(buf, "chapter", args)
	if err != nil {
		logger.Error("failed to execute template", "error", err)
		return err
	}
	
	// write Markdown content to file
	err = w.getWriter()(
		fmt.Sprintf("../../build/%s/%s/%d.html", book.Tran, book.BookCode, chapter.Number),
		buf.Bytes(),
		0644,
	)
	
	if err != nil {
		logger.Error("failed to write chapter", "error", err)
		return err
	}
	
	return nil
}
