package main

import (
	"bufio"
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type User struct {
	Id       int64
	Created  int64
	Updated  int64
	Username string
	Password []byte
}

func main() {

	//prep the input buffer
	in := bufio.NewReader(os.Stdin)

	//read in username and strip newline character from end
	fmt.Println("Enter Username:")
	username, _ := in.ReadString('\n')
	cleanUsername := username[:len(username)-1]

	//read in Password and strip newline
	fmt.Println("Enter Password:")
	password, _ := in.ReadString('\n')
	cleanPassword := password[:len(password)-1]
	//hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(cleanPassword), bcrypt.DefaultCost)

	/* BCrypt Example
	//hash the line and output the hash
	hashedLine, _ := bcrypt.GenerateFromPassword([]byte(cleanLine), bcrypt.DefaultCost)
	fmt.Println("Hashed Line:", hashedLine)

	//check the hash cost, just for comparison
	hashCost, _ := bcrypt.Cost(hashedLine)
	fmt.Println("Hash Cost:", hashCost)

	//verify hash against the original string
	if bcrypt.CompareHashAndPassword(hashedLine, []byte(cleanLine)) != nil {
		fmt.Println("Nope, it changed somehwere!?")
	} else {
		fmt.Println("Yep, still the same line.")
	}
	*/

	//db open
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Database Opened.")
	}

	//setup gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbmap.AddTableWithName(User{}, "user_test").SetKeys(true, "Id")

	//select user ID from username
	userId, _ := dbmap.SelectInt("select Id from user_test where Username=?", cleanUsername)

	//return user object from Id
	obj, _ := dbmap.Get(User{}, userId)
	user := obj.(*User)

	//verify stored hash against entered password
	if bcrypt.CompareHashAndPassword(user.Password, []byte(cleanPassword)) != nil {
		fmt.Println("Wrong Password")
	} else {
		fmt.Println("Right Password")
		fmt.Printf("Your UserId is %d.\n", userId)
	}
	defer db.Close()
}
