package main

// Import our dependencies. We'll use the standard http library as well as the gorilla router for this app
import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Here we are instantiating the gorilla/mux router
	r := mux.NewRouter()

	// On the default page we will simply serve our static index page.
	r.Handle("/", http.FileServer(http.Dir("./views/")))
	// We will setup our server so we can serve static assest like images, css from the /static/{file} route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Our API is going to consist of three routes
	// /status - which we will call to make sure that our API is up and running
	// /products - which will retrieve a list of products that the user can leave feedback on
	// /products/{slug}/feedback - which will capture user feedback on products
	r.Handle("/status", StatusHandler).Methods("GET")
	r.Handle("/products", ProductsHandler).Methods("GET")
	r.Handle("/products/{slug}/feedback", AddFeedbackHandler).Methods("POST")

	// Wrap the LoggingHandler function around our router so that the logger is called first on each route request
	http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))
}

// NotImplemented : Here we are implementing the NotImplemented handler. Whenever an API endpoint is hit
// we will simply return the message "Not Implemented"
var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})

// Product : We will first create a new type called Product
//  This type will contain information about boardgames */
type Product struct {
	ID          int
	Name        string
	Slug        string
	Description string
}

// We will create our catalog of boardgames and store them in a slice.
var products = []Product{
	Product{ID: 1, Name: "Cards Against Humanity", Slug: "cah", Description: "Cards Against Humanity is a party game for horrible people."},
	Product{ID: 2, Name: "Space Team", Slug: "space-team", Description: "A fast-paced, shouting card game where you work together as a team to repair a busted spaceship."},
	Product{ID: 3, Name: "Sonar", Slug: "sonar", Description: "You and your teammates control a state-of-the-art submarine and are trying to locate an enemy submarine in order to blow it out of the water before they can do the same to you."},
	Product{ID: 4, Name: "Codenames", Slug: "codenames", Description: "In Codenames, two teams compete to see who can make contact with all of their agents first."},
	Product{ID: 5, Name: "Dixit", Slug: "dixit", Description: "Every picture tells a story - but what story will your picture tell? Dixit is the lovingly illustrated game of creative guesswork, where your imagination unlocks the tale."},
	Product{ID: 6, Name: "Ticket To Ride", Slug: "ticket-to-ride", Description: "Ticket to Ride is a cross-country train adventure where players collect cards of various types of train cars that enable them to claim railway routes connecting cities in various countries around the world."},
}

// StatusHandler : The status handler will be invoked when the user calls the /status route
//  It will simply return a string with the message "API is up and running"
var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API is up and running"))
})

// ProductsHandler : The products handler will be called when the user makes a GET request to the /products endpoint.
//  This handler will return a list of products available for users to review
var ProductsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Here we are converting the slice of products to json
	payload, _ := json.Marshal(products)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})

// AddFeedbackHandler : The feedback handler will add either positive or negative feedback to the product
//  We would normally save this data to the database - but for this demo we'll fake it
//  so that as long as the request is successful and we can match a product to our catalog of products
//  we'll return an OK status.
var AddFeedbackHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var product Product
	vars := mux.Vars(r)
	slug := vars["slug"]

	for _, p := range products {
		if p.Slug == slug {
			product = p
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if product.Slug != "" {
		payload, _ := json.Marshal(product)
		w.Write([]byte(payload))
	} else {
		w.Write([]byte("Product Not Found"))
	}
})
