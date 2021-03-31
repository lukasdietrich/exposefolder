package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var indexFilenames = []string{
	"index.html",
	"index.htm",
}

type requestContext struct {
	w http.ResponseWriter
	r *http.Request

	path     string
	filename string
}

func newContext(foldername string, w http.ResponseWriter, r *http.Request) requestContext {
	path := filepath.Clean(r.RequestURI)
	filename := filepath.Join(foldername, path)

	return requestContext{
		w:        w,
		r:        r,
		path:     path,
		filename: filename,
	}
}

func makeHandler(foldername string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %q", r.RemoteAddr, r.Method, r.RequestURI)

		if err := handleRequest(newContext(foldername, w, r)); err != nil {
			respondWithStatus(w, http.StatusInternalServerError)
			log.Printf("\t=> %v", err)
		}
	})
}

func respondWithStatus(w http.ResponseWriter, status int) error {
	w.WriteHeader(status)
	_, err := fmt.Fprintln(w, http.StatusText(status))
	return err
}

func handleRequest(ctx requestContext) error {
	switch ctx.r.Method {
	case http.MethodGet:
		return handleGet(ctx)

	case http.MethodPost, http.MethodPut:
		return handlePost(ctx)

	default:
		return errors.New("ignored request because of the http method")
	}
}

func handleGet(ctx requestContext) error {
	info, err := os.Stat(ctx.filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return respondWithStatus(ctx.w, http.StatusNotFound)
		}

		return err
	}

	if info.IsDir() {
		for _, filename := range indexFilenames {
			filename = filepath.Join(ctx.filename, filename)

			if fileExists(filename) {
				ctx.filename = filename
				return serveFile(ctx)
			}
		}
		return serveFolder(ctx)
	} else {
		return serveFile(ctx)
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func serveFile(ctx requestContext) error {
	http.ServeFile(ctx.w, ctx.r, ctx.filename)
	return nil
}

func serveFolder(ctx requestContext) error {
	files, err := os.ReadDir(ctx.filename)
	if err != nil {
		return err
	}

	infos := make([]fs.FileInfo, len(files))
	for i, file := range files {
		infos[i], err = file.Info()
		if err != nil {
			return err
		}
	}

	return renderFolder(ctx.w, FolderData{Path: ctx.path, Entries: infos})
}

func handlePost(ctx requestContext) error {
	r, err := ctx.r.MultipartReader()
	if err != nil {
		return err
	}

	for {
		part, err := r.NextPart()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}

			return err
		}

		filename := part.FileName()
		if filename != "" {
			log.Printf("\tuploading %q to %q", filename, ctx.filename)
		}

		if ctx.r.Method == http.MethodPost && fileExists(filename) {
			log.Printf("\t%q already exists!", filename)
			return respondWithStatus(ctx.w, http.StatusConflict)
		}

		writeFile(part, filepath.Join(ctx.filename, filename))
	}
}

func writeFile(r io.ReadCloser, filename string) error {
	defer r.Close()

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = io.Copy(f, r)
	return err
}
