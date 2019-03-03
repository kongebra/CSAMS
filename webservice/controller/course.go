package controller

import (
	"fmt"
	"github.com/JohanAanesen/NTNU-Bachelor-Management-System-For-CS-Assignments/webservice/model"
	"github.com/JohanAanesen/NTNU-Bachelor-Management-System-For-CS-Assignments/webservice/shared/session"
	"github.com/JohanAanesen/NTNU-Bachelor-Management-System-For-CS-Assignments/webservice/shared/view"
	"github.com/shurcooL/github_flavored_markdown"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

//CourseGET serves class page to users
func CourseGET(w http.ResponseWriter, r *http.Request) {

	var course model.Course

	//check if request has an classID
	id := r.FormValue("id")
	if id == "" {
		//redirect to error pageinfo
		ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	//check that id is a number
	courseID, err := strconv.Atoi(id)
	if err != nil {
		//redirect to error pageinfo
		log.Println(err.Error())
		ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	//check if user is logged in
	if !session.GetUserFromSession(r).Authenticated {
		LoginGET(w, r)
		return
	}

	//get user
	user := session.GetUserFromSession(r)

	//course repo
	courseRepo := &model.CourseRepository{}

	//get info from db
	course, err = courseRepo.GetSingle(courseID)
	if err != nil {
		log.Println(err)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	assignmentRepo := model.AssignmentRepository{}
	assignments, err := assignmentRepo.GetAllFromCourse(courseID)
	if err != nil {
		log.Println(err)
	}

	//check if user is an participant of said class or a teacher
	participant := courseRepo.UserExistsInCourse(user.ID, courseID) || user.ID == course.Teacher
	if !participant {
		log.Println("user not participant of class")
		ErrorHandler(w, r, http.StatusUnauthorized)
		return
	}

	classmates := model.GetUsersToCourse(courseID)

	//all a-ok
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	md := []byte(course.Description)
	description := github_flavored_markdown.Markdown(md) //todo sanitize markdown

	v := view.New(r)
	v.Name = "course"

	v.Vars["Course"] = course
	v.Vars["User"] = user
	v.Vars["Classmates"] = classmates
	v.Vars["Description"] = template.HTML(description)
	v.Vars["Assignments"] = assignments

	v.Render(w)
}

//CourseListGET serves class list page to users
func CourseListGET(w http.ResponseWriter, r *http.Request) {

	//check if request has an classID
	if r.Method == http.MethodGet {

		id := r.FormValue("id")

		if id == "" {
			//redirect to error pageinfo
			ErrorHandler(w, r, http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "ID is %s\n", id)
	}
	//check if user is an participant of said class or a teacher

	//get classlist from db

	//parse info to html template
	temp, err := template.ParseFiles("template/test.html")
	if err != nil {
		log.Fatal(err)
		ErrorHandler(w, r, http.StatusInternalServerError)
	}

	temp.Execute(w, nil)
}
