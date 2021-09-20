package localize

import (
	"embed"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Localize struct {
	fsys         fs.FS
	language     *language.Tag
	bundle       *i18n.Bundle
	translations []i18n.MessageFile
	format       string
	path         string
}

var (
	defaultLocales  embed.FS
	defaultLanguage = &language.English
)

type Config struct {
	fsys     fs.FS
	language *language.Tag
	format   string
	path     string
}

func GetDefaultLanguage() *language.Tag {
	return defaultLanguage
}

func New(cfg *Config) (*Localize, error) {
	if cfg == nil {
		cfg = &Config{}
	}

	wd, err1 := os.Getwd()
	if err1 != nil {
		fmt.Errorf("Failed to get wd: %s", err1)
	}

	if cfg.fsys == nil {
		cfg.path = filepath.Join(wd, "localize", "locales")
		cfg.fsys = os.DirFS(cfg.path)
		cfg.format = "toml"
	}
	if cfg.language == nil {
		cfg.language = GetDefaultLanguage()
	}
	if cfg.format == "" {
		cfg.format = "toml"
	}

	bundle := i18n.NewBundle(*cfg.language)
	loc := &Localize{
		fsys:     cfg.fsys,
		language: cfg.language,
		bundle:   bundle,
		format:   cfg.format,
		path:     cfg.path,
	}

	err := loc.load()
	return loc, err
}

func (l *Localize) load() error {
	return fs.WalkDir(l.fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		return l.MustLocalizeFile(l.fsys, path)
	})
}

func (l *Localize) MustLocalizeFile(fsys fs.FS, path string) (err error) {
	buf, err := fs.ReadFile(fsys, path)
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

func (l *Localize) GetTranslations() []i18n.MessageFile {
	return l.translations
}
