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

type Project struct {
	Id           int
	Title        string
	StartDate    string
	EndDate      string
	Duration     string
	Description  string
	Technologies []string
}

var Projects = []Project{
	// {
	// 	Title:       "Dumbways Mobile App - 2023",
	// 	Duration:    "2 bulan",
	// 	Description: "This is a longer card with supporting text below as a natural lead-in to additional content. This content is a little bit longer.",
	// },
}

func main() {
	router := mux.NewRouter()

	router.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/contact-me", contactMe).Methods("GET")
	router.HandleFunc("/add-project", addProject).Methods("GET")
	router.HandleFunc("/add-data-project", addDataProject).Methods("POST")
	router.HandleFunc("/detail-project/{id}", detProject).Methods("GET")
	router.HandleFunc("/edit-project/{id}", editProject).Methods("GET")
	router.HandleFunc("/update-project/{id}", updateProject).Methods("POST")
	router.HandleFunc("/delete-project/{id}", delProject).Methods("GET")

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

	resp := map[string]interface{}{
		"Data":     Data,
		"Projects": Projects,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
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

func detProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var tmpl, err = template.ParseFiles("views/detail-project.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	DetailProject := Project{}

	for i, data := range Projects {
		if i == id {
			DetailProject = Project{
				Title:        data.Title,
				StartDate:    data.StartDate,
				EndDate:      data.EndDate,
				Description:  data.Description,
				Technologies: data.Technologies,
			}
		}
	}

	resp := map[string]interface{}{
		"Data":    Data,
		"Project": DetailProject,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func editProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var tmpl, err = template.ParseFiles("views/edit-project.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	EditProject := Project{}

	iNode := false
	iReact := false
	iNext := false
	iTypescript := false

	for _, data := range EditProject.Technologies {
		if data == "Node JS" {
			iNode = true
		}
		if data == "React JS" {
			iReact = true
		}
		if data == "Next JS" {
			iNext = true
		}
		if data == "Typescript" {
			iTypescript = true
		}
	}

	for i, data := range Projects {
		if i == id {
			EditProject = Project{
				Title:       data.Title,
				StartDate:   data.StartDate,
				EndDate:     data.EndDate,
				Description: data.Description,
			}
		}
	}

	resp := map[string]interface{}{
		"Data":       Data,
		"Project":    EditProject,
		"Node":       iNode,
		"React":      iReact,
		"Next":       iNext,
		"Typescript": iTypescript,
		"Id":         id,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func addDataProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal()
	}

	title := r.PostForm.Get("title")
	startDate := r.PostForm.Get("start-date")
	endDate := r.PostForm.Get("end-date")
	description := r.PostForm.Get("description")
	technologies := r.Form["tech"]

	var newProject = Project{
		Title:        title,
		StartDate:    startDate,
		EndDate:      endDate,
		Description:  description,
		Technologies: technologies,
	}

	Projects = append(Projects, newProject)

	fmt.Println(newProject)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	Projects = append(Projects[:id], Projects[id+1:]...)

	err := r.ParseForm()
	if err != nil {
		log.Fatal()
	}

	title := r.PostForm.Get("title")
	startDate := r.PostForm.Get("start-date")
	endDate := r.PostForm.Get("end-date")
	description := r.PostForm.Get("description")
	technologies := r.Form["tech"]

	var newProject = Project{
		Title:        title,
		StartDate:    startDate,
		EndDate:      endDate,
		Description:  description,
		Technologies: technologies,
	}

	Projects = append(Projects, newProject)

	fmt.Println(newProject)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func delProject(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Chace-control", "no-cache, no-store, must-revalidate")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	Projects = append(Projects[:id], Projects[id+1:]...)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
