package middlewares

import (
	"embed"

	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

//go:embed i18n/locale/*
var fs embed.FS

func I18nMiddleware() gin.HandlerFunc {
	return ginI18n.Localize(ginI18n.WithBundle(&ginI18n.BundleCfg{
		DefaultLanguage:  language.Persian,
		FormatBundleFile: "yaml",
		AcceptLanguage:   []language.Tag{language.English, language.Persian},
		RootPath:         "i18n/locale/",
		UnmarshalFunc:    yaml.Unmarshal,
		// After commenting this line, use defaultLoader
		// it will be loaded from the file
		Loader: &ginI18n.EmbedLoader{
			FS: fs,
		},
	}))
}
