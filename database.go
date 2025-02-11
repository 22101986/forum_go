package main

import (
	"database/sql"
	"fmt"
	"log"
)

func BDD() { // create database and create all table
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	DefineTables(db)
}

func DefineTables(db *sql.DB) { // define and create all tables
	usersTable := `CREATE TABLE IF NOT EXISTS users (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"uuid" VARCHAR(36) NOT NULL UNIQUE,
			"name" VARCHAR(255) NOT NULL UNIQUE,
			"email" VARCHAR(255) NOT NULL UNIQUE,
			"password" VARCHAR(72) NOT NULL,
			"picture" BLOB,
			"role_id" INTEGER,
			"is_deleted" BOOL,
			"is_external" INT,
			FOREIGN KEY(role_id) REFERENCES roles(id)
		);`
	postsTable := `CREATE TABLE IF NOT EXISTS posts (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,        
			"title" VARCHAR(255) NOT NULL,
			"content" TEXT NOT NULL,
			"date" VARCHAR(32) NOT NULL,
			"user_id" INTEGER NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);`
	commentsTable := `CREATE TABLE IF NOT EXISTS comments (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,        
			"content" TEXT NOT NULL,
			"date" VARCHAR(32) NOT NULL,
			"user_id" INTEGER, 
			"post_id" INTEGER NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
		);`
	categoriesTable := `CREATE TABLE IF NOT EXISTS categories (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,        
			"name" VARCHAR(32) NOT NULL UNIQUE
		);`
	likePostTable := `CREATE TABLE IF NOT EXISTS likepost (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,    
			"date" VARCHAR(32) NOT NULL,
			"user_id" INTEGER NOT NULL,
			"post_id" INTEGER NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
		);`
	dislikePostTable := `CREATE TABLE IF NOT EXISTS dislikepost (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			"date" VARCHAR(32) NOT NULL,
			"user_id" INTEGER NOT NULL,
			"post_id" INTEGER NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
		);`
	likeComTable := `CREATE TABLE IF NOT EXISTS likecom (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"date" VARCHAR(32) NOT NULL,        
			"user_id" INTEGER NOT NULL,
			"comments_id" INTEGER NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(comments_id) REFERENCES comments(id) ON DELETE CASCADE
		);`
	dislikeComTable := `CREATE TABLE IF NOT EXISTS dislikecom (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"date" VARCHAR(32) NOT NULL,        
			"user_id" INTEGER NOT NULL,
			"comments_id" INTEGER NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(comments_id) REFERENCES comments(id) ON DELETE CASCADE
		);`
	catPostRelTable := `CREATE TABLE IF NOT EXISTS catpostrel (
			"cat_id" INTEGER NOT NULL,
			"post_id" INTEGER NOT NULL,
			FOREIGN KEY(cat_id) REFERENCES categories(id),
			FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
		);`
	roleTable := `CREATE TABLE IF NOT EXISTS roles (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,        
			"name" VARCHAR(32) NOT NULL UNIQUE
		);`
	typeTable := `CREATE TABLE IF NOT EXISTS types (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,        
			"name" VARCHAR(32) NOT NULL UNIQUE
		);`
	blobTable := `CREATE TABLE IF NOT EXISTS blob (
			"picture" MEDIUMBLOB,
			"post_id" INTEGER NOT NULL UNIQUE,
			FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
		);`
	notifTable := `CREATE TABLE IF NOT EXISTS notifs (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"date" VARCHAR(32) NOT NULL,
			"type_id" INTEGER NOT NULL,
			"user_id_to" INTEGER NOT NULL,
			"user_id_from" INTEGER,
			"post_id" INTEGER,
			"comments_id" INTEGER,
			FOREIGN KEY(type_id) REFERENCES types(id),
			FOREIGN KEY(user_id_to) REFERENCES users(id),
			FOREIGN KEY(user_id_from) REFERENCES users(id),
			FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY(comments_id) REFERENCES comments(id) ON DELETE CASCADE
		);`

	// create all tables with above definitions
	CreateTable(db, usersTable, "users")
	CreateTable(db, postsTable, "posts")
	CreateTable(db, commentsTable, "comments")
	CreateTable(db, categoriesTable, "categories")
	CreateTable(db, likePostTable, "like_post")
	CreateTable(db, dislikePostTable, "dislike_post")
	CreateTable(db, likeComTable, "like_comments")
	CreateTable(db, dislikeComTable, "dislike_comments")
	CreateTable(db, catPostRelTable, "cat_post_rel")
	CreateTable(db, roleTable, "roles")
	CreateTable(db, typeTable, "types")
	CreateTable(db, blobTable, "blob")
	CreateTable(db, notifTable, "notif")
}

func CreateTable(db *sql.DB, createTableSQL string, tableName string) { //create one table already defined
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Creation table Failed %s : %v", tableName, err)
	}
	fmt.Printf("Table %s allready exist.\n", tableName)
}

func InsertNamesInDB(db *sql.DB, chosenNames []string, sqlExecQuery string) error { //fill references tables
	var err error
	for _, name := range chosenNames {
		_, err := db.Exec(sqlExecQuery, name)
		if err != nil {
			return err
		}
	}
	return err
}
