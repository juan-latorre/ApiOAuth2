package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

var (
	oauthConfig = &oauth2.Config{
		ClientID:     "TU_ID_DE_CLIENTE",
		ClientSecret: "TU_SECRETO_DE_CLIENTE",
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"perfil", "correo"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:8080/auth",
			TokenURL: "http://localhost:8080/token",
		},
	}
)

func main() {
	http.HandleFunc("/auth", handleAuth)
	http.HandleFunc("/auth/callback", handleCallback)
	http.HandleFunc("/token", handleToken)
	http.HandleFunc("/protected", handleProtected)

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, oauthConfig.AuthCodeURL("state"), http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Error al intercambiar token", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Token de Acceso: %s\n", token.AccessToken)
	fmt.Fprintf(w, "Token de Actualización: %s\n", token.RefreshToken)
	fmt.Fprintf(w, "Vencimiento: %s\n", token.Expiry.Format(time.RFC3339))
}

func handleToken(w http.ResponseWriter, r *http.Request) {
	// Implementar la lógica del endpoint de tokens para generar tokens de acceso
	// Por propósitos de demostración, puedes simplemente retornar un token codificado
	accessToken := "TU_TOKEN_DE_ACCESO"

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"access_token": "%s", "token_type": "bearer"}`, accessToken)
}

func handleProtected(w http.ResponseWriter, r *http.Request) {
	// Implementar la lógica para proteger recursos utilizando tokens de acceso
	// Por propósitos de demostración, puedes simplemente verificar si la solicitud contiene un token de acceso válido
	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		http.Error(w, "No Autorizado", http.StatusUnauthorized)
		return
	}

	// Verificar el token de acceso (por ejemplo, validar contra una tienda de tokens o decodificar JWT)
	// Por propósitos de demostración, simplemente comprobaremos si coincide con el token codificado
	accessTokenEsperado := "TU_TOKEN_DE_ACCESO"
	if accessToken != "Bearer "+accessTokenEsperado {
		http.Error(w, "No Autorizado", http.StatusUnauthorized)
		return
	}

	// Acceso concedido
	fmt.Fprintln(w, "¡Recurso Protegido Accedido!")
}
