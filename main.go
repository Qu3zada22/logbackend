package main

import (
    "log"
    "net/http"

    "database/sql"
    "encoding/json"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"

	"myapp/handlers"
    "myapp/utils"

    "myapp/models"
    

)

func main() {
    // Conectar a la base de datos
    db, err := setupDatabase("../db/users.db")
    if err != nil {
        log.Fatal("CRITICAL: No se pudo conectar a la base de datos:", err)
    }
    defer db.Close() // Asegurar que se cierre al final

    // Crear router Chi
    r := chi.NewRouter()

    // Middlewares
    r.Use(middleware.Logger)    // Loggea cada request
    r.Use(middleware.Recoverer) // Recupera de panics
    r.Use(configureCORS())      // Aplica nuestra configuración CORS

    // Rutas Públicas (sin autenticación requerida inicialmente)
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("API de Login v1.0"))
    })
    r.Post("/register", handlers.PostRegisterHandler(db))
    r.Post("/login", handlers.PostLoginHandler(db))

    // Ruta para obtener datos de usuario (protegida más tarde con JWT)
    // Por ahora, cualquiera puede acceder si conoce el ID
    //r.Get("/users/{userID}", handlers.GetUserHandler(db))
    // --- Rutas Protegidas ---
    r.Group(func(r chi.Router) {
        r.Use(utils.JwtAuthMiddleware(db))

        r.Post("/auth/logout", handlers.PostLogoutHandler(db))
        r.Get("/users/profile", getUserProfileHandler(db))
    })

    // Ruta pública opcional
    r.Get("/users/{userID}", getUserProfileHandler(db)) 

    port := ":3000"
    log.Printf("Servidor escuchando en puerto %s", port)
    log.Fatal(http.ListenAndServe(port, r))
}

// --- Nuevo Handler para Perfil ---
func getUserProfileHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID, ok := r.Context().Value("userID").(int)
        if !ok || userID == 0 {
            http.Error(w, `{"error": "No se pudo obtener ID de usuario del token"}`, http.StatusInternalServerError)
            return
        }

        var userResp models.LoginSuccessData
        err := db.QueryRow("SELECT id, username FROM users WHERE id = ?", userID).Scan(&userResp.UserID, &userResp.Username)
        if err != nil {
            if err == sql.ErrNoRows {
                http.Error(w, `{"error": "Usuario del token no encontrado"}`, http.StatusNotFound)
            } else {
                log.Printf("Error consultando perfil para user %d: %v", userID, err)
                http.Error(w, `{"error": "Error interno del servidor"}`, http.StatusInternalServerError)
            }
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(userResp)
    }
}