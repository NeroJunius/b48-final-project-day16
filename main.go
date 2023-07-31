package main

import (
	"batch48/connection"
	"batch48/middleware"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Projects struct {
	ID          int
	ProjectName string
	Author      string
	AuthorID int

	StartDateFormat string
	EndDateFormat   string
	DurationFormat  string

	DescriptionProject string
	NodeJS             bool
	ReactJS            bool
	NextJS             bool
	TypeScript         bool
	Img                string

	StartDate time.Time
	EndDate   time.Time
	Duration  time.Duration
}

type User struct {
	ID      int
	Name     string
	Email    string
	Password string
}

type UserLoginSession struct {
	IsLogin bool
	Name    string
}

var userLoginSession = UserLoginSession {}

func main() {
	e := echo.New()

	connection.DatabaseConnect()

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session"))))

	e.Static("/assets", "assets")
	e.Static("/uploads", "uploads")

	// project //
	e.GET("/", Home)
	e.GET("/contactMe", contactMe)
	e.GET("/testimonial", testimonials)
	e.GET("/createProject", projectPage)
	e.GET("/projectDetail/:id", projectDetail)
	e.GET("/editProject/:id", editProject)

	// processing //
	e.POST("/add-project", middleware.UploadFile(addProject))
	e.POST("/edit-project/:id", middleware.UploadFile(editedProject))
	e.POST("/delete-project/:id", DeleteProject)

	// LOG IN //
	e.GET("/login-page", LoginPage)
	e.POST("/logged-in", LoggedIn)

	// LOG IN //
	e.POST("/log-out", LoggedOut)

	// REGISTER IN //
	e.GET("/register-page", RegisterPage)
	e.POST("/registered", Registered)

	fmt.Println("server started on port 5900")
	e.Logger.Fatal(e.Start("localhost:5900"))
}

// List Fungsi GET Project nya /

func Home(c echo.Context) error {

	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userLoginSession.IsLogin = false
	} else {
		userLoginSession.IsLogin = sess.Values["isLogin"].(bool)
		userLoginSession.Name = sess.Values["name"].(string)
	}

	tmpl, err := template.ParseFiles("tabs/index.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	var results []Projects

	if sess.Values["isLogin"] != true {

		userLoginSession.IsLogin = false
		
		data, errorData := connection.Conn.Query(context.Background(),
	 	"SELECT id, author, project_title, start_date, end_date, description, node_js, react_js, next_js, type_script, image  FROM tb_projects")

		if errorData != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": errorData.Error()})
		}
		
		for data.Next() {
			var each = Projects{}

			err := data.Scan(&each.ID, &each.Author, &each.ProjectName, &each.StartDate, &each.EndDate, &each.DescriptionProject, &each.NodeJS, &each.ReactJS, &each.NextJS, &each.TypeScript, &each.Img )

			if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
			}

			each.Duration = each.EndDate.Sub(each.StartDate)
			each.DurationFormat = DurationFormat(each.Duration)
			results = append(results, each)
		}
	} else {
		userLoginSession.IsLogin = true

		id := sess.Values["id"].(int)

		dataProjects, errProjects := connection.Conn.Query(context.Background(), 
		"SELECT tb_projects.id, tb_projects.project_title, tb_projects.start_date, tb_projects.end_date,tb_projects.description, tb_projects.node_js, tb_projects.react_js, tb_projects.next_js, tb_projects.type_script, tb_projects.image, tb_user.name AS author FROM tb_projects JOIN tb_user ON tb_projects.author_id = tb_user.id WHERE tb_user.id=$1 ORDER BY tb_projects.id DESC", id)
		if errProjects != nil {
		 return c.JSON(http.StatusInternalServerError, err.Error())
		}
		for dataProjects.Next() {
		 var each = Projects{}
	   
		 err := dataProjects.Scan(&each.ID, &each.ProjectName, &each.StartDate, &each.EndDate, &each.DescriptionProject, &each.NodeJS, &each.ReactJS, &each.NextJS, &each.TypeScript, &each.Img, &each.Author)
		 if err != nil {
		  return c.JSON(http.StatusInternalServerError, err.Error())
		 }
		//  // each.Image = ""
		//  each.Start = each.StartDate.Format("2006-01-02")
		//  each.End = each.EndDate.Format("2006-01-02")
		each.Duration = each.EndDate.Sub(each.StartDate)
		each.DurationFormat = DurationFormat(each.Duration)
		results = append(results, each)
		}
	   }

	// fmt.Println(results)

	projects := map[string]interface{}{
		"Project": results,
		"FlashStatus": sess.Values["status"],
		"FlashMessage": sess.Values["message"],
		"UserLoginSession": userLoginSession,
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	return tmpl.Execute(c.Response(), projects)
}

