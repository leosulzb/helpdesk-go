package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	dsn := "lesbarros:Cookie@leo18@tcp(127.0.0.1:3306)/helpdesk"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	fmt.Println("Conexão com o banco de dados MySQL estabelecida com sucesso!")

}
