package translation

type LANG string

const (
	EN LANG = "EN"
	VI LANG = "VI"
)

type translation map[LANG]string

var translations = map[string]map[LANG]string{
	"Chapter": {
		VI: "Chương",
		EN: "Chapter",
	},
}

func Get(key string, lang LANG) string {
	if value, ok := translations[key][lang]; ok {
		return value
	}
	
	return "<translation-not-found: " + key + ">"
}
