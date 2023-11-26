package main

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Print("Entrez l'utilisateur de la base de données : ")
	var username string
	_, err := fmt.Scanf("%s", &username)
	if err != nil {
		log.Println("Erreur lors de la lecture du nom d'utilisateur :", err)
		return
	}

	fmt.Print("Entrez le mot de passe de la base de données : ")
	bytePassword, err := terminal.ReadPassword(0)
	if err != nil {
		log.Println("Erreur lors de la lecture du mot de passe :", err)
		return
	}
	password := string(bytePassword)

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/my_database", username, password))
	if err != nil {
		log.Panicln(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println("erreur")
		}
	}(db)

	fmt.Print("Dans quel dossier souhaitez-vous créer la base de données ? ")
	var folder string
	_, err = fmt.Scanf("%s", &folder)
	if err != nil {
		return
	}

	// Merci Antoine LEVEUGLE de m'avoir dit de faire du regex ^^
	regex := regexp.MustCompile("^[a-zA-Z0-9_\\-/.]+$")
	if !regex.MatchString(folder) {
		fmt.Println("Le chemin du dossier doit être composé uniquement de lettres, de chiffres, de tirets, de underscores, de points et de barres obliques.")
		return
	}

	fmt.Print("Quel nom souhaitez-vous donner à la base de données ? ")
	var database string
	_, err = fmt.Scanf("%s", &database)
	if err != nil {
		return
	}

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err := os.MkdirAll(folder, 0777)
		if err != nil {
			return
		}
	}

	_, err = db.Exec("CREATE DATABASE `" + database + "`")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("La base de données n'a pas été créée correctement.")
		return
	}

	fmt.Println("Base de données créée")
}
