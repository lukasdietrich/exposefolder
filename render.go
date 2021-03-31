package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"path"
	"sort"
	"strconv"
	"time"
)

var (
	//go:embed *.html
	templateFiles embed.FS
	templates     = template.Must(template.New("").Funcs(templateFuncs).ParseFS(templateFiles, "*"))
	templateFuncs = template.FuncMap{
		"dir":         path.Dir,
		"base":        path.Base,
		"join":        path.Join,
		"fmtBytes":    fmtBytes,
		"fmtTime":     fmtTime,
		"fmtFilename": fmtFilename,
		"sortEntries": sortEntries,
	}
)

func fmtBytes(b int64) string {
	const (
		kilo = 1 << 10
		mega = 1 << 20
		giga = 1 << 30
	)

	switch {
	case b >= giga:
		return fmt.Sprintf("%.1fG", float64(b)/giga)
	case b >= mega:
		return fmt.Sprintf("%.1fM", float64(b)/mega)
	case b >= kilo:
		return fmt.Sprintf("%.1fK", float64(b)/kilo)
	default:
		return strconv.FormatInt(b, 10)
	}
}

func fmtTime(t time.Time) string {
	return t.Format(time.UnixDate)
}

func fmtFilename(info fs.FileInfo) string {
	name := info.Name()

	if info.IsDir() {
		name += "/"
	}

	return name
}

func sortEntries(entries []fs.FileInfo) []fs.FileInfo {
	sort.Slice(entries, func(i, j int) bool {
		a, b := entries[i], entries[j]
		return (a.IsDir() == b.IsDir() && a.Name() < b.Name()) || a.IsDir()
	})

	return entries
}

type FolderData struct {
	Path    string
	Entries []fs.FileInfo
}

func renderFolder(w io.Writer, data FolderData) error {
	return templates.ExecuteTemplate(w, "folder.html", data)
}
