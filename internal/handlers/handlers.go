package handlers

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		http.Error(w, "Not found: "+r.URL.Path, http.StatusNotFound)
		return
	}

	if _, err := os.Stat("index.html"); err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, "index.html")
}

func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	var file multipart.File
	var header *multipart.FileHeader
	var err error

	fieldNames := []string{"file", "myFile", "upload", "document"}
	for _, fieldName := range fieldNames {
		file, header, err = r.FormFile(fieldName)
		if err == nil {
			break
		}
	}

	if err != nil {
		http.Error(w, "No file uploaded: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(fileData) == 0 {
		http.Error(w, "File is empty", http.StatusBadRequest)
		return
	}

	convertedData, err := service.ConvertText(string(fileData))
	if err != nil {
		http.Error(w, "Error converting: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".txt"
	}

	timestamp := time.Now().UTC().Format("20060102_150405.000")
	outputFilename := "converted_" + timestamp + ext

	outputPath := filepath.Join("uploads", outputFilename)
	if err := os.MkdirAll("uploads", 0755); err != nil {
		http.Error(w, "Error creating uploads directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := os.WriteFile(outputPath, []byte(convertedData), 0644); err != nil {
		http.Error(w, "Error writing file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(convertedData))
}
