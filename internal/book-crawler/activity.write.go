package book_crawler

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	
	"go.temporal.io/sdk/activity"
	"temporal-crawler/internal/resources"
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

type writer func(name string, data []byte, perm os.FileMode) error

func NewResultWriter(writer writer, dir string) *ResultWriter {
	if writer == nil {
		writer = os.WriteFile
	}
	
	return &ResultWriter{
		baseDir:          dir,
		createMissingDir: true,
		writer:           writer,
	}
}

type ResultWriter struct {
	baseDir          string
	createMissingDir bool
	writer           writer
}

func (w *ResultWriter) getWriter() writer {
	if w.writer != nil {
		return w.writer
	}
	
	return os.WriteFile
}

func (w *ResultWriter) WriteResultActivity(ctx context.Context, book BookInfo, chapter ChapterInfo) error {
	logger := activity.GetLogger(ctx)
	t, err := template.New("chapter").Parse(resources.ChapterTemplateFile)
	
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
	
	writePath, err := w.getWritingPath("%s/build/static/%s/%s/%d.md", w.baseDir, book.Tran, book.BookCode, chapter.Number)
	if err != nil {
		logger.Error("failed to get writing path", "error", err)
		return err
	}
	
	logger.Info("writing chapter", "path", writePath)
	err = w.getWriter()(writePath, buf.Bytes(), 0644)
	
	if err != nil {
		logger.Error("failed to write chapter", "error", err)
		return err
	}
	
	return nil
}

func (w *ResultWriter) getWritingPath(format string, a ...any) (string, error) {
	writePath := fmt.Sprintf(format, a...)
	
	if w.createMissingDir {
		if _, err := os.Stat(writePath); os.IsNotExist(err) {
			err := os.MkdirAll(filepath.Dir(writePath), 0755)
			if err != nil {
				return "", err
			}
		}
	}
	
	return writePath, nil
}
