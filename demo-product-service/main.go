package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"ranveer/common"
	"strconv"
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

//type Person struct {
//	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
//	FirstName string `json:"firstname,omitempty" bson:"firstname,omitempty"`
//	LastName string `json:"lastname,omitempty" bson:"lastname,omitempty"`
//}



var client *mongo.Client



type ProductConfiguration struct {
	Categories []string `json:"categories"`
}


//
func registerServiceWithConsul() {
	log.Println("consul registration for product service started...")
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	ip, err :=  common.IPAddr()
	if err != nil {
		log.Fatalf("could not determine IP address to register this service with... %v", err)
	}
	registration := new(consulapi.AgentServiceRegistration)

	registration.ID = "product-service"
	registration.Name = "product-service"
	//address := hostname()
	registration.Address = ip.String()
	port, err := strconv.Atoi(port()[1:len(port())])
	if err != nil {
		log.Fatalln(err)
	}
	registration.Port = port
	registration.Check = new(consulapi.AgentServiceCheck)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck", ip.String(), port)
	registration.Check.Interval = "5s"
	registration.Check.Timeout = "5s"
	consul.Agent().ServiceRegister(registration)
	log.Println(" product service registered to consul...")
}

func main() {
	 registerServiceWithConsul()

	log.Println("starting product service ...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ =  mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	router := mux.NewRouter()
	router.HandleFunc("/product", CreateProductEndpoint).Methods("POST")
	router.HandleFunc("/products", GetProductsEndpoint).Methods("GET")
	router.HandleFunc("/product/{id}", GetProductEndpoint).Methods("GET")
	router.HandleFunc("/product/{id}", DeleteProductByID).Methods("DELETE")
	router.HandleFunc("/healthcheck", healthcheck)
	//router.HandleFunc("/product-configuration", Configuration)
	//http.ListenAndServe(":8888", router)


	fmt.Printf("product service is up on port: %s", port())
	http.ListenAndServe(port(), router)
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, `product service is good`)
	r.Body.Close()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func port() string {
	p := os.Getenv("PRODUCT_SERVICE_PORT")
	if len(strings.TrimSpace(p)) == 0 {
		return ":8888"
	}
	return fmt.Sprintf(":%s", p)
}

func DeleteProductByID(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type","application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	//log.Println("dd id = "+id.String())
	//var product Product

	collection := client.Database("go").Collection("product")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	col,err := collection.DeleteOne(ctx, common.Product{ID: id})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() +`"}`))
		return
	}

	if col.DeletedCount > 0 {
		json.NewEncoder(response).Encode("deleted product id "+params["id"])
	} else {
		json.NewEncoder(response).Encode(params["id"] + " not found")
	}

}



func CreateProductEndpoint (response http.ResponseWriter, request *http.Request){
	response.Header().Add("content-type","application/json")
	var product common.Product
	_ = json.NewDecoder(request.Body).Decode(&product)
	collection := client.Database("go").Collection("product")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, product)
	json.NewEncoder(response).Encode(result)

}
func GetProductsEndpoint (response http.ResponseWriter, request *http.Request){
	response.Header().Add("content-type","application/json")

	var products []common.Product
	collection := client.Database("go").Collection("product")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"`+err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var product common.Product
		cursor.Decode(&product)
		products = append(products, product)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"`+ err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(products)
}

func GetProductEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type","application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var product common.Product
	//_ = json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("go").Collection("product")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, common.Product{ID: id}).Decode(&product)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() +`"}`))
		return
	}
	json.NewEncoder(response).Encode(product)

}



