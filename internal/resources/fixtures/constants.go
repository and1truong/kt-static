package fixtures

import (
	_ "embed"
)

//go:embed book-info.jude.json
var BookInfoJudeJson []byte

//go:embed book-info.nahum.json
var BookInfoNahumJson []byte

//go:embed chapter.vi.html
var ChapterSampleViHtml []byte

//go:embed chapter.vi.json
var ChapterSampleViJson []byte

//go:embed chapter.en.html
var ChapterSampleEnHtml []byte

//go:embed book-info.jude.json
var BookSampleJudeJson []byte

//go:embed fetch.translation.VI1934.html
var TranslationSampleVI1934Html []byte
