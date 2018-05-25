package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

  "time"

"github.com/gorilla/mux"
	_ "github.com/denisenkom/go-mssqldb"
	//"github.com/gorilla/mux"
)

type Estructura struct {
	ProductID         int      `json:"ProductID,omitempty"`
	ProductName       string   `json:"ProductName,omitempty"`
	Stock             int      `json:"Stock,omitempty"`
	QuantitySold      int      `json:"QuantitySold,omitempty"`
	LastSoldDate      time.Time     `json:"LastSoldDate,omitempty"`
	BestCustomer      string   `json:"BestCustomer,omitempty"`

}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/inventory/stock", obtenerResultados).Methods("GET")


	log.Fatal(http.ListenAndServe(":5000", router))
}



func hola(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "No sirve")
}

func obtenerResultados(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mssql", "server=192.168.1.115\\devbc; database = AdventureWorks2014; user id=sa; password=1234")
	if err != nil {
		fmt.Println("From Open() attempt: " + err.Error())
	}
	slcstrContent, err := GetContent(db, 1)
	if err != nil {
		fmt.Println("(sqltest) Error en el contenido: " + err.Error())
	}

	json.NewEncoder(w).Encode(slcstrContent)
}

func GetContent(db *sql.DB, op int) ([]Estructura, error) {

	var slcstrContent []Estructura

	var discount string

	if op == 1 {
		discount = `
		SELECT top 10
		P.ProductID,
		P.Name as ProductName,
		SUM(PIV.Quantity) AS Stock,
		SUM(SOD.OrderQty) AS QuantitySold,
		MAX(SOH.OrderDate) AS LastSoldDate,
		(select TOP 1 pe.LastName + ' ' + pe.FirstName
		from Sales.SalesOrderDetail sd
		inner join Sales.SalesOrderHeader soh on sd.SalesOrderID =soh.SalesOrderID
		inner join Sales.Customer c on soh.CustomerID=c.CustomerID
		inner join Person.Person pe on c.PersonID=pe.BusinessEntityID
		where ProductID = P.ProductID
		GROUP BY pe.LastName + ' ' + pe.FirstName
		ORDER BY COUNT(1) DESC) as BestCustomer
		FROM Production.ProductInventory PIV
		INNER JOIN Production.Product P ON PIV.ProductID=P.ProductID
		INNER JOIN Sales.SalesOrderDetail SOD ON P.ProductID = SOD.ProductID
		INNER JOIN Sales.SalesOrderHeader SOH ON SOD.SalesOrderID = SOH.SalesOrderID
		GROUP By  P.ProductID, P.Name
		ORDER By  Stock asc,QuantitySold desc`
	}
	rows, err := db.Query(discount)

	if err != nil {
		return slcstrContent, err
	}
	defer rows.Close()

	for rows.Next() {

		var ProductID, Stock, QuantitySold  int
		var ProductName, BestCustomer string
		var LastSoldDate time.Time


		err := rows.Scan(&ProductID,&ProductName, &Stock, &QuantitySold,&LastSoldDate,&BestCustomer)

		if err != nil {
			return slcstrContent, err
		}

		var producto Estructura

		producto.ProductID = ProductID
		producto.ProductName = ProductName
		producto.Stock = Stock
		producto.QuantitySold = QuantitySold
		producto.LastSoldDate = LastSoldDate
		producto.BestCustomer = BestCustomer

		slcstrContent = append(slcstrContent, producto)

	}
	return slcstrContent, nil
}

func obtenerResultadosDos(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mssql", "server=192.168.1.115\\devbc; database = AdventureWorks2014; user id=sa; password=1234")
	if err != nil {
		fmt.Println("From Open() attempt: " + err.Error())
	}
	slcstrContent, err := GetContent(db, 2)
	if err != nil {
		fmt.Println("(sqltest) Error en el contenido: " + err.Error())
	}

	json.NewEncoder(w).Encode(slcstrContent)
}
