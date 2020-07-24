package auth

import (
	"bufio"
	"dbus-api/pkg/common"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Database struct {
	authorizedUsers map[string]string
}

func (db *Database) isAuthorized(userName string, plainPassword string) error {
	hash, userExists := db.authorizedUsers[userName]
	if !userExists {
		return errors.New("user does not exist")
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainPassword))
}

func (db *Database) BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			common.GenerateProblemResponse(http.StatusUnauthorized, "malformed auth")
			return
		}
		authErr := db.isAuthorized(username, password)
		if authErr != nil {
			log.Println("rejected user: " + username)
			common.GenerateProblemResponse(http.StatusUnauthorized, "user is unauthorized")
			return
		}
		log.Println("authorized user: " + username)
		next.ServeHTTP(w, r)
	})
}

//NewDatabase creates a new authorization database that allows users to use the api. The file must be an httpd formatted
//basic auth file in the form of <user>:<bcrypt password hash> with newlines separating each user.
func NewDatabase(pathToAuthFile string) (Database, error) {
	file, openErr := os.Open(pathToAuthFile)
	if openErr != nil {
		return Database{}, openErr
	}
	defer file.Close()
	db := Database{map[string]string{}}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		splitString := strings.Split(scanner.Text(), ":")
		if len(splitString) != 2 {
			return Database{}, errors.New("malformed httpd password file")
		}
		username := splitString[0]
		passwordHash := splitString[1]
		db.authorizedUsers[username] = passwordHash
	}
	return db, nil
}
