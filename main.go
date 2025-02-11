package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"forum/methods"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type StructSquare struct { // isolated methods struct
	Users    *methods.UserMethod
	Posts    *methods.PostMethod
	Comments *methods.CommentMethod
	Likes    *methods.LikeMethod
	Notifs   *methods.NotifMethod
	Blobs    *methods.BlobMethod
}
type StructData struct { // data struct
	AllUsers  []*methods.User
	AllPosts  []*methods.Post
	AllCats   []*methods.Categories
	AllNotifs []*methods.Notif
}

func main() {
	err := loadEnvFile("local.env") // Open local.env in order to extract data
	if err != nil {
		log.Fatalf("Loading error on file local.env: %v", err)
	}
	fmt.Println("File local.env correctly opened")

	BDD() // create database and tables

	db, err := sql.Open("sqlite3", DBPath) // open database for nexts functions
	if err != nil {
		log.Fatalf("Loading error on the database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`PRAGMA foreign_keys = ON;`)
	if err != nil {
		log.Fatalf("Error on Pragmas: %v", err)
	}

	sSquared := CreateStructSquare(db) // create struct for methods
	sData := CreateStructData()        // create struct for all datas

	InsertNamesInDB(db, []string{"Tips", "Ponds", "Fishing spots", "Catches", "Boats", "Crustaceans", "Mollusks", "Fishes"}, `INSERT INTO categories (name) VALUES (?)`)
	InsertNamesInDB(db, []string{"Drowned", "Classic", "Moderator", "Administrator"}, `INSERT INTO roles (name) VALUES (?)`)
	InsertNamesInDB(db, []string{"likepost", "dislikepost", "likecom", "dislikecom", "comonpost", "askmod", "reportpost", "reportcom", "adminanswer+", "adminanswer-", "admindelete+", "admindelete-"}, `INSERT INTO types (name) VALUES (?)`)

	FillingStruct(nil, nil, &sSquared, &sData) // fill the data struct with database dataP

	window := NewWindow(CountMaxReq, WindowSize) // Creating sliding window for rate limiting
	fmt.Println("Sliding Window created")

	ServerCreate(sSquared, sData, CleanTLS(), window) // build TLS structure, and then launch server
}

func CleanTLS() *tls.Config {
	cert, err := tls.LoadX509KeyPair(CertPath, KeyPath)
	if err != nil {
		log.Fatalf("Loading error in pairing certificates: %v", err)
	}
	certConfig := []tls.Certificate{cert}

	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS13,
		Certificates:             certConfig,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
		},
	}
	return tlsConfig
}

