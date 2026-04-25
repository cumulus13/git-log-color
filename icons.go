package main

import "strings"

func defaultIcons() map[string]string {
	return map[string]string{
		"success":      "✅ ",
		"error":        "❌ ",
		"warning":      "⚠️ ",
		"info":         "ℹ️ ",
		"git":          "🔧 ",
		"github":       "🐙 ",
		"version":      "📦 ",
		"tag":          "🏷️ ",
		"push":         "🚀 ",
		"pull":         "⬇️ ",
		"clone":        "📥 ",
		"commit":       "💾 ",
		"remote":       "🌐 ",
		"folder":       "📁 ",
		"file":         "📄 ",
		"lock":         "🔒 ",
		"unlock":       "🔓 ",
		"search":       "🔍 ",
		"clipboard":    "📋 ",
		"config":       "⚙️ ",
		"time":         "⏰ ",
		"notification": "🔔 ",
		"download":     "⬇️ ",
		"release":      "🎯 ",
		"username":     "🙎 ",
		"owner":        "🚹 ",
		"node":         "📮 ",
		"key":          "🔑 ",
		"date":         "📅 ",
		"email":        "📧 ",
		"link":         "🔗 ",
		"mark":         "🔖 ",
		"find":         "🔎 ",
		"block":        "🧱 ",
		"gitignore":    "🌿 ",
		"diff":         "⚓ ",
		"color":        "🎨 ",
		"linenumber":   "🥦 ",
		"pending":      "📧 ",
		"cancel":       "⛔ ",
		"chapter":      "🎬 ",
		"rust":         "🦀 ",
		"dotnet":       "C++ ",
		"golang":       "🇬🇴 ",
		"nsis":         "🧪 ",
		"test":         "👌 ",
		"python":       "🐍 ",
		"any":          "🍀 ",
	}
}

type IconSet struct {
	values  map[string]string
	enabled bool
}

func NewIconSet(cfg Config) IconSet {
	return IconSet{
		values:  cfg.Icons,
		enabled: cfg.UI.IconsEnabled,
	}
}

func (s IconSet) Get(name string) string {
	if !s.enabled {
		return ""
	}
	return s.values[strings.ToLower(name)]
}
