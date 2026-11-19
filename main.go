package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Nom struct {
	Name string `json:"name"`
}
type Type struct {
	Type Nom `json:"type"`
}
type Pokemon struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Types     []Type  `json:"types"`
	Height    float32 `json:"height"`
	Weight    float32 `json:"weight"`
	Abilities []struct {
		Ability Nom `json:"ability"`
	} `json:"abilities"`
}
type ListePokémon struct {
	Result []Pokemon `json:"results"`
}

type Traduction struct {
	Name     string `json:"name"`
	Language Nom    `json:"language"`
}
type Espece struct {
	Names []Traduction `json:"names"`
}

type PagePokemon struct {
	Pokemon Pokemon
	NomFr   string
}

func accueil(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("erreur de chargement du template : %v", err)
		http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
		return
	}

	pokémons := getPokemons()

	if err := tmpl.Execute(w, pokémons); err != nil {
		log.Printf("erreur de rendu du template : %v", err)
	}
}

func getPokemons() []Pokemon {
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

	for i := range liste.Result {
		liste.Result[i].ID = i + 1
	}
	return liste.Result
}

func getPokemon(id string) Pokemon {
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + id)
	if err != nil {
		fmt.Println("Erreur :", err)
		return Pokemon{}
	}
	defer resp.Body.Close()
	var pokemon Pokemon
	json.NewDecoder(resp.Body).Decode(&pokemon)
	return pokemon
}

func getEspece(id string) Espece {
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon-species/" + id)
	if err != nil {
		fmt.Println("Erreur :", err)
		return Espece{}
	}
	defer resp.Body.Close()
	var espece Espece
	json.NewDecoder(resp.Body).Decode(&espece)
	return espece
}

func nomFrancais(espece Espece) string {
	for _, traduction := range espece.Names {
		if traduction.Language.Name == "fr" {
			return traduction.Name
		}
	}
	return ""
}

func fiche(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	page := PagePokemon{
		Pokemon: getPokemon(id),
		NomFr:   nomFrancais(getEspece(id)),
	}
	tmpl := template.Must(template.ParseFiles("templates/fiche.html"))
	tmpl.Execute(w, page)
}

func main() {
	http.HandleFunc("/", accueil)
	http.HandleFunc("/pokemon/{id}", fiche)
	http.Handle("/static/", http.FileServer(http.Dir(".")))

	log.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
