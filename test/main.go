package main

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/crewjam/saml/samlsp"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	SSoEnabled bool `json:"SSoEnabled"`
}

func main() {
	//r := gin.New()

	keyPair, err := tls.LoadX509KeyPair("myservice.cert", "myservice.key")
	if err != nil {
		log.Fatal(err)
		return
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		log.Fatal(err)
		return
	}

	idpMetadataURL, err := url.Parse("https://samltest.id/saml/idp") // keycloak url or onelogin
	if err != nil {
		log.Fatal(err)
		return
	}
	idpMetadata, err := samlsp.FetchMetadata(context.Background(), http.DefaultClient, *idpMetadataURL)
	if err != nil {
		log.Fatal(err)
		return
	}

	rootURL, err := url.Parse("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
		return
	}

	samlSP, _ := samlsp.New(samlsp.Options{
		URL:         *rootURL,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		IDPMetadata: idpMetadata,
	})

	//r.Any("/saml/*action", gin.WrapH(samlSP))
	//r.Use(adapter.Wrap(samlSP.RequireAccount))
	//r.GET("/login", Login)
	//err = r.Run(":8080")
	//if err != nil {
	//	return
	//}
	app := http.HandlerFunc(hello)
	http.Handle("/hello", RequireAccount(app, samlSP))
	http.Handle("/saml/", samlSP)
	http.ListenAndServe(":8080", nil)

}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func Login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "hello")
}

func RequireAccount(handler http.Handler, m *samlsp.Middleware) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Apply conditional check here.
		session, err := m.Session.GetSession(r)
		if session != nil {
			r = r.WithContext(samlsp.ContextWithSession(r.Context(), session))
			handler.ServeHTTP(w, r)
			return
		}
		if err == samlsp.ErrNoSession {
			m.HandleStartAuthFlow(w, r)
			return
		}

		m.OnError(w, r, err)
		return
	})
}
