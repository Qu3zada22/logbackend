package main

import (
    "database/sql"
    "log"
    _ "github.com/mattn/go-sqlite3"
)


// nuestros modelos, estos representan las tablas de la base de datos. Los pueden poner en un archivo models.go

type UserModel struct {
    Username string `json:"username"`
    Password string `json:"password"` // Needed from the client request
}

// setupDatabase inicializa la conexión a la BD
func setupDatabase(dbPath string) (*sql.DB, error) {
    log.Printf("Conectando a la base de datos en: %s", dbPath)
    // Nota: Con modernc.org/sqlite necesitamos el prefijo "file:" y el nombre del driver es "sqlite"
    db, err := sql.Open("sqlite", "file:"+dbPath+"?_foreign_keys=on")
    if err != nil {
        return nil, err
    }

    // Es buena idea hacer ping para verificar la conexión inmediatamente
    if err = db.Ping(); err != nil {
        db.Close() // Cerrar si el ping falla
        return nil, err
    }
    log.Println("Base de datos conectada exitosamente.")
    // Podríamos añadir aquí la creación de tablas si no existen
    // _, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (...)`)

    return db, nil
}