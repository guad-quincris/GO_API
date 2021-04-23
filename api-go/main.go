package main

import (
	"encoding/json" //El paquete json implementa la codificación y decodificación de JSON
	"fmt"           //El paquete fmt implementa E / S formateado con funciones análogas a printf y scanf de C
	"io/ioutil"
	"log"      // Define un tipo, Logger, con métodos para formatear la salida
	"net/http" //El paquete http proporciona implementaciones de servidor y cliente HTTP.
	"strconv"  //El paquete strconv implementa conversiones hacia y desde representaciones de cadenas de tipos de datos básicos

	"github.com/gorilla/mux" //"multiplexor de solicitud HTTP"
)

//EStablecemos nuestro struct que en otros lenguajes se conoce como clase, modelo o interface
type ticket struct {
	Id         int    `json:"id"`
	User       string `json:"user"`
	StartDate  string `json:"startDate"`
	UpdateDate string `json:"updateDate"`
	Status     string `json:"status"`
}

type allTickets []ticket

//Base de datos de Prueba o data dumie

var tickets = allTickets{
	{
		Id:         111,
		User:       "user1",
		StartDate:  "inicio",
		UpdateDate: "actualizacion",
		Status:     "abierto",
	},
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Wecome the my GO API!")
}

//Funcion para listar todos los tickets

func getTickets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Tipe", "application/json")
	json.NewEncoder(w).Encode(tickets)
}

// funcion para listar o llamar un ticket por su ID
func getTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ticketId, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Not valid Id")
		return
	}

	for _, ticket := range tickets {
		if ticket.Id == ticketId {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ticket)
		}
	}

}

// funcion para crear un nuevo ticket

func createTickets(w http.ResponseWriter, r *http.Request) {
	var newTicket ticket
	reqBody, err := ioutil.ReadAll(r.Body) //Recibo la informacion que el cliente va a estar enviando al servidor

	// Condicional para detectar errores
	if err != nil {
		fmt.Fprintf(w, "Insert a Valid Data")
	}

	json.Unmarshal(reqBody, &newTicket)  // asigno la informacion que estoy recibiendo a la variable newTicket
	newTicket.Id = len(tickets) + 111    // creamos un ID para el nuevo ticket
	tickets = append(tickets, newTicket) // agrego newTicket a la lista ya existente de tickets

	w.Header().Set("Content-Tipe", "application/json") // expecificamos el tipo de contenido que vamos a enviar en la peticion
	w.WriteHeader(http.StatusOK)                       //envio de un codigo de estado
	json.NewEncoder(w).Encode(newTicket)               //Respondemos al cliente con el ticket que se acaba de crear
}

//funcion para eliminar un ticked deacuerdo a in ID

func deleteTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ticketId, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Not valid Id")
		return
	}

	for index, ticket := range tickets {
		if ticket.Id == ticketId {
			tickets = append(tickets[:index], tickets[index+1:]...)
			fmt.Fprintf(w, "El Ticket con el ID %v Se ha eliminado correctamente", ticketId)
		}
	}

}

// funcion para actualizar la informacion de un ticket seleccionado por ID

func updateTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ticketId, err := strconv.Atoi(vars["id"])
	var updatedTicket ticket

	if err != nil {
		fmt.Fprintf(w, "Not valid Id")
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a Valid Data")
	}
	json.Unmarshal(reqBody, &updatedTicket)

	for index, ticket := range tickets {
		if ticket.Id == ticketId {
			tickets = append(tickets[:index], tickets[index+1:]...)
			updatedTicket.Id = ticketId
			tickets = append(tickets, updatedTicket)
			fmt.Fprintf(w, "El Ticket con el ID %v Se ha Actualizado correctamente", ticketId)
		}
	}

}

func main() {

	//Establecemos el enrutador  que funcionara en el puerto: 3000

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tickets", getTickets).Methods("GET")
	router.HandleFunc("/tickets", createTickets).Methods("POST")
	router.HandleFunc("/tickets/{id}", getTicket).Methods("GET")
	router.HandleFunc("/tickets/{id}", deleteTicket).Methods("DELETE")
	router.HandleFunc("/tickets/{id}", updateTicket).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000", router))

}
