package prop

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type dbProp struct {
	DbHost           string
	DbPort           string
	DbUser           string
	DbPassword       string
	DbName           string
	DbSslMode        string
	DbSelectLimit    int
	DbMaxConns       int32
	DbWaitAfterQuery time.Duration
}

const (
	EQUAL     string = " = "
	NOTEQUAL  string = " <> "
	GREATER   string = " > "
	LESS      string = " < "
	GREATERE  string = " >= "
	LESSE     string = " <= "
	LIKE      string = " LIKE "
	NOTLIKE   string = " NOT LIKE "
	IN        string = " IN "
	NOTIN     string = " NOT IN "
	ISNULL    string = " IS NULL "
	ISNOTNULL string = " IS NOT NULL "
)

var (
	dbPropInstance *dbProp
	dbPropOnce     sync.Once
)

func GetDbProp() *dbProp {
	dbPropOnce.Do(func() {
		dbPropInstance = &dbProp{}
		setDbProp(dbPropInstance)
	})
	return dbPropInstance
}

func setDbProp(prop *dbProp) {
	// DbHost
	if v, ok := os.LookupEnv("DB_HOST"); ok {
		prop.DbHost = v
	} else {
		log.Fatal("DB_HOST not found")
	}

	// DbPort
	if v, ok := os.LookupEnv("DB_PORT"); ok {
		prop.DbPort = v
	} else {
		log.Fatal("DB_PORT not found")
	}

	// DbUser
	if v, ok := os.LookupEnv("DB_USER"); ok {
		prop.DbUser = v
	} else {
		log.Fatal("DB_USER not found")
	}

	// DbPassword
	if v, ok := os.LookupEnv("DB_PASSWORD"); ok {
		prop.DbPassword = v
	} else {
		log.Fatal("DB_PASSWORD not found")
	}

	// DbName
	if v, ok := os.LookupEnv("DB_NAME"); ok {
		prop.DbName = v
	} else {
		log.Fatal("DB_NAME not found")
	}

	// DbSqlSslMode
	if v, ok := os.LookupEnv("DB_SSL_MODE"); ok {
		if v == "disable" || v == "require" || v == "verify-ca" || v == "verify-full" {
			prop.DbSslMode = v
		} else {
			log.Fatalf("DB_SSL_MODE invalid value: %s", v)
		}
		prop.DbSslMode = v
	} else {
		prop.DbSslMode = "require"
	}

	// DbSelectLimit
	if v, ok := os.LookupEnv("DB_SELECT_LIMIT"); ok {
		if v, err := strconv.Atoi(v); err == nil {
			prop.DbSelectLimit = v
		} else {
			log.Fatalf("DB_SELECT_LIMIT invalid value: %s", v)
		}
	} else {
		prop.DbSelectLimit = 5000
	}

	// DbMaxConns
	if v, ok := os.LookupEnv("DB_MAX_CONNS"); ok {
		if v, err := strconv.Atoi(v); err == nil {
			prop.DbMaxConns = int32(v)
		} else {
			log.Fatalf("DB_MAX_CONNS invalid value: %s", v)
		}
	} else {
		prop.DbMaxConns = 30
	}

	// DbWaitAfterQuery
	if v, ok := os.LookupEnv("DB_WAIT_AFTER_QUERY"); ok {
		if v, err := strconv.Atoi(v); err == nil {
			prop.DbWaitAfterQuery = time.Millisecond * time.Duration(v)
		} else {
			log.Fatalf("DB_WAIT_AFTER_QUERY invalid value: %s", v)
		}
	} else {
		prop.DbWaitAfterQuery = time.Millisecond * 10
	}
}
