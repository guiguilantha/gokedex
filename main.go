package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Pokemon struct {
	Name string `json:"name"`
}

type ListePokémon struct {
	Result []Pokemon `json:"results"`
}

func accueil(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("erreur de chargement du template : %v", err)
		http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
		return
	}

	pokémons := getPokemon()
	fmt.Println(len(pokémons), "Pokémon")

	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("erreur de rendu du template : %v", err)
	}
}
func getPokemon() []Pokemon {
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon?limit=151")
	if err != nil {
		fmt.Println("erreur :", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("erreur API :", resp.Status)
		return nil
	}

	var liste ListePokémon
	if err := json.NewDecoder(resp.Body).Decode(&liste); err != nil {
		fmt.Println("erreur de décodage :", err)
		return nil
	}

	return liste.Result
}
func main() {
	http.HandleFunc("/", accueil)

	log.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