func contactMe(c echo.Context) error {

	// ses //
	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userLoginSession.IsLogin = false
	} else {
		userLoginSession.IsLogin = sess.Values["isLogin"].(bool)
		userLoginSession.Name = sess.Values["name"].(string)
	}

	projects := map[string]interface{}{
		"UserLoginSession": userLoginSession,
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	var tmpl, err = template.ParseFiles("tabs/contact.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), projects)
}

func projectPage(c echo.Context) error {

	// ses //
	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userLoginSession.IsLogin = false
	} else {
		userLoginSession.IsLogin = sess.Values["isLogin"].(bool)
		userLoginSession.Name = sess.Values["name"].(string)
	}

	projects := map[string]interface{}{
		"UserLoginSession": userLoginSession,
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())
	// end ses //

	var tmpl, err = template.ParseFiles("tabs/project.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), projects)
}

func testimonials(c echo.Context) error {

	// ses //
	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userLoginSession.IsLogin = false
	} else {
		userLoginSession.IsLogin = sess.Values["isLogin"].(bool)
		userLoginSession.Name = sess.Values["name"].(string)
	}

	projects := map[string]interface{}{
		"UserLoginSession": userLoginSession,
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())
	// end ses //

	var tmpl, err = template.ParseFiles("tabs/testimonial.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), projects)
}

func projectDetail(c echo.Context) error {
	Id := c.Param("id")
	idToInt, _ := strconv.Atoi(Id)

	projectDetail := Projects {}

	err := connection.Conn.QueryRow(context.Background(), "SELECT project_title, start_date, end_date, description, node_js, react_js, next_js, type_script, image, author FROM tb_projects WHERE id=$1", idToInt).Scan(
		&projectDetail.ProjectName, &projectDetail.StartDate, &projectDetail.EndDate, &projectDetail.DescriptionProject, &projectDetail.NodeJS, &projectDetail.ReactJS, &projectDetail.NextJS, &projectDetail.TypeScript, &projectDetail.Img, &projectDetail.Author)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	projectDetail.StartDateFormat = projectDetail.StartDate.Format("2006-01-02")
	projectDetail.EndDateFormat = projectDetail.EndDate.Format("2006-01-02")
	
	projectDetail.Duration = projectDetail.EndDate.Sub(projectDetail.StartDate)
	projectDetail.DurationFormat = DurationFormat(projectDetail.Duration)

	projectDetail = Projects {

		ProjectName : projectDetail.ProjectName,
		DurationFormat: projectDetail.DurationFormat,
		DescriptionProject : projectDetail.DescriptionProject,
		Img: projectDetail.Img,
		StartDateFormat: projectDetail.StartDateFormat,
		EndDateFormat: projectDetail.EndDateFormat,
		NodeJS: projectDetail.NodeJS,
		ReactJS: projectDetail.ReactJS,
		NextJS: projectDetail.NextJS,
		TypeScript: projectDetail.TypeScript,
		Author: projectDetail.Author,

	}

	tmpl, err := template.ParseFiles("tabs/project-detail.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userLoginSession.IsLogin = false
	} else {
		userLoginSession.IsLogin = true
		userLoginSession.Name = sess.Values["name"].(string)
	}

	data := map[string]interface{}{ // interface -> tipe data apapun
		"Id":   Id,
		"UserLoginSession" : userLoginSession,
		"Project": projectDetail,
	}

	return tmpl.Execute(c.Response(), data)
}

func editProject(c echo.Context) error {

	// ses //
	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userLoginSession.IsLogin = false
	} else {
		userLoginSession.IsLogin = sess.Values["isLogin"].(bool)
		userLoginSession.Name = sess.Values["name"].(string)
	}

	id, _ := strconv.Atoi(c.Param("id"))

	var Previous_Data = Projects{}

	err := connection.Conn.QueryRow(context.Background(),
		"SELECT project_title, start_date, end_date, description, node_js, react_js, next_js, type_script, image FROM tb_projects WHERE id=$1", id).Scan(&Previous_Data.ProjectName, &Previous_Data.StartDate, &Previous_Data.EndDate, &Previous_Data.DescriptionProject, &Previous_Data.NodeJS, &Previous_Data.ReactJS, &Previous_Data.NextJS, &Previous_Data.TypeScript, &Previous_Data.Img)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	data := map[string]interface{}{
		"ID":            id,
		"Previous_Data": Previous_Data,
	}

	tmpl, errTemp := template.ParseFiles("tabs/edit-form.html")

	if errTemp != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), data)
}

