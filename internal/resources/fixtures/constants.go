package fixtures

import (
	_ "embed"
)

//go:embed book-info.jude.json
var BookInfo_VI_JudeJson []byte

//go:embed book-info.nahum.json
var BookInfo_VI_NahumJson []byte

//go:embed chapter.vi.html
var ChapterSample_VI_Html []byte

//go:embed chapter.vi.json
var ChapterSample_VI_Json []byte

//go:embed chapter.en.html
var ChapterSample_EN_Html []byte

//go:embed book-info.jude.json
var BookSample_VI_JudeJson []byte

//go:embed fetch.translation.VI1934.html
var TranslationSample_VI1934_Html []byte
