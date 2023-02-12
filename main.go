package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"personal-web/connection"
	"personal-web/middleware"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type MetaData struct {
	Title     string
	IsLogin   bool
	Id        int
	Username  string
	FlashData string
}

var Data = MetaData{
	Title: "Personal Web |",
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

type Project struct {
	Id              int
	Title           string
	StartDate       time.Time
	EndDate         time.Time
	StartFormatDate string
	EndFormatDate   string
	DurationDate    string
	Description     string
	Technologies    []string
	Image           string
	Author          string
	IsLogin         bool
}

func main() {
	//deklarasi new routermiddleware
	router := mux.NewRouter()

	//connect ke database
	connection.DbConnect()

	//static folder
	router.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	router.PathPrefix("/uploads").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	//handling URL
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/add-project", addProject).Methods("GET")
	router.HandleFunc("/add-data-project", middleware.UploadFile(addDataProject)).Methods("POST")
	router.HandleFunc("/edit-project/{id}", editProject).Methods("GET")
	router.HandleFunc("/update-data-project/{id}", updateProject).Methods("POST")
	router.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")
	router.HandleFunc("/detail-project/{id}", detailProject).Methods("GET")
	router.HandleFunc("/contact-me", contactMe).Methods("GET")
	router.HandleFunc("/login", formLogin).Methods("GET")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/logout", logout).Methods("POST")
	router.HandleFunc("/register", formRegister).Methods("GET")
	router.HandleFunc("/register", register).Methods("POST")

	fmt.Println("Server running on port 5000")
	http.ListenAndServe("localhost:5000", router)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var result []Project

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false

		rows, _ := connection.Conn.Query(context.Background(), "select tb_project.id, title, description, image, start_date, end_date, technologies, tb_user.name as author from tb_project left join tb_user on tb_project.author_id = tb_user.id order by id desc;")

		for rows.Next() {
			var each = Project{}

			var err = rows.Scan(&each.Id, &each.Title, &each.Description, &each.Image, &each.StartDate, &each.EndDate, &each.Technologies, &each.Author)

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			var store = sessions.NewCookieStore([]byte("SESSION_ID"))
			session, _ := store.Get(r, "SESSION_ID")

			if session.Values["IsLogin"] != true {
				each.IsLogin = false
			} else {
				each.IsLogin = session.Values["IsLogin"].(bool)
			}

			each.DurationDate = getDuration(each.StartDate, each.EndDate)

			result = append(result, each)
		}
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Name"].(string)
		Data.Id = session.Values["ID"].(int)

		rows, _ := connection.Conn.Query(context.Background(), "select tb_project.id, title, description, image, start_date, end_date, technologies, tb_user.name as author from tb_project left join tb_user on tb_project.author_id = tb_user.id where author_id=$1 order by id desc;", Data.Id)

		for rows.Next() {
			var each = Project{}

			var err = rows.Scan(&each.Id, &each.Title, &each.Description, &each.Image, &each.StartDate, &each.EndDate, &each.Technologies, &each.Author)

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			var store = sessions.NewCookieStore([]byte("SESSION_ID"))
			session, _ := store.Get(r, "SESSION_ID")

			if session.Values["IsLogin"] != true {
				each.IsLogin = false
			} else {
				each.IsLogin = session.Values["IsLogin"].(bool)
			}

			each.DurationDate = getDuration(each.StartDate, each.EndDate)

			result = append(result, each)
		}
	}

	resp := map[string]interface{}{
		"Data":    Data,
		"Project": result,
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

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Name"].(string)
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func formLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/login.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Name"].(string)
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	user := User{}

	err = connection.Conn.QueryRow(context.Background(), "select * from tb_user where email=$1", email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	session.Values["IsLogin"] = true
	session.Values["Name"] = user.Name
	session.Values["ID"] = user.Id
	session.Options.MaxAge = 10800

	session.AddFlash("Login success!", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func logout(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func formRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/register.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Name"].(string)
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal()
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "insert into tb_user (name, email, password) values ($1, $2, $3)", name, email, passwordHash)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
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

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Name"].(string)
	}

	DetailProject := Project{}

	err = connection.Conn.QueryRow(context.Background(), "select * from tb_project where id=$1", id).Scan(&DetailProject.Id, &DetailProject.Title, &DetailProject.Description, &DetailProject.Image, &DetailProject.Author, &DetailProject.StartDate, &DetailProject.EndDate, &DetailProject.Technologies)

	DetailProject.StartFormatDate = DetailProject.StartDate.Format("2 January 2006")
	DetailProject.EndFormatDate = DetailProject.EndDate.Format("2 January 2006")

	DetailProject.DurationDate = getDuration(DetailProject.StartDate, DetailProject.EndDate)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	resp := map[string]interface{}{
		"Data":    Data,
		"Project": DetailProject,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/add-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Name"].(string)
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
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
	imageContext := r.Context().Value("dataFile")
	image := imageContext.(string)

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")
	author := session.Values["ID"].(int)

	_, err = connection.Conn.Exec(context.Background(), "insert into tb_project (title, description, image, author_id, start_date, end_date, technologies) values ($1, $2, $3, $4, $5, $6, $7)", title, description, image, author, startDate, endDate, technologies)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
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

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Name"].(string)
	}

	EditProject := Project{}

	err = connection.Conn.QueryRow(context.Background(), "select * from tb_project where id=$1", id).Scan(&EditProject.Id, &EditProject.Title, &EditProject.Description, &EditProject.Image, &EditProject.Author, &EditProject.StartDate, &EditProject.EndDate, &EditProject.Technologies)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	startDate := EditProject.StartDate.Format("2006-01-02")
	endDate := EditProject.EndDate.Format("2006-01-02")

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

	resp := map[string]interface{}{
		"Data":       Data,
		"Project":    EditProject,
		"StartDate":  startDate,
		"EndDate":    endDate,
		"NodeJS":     iNode,
		"ReactJS":    iReact,
		"NextJS":     iNext,
		"Typescript": iTypescript,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal()
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	title := r.PostForm.Get("title")
	startDate := r.PostForm.Get("start-date")
	endDate := r.PostForm.Get("end-date")
	description := r.PostForm.Get("description")
	technologies := r.Form["tech"]

	_, err = connection.Conn.Exec(context.Background(), "update tb_project set title=$1, description=$2, start_date=$3, end_date=$4, technologies=$5 where id=$6", title, description, startDate, endDate, technologies, id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Conn.Exec(context.Background(), "delete from tb_project where id=$1", id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func getDuration(s, e time.Time) string {
	diff := e.Sub(s).Hours()

	Month := math.Floor(diff / (30 * 24))
	Week := math.Floor(diff / (7 * 24))
	Day := math.Floor(diff / (24))

	var DurationDate string

	if Month > 1 {
		DurationDate = strconv.FormatFloat(Month, 'f', 0, 64) + " months"
	} else if Month > 0 {
		DurationDate = strconv.FormatFloat(Month, 'f', 0, 64) + " month"
	} else {
		if Week > 1 {
			DurationDate = strconv.FormatFloat(Week, 'f', 0, 64) + " weeks"
		} else if Week > 0 {
			DurationDate = strconv.FormatFloat(Week, 'f', 0, 64) + " week"
		} else {
			if Day > 1 {
				DurationDate = strconv.FormatFloat(Day, 'f', 0, 64) + " days"
			} else if Day > 0 {
				DurationDate = strconv.FormatFloat(Day, 'f', 0, 64) + " day"
			}
		}
	}

	return DurationDate
}