// server function with almost all handlers
func ServerCreate(sSquared StructSquare, sData StructData, tlsConfig *tls.Config, window *Window) {
	// No error/logout handler directly, see handlersBasic.go/handlersLog.go for this
	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, IndexHandler, w, r, &sSquared, &sData)
	}
	registerHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, RegisterHandler, w, r, &sSquared, &sData)
	}
	loginHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, LoginHandler, w, r, &sSquared, &sData)
	}
	logoutHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, LogoutHandler, w, r, &sSquared, &sData)
	}
	// googleCallbackHandler := func(w http.ResponseWriter, r *http.Request) {
	// 	VerifReq(window, GoogleCallbackHandler, w, r, &sSquared, &sData)
	// }
	// googleLoginHandler := func(w http.ResponseWriter, r *http.Request) {
	// 	VerifReq(window, GoogleLoginHandler, w, r, &sSquared, &sData)
	// }
	// discordCallbackHandler := func(w http.ResponseWriter, r *http.Request) {
	// 	VerifReq(window, DiscordCallbackHandler, w, r, &sSquared, &sData)
	// }
	// discordLoginHandler := func(w http.ResponseWriter, r *http.Request) {
	// 	VerifReq(window, DiscordLoginHandler, w, r, &sSquared, &sData)
	// }
	// githubCallbackHandler := func(w http.ResponseWriter, r *http.Request) {
	// 	VerifReq(window, GithubCallbackHandler, w, r, &sSquared, &sData)
	// }
	// githubLoginHandler := func(w http.ResponseWriter, r *http.Request) {
	// 	VerifReq(window, GithubLoginHandler, w, r, &sSquared, &sData)
	// }
	newPostHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, NewPostHandler, w, r, &sSquared, &sData)
	}
	detailPostHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, DetailPostHandler, w, r, &sSquared, &sData)
	}
	deletePostHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, DeletePostHandler, w, r, &sSquared, &sData)
	}
	editPostHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, EditPostHandler, w, r, &sSquared, &sData)
	}
	deleteCommentHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, DeleteCommentHandler, w, r, &sSquared, &sData)
	}
	editCommentHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, EditCommentHandler, w, r, &sSquared, &sData)
	}
	myProfileHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, ProfileHandler, w, r, &sSquared, &sData)
	}
	myActivityHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, ActivityProfileHandler, w, r, &sSquared, &sData)
	}
	myPostsHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, PostProfileHandler, w, r, &sSquared, &sData)
	}
	myLikesHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, LikeProfileHandler, w, r, &sSquared, &sData)
	}
	deleteProfileHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, DeleteProfileHandler, w, r, &sSquared, &sData)
	}
	likePostHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, LikePostHandler, w, r, &sSquared, &sData)
	}
	dislikePostHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, DislikePostHandler, w, r, &sSquared, &sData)
	}
	likeComHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, LikeComHandler, w, r, &sSquared, &sData)
	}
	dislikeComHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, DislikeComHandler, w, r, &sSquared, &sData)
	}
	notifHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, NotifHandler, w, r, &sSquared, &sData)
	}
	imagePostHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, ImagePostHandler, w, r, &sSquared, &sData)
	}
	profilePictureHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, ProfilePictureHandler, w, r, &sSquared, &sData)
	}
	modoHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, ModoHandler, w, r, &sSquared, &sData)
	}
	responseHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, ResponseHandler, w, r, &sSquared, &sData)
	}
	reportHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, ReportHandler, w, r, &sSquared, &sData)
	}
	adminHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, AdminHandler, w, r, &sSquared, &sData)
	}
	demoteHandler := func(w http.ResponseWriter, r *http.Request) {
		VerifReq(window, DemoteHandler, w, r, &sSquared, &sData)
	}

	mux := http.NewServeMux() // Mux for multiple handlers
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.Handle("/pictures/", http.StripPrefix("/pictures/", http.FileServer(http.Dir("./pictures"))))
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/login", loginHandler)
	// mux.HandleFunc("/login/google", googleLoginHandler)
	// mux.HandleFunc("/callback/google", googleCallbackHandler)
	// mux.HandleFunc("/login/github", githubLoginHandler)
	// mux.HandleFunc("/callback/github", githubCallbackHandler)
	// mux.HandleFunc("/login/discord", discordLoginHandler)
	// mux.HandleFunc("/callback/discord", discordCallbackHandler)
	mux.HandleFunc("/logout", logoutHandler)
	mux.HandleFunc("/newPost", newPostHandler)
	mux.HandleFunc("/detailPost", detailPostHandler)
	mux.HandleFunc("/deletePost", deletePostHandler)
	mux.HandleFunc("/editPost", editPostHandler)
	mux.HandleFunc("/deleteComment", deleteCommentHandler)
	mux.HandleFunc("/editComment", editCommentHandler)
	mux.HandleFunc("/myProfile", myProfileHandler)
	mux.HandleFunc("/myActivity", myActivityHandler)
	mux.HandleFunc("/myPosts", myPostsHandler)
	mux.HandleFunc("/myLikes", myLikesHandler)
	mux.HandleFunc("/deleteProfile", deleteProfileHandler)
	mux.HandleFunc("/likePost", likePostHandler)
	mux.HandleFunc("/dislikePost", dislikePostHandler)
	mux.HandleFunc("/likeCom", likeComHandler)
	mux.HandleFunc("/dislikeCom", dislikeComHandler)
	mux.HandleFunc("/notifications", notifHandler)
	mux.HandleFunc("/imagePost", imagePostHandler)
	mux.HandleFunc("/profilePicture", profilePictureHandler)
	mux.HandleFunc("/askModo", modoHandler)
	mux.HandleFunc("/adminAnswer", responseHandler)
	mux.HandleFunc("/report", reportHandler)
	mux.HandleFunc("/admin", adminHandler)
	mux.HandleFunc("/demote", demoteHandler)

	server := &http.Server{
		Addr:              ServerPort,       //adresse du server
		Handler:           mux,              // listes des handlers
		ReadHeaderTimeout: 10 * time.Second, // temps autorisé pour lire les headers
		WriteTimeout:      10 * time.Second, // temps maximum d'écriture de la réponse
		IdleTimeout:       30 * time.Second, // temps maximum entre deux rêquetes
		MaxHeaderBytes:    1 << 20,          // 1 MB // maximum de bytes que le serveur va lire
		TLSConfig:         tlsConfig,
	}

	log.Println("https://localhost:8080")
	if err := server.ListenAndServeTLS(CertPath, KeyPath); err != nil { // open server with auto-signed SSL certificate
		log.Fatalf("Listen&Serve/TLS error: %v", err)
	}
}
