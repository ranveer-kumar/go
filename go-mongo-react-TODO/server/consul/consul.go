package consul

import (
	"fmt"
	"log"
	"ranveer/common"

	// "strconv"

	// "net/http"

	consulapi "github.com/hashicorp/consul/api"
)

func RegisterServiceWithConsul() {
	log.Println("inside consule registration..")
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}
	ip, err := common.IPAddr()
	registration := new(consulapi.AgentServiceRegistration)

	registration.ID = "go-mongo-crud"
	registration.Name = "go-mongo-crud"
	//address := hostname()
	registration.Address = ip.String()
	// p, err := strconv.Atoi(port()[1:len(port())])
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// registration.Port = p
	log.Println("ip=", ip.String())
	registration.Port = 8080
	registration.Check = new(consulapi.AgentServiceCheck)
	// registration.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck", ip.String(), p)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck", ip.String(), 8080)
	registration.Check.Interval = "5s"
	registration.Check.Timeout = "3s"
	consul.Agent().ServiceRegister(registration)
	log.Println("service registered..")
}

// func healthcheck(w http.ResponseWriter, r *http.Request) {
// 	//fmt.Fprintf(w, `product service is good`)
// 	r.Body.Close()
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("OK"))
// }
