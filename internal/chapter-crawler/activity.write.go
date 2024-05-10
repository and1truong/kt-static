package chapter_crawler

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	
	"go.temporal.io/sdk/activity"
	"temporal-crawler/internal"
	"temporal-crawler/internal/activities"
	"temporal-crawler/internal/entity"
	"temporal-crawler/internal/resources"
	"temporal-crawler/internal/resources/translation"
)

type (
	templateArgs struct {
		BookName        string
		ChapterNumber   int
		TranslationName string
		LanguageName    translation.LANG
		Slug            string
		Content         string
		AudioFiles      string
	}
)

func NewResultWriter(writer internal.FileWriter, createMissingDir bool) *ResultWriter {
	if writer == nil {
		writer = os.WriteFile
	}
	
	return &ResultWriter{
		createMissingDir: createMissingDir,
		writer:           writer,
	}
}

type ResultWriter struct {
	createMissingDir bool
	writer           internal.FileWriter
}

func (w *ResultWriter) getWriter() internal.FileWriter {
	if w.writer != nil {
		return w.writer
	}
	
	return os.WriteFile
}

func (w *ResultWriter) WriteResultActivity(ctx context.Context, book entity.BookInfo, chapter entity.ChapterInfo) error {
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
		LanguageName:    translation.Translations[book.Tran],
		Slug:            fmt.Sprintf("/%s/%d/%s", book.BookCode, chapter.Number, book.Tran),
		Content:         chapter.String(),
		AudioFiles:      internal.JsonifyStringSlice(chapter.AudioLinks),
	}
	
	buf := bytes.NewBufferString("")
	err = t.ExecuteTemplate(buf, "chapter", args)
	if err != nil {
		logger.Error("failed to execute template", "error", err)
		return err
	}
	
	writePath, err := w.getWritingPath("%s/build/static/%s/%s/%d.md", activities.GetBaseDir(), book.Tran, book.BookCode, chapter.Number)
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
