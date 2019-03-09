package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"routter"
	"strings"
)

type Data struct {
	Name       string
	TypeOfFile string `json:"Type"`
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func host(w http.ResponseWriter, r *http.Request, k routter.Keys) {
	log.Println("host--------------------------------")
	http.ServeFile(w, r, "data/"+k["user"]+"/"+k["path"])
}

func user(w http.ResponseWriter, r *http.Request, k routter.Keys) {
	log.Println("user--------------------------------")
	cookie := http.Cookie{Name: "username", Value: k["user"], MaxAge: 111111111, Path: "/"}
	http.SetCookie(w, &cookie)
}

func static(w http.ResponseWriter, r *http.Request, k routter.Keys) {
	log.Println("static--------------------------------")
	http.ServeFile(w, r, "static/"+k["path"])
}

func dashboard(w http.ResponseWriter, r *http.Request, _ routter.Keys) {
	log.Println("dashboard--------------------------------")
	log.Println(r.RemoteAddr, "\n", r.Header.Get("X-FORWARDED-FOR"))
	http.ServeFile(w, r, "static/dashboard.html")
}

func get(w http.ResponseWriter, r *http.Request, _ routter.Keys) {
	log.Println("get-------------------------------- ", r.URL.Path)
	urlquery, _ := url.ParseQuery(r.URL.RawQuery)
	fmt.Println(urlquery)
	fmt.Println(len(urlquery["path"]))

	cookie, err := r.Cookie("username")
	if err != nil {
		fmt.Fprint(w, err)
	}

	path := cookie.Value
	log.Println(path)
	if len(urlquery["path"]) != 0 {
		path += "/" + urlquery["path"][0]
		log.Println(path)
	}
	dir, err := ioutil.ReadDir("data/" + path)
	if err != nil {
		log.Println(err)
	}
	s := make([]Data, 0)
	var ext []string
	if path == cookie.Value {
		for i := range dir {
			s = append(s, Data{dir[i].Name(), "md-cloud"})
		}
	} else {
		for i := range dir {
			if dir[i].IsDir() {
				s = append(s, Data{dir[i].Name(), "md-folder"})
			}
		}
		for i := range dir {
			if !dir[i].IsDir() {
				ext = strings.Split(dir[i].Name(), ".")
				if len(ext) <= 1 {
					s = append(s, Data{dir[i].Name(), "md-document"})
				} else {
					switch ext[len(ext)-1] {
					case "htm", "html":
						s = append(s, Data{dir[i].Name(), "logo-html5"})
					case "js", "jsx":
						s = append(s, Data{dir[i].Name(), "logo-javascript"})
					case "ts", "coffee", "vue", "php", "go", "sh":
						s = append(s, Data{dir[i].Name(), "md-code"})
					case "css", "sass":
						s = append(s, Data{dir[i].Name(), "md-color-palette"})
					case "ttf", "otf", "woff":
						s = append(s, Data{dir[i].Name(), "md-open"})
					case "zip", "tar", "rar", "tar.xz", "tar.bz", "tar.bz2", "7z", "gz", "gzip", "jar", "tgz", "tar-gz", "xz":
						s = append(s, Data{dir[i].Name(), "md-filing"})
					default:
						if strings.Contains(mime.TypeByExtension("."+ext[len(ext)-1]), "image") {
							s = append(s, Data{dir[i].Name(), "md-image"})
						} else if strings.Contains(mime.TypeByExtension("."+ext[len(ext)-1]), "video") {
							s = append(s, Data{dir[i].Name(), "md-videocam"})
						} else if strings.Contains(mime.TypeByExtension("."+ext[len(ext)-1]), "text") {
							s = append(s, Data{dir[i].Name(), "md-list-box"})
						} else {
							s = append(s, Data{dir[i].Name(), "md-document"})
						}
					}

				}
			}
		}
	}
	fmt.Println(s)
	sendData, err := json.Marshal(&s)
	fmt.Println(string(sendData))
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(sendData)
	if err != nil {
		log.Println(err)
	}
}

func mainPage(w http.ResponseWriter, r *http.Request, _ routter.Keys) {
	log.Println("main--------------------------------")
	http.ServeFile(w, r, "static/index.html")
}

func login(w http.ResponseWriter, r *http.Request, _ routter.Keys) {
	log.Println("login--------------------------------")
	http.ServeFile(w, r, "static/login.html")
}

func del(w http.ResponseWriter, r *http.Request, _ routter.Keys) {
	log.Println(r.Method)
	var name struct {
		Files []string `json:"FileNames"`
	}
	err := json.NewDecoder(r.Body).Decode(&name)
	log.Println("delete.....", err, name.Files)
	for _, e := range name.Files {
		cookie, err := r.Cookie("username")
		if err != nil {
			fmt.Fprint(w, err)
		}

		e = "data/" + cookie.Value + "/" + e
		err = os.Remove(e)
		if err != nil {
			log.Println(err)
			_ = os.RemoveAll(e + "/")
		}
	}
}

func upload(w http.ResponseWriter, r *http.Request, _ routter.Keys) {
	file, handler, err := r.FormFile("Files")
	check(err)
	defer file.Close()
	cookie, err := r.Cookie("username")
		if err != nil {
			fmt.Fprint(w, err)
		}
	dir := "data/" + cookie.Value + "/" + r.FormValue("Directory")
	log.Println(dir+handler.Filename)
	f, err := os.Create(dir+handler.Filename)
	check(err)
	defer f.Close()
	_, err = io.Copy(f, file)
	check(err)
}

func createProject(w http.ResponseWriter, r *http.Request, _ routter.Keys) {
	cookie, err := r.Cookie("username")
		if err != nil {
			fmt.Fprint(w, err)
	}
	name := r.FormValue("Name")
	name = strings.Replace(name, "/", "_slash_", -1)
	dir := "data/" + cookie.Value + "/" + name
	log.Println(">   ", dir)
	err = os.Mkdir(dir, 644)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	r := routter.New()
	r.Host = HOST
	r.Add(r.Host+"/login", login)
	r.Add(r.Host+"/get", get)
	r.Add(r.Host+"/upload", upload)
	r.Add(r.Host+"/delete", del)
	r.Add(r.Host+"/create", createProject)
	r.Add(r.Host+"/user/$user", user)
	r.Add(r.Host+"/dashboard", dashboard)
	r.Add("$user."+r.Host+"/$path:", host)
	r.Add(r.Host+"/static/$path:", static)
	r.Add(r.Host, mainPage)
	// TODO: проверить роутер
	r.Run(":8080")
}
