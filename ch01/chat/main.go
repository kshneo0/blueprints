package main

import (
	"blueprints/ch01/trace"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

// Avatar 구현을 활성화한다.
var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar,
}

//tmpl은 하나의 템플릿을 나타냄
type templateHandler struct {
	filename	string
	templ		*template.Template
}

// ServeHTTP가 HTTP 요청을 처리한다.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if t.templ == nil {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	}

	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)
}

var host = flag.String("host", ":8080", "The host of the application.")

func main(){

	flag.Parse() // parse the flags

	// setup gomniauth
	gomniauth.SetSecurityKey("AIzaSyBuBkFOY5kcPV5O0pPU1y9LIJKPGYeUUOk")
	gomniauth.WithProviders(
		google.New("735592233684-snegp4sbkk04s1l7b5kpjaokdd7kr3bi.apps.googleusercontent.com", "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", "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	http.Handle("/chat",  MustAuth(&templateHandler{filename:"chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room",r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, & http.Cookie{
			Name: "auth",
			Value:"",
			Path: "/",
			MaxAge:-1,
		})
		w.Header().Set("Location","/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/upload", &templateHandler{filename:"upload.html"})
	http.HandleFunc("/uploader",uploaderHandler)
	http.Handle("/avatars/",
		http.StripPrefix("/avatars/",
			http.FileServer(http.Dir("./avatars"))))

	//방을 가져옴
	go r.run()

	//웹서버 시작
	log.Println("Starting web server on", *host)
	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}