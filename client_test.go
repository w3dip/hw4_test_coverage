package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

type Root struct {
	XMLName xml.Name `xml:"root"`
	Text    string   `xml:",chardata"`
	Rows    []Row    `xml:"row"`
}

type Row struct {
	Text          string `xml:",chardata"`
	ID            string `xml:"id"`
	Guid          string `xml:"guid"`
	IsActive      string `xml:"isActive"`
	Balance       string `xml:"balance"`
	Picture       string `xml:"picture"`
	Age           string `xml:"age"`
	EyeColor      string `xml:"eyeColor"`
	FirstName     string `xml:"first_name"`
	LastName      string `xml:"last_name"`
	Gender        string `xml:"gender"`
	Company       string `xml:"company"`
	Email         string `xml:"email"`
	Phone         string `xml:"phone"`
	Address       string `xml:"address"`
	About         string `xml:"about"`
	Registered    string `xml:"registered"`
	FavoriteFruit string `xml:"favoriteFruit"`
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("AccessToken")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	limit := r.FormValue("limit")

	query := r.FormValue("query")
	if query == "fail" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if query == "wrong_body" {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"status": 400}`)
		return
	}

	offset := r.FormValue("offset")
	if offset == "100" {
		w.WriteHeader(http.StatusBadRequest)
		error := SearchErrorResponse{
			Error: "UnknownBadRequest",
		}
		error_result, err := json.Marshal(error)
		if err != nil {
			panic(err)
		}
		io.WriteString(w, string(error_result))
		return
	}

	if offset == "101" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if offset == "102" {
		//w.Header().Set("Location", "")
		//return
		panic("Server failed")
	}

	order := r.FormValue("order_field")
	switch order {
	case "Id", "Age", "Name":
		//поиск по полям
		break
	case "":
		//поиск по Name
		break
	case "About":
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		break
	default:
		w.WriteHeader(http.StatusBadRequest)
		error := SearchErrorResponse{
			Error: "ErrorBadOrderField",
		}
		error_result, err := json.Marshal(error)
		if err != nil {
			panic(err)
		}
		io.WriteString(w, string(error_result))
		return
		//panic("Wrong order type")
	}

	file, err := os.Open("dataset.xml")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	xmlData, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	users := make([]User, 0)
	v := new(Root)
	error := xml.Unmarshal(xmlData, &v)
	if error != nil {
		fmt.Printf("error: %v", err)
		return
	}
	for _, row := range v.Rows {
		id, err := strconv.Atoi(row.ID)
		if err != nil {
			// handle error
			fmt.Println(err)
		}

		age, err := strconv.Atoi(row.Age)
		if err != nil {
			// handle error
			fmt.Println(err)
		}
		users = append(users, User{
			Id:     id,
			Name:   row.FirstName + " " + row.LastName,
			Age:    age,
			About:  row.About,
			Gender: row.Gender,
		})
	}
	w.WriteHeader(http.StatusOK)

	new_users := make([]User, 0)
	if limit == "21" {
		for indx := 0; indx < 21; indx++ {
			new_users = append(new_users, users[indx])
		}
	} else {
		for indx := 0; indx < len(users); indx++ {
			new_users = append(new_users, users[indx])
		}
	}
	result, err := json.Marshal(new_users)
	if err != nil {
		panic(err)
	}
	io.WriteString(w, string(result))
}

func TestFindUsers(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	client := &SearchClient{
		URL:         ts.URL,
		AccessToken: "some token",
	}

	req := SearchRequest{
		Limit:      10,
		Offset:     0,
		Query:      "Name",
		OrderField: "",
		OrderBy:    0,
	}

	result, err := client.FindUsers(req)
	if err != nil {
		t.Error("Was error", err)
	}
	if result == nil {
		t.Error("Result is empty")
	}
	ts.Close()
}

func TestFindUsersLimitLessZero(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	client := &SearchClient{
		URL:         ts.URL,
		AccessToken: "some token",
	}

	req := SearchRequest{
		Limit:      -1,
		OrderField: "",
	}
	client.FindUsers(req)
	ts.Close()
}

func TestFindUsersLimitGreaterMax(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	client := &SearchClient{
		URL:         ts.URL,
		AccessToken: "some token",
	}

	req := SearchRequest{
		Limit:      26,
		OrderField: "",
	}
	client.FindUsers(req)
	ts.Close()
}

func TestFindUsersOffsetLessZero(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	client := &SearchClient{
		URL:         ts.URL,
		AccessToken: "some token",
	}

	req := SearchRequest{
		Offset:     -1,
		OrderField: "",
	}
	client.FindUsers(req)
	ts.Close()
}

func TestFindUsersTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	client := &SearchClient{
		URL:         ts.URL,
		AccessToken: "some token",
	}

	req := SearchRequest{
		Offset:     0,
		OrderField: "About",
	}
	client.FindUsers(req)
	ts.Close()
}

func TestFindUsersUnauthorized(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	client := &SearchClient{
		URL:         ts.URL,
		AccessToken: "",
	}

	req := SearchRequest{
		Offset:     0,
		OrderField: "About",
	}
	client.FindUsers(req)
	ts.Close()
}

func TestFindUsersServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	client := &SearchClient{
		URL:         ts.URL,
		AccessToken: "some toke",
	}

	req := SearchRequest{
		Offset: 0,
		Query:  "fail",
	}
	client.FindUsers(req)
	ts.Close()
}

func TestFindUsersBadRequestWithCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	client := &SearchClient{
		URL:         ts.URL,
		AccessToken: "some token",
	}

	req := SearchRequest{
		Offset:     0,
		OrderField: "Gender",
	}
	client.FindUsers(req)
	ts.Close()
}

func TestFindUsersBadRequestWrongOffset(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	client := &SearchClient{
		URL:         ts.URL,
		AccessToken: "some token",
	}

	req := SearchRequest{
		Offset:     100,
		OrderField: "Gender",
	}
	client.FindUsers(req)
	ts.Close()
}

func TestFindUsersBadRequestCommonCase(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	client := &SearchClient{
		URL:         ts.URL,
		AccessToken: "some token",
	}

	req := SearchRequest{
		Offset:     101,
		OrderField: "Gender",
	}
	client.FindUsers(req)
	ts.Close()
}

func TestFindUsersServerFailed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	client := &SearchClient{
		URL:         ts.URL,
		AccessToken: "some token",
	}

	req := SearchRequest{
		Offset:     102,
		OrderField: "Gender",
	}
	client.FindUsers(req)
	ts.Close()
}

func TestFindUsersFailedParsing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	client := &SearchClient{
		URL:         ts.URL,
		AccessToken: "some token",
	}

	req := SearchRequest{
		Limit:      10,
		Offset:     0,
		Query:      "wrong_body",
		OrderField: "",
		OrderBy:    0,
	}
	client.FindUsers(req)
	ts.Close()
}

func TestFindUsersLimitEqualsSize(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	client := &SearchClient{
		URL:         ts.URL,
		AccessToken: "some token",
	}

	req := SearchRequest{
		Limit:      20,
		Offset:     0,
		Query:      "Name",
		OrderField: "",
		OrderBy:    0,
	}

	client.FindUsers(req)
	ts.Close()
}
