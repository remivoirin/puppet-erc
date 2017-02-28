package main

import (
	"net/http"
	"regexp"
	"strconv"
	"github.com/labstack/echo"
)

// basic JSON response and explanation
type Basicresp struct {
	Completed	bool	`json:"completed"`
	Comment		string	`json:"comment"`
}

// entry types correlated to DB
type Newentry struct {
	Host_regex	string	`json:"host_regex"`
	Role		string	`json:"role"`
	Comment		string	`json:"comment"`
}

type Fullentry struct {
	Id		string	`json:"id"`
	Host_regex	string	`json:"host_regex"`
	Role		string	`json:"role"`
	Comment		string	`json:"comment"`
}

func Index(c echo.Context) error {
	return c.JSON(http.StatusOK, "Welcome")
}

func Initialize(c echo.Context) error {
	result, reason := Dbinitialize()
	resp := Basicresp{Completed: result, Comment: reason}
	if result == false {
		return c.JSON(http.StatusBadRequest, resp)
	}
	return c.JSON(http.StatusOK, resp)
}

func Insert(c echo.Context) error {
	entry := &Newentry{}

	// check if input JSON fits the Newentry struct
	if err := c.Bind(entry); err != nil {
		return err
	}

	// call the actual insert into DB
	result, reason := Dbinsert(entry.Host_regex, entry.Role, entry.Comment)
	resp := Basicresp{Completed: result, Comment: reason}
	if result == false {
		return c.JSON(http.StatusBadRequest, resp)
	}
	return c.JSON(http.StatusCreated, resp)
}

func List(c echo.Context) error {
	var entries []Fullentry

	// call the actual list
	entries = Dblist()
	return c.JSON(http.StatusOK, entries)
}

func Deletebyid(c echo.Context) error {
	// simple type validation of parameter
	myid, err := strconv.Atoi(c.Param("deleteid"))
	if err != nil {
		return err
	}

	// call the actual delete
	result, reason := Dbdeletebyid(myid)
	resp := Basicresp{Completed: result, Comment: reason}
	if result == false {
		return c.JSON(http.StatusBadRequest, resp)
	}

	// return HTTP 204 if the record is successfully deleted
	return c.NoContent(http.StatusNoContent)
}

// fulltext output for curl / Puppet agents
func Getrolebyhostname(c echo.Context) error {
	var entries []Fullentry

	// get all entries
	entries = Dblist()

	// check hostname against roles
	for _, entry := range entries {
		myregex, err := regexp.Compile(entry.Host_regex)
		if err != nil {
			return err
		}
		if myregex.MatchString(c.Param("hostname")) {
			return c.String(http.StatusOK, entry.Role)
		}
	}

	// return HTTP 500 because default should have poped up eventually
	return c.NoContent(http.StatusInternalServerError)
}
