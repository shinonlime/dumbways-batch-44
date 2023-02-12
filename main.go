package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var Data = map[string]interface{}{
	"Title": "Personal Web |",
}

func main() {
	router := mux.NewRouter()

	router.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/contact-me", contactMe).Methods("GET")
	router.HandleFunc("/add-project", addProject).Methods("GET")
	router.HandleFunc("/add-data-project", addDataProject).Methods("POST")
	router.HandleFunc("/detail-project/{id}", detailProject).Methods("GET")

	fmt.Println("Server running on port 5000")
	http.ListenAndServe("localhost:5000", router)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func contactMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact-me.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/add-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func detailProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var tmpl, err = template.ParseFiles("views/detail-project.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	resp := map[string]interface{}{
		"Data": Data,
		"Id":   id,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func addDataProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal()
	}

	fmt.Println(r.PostForm.Get("title"))
	fmt.Println(r.PostForm.Get("start-date"))
	fmt.Println(r.PostForm.Get("end-date"))
	fmt.Println(r.PostForm.Get("description"))
	fmt.Println(r.Form["technologies"])

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}
