package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/jinzhu/gorm/dialects/postgres" // using postgres sql
)

type message struct {
	StatusMessage int
	Message       interface{}
}

//post contains the struct to insert to DB
type Comments struct {
	Postid           uuid.UUID `db:"postid"`
	Comment          string    `json:"comment",db:"comment"`
	Commented_user   string    `json:"commented_user",db:"commented_user"`
	Comment_reaction string    `json:"commented_reaction",db:"comment_user"`
	Parent_path      string    `db:"parent_path"`
	Created_at       time.Time `json:"created_at",db:"created_at"`
}
type Post struct {
	Postid           uuid.UUID `db:"postid"`
	Postname         string    `json:"postname", db:"postname"`
	Createdby        string    `json:"createdby", db:"createdby"`
	Comment          string    `gorm:"-"`
	Comment_reaction string    `gorm:"-"`
	Parent_path      string    `gorm:"-"`
}

//Users contains a struct to read the username and password from the request body

type Users struct {
	Username string `json:"username", db:"username"`
	Password string `json:"password", db:"password"`
}

type M map[string]interface{}

//Register stores the user in the database
func Register(w http.ResponseWriter, r *http.Request) {
	creds := Users{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	if err != nil {
		msg := "cant generate password"
		http.Error(w, msg, http.StatusInternalServerError)
	}
	tx := db.Begin()
	creds.Password = string(hashedPassword)

	if err := tx.Create(&creds); err.Error != nil {
		log.Println(err.Error.Error())

		tx.Rollback()
		http.Error(w, "the user is taken", http.StatusInternalServerError)
		return
	}
	tx.Commit()
}

//AddComments adds first level comments on posts
func AddComments(w http.ResponseWriter, r *http.Request) {

	validated, user := validateToken(w, r)
	if !validated {
		http.Error(w, "could not validate the jwt", http.StatusBadRequest)
		return
	}
	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		log.Printf("could not decode the request body :%v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if post.Comment == "" {
		log.Println("no comment to add")
		http.Error(w, "no comment to add", http.StatusBadRequest)
		return
	}
	// here we are trying to get the unique id to get the comments
	var result Post
	rows := db.Table("posts").Select("*").Where("postname = ? and createdby = ?", post.Postname, post.Createdby).Row()

	err = rows.Scan(&result.Postid, &result.Postname, &result.Createdby)
	if err != nil {
		//set error code
		msg := "The row does not exists"
		http.Error(w, msg, http.StatusBadRequest)
		log.Printf("%v: %v", msg, err)

		return
	}
	var commentCount int
	row := db.DB().QueryRow("SELECT count(comment)from comments where parent_path  ~  $1", strings.ReplaceAll(result.Postid.String(), "-", "")+".*{1,1}")
	err = row.Scan(&commentCount)
	if err != nil {
		msg := "could not get the count"
		http.Error(w, msg, http.StatusBadRequest)
		log.Printf("%v:%v", msg, err)
		return
	}

	var comment Comments

	comment.Postid = result.Postid
	comment.Comment = post.Comment
	comment.Comment_reaction = post.Comment_reaction
	comment.Commented_user = user
	comment.Parent_path = strings.ReplaceAll(result.Postid.String(), "-", "") + "." + fmt.Sprint(commentCount+1)
	comment.Created_at = time.Now()
	log.Printf("%+v", comment)
	row = db.Table("comments").Select("*").Where("parent_path = ?", comment.Parent_path).Row()
	if row != nil {
		comment.Parent_path = strings.ReplaceAll(result.Postid.String(), "-", "") + "." + fmt.Sprint(commentCount+2)
	}
	tx := db.Begin()
	if err := tx.Create(&comment); err.Error != nil {
		tx.Rollback()
		log.Println(err)
		http.Error(w, err.Error.Error(), http.StatusBadRequest)
		return
	}
	tx.Commit()
}

//DeleteComments will only delete comments if the user is the
//owner of the posts or the user is owner of the comments
func DeleteComments(w http.ResponseWriter, r *http.Request) {
	validated, user := validateToken(w, r)
	if !validated {
		http.Error(w, "could not validate the jwt", http.StatusBadRequest)
		return
	}
	var posts Post
	err := json.NewDecoder(r.Body).Decode(&posts)
	if err != nil {
		log.Println("could not decode the request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result Comments
	var row *sql.Row
	row = db.Table("comments").Select("*").Where("parent_path=?", posts.Parent_path).Row()
	err = row.Scan(&result.Postid, &result.Comment, &result.Comment_reaction, &result.Commented_user, &result.Created_at, &result.Parent_path)
	if err != nil {
		//set error code
		msg := "The row does not exists in comments"
		http.Error(w, msg, http.StatusBadRequest)
		log.Printf("%v:%v", msg, err)
		return
	}

	commented_user := result.Commented_user
	row = db.Table("posts").Select("*").Where("postid=?", result.Postid).Row()
	err = row.Scan(&posts.Postid, &posts.Postname, &posts.Createdby)
	if err != nil {
		//set error code
		msg := "The row does not exists in posts"
		http.Error(w, msg, http.StatusBadRequest)
		log.Println("%v:%v", msg, err)
		return
	}
	created_user := posts.Createdby

	if !(user == created_user || user == commented_user) {
		http.Error(w, "user not authorized to delete the comment", http.StatusBadRequest)
		return
	}
	db.Exec("delete from comments where parent_path ~ ?", result.Parent_path)
	db.Exec("delete from comments where parent_path ~ ?", result.Parent_path+".*{1,10000}")
	response := fmt.Sprintf("deleted %v and all its sub comments if any existed on post %v", result.Comment, result.Postid)
	msg := message{StatusMessage: http.StatusOK, Message: response}
	json.NewEncoder(w).Encode(msg)

}

//Delete Post deletes the post only if the post is created by the same user
func DeletePosts(w http.ResponseWriter, r *http.Request) {
	validated, user := validateToken(w, r)
	if !validated {
		http.Error(w, "could not validate the jwt", http.StatusBadRequest)
		return
	}
	var posts Post
	type Data struct {
		Post_id string
	}
	var data Data
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Printf("could not decode the request body :%v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fromString, err := uuid.FromString(data.Post_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	posts.Postid = fromString

	row := db.Table("posts").Select("*").Where("postid=? ", posts.Postid).Row()
	err = row.Scan(&posts.Postid, &posts.Postname, &posts.Createdby)

	if err != nil || err == sql.ErrNoRows {
		//set error code
		msg := "current user not authorized to delete the post"
		http.Error(w, msg, http.StatusBadRequest)
		log.Printf("%v:%v", msg, err)
		return
	}
	if user != posts.Createdby {
		http.Error(w, "user not authorized to delete the post", http.StatusBadRequest)
		return
	}
	db.Exec("delete from posts where postid =?", fromString)
	response := fmt.Sprintf("delete post %v from database", posts.Postid)
	msg := message{StatusMessage: http.StatusOK, Message: response}
	json.NewEncoder(w).Encode(msg)
}

//GetPosts returns the postID for the user to trigger delete Request
func GetPosts(w http.ResponseWriter, r *http.Request) {
	validated, _ := validateToken(w, r)
	if !validated {
		http.Error(w, "could not validate the jwt", http.StatusBadRequest)
		return
	}
	var post Post

	rows, err := db.Table("posts").Select("*").Rows()
	if err != nil {
		log.Printf("Failed to get Rows %v", err)
		http.Error(w, "could not gets posts for particular user", http.StatusBadRequest)
	}
	var myMapSlice []M
	for rows.Next() {
		err := rows.Scan(&post.Postid, &post.Postname, &post.Createdby)
		if err != nil {
			//set error code
			msg := "could not fetch the value from db"
			http.Error(w, msg, http.StatusBadRequest)
			log.Printf("%v :%v", msg, err)
			return
		}
		m1 := M{"post_id": post.Postid, "post_name": post.Postname, "created_by": post.Createdby}
		myMapSlice = append(myMapSlice, m1)
	}
	msg := message{StatusMessage: http.StatusOK, Message: myMapSlice}
	json.NewEncoder(w).Encode(msg)
}

//Interact adds sub comments
func Interact(w http.ResponseWriter, r *http.Request) {
	validated, user := validateToken(w, r)
	if !validated {
		http.Error(w, "could not validate the jwt", http.StatusBadRequest)
		return
	}
	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		log.Printf("could not decode the request body :%v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if post.Comment == "" {
		log.Println("nothing to add")
		http.Error(w, "nothing to add", http.StatusBadRequest)
		return

	}
	//here we are forcing the user to add postname and createdby because if there were a UI
	// one would go and choose the same way
	var result Post
	rows := db.Table("posts").Select("*").Where("postname = ? and createdby = ?", post.Postname, post.Createdby).Row()

	err = rows.Scan(&result.Postid, &result.Postname, &result.Createdby)
	if err != nil {
		//set error code
		msg := "The row does not exists"
		http.Error(w, msg, http.StatusBadRequest)
		log.Printf("%v: %v", msg, err)

		return
	}
	var commentCount int
	row := db.DB().QueryRow("SELECT count(comment)from comments where parent_path  ~  $1", post.Parent_path+".*{1,1}")
	err = row.Scan(&commentCount)
	if err != nil {
		msg := "could not get the count"
		http.Error(w, msg, http.StatusBadRequest)
		log.Printf("%v:%v", msg, err)
		return
	}
	var comment Comments
	comment.Postid = result.Postid
	comment.Comment = post.Comment
	comment.Comment_reaction = post.Comment_reaction
	comment.Commented_user = user
	comment.Parent_path = post.Parent_path + "." + fmt.Sprint(commentCount+1)
	comment.Created_at = time.Now()
	log.Printf("%+v", comment)
	tx := db.Begin()
	if err := tx.Create(&comment); err.Error != nil {
		tx.Rollback()
		log.Println(err)
		return
	}
	tx.Commit()
	response := fmt.Sprintf("added reply %v to the parent comment", comment.Comment)
	msg := message{http.StatusOK, response}
	json.NewEncoder(w).Encode(msg)
}

//Add reaction adds reaction such as like, dislike into the DB
func AddReaction(w http.ResponseWriter, r *http.Request) {
	validated, _ := validateToken(w, r)
	if !validated {
		http.Error(w, "could not validate the jwt", http.StatusBadRequest)
		return
	}
	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		log.Printf("could not decode the request body :%v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var comment Comments
	comment.Comment_reaction = post.Comment_reaction
	//comment.Parent_path = post.Parent_path
	db.Model(&comment).Where("parent_path= ?", post.Parent_path).Updates(comment)
	response := fmt.Sprintf("added reaction %v on parent_path %v", comment.Comment_reaction, post.Parent_path)
	msg := message{StatusMessage: http.StatusOK, Message: response}
	json.NewEncoder(w).Encode(msg)
}

//GetInteraction returns top level comments and reply count
//and if user wants to reply he can interact with that comment
//with resulting post id
func GetInteraction(w http.ResponseWriter, r *http.Request) {
	validated, _ := validateToken(w, r)
	if !validated {
		http.Error(w, "could not validate the jwt", http.StatusBadRequest)
		return
	}
	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		log.Println("could not decode the request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result Post
	row := db.Table("posts").Select("*").Where("postname = ? and createdby = ?", post.Postname, post.Createdby).Row()

	err = row.Scan(&result.Postid, &result.Postname, &result.Createdby)
	if err != nil {
		//set error code
		msg := "The row does not exists"
		http.Error(w, msg, http.StatusBadRequest)
		log.Println(msg)
		return
	}

	var comment Comments
	var myMapSlice []M
	var rows *sql.Rows
	if post.Parent_path == "" {
		rows, err = db.Table("comments").Select("*").Where("parent_path  ~  $1  ", strings.ReplaceAll(result.Postid.String(), "-", "")+".*{1,1}").Rows()
	} else {
		rows, err = db.Table("comments").Select("*").Where("parent_path  ~  $1  ", post.Parent_path+".*{1,1}").Rows()
	}

	if err != nil {
		msg := "count not find the comments for given posts"
		http.Error(w, msg, http.StatusBadRequest)
		log.Printf("%v :%v", msg, err)
		return
	}
	var commentCount int
	for rows.Next() {
		err := rows.Scan(&comment.Postid, &comment.Comment, &comment.Comment_reaction, &comment.Commented_user, &comment.Created_at, &comment.Parent_path)
		if err != nil {
			//set error code
			msg := "could not fetch the value from db"
			http.Error(w, msg, http.StatusBadRequest)
			log.Printf("%v :%v", msg, err)
			return
		}
		row := db.DB().QueryRow("SELECT count(comment)from comments where parent_path  ~  $1", comment.Parent_path+".*{1,1}")
		err = row.Scan(&commentCount)
		if err != nil {
			msg := "could not get the count"
			http.Error(w, msg, http.StatusBadRequest)
			log.Printf("%v:%v", msg, err)
			return
		}
		m1 := M{"comment": comment.Comment, "commented_by": comment.Commented_user, "comment_reaction": comment.Comment_reaction, "parent_path": comment.Parent_path, "reply": commentCount}
		myMapSlice = append(myMapSlice, m1)
	}

	msg := message{StatusMessage: http.StatusOK, Message: myMapSlice}
	json.NewEncoder(w).Encode(msg)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	validated, username := validateToken(w, r)
	if !validated {
		http.Error(w, "could not validate the jwt", http.StatusBadRequest)
		return
	}
	uuid := uuid.NewV4()

	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		log.Println("could not decode the request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	post.Createdby = username
	post.Postid = uuid

	tx := db.Begin()
	if err := tx.Create(&post); err.Error != nil {
		tx.Rollback()
		log.Println(err)
		http.Error(w, "post already exists", http.StatusBadRequest)
		return
	}

	var comment Comments
	comment.Postid = uuid
	comment.Parent_path = strings.ReplaceAll(uuid.String(), "-", "")
	comment.Created_at = time.Now()

	if err := tx.Create(&comment); err.Error != nil {
		tx.Rollback()
		log.Println(err)
		return
	}
	tx.Commit()

	log.Printf("post %v created by %v added to database", post.Postname, post.Createdby)
	msg := message{StatusMessage: http.StatusOK, Message: "added post to database"}
	json.NewEncoder(w).Encode(msg)

}
