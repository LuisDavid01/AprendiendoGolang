package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//se inicializa el cuerpo del almub
type album struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}
//Por defecto creamos 3 albumes
var albums = []album{
    {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
    {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
    {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}
//Muestra todos los albunes
func getAlbums(c *gin.Context){
	c.IndentedJSON(http.StatusOK, albums)
}

//crea un album nuevo
func createAlbum(c *gin.Context){
	var newAlbum  album

	//BindJSON carga el body del request con el struct album
	if err := c.BindJSON(&newAlbum); err != nil{
		return
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumById(c *gin.Context){
	id := c.Param("id")

	for _, a := range albums{
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}

	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"mensaje": "Album no encontrado"})

}
//encuentra el index del album en memoria
func findIndex(id string) int {
	for i, album := range albums {
		if album.ID == id{
			return i
		}
	}	
		return -1
}
//borrar un album
func deleteAlbumById(c *gin.Context){
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"mensaje": "Coloque un id porfavor"})
		return
	}
	idx := findIndex(id)

	if idx == -1{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"mensaje": "album no encontrado"})
		return
	}
	deletedAlbum := albums[idx]
	albums = append(albums[:idx], albums[idx + 1:]... )

	c.IndentedJSON(http.StatusGone, gin.H{"album eliminado": deletedAlbum })



}
//editar un album
func editAlbum(c *gin.Context){
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"mensaje": "Coloque un id porfavor"})
		return
	}
	idx := findIndex(id)

	if idx == -1{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"mensaje": "album no encontrado"})
		return
	}
	//se carga el album del body
	var newAlbum  album

	//BindJSON carga el body del request con el struct album
	if err := c.BindJSON(&newAlbum); err != nil{
		return
	}
	//se actualiza el album	
	albums[idx] = newAlbum
	c.IndentedJSON(http.StatusOK, gin.H{"Album-editado": newAlbum})

	




}



func main(){
	//se inicializa el enrutador del api
	router := gin.Default()

	//rutas
	router.GET("/albums", getAlbums)
	router.POST("/albums", createAlbum)
	router.GET("/albums/:id", getAlbumById)
	router.DELETE("/albums/:id", deleteAlbumById)
	router.PUT("/albumns/:id", editAlbum)
	//el servidor esta corriendo
	router.Run("localhost:8080")
}
