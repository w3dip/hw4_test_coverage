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
)

//type TestCase struct {
//	ID      string
//	Result  *CheckoutResult
//	IsError bool
//}
//
//type CheckoutResult struct {
//	Status  int
//	Balance int
//	Err     string
//}

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
			//os.Exit(2)
		}

		age, err := strconv.Atoi(row.Age)
		if err != nil {
			// handle error
			fmt.Println(err)
			//os.Exit(2)
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
	result, err := json.Marshal(users)
	if err != nil {
		panic(err)
	}
	io.WriteString(w, string(result))

	//key := r.FormValue("id")
	//switch key {
	//case "42":
	//	w.WriteHeader(http.StatusOK)
	//	io.WriteString(w, `{"status": 200, "balance": 100500}`)
	//case "100500":
	//	w.WriteHeader(http.StatusOK)
	//	io.WriteString(w, `{"status": 400, "err": "bad_balance"}`)
	//case "__broken_json":
	//	w.WriteHeader(http.StatusOK)
	//	io.WriteString(w, `{"status": 400`) //broken json
	//case "__internal_error":
	//	fallthrough
	//default:
	//	w.WriteHeader(http.StatusInternalServerError)
	//}
}

//func parseXml() {
//	logins := make([]string, 0)
//	v := new(Users)
//	err := xml.Unmarshal(xmlData, &v)
//	if err != nil {
//		fmt.Printf("error: %v", err)
//		return
//	}
//	for _, u := range v.List {
//		logins = append(logins, u.Login)
//	}
//}

//type Cart struct {
//	PaymentApiURL string
//}
//
//func (c *Cart) Checkout(id string) (*CheckoutResult, error) {
//	url := c.PaymentApiURL + "?id=" + id
//	resp, err := http.Get(url)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//
//	data, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	result := &CheckoutResult{}
//
//	err = json.Unmarshal(data, result)
//	if err != nil {
//		return nil, err
//	}
//	return result, nil
//}

func TestCartCheckout(t *testing.T) {
	//cases := []TestCase{
	//	TestCase{
	//		ID: "42",
	//		Result: &CheckoutResult{
	//			Status:  200,
	//			Balance: 100500,
	//			Err:     "",
	//		},
	//		IsError: false,
	//	},
	//	TestCase{
	//		ID: "100500",
	//		Result: &CheckoutResult{
	//			Status:  400,
	//			Balance: 0,
	//			Err:     "bad_balance",
	//		},
	//		IsError: false,
	//	},
	//	TestCase{
	//		ID:      "__broken_json",
	//		Result:  nil,
	//		IsError: true,
	//	},
	//	TestCase{
	//		ID:      "__internal_error",
	//		Result:  nil,
	//		IsError: true,
	//	},
	//}

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
	//for caseNum, item := range cases {
	//	c := &Cart{
	//		PaymentApiURL: ts.URL,
	//	}
	//	result, err := c.Checkout(item.ID)
	//
	//	if err != nil && !item.IsError {
	//		t.Errorf("[%d] unexpected error: %#v", caseNum, err)
	//	}
	//	if err == nil && item.IsError {
	//		t.Errorf("[%d] expected error, got nil", caseNum)
	//	}
	//	if !reflect.DeepEqual(item.Result, result) {
	//		t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, item.Result, result)
	//	}
	//}
	ts.Close()
}
