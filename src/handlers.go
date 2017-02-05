package munhasir

import (
	"net/http"
	"html/template"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not allowed.", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "text/html")

	var funcMap = template.FuncMap{
		"titleF": func() string {
			return "munhasir"
		},

	}

	template, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	info := make(map[string]string)
	info["title"] = "munhasir"
	info["content"] = "A platform to keeping diaries for those who are cautious(or paranoid)."

	err = template.Execute(w, info)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}