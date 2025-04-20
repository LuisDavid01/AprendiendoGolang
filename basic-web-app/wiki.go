package main

import (
	//"fmt"
	"errors"
	"regexp"
	"os"
	"net/http"
	"log"
	"html/template"
)
//La estructura de la pagina
type Page struct 
	{
		Title string
		Body []byte
	}

var templates = template.Must(template.ParseFiles("./tmpl/edit.html", "./tmpl/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
//Funcion que guarda una pagina localmente como texto
//se guarda en la carpeta data
func (p *Page) save() error{
	filename := "./data/" +p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error){
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid page title")

	}

	return m[2], nil
}

//Carga una pagina por nombre
func loadPage(title string) (*Page, error) {
	//carga la pagina de la carpeta data
	filename := "./data/"+title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil{
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
//Renderiza un template html con la informacion del body
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page){
	err := templates.ExecuteTemplate(w, tmpl + ".html" , p) 
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
}
//vista de la wiki
func viewHandler(w http.ResponseWriter, r *http.Request, title string){


	p, err := loadPage(title)

	if err != nil{
		http.Redirect(w, r, "/edit/"+ title, http.StatusFound)
		return
	} 
	renderTemplate(w, "view" ,p)
}

//vista de editar pagina
func editHandler(w http.ResponseWriter, r *http.Request, title string) {

    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
	renderTemplate(w, "edit", p)
}

//guardar pagina 
func saveHandler(w http.ResponseWriter, r *http.Request , title string) {

    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w,r)
			return
		}
		fn(w,r,m[2])
	}

}

func main(){
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	//se inicia el servidor en el puerto 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}

