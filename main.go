package main

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
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
	r := gin.New()

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

	idpMetadataURL, err := url.Parse("https://samltest.id/saml/idp")
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

	r.Any("/saml/*action", gin.WrapH(samlSP))

	r.Use(SsoEnabledNew(r, samlSP))
	//	r.Use(CheckRequest())
	r.GET("/login", Login)
	//
	//ctx := gin.Context{}
	//
	//req, _ := ctx.Get("key")
	//fmt.Printf("req---%+v\n", req)
	//if req == true {
	//	fmt.Println("hello----")
	//} else {
	//r.GET("/hello", func(ctx *gin.Context) {
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"message": "hello world",
	//	})
	//})

	err = r.Run(":8080")
	if err != nil {
		return
	}
}

func Login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "hello")
}

func CheckRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req LoginRequest

		err := json.NewDecoder(ctx.Request.Body).Decode(&req)
		if err != nil {
			log.Fatal(err)
			return
		}
		if req.SSoEnabled == true {
			ctx.Set("key", true)
			ctx.Next()
		} else {
			ctx.Set("key", false)
			ctx.Next()
		}
	}
}

//func SsoEnabled2(r *gin.Engine, samlSP *samlsp.Middleware) gin.HandlerFunc {
//	return adapter.Wrap()
//}

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


func SsoEnabledNew(r *gin.Engine, samlSP *samlsp.Middleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		req.SSoEnabled = true

		// SSO enabled for user
		if req.SSoEnabled == true {
			//adapter.Wrap(samlSP.RequireAccount)
			session, err := samlSP.Session.GetSession(c.Request)
			//fmt.Printf("context request before---%+v\n", c.Request)
			//
			if session != nil {
			//	c.Request = c.Request.WithContext(samlsp.ContextWithSession(c.Request.Context(), session))
			//	if c.Request.URL.Path == samlSP.ServiceProvider.MetadataURL.Path {
			//		samlSP.ServeMetadata(c.Writer, c.Request)
			//		return
			//	}
			//	if c.Request.URL.Path == samlSP.ServiceProvider.AcsURL.Path {
			//		samlSP.ServeACS(c.Writer, c.Request)
			//		return
			//	}
			//	http.NotFoundHandler().ServeHTTP(c.Writer, c.Request)
			//	c.Next()
			//	return
			}
			if err == samlsp.ErrNoSession {
				samlSP.HandleStartAuthFlow(c.Writer, c.Request)
				c.Abort()
				return
			}
			samlSP.OnError(c.Writer, c.Request, err)
			c.Next()
			return
		} else {
			fmt.Println("sso is not enabled for user")
		}
	}
}

func Test(_ *gin.Context) (string, error) {
	e := fmt.Sprintf("code is running")
	return e, nil
}


func SsoEnabled(r *gin.Engine, samlSP *samlsp.Middleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("middleware ssoenables is called")
		var req LoginRequest
		req.SSoEnabled = true
		//err := json.NewDecoder(c.Request.Body).Decode(&req)
		//if err != nil {
		//	log.Fatal(err)
		//	return
		//}
		if req.SSoEnabled == true {

			//session, err := samlSP.Session.GetSession(r)
			//if session != nil {
			//	c.Request = c.Request.WithContext(samlsp.ContextWithSession(c.Request.Context(), session))
			//	r.ServeHTTP(w, r)
			//	return
			//}
			//if err == ErrNoSession {
			//	m.HandleStartAuthFlow(w, r)
			//	return
			//}
			//
			//m.OnError(w, r, err)
			//return
			//keyPair, err := tls.LoadX509KeyPair("myservice.cert", "myservice.key")
			//if err != nil {
			//	log.Fatal(err)
			//	return
			//}
			//keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
			//if err != nil {
			//	log.Fatal(err)
			//	return
			//}
			//
			//idpMetadataURL, err := url.Parse("https://samltest.id/saml/idp")
			//if err != nil {
			//	log.Fatal(err)
			//	return
			//}
			//idpMetadata, err := samlsp.FetchMetadata(c, http.DefaultClient, *idpMetadataURL)
			//if err != nil {
			//	log.Fatal(err)
			//	return
			//}
			//
			//rootURL, err := url.Parse("http://localhost:8080")
			//if err != nil {
			//	log.Fatal(err)
			//	return
			//}
			//
			//samlSP, _ := samlsp.New(samlsp.Options{
			//	URL:         *rootURL,
			//	Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
			//	Certificate: keyPair.Leaf,
			//	IDPMetadata: idpMetadata,
			//})

			session, err := samlSP.Session.GetSession(c.Request)

			fmt.Printf("session--%+v\n", session)
			fmt.Printf("context request before---%+v\n", c.Request)

			if session != nil {
				c.Request = c.Request.WithContext(samlsp.ContextWithSession(c.Request.Context(), session))
				r.ServeHTTP(c.Writer, c.Request)
				//return

				//fmt.Printf("request after---%+v\n", c.Request)
				//fmt.Printf("metadataurl--%+v\n", samlSP.ServiceProvider.MetadataURL.Path)
				//if c.Request.URL.Path == samlSP.ServiceProvider.MetadataURL.Path {
				//	samlSP.ServeMetadata(c.Writer, c.Request)
				//	return
				//}
				//
				//fmt.Printf("Acs url--%+v\n", samlSP.ServiceProvider.AcsURL.Path)
				//fmt.Printf("url--path--%+v\n", c.Request.URL.Path)
				//
				//if c.Request.URL.Path == samlSP.ServiceProvider.AcsURL.Path {
				//	samlSP.ServeACS(c.Writer, c.Request)
				//
				//	return
				//}
				//http.NotFoundHandler().ServeHTTP(c.Writer, c.Request)
				return
			}
			if err == samlsp.ErrNoSession {

				fmt.Printf("status----%+v\n", c.Writer.Status())
				samlSP.HandleStartAuthFlow(c.Writer, c.Request)

				return
			}

			samlSP.OnError(c.Writer, c.Request, err)

			return

		} else {
			test, err := Test(c)
			if err != nil {
				return
			}

			fmt.Printf("test--%+v\n", test)
			return
		}
	}
}
