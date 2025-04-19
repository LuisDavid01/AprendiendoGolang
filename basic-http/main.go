package main 
import (
    "fmt"
    "net/http"
	"github.com/gorilla/mux"
)

func main(){
	r := mux.NewRouter()


	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        title := vars["title"]
        page := vars["page"]
        fmt.Fprintf(w, "Here is your book: %s, page: %s", title, page)
    })	/*
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprint(w, "welcome to my website")
	})
	fs := http.FileServer(http.Dir("static/"))
	
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	*/
	http.ListenAndServe(":80", r)

}
