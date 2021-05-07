package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

type UserHandler struct {
	DB *sql.DB
}

func main() {

	fmt.Println("Welcome to API")

	e := echo.New()
	g := e.Group("/admin")

	// this logs the server interaction
	g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}]  ${status}  ${method} ${host}${path} ${latency_human}` + "\n",
	}))

	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// check in the DB
		if username == "jack" && password == "1234" {
			return true, nil
		}

		return true, nil
	}))

	g.GET("/main", mainAdmin)

	//กำหนด Route ก่อนเลย พร้อมให้ call ไปยัง func ต่างๆ
	h := UserHandler{}
	h.Initialize()

	e.GET("/users", h.GetUser)

	e.GET("/GetUser/:email", h.Getusers)

	e.POST("/user/add", h.PostUser)

	e.PUT("/user/update/:id", h.PutUser)

	e.DELETE("/user/delete/:id", h.DelUser)

	e.Logger.Fatal(e.Start(":9999"))
}

//ให้เชื่อมต่อฐานข้อมูลเมื่อ Initialize
func (h *UserHandler) Initialize() {
	connStr := "postgres://postgres:sanchai@34.87.1.245/demolab?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	h.DB = db
}

type Userhotel struct {
	Id            uint   `sql:"primary_key" json:"id" form:"id" query:"id"`
	Email         string `json:"email" form:"email" query:"email"`
	Password      string `json:"password" form:"password" query:"password"`
	Userpermis_id int    `json:"userpermis_id" form:"userpermis_id" query:"userpermis_id"`
}

func mainAdmin(c echo.Context) error {

	return c.String(http.StatusOK, "Are you admin")
}

func (h *UserHandler) GetUser(c echo.Context) error {

	rows, err := h.DB.Query("SELECT id , email, password, userpermis_id from userhotel order by id")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var userhotels []Userhotel

	for rows.Next() {
		var data Userhotel
		err := rows.Scan(&data.Id, &data.Email, &data.Password, &data.Userpermis_id)
		if err != nil {
			log.Fatal("Scan failed:", err.Error())
		}
		userhotels = append(userhotels, data)
	}

	return c.JSON(http.StatusOK, userhotels)

}

func (h *UserHandler) Getusers(c echo.Context) error {

	id := c.Param("email")
	var data Userhotel
	rows, err := h.DB.Prepare("SELECT id , email, password, userpermis_id from userhotel WHERE email = $1")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var userhotels []Userhotel
	
	
		err = rows.QueryRow(id).Scan(&data.Id, &data.Email, &data.Password, &data.Userpermis_id)
		if err != nil {
			log.Fatal("Scan failed:", err.Error())
		}
		userhotels = append(userhotels, data)
	

	return c.JSON(http.StatusOK, userhotels)
}

func (h *UserHandler) PostUser(c echo.Context) (err error) {

	data := new(Userhotel)

	if err = c.Bind(data); err != nil {
		return err
	}

	log.Println(data)

	stmt, err := h.DB.Prepare("INSERT INTO userhotel(email,password,userpermis_id) VALUES($1,$2,$3)")

	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}

	_, err = stmt.Exec(data.Email, data.Password, data.Userpermis_id)

	fmt.Println("Email :", data.Email)
	fmt.Println("Password :", data.Password)

	if err != nil {
		log.Fatal("DATABASE INSERT Error :", err.Error())
	}

	defer stmt.Close()

	return c.JSON(http.StatusOK, data)
}

func (h *UserHandler) PutUser(c echo.Context) (err error) {

	id := c.Param("id")

	data := new(Userhotel)

	if err := c.Bind(data); err != nil {
		return err
	}

	stmt, err := h.DB.Prepare("UPDATE userhotel SET email=$1 , password=$2 , userpermis_id=$3 WHERE id = $4 ")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(data.Email, data.Password, data.Userpermis_id, id)

	if err != nil {
		return err
	}

	defer stmt.Close()

	return c.JSON(http.StatusOK, data)
}

func (h *UserHandler) DelUser(c echo.Context) (err error) {

	id := c.Param("id")

	stmt, err := h.DB.Prepare("DELETE FROM userhotel WHERE id = $1")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)

	if err != nil {
		log.Println("Database DELETE failed:", err.Error())
	}
	defer stmt.Close()

	return c.JSON(http.StatusOK, "OK")

}
