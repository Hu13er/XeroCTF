package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	sqlfile = "users.db"
)

var (
	db    *sql.DB
	mutex sync.Mutex

	users = map[string]struct {
		id        int
		plainPass string
		desc      string
	}{
		"admin0": {
			id:        0,
			plainPass: "Nimbus2000",
			desc:      "admin0 loves harry; hehe",
		},
		"admin0.0": {
			id:        1,
			plainPass: "nimbusTwoThousands",
			desc:      "admin0 loves harry; shampoo",
		},
		"admin1": {
			id:        2,
			plainPass: "fatFang123",
			desc:      "admin1 loves ron; hehe",
		},
		"admin2": {
			id:        3,
			plainPass: "fanghateexams",
			desc:      "nothing to say about admin2! shampoo",
		},
	}
)

func login(name, pass string) error {
	mutex.Lock()
	defer mutex.Unlock()

	hashed := md5.Sum([]byte(pass))
	hx := hex.EncodeToString(hashed[:])

	rows := db.QueryRow(`SELECT name FROM users WHERE name=? AND pass=?`, name, hx)
	return rows.Scan(new(string))
}

type nameDesc struct {
	name string
	desc string
}

func desc(name string) ([]nameDesc, error) {
	mutex.Lock()
	defer mutex.Unlock()

	// *SQLInjection*
	q := fmt.Sprintf(`SELECT name, desc FROM users WHERE name = "%s"`,
		name)
	log.Println(q)

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make([]nameDesc, 0)
	for rows.Next() {
		var nd nameDesc
		err = rows.Scan(&nd.name, &nd.desc)
		if err != nil {
			return nil, err
		}
		data = append(data, nd)
	}
	log.Println(data)
	return data, nil
}

func ensureTables() error {
	mutex.Lock()
	defer mutex.Unlock()
	// Drop All tables:
	tablesNames := make([]string, 0)
	if err := func() error {
		rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type is 'table'`)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var name string
			err = rows.Scan(&name)
			if err != nil {
				return err
			}
			tablesNames = append(tablesNames, name)
		}
		return nil
	}(); err != nil {
		return err
	}
	for _, tn := range tablesNames {
		//log.Println("Drop table", tn)
		_, err := db.Exec(fmt.Sprintf(`DROP TABLE %s;`, tn))
		if err != nil {
			log.Panicln(err)
		}
	}

	//log.Println(`Create Table "users"`)
	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
	id INTEGER NOT NULL PRIMARY KEY,
	name TEXT NOT NULL,
	pass TEXT NOT NULL,
	desc TEXT NOT NULL
	);`); err != nil {
		log.Fatalln("Could not ensure table useres:", err)
	}

	err := func() error {
		for name, pd := range users {
			//log.Println("Inserting user", name)
			hashed := md5.Sum([]byte(pd.plainPass))
			hx := hex.EncodeToString(hashed[:])
			_, err := db.Exec(`INSERT INTO users VALUES(?,?,?,?);`, pd.id, name, hx, pd.desc)
			if err != nil {
				return err
			}
		}
		return nil
	}()
	if err != nil {
		return err
	}

	return nil
}

func init() {
	var err error
	db, err = sql.Open("sqlite3", sqlfile)
	if err != nil {
		log.Fatalln("Could not open db:", err)
	}

	err = ensureTables()
	if err != nil {
		log.Fatalln(err)
	}

	t := time.NewTicker(5 * time.Second)
	go func() {
		for range t.C {
			// Reset tables:
			// cause every1 can fuck this db
			if err := ensureTables(); err != nil {
				// K8s should restart service
				log.Fatalln(err)
			}
		}
	}()
}
