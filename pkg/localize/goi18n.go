package localize

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"sort"

	"gopkg.in/yaml.v2"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Goi18n struct {
	files        fs.FS
	language     *language.Tag
	bundle       *i18n.Bundle
	translations []i18n.MessageFile
	format       string
	path         string
}

var (
	//go:embed locales
	defaultLocales  embed.FS
	defaultLanguage = &language.English
)

type Config struct {
	files    fs.FS
	language *language.Tag
	format   string
	path     string
}

func GetDefaultLocales() fs.FS {
	return defaultLocales
}

func GetDefaultLanguage() *language.Tag {
	return defaultLanguage
}

func New(cfg *Config) (*Goi18n, error) {
	if cfg == nil {
		cfg = &Config{}
	}
	if cfg.files == nil {
		cfg.files = GetDefaultLocales()
		cfg.path = "locales"
		cfg.format = "toml"
	}
	if cfg.language == nil {
		cfg.language = GetDefaultLanguage()
	}
	if cfg.format == "" {
		cfg.format = "toml"
	}

	bundle := i18n.NewBundle(*cfg.language)
	loc := &Goi18n{
		files:    cfg.files,
		language: cfg.language,
		bundle:   bundle,
		format:   cfg.format,
		path:     cfg.path,
	}

	err := loc.load()
	return loc, err
}

// walk the file system and load each file into memory
func (l *Goi18n) load() error {
	return fs.WalkDir(l.files, l.path, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return l.MustLocalizeFile(l.files, path)
	})
}

func (l *Goi18n) MustLocalizeFile(files fs.FS, path string) (err error) {
	buf, err := fs.ReadFile(files, path)
	if err != nil {
		return err
	}
	fileext := fmt.Sprintf("%v.%v", l.language.String(), l.format)
	var unmarshalFunc i18n.UnmarshalFunc
	switch l.format {
	case "toml":
		unmarshalFunc = toml.Unmarshal
	case "yaml", "yml":
		unmarshalFunc = yaml.Unmarshal
	case "json":
		unmarshalFunc = json.Unmarshal
	default:
		return fmt.Errorf("unsupported format \"%v\"", l.format)
	}

	l.bundle.RegisterUnmarshalFunc(l.format, unmarshalFunc)
	file, err := l.bundle.ParseMessageFileBytes(buf, fileext)
	if err != nil {
		return err
	}

	l.translations = append(l.translations, *file)

	return nil
}

func (l *Goi18n) GetTranslations() []i18n.MessageFile {
	return l.translations
}

func (l *Goi18n) FindTranslationString(key string) {
	for _, file := range l.translations {
		fmt.Println(key)
		id := sort.Search(len(file.Messages), func(i int) bool {
			return file.Messages[i].ID == key
		})

		if id < len(file.Messages) && file.Messages[id].ID == key {
			fmt.Println(id)
		}
	}
}
