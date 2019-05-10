package main

//Try Scan
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"


	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
)

type GlobalQuote struct {
  ID string `json:"id"`
	Symbol string    `json:"Symbol"`
	Price int    `json:"Price"`
}

type Status struct {
	Table       string `json:"table"`
	Recordcount int    `json:"recordCount"`
}

var status []Status

//Get All Table Items
func getTableItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	svc := dynamodb.New(sess)

	params := &dynamodb.ScanInput{
		TableName: aws.String("AppleStock"),
	}

	result, err := svc.Scan(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	obj := []GlobalQuote{}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &obj)
	if err != nil {
		fmt.Printf("failed to unmarshal Query result items, %v", err)
	}

	json.NewEncoder(w).Encode(obj)
}

func getTableInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func main() {

	//Init Router
	r := mux.NewRouter()

	status = append(status, Status{Table: "AppleStock", Recordcount: 3})

	//Route Handlers / Endpoints
	r.HandleFunc("/tdominic/all", getTableItems).Methods("GET")
	r.HandleFunc("/tdominic/status", getTableInfo).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
