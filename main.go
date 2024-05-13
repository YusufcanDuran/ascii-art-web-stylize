package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"strings"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed ERROR 405", http.StatusMethodNotAllowed)
		return
	}
	text := strings.TrimSpace(r.FormValue("text"))
	banner := r.FormValue("banner")
	if text == "" || banner == "" {
		http.Error(w, "Missing text or banner ERROR 400 BAD REQUEST", http.StatusBadRequest)
		return
	}
	turkcechar := "ışğüöçİŞÖÇĞÜ"
	if strings.ContainsAny(text, turkcechar) {
		fmt.Fprint(w, "türkçe karakter girildi. ERROR 400 BAD REQUEST")
		return
	}
	asciiArt, err := generateAsciiArt(text, banner)
	if err != nil {
		http.Error(w, "Error generating ASCII Art ERROR 500", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%s", asciiArt)
}
func generateAsciiArt(text, banner string) (string, error) {
	// Run argument.go script passing 'text' and 'banner' as arguments
	cmd := exec.Command("go", "run", "argument.go", text, banner)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running command: %v", err)
	}
	return string(output), nil
}
func main() {
	http.Handle("/templates/style.css", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/generate", asciiArtHandler)
	http.ListenAndServe(":8080", nil)
}