// LIST QUERY PROJECT //

// time //
func DurationFormat(Duration time.Duration) string {
	if Duration <= 24*time.Hour {
		return "Less than a day"
	}

	Days := int(Duration.Hours() / 24)
	Weeks := Days / 7
	Months := Days / 30
	Years := Months / 12

	if Years > 1 {
		return fmt.Sprintf("%d years", Years)
	} else if Years == 1 {
		return "A year"
	} else if Months > 1 {
		return fmt.Sprintf("%d months", Months)
	} else if Months == 1 {
		return "A month"
	} else if Weeks > 1 {
		return fmt.Sprintf("%d weeks", Weeks)
	} else if Weeks == 1 {
		return "A week"
	} else if Days > 1 {
		return fmt.Sprintf("%d days", Days)
	} else {
		return "A day"
	}
}

// func duration (startDate time.Time, endDate time.Time) string {
// 	duration := endDate.Sub(startDate)
// 	days := int(duration.Hours() / 24)
// 	weeks := days / 7
// 	months := days / 30

// 	if months > 12 {
// 		return strconv.Itoa(months/12) + " Year"
// 	}
// 	if months > 0 {
// 		return strconv.Itoa(months) + " Months"
// 	}
// 	if weeks > 0 {
// 		return strconv.Itoa(weeks) + " Weeks"
// 	}
// 	return strconv.Itoa(days) + " days"
// }

