package storage

import (
	"encoding/json"
	"log"

	"github.com/fatihdumanli/cnote/pkg/oauthv2"
	"github.com/fatihdumanli/cnote/pkg/onenote"
	"github.com/xujiajun/nutsdb"
)

const BUCKET = "cnote"

//expired, doesnt exist
//TODO: is it really the job of storage package to check if the token has expired?
func CheckToken() (oauthv2.OAuthToken, TokenStatus) {

	var token oauthv2.OAuthToken

	db, closer, err := openDb()
	defer closer()

	//TODO: perhaps a better handling is required.
	if err != nil {
		return token, DoesntExist
	}

	var e *nutsdb.Entry

	fnGet := func(tx *nutsdb.Tx) error {
		key := []byte(TOKEN_KEY)

		if e, err = tx.Get(BUCKET, key); err != nil {
			return err
		}

		return nil
	}
	err = db.View(fnGet)

	//TODO: perhaps a better handling is required.
	if err != nil {
		return token, DoesntExist
	}

	err = json.Unmarshal(e.Value, &token)

	//TODO: perhaps a better handling is required.
	if err != nil {
		log.Fatal(err)
		return token, DoesntExist
	}

	//check if it expired
	if token.IsExpired() {
		return token, Expired
	} else {
		return token, Valid
	}
}

func StoreToken(t interface{}) error {

	if _, ok := t.(oauthv2.OAuthToken); !ok {
		return InvalidTokenType
	}

	//open nuts db
	db, closer, err := openDb()
	defer closer()

	if err != nil {
		log.Fatal(err)
		return err
	}

	//convert the token into bytes
	bytes, err := json.Marshal(t)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fnUpdate := func(tx *nutsdb.Tx) error {
		key := []byte(TOKEN_KEY)
		val := bytes
		if err := tx.Put(BUCKET, key, val, 0); err != nil {
			log.Fatal(err)
			return err
		}
		return nil
	}

	//save the token
	err = db.Update(fnUpdate)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

//TODO: Complete
func SaveAlias(a onenote.AliasName, n onenote.NotebookName, s onenote.SectionName) error {

	db, closer, err := openDb()
	_ = db
	defer closer()

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

//NOTE: Currently is a mock
func GetAlias(a string) (onenote.Alias, bool) {
	notebookname := "Fatih adlı kişinin Not Defteri"
	sectionname := "Go"

	if a == "go" {
		return onenote.Alias{
			Notebook: onenote.NotebookName(notebookname),
			Section:  onenote.SectionName(sectionname),
		}, true
	} else {
		return onenote.Alias{}, false
	}
}

//opens the nuts db
//returns nuts db, closer and an error
//call closer to clean up the resources.
func openDb() (*nutsdb.DB, func() error, error) {
	opts := nutsdb.DefaultOptions
	opts.Dir = "/tmp/cnotedb"
	db, err := nutsdb.Open(opts)

	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	return db, db.Close, nil
}