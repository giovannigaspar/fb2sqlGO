// If you don't need some of the codes, you can just comment it. For example, if
// you don't need triggers, just comment the code block.


package main
import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/nakagami/firebirdsql"
)


// Tables' name (the same from the other database)
// Edit this list with your tables' name.
var tables = []string{
    "table_1",
    "table_2"}


func main() {
	// Firebird -> Source
	connOrigem := "user:password@ip:port/dbname"
	// PostgreSQL -> destination
	connDestino := "user=user password=pass dbname=dbname"

	fmt.Printf("Connecting to source database . . . ")
	dbOrigem, err := sql.Open("firebirdsql", connOrigem)
	defer dbOrigem.Close()
	if err != nil {
		fmt.Printf(err.Error() + "\n")
		return
	}
	fmt.Printf("Connected\n")

	fmt.Printf("Connecting to destination database . . . ")
	dbDestino, err := sql.Open("postgres", connDestino)
	defer dbDestino.Close()
	if err != nil {
		fmt.Printf(err.Error() + "\n")
		return
	}
	fmt.Printf("Connected\n")

	fmt.Printf("Creating tables . . . ")
	err = executeSQLFile(dbDestino, "sql/tables.sql")
	if err != nil {
		fmt.Printf(err.Error() + "\n")
		return
	}
	fmt.Printf("Done\n")

	for i := range tables {
		err = copyTables(dbOrigem, dbDestino, tables[i])
		if err != nil { fmt.Printf(err.Error() + "\n") }
	}

	fmt.Printf("Creating triggers . . . ")
	err = executeTriggersFile(dbDestino, "sql/triggers.sql")
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	fmt.Printf("Done\n")

	fmt.Printf("Making additional configurations . . . ")
	err = executeSQLFile(dbDestino, "sql/configs.sql")
	if err != nil {
		fmt.Printf(err.Error() + "\n")
		return
	}
	fmt.Printf("Done\n")

	fmt.Println("Process terminated")
}

func executeTriggersFile(db *sql.DB, filename string) (error){
	file, err := ioutil.ReadFile(filename)
	if err != nil { return err }

	requests := strings.Split(string(file), "/*SPLITHERE*/")
	tx, err := db.Begin()
	if err != nil { return err }
	for _, request := range requests {
		_, err := tx.Exec(request)
		if err != nil { return err }
	}

	tx.Commit()
	return nil
}

func executeSQLFile(db *sql.DB, filename string) (error){
	file, err := ioutil.ReadFile(filename)
	if err != nil { return err }

	requests := strings.Split(string(file), ";")
	tx, err := db.Begin()
	if err != nil { return err }
	for _, request := range requests {
		_, err := tx.Exec(request)
		if err != nil { return err }
	}

	tx.Commit()
	return nil
}

func copyTables(dbOrigem *sql.DB, dbDestino *sql.DB, nomeTabela string) (error){
	fmt.Printf("Copying table %s . . . ", nomeTabela)

	txOrigem, err := dbOrigem.Begin()
	if err != nil { return err }
	rowsOrigem, err := txOrigem.Query("SELECT * FROM " + nomeTabela)
	if err != nil { return err }
	colsOrigem, err := rowsOrigem.Columns()
	if err != nil { return err }

	txDestino, err := dbDestino.Begin()
	if err != nil { return err }

	for rowsOrigem.Next() {
		columnsOrigem := make([]interface{}, len(colsOrigem))
		columnOrigemPointers := make([]interface{}, len(colsOrigem))
		for i := range columnsOrigem {
			columnOrigemPointers[i] = &columnsOrigem[i]
		}

		if err := rowsOrigem.Scan(columnOrigemPointers...); err != nil { return err }
		m := make(map[string]interface{})

		var keys []string
		var values []string
		for i, colName := range colsOrigem {
			val := columnOrigemPointers[i].(*interface{})
			m[colName] = *val

			strVal := fmt.Sprint(*val)
			if strVal != `<nil>` {
				strVal = parseData(strVal)
				keys = append(keys, colName)
				values = append(values, strVal)
			}
		}

		insertQuery := createInsertSQLQuery(nomeTabela, keys, values)
		_, err = txDestino.Exec(insertQuery)
		if err != nil { return err }
	}

	txOrigem.Commit()
	txDestino.Commit()
	fmt.Printf("Done\n")
	return nil
}

func createInsertSQLQuery(tableName string, keys []string, values []string) (string) {
	var sqlQuery bytes.Buffer

	sqlQuery.WriteString(`INSERT INTO `)
	sqlQuery.WriteString(tableName)
	sqlQuery.WriteString(`(`)
	sqlQuery.WriteString(strings.Join(keys, ","))
	sqlQuery.WriteString(`)`)
	sqlQuery.WriteString(` values('`)
	sqlQuery.WriteString(strings.Join(values, "','"))
	sqlQuery.WriteString(`')`)

	return sqlQuery.String()
}

// If you need any additional data parse, put it here.
func parseData(s string) (string) {
    /* Examples:
	if strings.Contains(s, `+0000 UTC`) {
		s = s[:len(s)-10]
	}
	if strings.Contains(s, `0000-01-01`) {
		s = s[11:len(s)]
	}
	if strings.Contains(s, `'`) {
		s = strings.Replace(s, `'`, `"`, -1)
	}*/
	return s
}
