package main

// Import our dependencies. We'll use the standard http library as well as the gorilla router for this app
import (
  "net/http"
  "github.com/gorilla/mux"
)

func main() {
  // Here we are instantiating the gorilla/mux router
  r := mux.NewRouter()

  // On the default page we will simply serve our static index page.
  r.Handle("/", http.FileServer(http.Dir("./views/")))
  // We will setup our server so we can serve static assest like images, css from the /static/{file} route
  r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

  // Our application will run on port 3000. Here we declare the port and pass in our router.
  http.ListenAndServe(":3000", r)
}