// buat project nya //
func addProject(c echo.Context) error {
	
	ProjectName := c.FormValue("projectName")
	StartDate := c.FormValue("startDate")
	EndDate := c.FormValue("endDate")
	DescriptionProject := c.FormValue("projectDescription")

	var NodeJS bool
	if c.FormValue("nodeJS") == "yes" {
		NodeJS = true
	}
	var NextJS bool
	if c.FormValue("nextJS") == "yes" {
		NextJS = true
	}
	var ReactJS bool
	if c.FormValue("reactJS") == "yes" {
		ReactJS = true
	}
	var TypeScript bool
	if c.FormValue("typeScript") == "yes" {
		TypeScript = true
	}

	Img := c.Get("dataFile").(string)

	sess, _ := session.Get("session", c)

	_, err := connection.Conn.Exec(context.Background(),
		"INSERT INTO tb_projects (project_title, start_date, end_date, description, node_js, react_js, next_js, type_script, image, author_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		ProjectName, StartDate, EndDate, DescriptionProject, NodeJS, ReactJS, NextJS, TypeScript, Img, sess.Values["id"].(int))

	fmt.Println("id:", sess.Values["id"])

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	fmt.Println(ProjectName, StartDate, EndDate, DescriptionProject, NodeJS, ReactJS, NextJS, TypeScript, Img)

	// dataProjects = append(dataProjects, createProject)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

// edit projectnya //
func editedProject(c echo.Context) error {
	id := c.Param("id")
	idToInt, _ := strconv.Atoi(id)
	sess, _ := session.Get("session", c)

	ProjectName := c.FormValue("projectName")
	StartDate := c.FormValue("startDate")
	EndDate := c.FormValue("endDate")
	DescriptionProjects := c.FormValue("projectDescription")

	var NodeJS bool
	if c.FormValue("nodeJS") == "yes" {
		NodeJS = true
	}
	var NextJS bool
	if c.FormValue("nextJS") == "yes" {
		NextJS = true
	}
	var ReactJS bool
	if c.FormValue("reactJS") == "yes" {
		ReactJS = true
	}
	var TypeScript bool
	if c.FormValue("typeScript") == "yes" {
		TypeScript = true
	}

	author := sess.Values["id"].(int)

	projectEdit := Projects{}

	erredit := connection.Conn.QueryRow(context.Background(), "SELECT image FROM tb_projects WHERE id=$1", idToInt).Scan(&projectEdit.Img)
	if erredit != nil {
		return c.JSON(http.StatusInternalServerError, erredit.Error())
	}

	img := c.Get("dataFile").(string)

	_, err := connection.Conn.Exec(context.Background(),
		"UPDATE tb_projects SET project_title=$1, start_date=$2, end_date=$3, description=$4, node_js=$5, react_js=$6, next_js=$7, type_script=$8, image=$9, author_id=$10 WHERE id=$11",
		ProjectName, StartDate, EndDate, DescriptionProjects, NodeJS, ReactJS, NextJS, TypeScript, img, author ,id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

// delete project //
func DeleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// dataProjects = append(dataProjects[:id], dataProjects[id+1:]...)

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	fmt.Println("Berhasil menghapus project!")

	return c.Redirect(http.StatusMovedPermanently, "/")
}

// AUTH & SESSION //
// list log in //
func LoginPage(c echo.Context) error {
	// ses //
	sess, _ := session.Get("session", c)

	flash := map[string]interface{}{
		"FlashStatus":  sess.Values["status"],
		"FlashMessage": sess.Values["message"],
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	// access page

	tmpl, err := template.ParseFiles("tabs/auth & session/login-page.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), flash)
}

func LoggedIn(c echo.Context) error {

	LoginEmail := c.FormValue("email")
	LoginPass := c.FormValue("password")

	LoginUser := User{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE email=$1", LoginEmail).Scan(&LoginUser.ID, &LoginUser.Name, &LoginUser.Email, &LoginUser.Password)
	if err != nil {
		return redirectWithMessage(c, "Email/Password Incorrect!", false, "/login-page")
	}

	errPassword := bcrypt.CompareHashAndPassword([]byte(LoginUser.Password), []byte(LoginPass))
	if errPassword != nil {
		return redirectWithMessage(c, "Email/Password Incorrect", false, "/login-page")
	}

	// ses //
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = 10800 // 3 JAM -> berapa lama expired
	sess.Values["message"] = "Login success!"
	sess.Values["status"] = true
	sess.Values["name"] = LoginUser.Name
	sess.Values["email"] = LoginUser.Email
	sess.Values["id"] = LoginUser.ID
	sess.Values["isLogin"] = true
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/")
}

// list register //
func RegisterPage(c echo.Context) error {
	// ses //
	sess, errSess := session.Get("session", c)
	if errSess != nil {
		return c.JSON(http.StatusInternalServerError, errSess.Error())
	}

	flash := map[string]interface{}{
		"FlashMessage": sess.Values["message"], // "Register berhasil"
		"FlashStatus":  sess.Values["status"],  // true
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	tmpl, err := template.ParseFiles("tabs/auth & session/register-page.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return tmpl.Execute(c.Response(), flash)
}

func Registered(c echo.Context) error {

	RegistUsername := c.FormValue("RegistUsername")
	RegistEmail := c.FormValue("email")
	RegistPass := c.FormValue("password")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(RegistPass), 10)

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES ($1, $2, $3)", RegistUsername, RegistEmail, hashedPassword)

	if err != nil {
		redirectWithMessage(c, "Register failed, please try again!", false, "/register-page")
	}

	fmt.Println(RegistUsername, RegistEmail, RegistPass)

	return redirectWithMessage(c, "Register success!", true, "/login-page")
}

// LOG OUT //
func LoggedOut(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	return redirectWithMessage(c, "Log out success!", true, "/login-page")
}

func redirectWithMessage(c echo.Context, message string, status bool, redirectPath string) error {

	sess, errSess := session.Get("session", c)

	if errSess != nil {
		return c.JSON(http.StatusInternalServerError, errSess.Error())
	}

	sess.Values["message"] = message
	sess.Values["status"] = status
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusMovedPermanently, redirectPath)
}
