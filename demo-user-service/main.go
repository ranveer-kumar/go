package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"ranveer/common"
	"strconv"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

func registerServiceWithConsul() {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}
	ip, err := common.IPAddr()
	registration := new(consulapi.AgentServiceRegistration)

	registration.ID = "user-service"
	registration.Name = "user-service"
	//address := hostname()
	registration.Address = ip.String()
	p, err := strconv.Atoi(port()[1:len(port())])
	if err != nil {
		log.Fatalln(err)
	}
	registration.Port = p
	registration.Check = new(consulapi.AgentServiceCheck)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck", ip.String(), p)
	registration.Check.Interval = "5s"
	registration.Check.Timeout = "3s"
	consul.Agent().ServiceRegister(registration)
}

func lookupServiceWithConsul(serviceName string) (string, error) {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		return "", err
	}
	services, err := consul.Agent().Services()
	if err != nil {
		return "", err
	}

	//var consulServices = nil
	//consulServices = nil
	log.Println("##################")
	for k := range services {
		log.Println(strings.ToLower(k))

	}

	srvc := services["product-service"]
	address := srvc.Address
	port := srvc.Port
	return fmt.Sprintf("http://%s:%v", address, port), nil
}

func main() {
	registerServiceWithConsul()
	http.HandleFunc("/healthcheck", healthcheck)
	http.HandleFunc("/api/product/products", UserProduct)
	fmt.Printf("user service is up on port: %s", port())
	http.ListenAndServe(port(), nil)
}


func healthcheck(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, `product service is good`)
	r.Body.Close()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func UserProduct(w http.ResponseWriter, r *http.Request) {
	//p := []product{}
	var product []common.Product
	url, err := lookupServiceWithConsul("user-service")
	fmt.Println("URL: ", url)
	if err != nil {
		fmt.Fprintf(w, "Error. %s", err)
		return
	}
	client := &http.Client{}
	resp, err := client.Get(url + "/products")
	if err != nil {
		fmt.Fprintf(w, "Error. %s", err)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		fmt.Fprintf(w, "Error. %s", err)
		return
	}
	json.NewEncoder(w).Encode(&product)
}

func port() string {
	p := os.Getenv("USER_SERVICE_PORT")
	if len(strings.TrimSpace(p)) == 0 {
		return ":8999"
	}
	return fmt.Sprintf(":%s", p)
}
