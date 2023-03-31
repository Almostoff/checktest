package db

import (
	"database/sql"
	"fmt"
	"log"
	"wbL0/server_sub/config"
	"wbL0/server_sub/entity"
)

type DBService struct {
	db *sql.DB
}

func NewDB(database *sql.DB) *DBService {
	return &DBService{db: database}
}

func InitDBConn(cfg config.Config) (*DBService, error) {
	dbConn := DBService{}
	var err error
	connStr := "user=" + cfg.Database.Username + " password=" + cfg.Database.Password + " dbname=" + cfg.Database.DBname + " sslmode=disable"
	dbConn.db, err = sql.Open(cfg.Database.DriverName, connStr)
	if err != nil {
		return &DBService{}, err
	}
	return &dbConn, err
}

func (dbService *DBService) Close() error {
	err := dbService.db.Close()
	return err
}

func (dbService *DBService) SaveOrder(jsonData *entity.DataItem) (sql.Result, error) {
	//fmt.Println("CHECK FROM db.go jsonData: ", jsonData)
	result, err := dbService.db.Exec(`insert into orders (id, orderdata) values ($1, $2)`, jsonData.ID, jsonData.OrderData)
	//fmt.Println("qqqq")
	if err == nil {
		log.Println("New data in database stored: ", jsonData)
	}
	//fmt.Println("CHECK result after exec db: ", result)
	return result, err
}

func (dbService *DBService) GetAllOrders() ([]entity.DataItem, error) {
	rows, err := dbService.db.Query("select * from orders")
	fmt.Println("CHECK rows from db.go: ", rows)
	rowItem := entity.DataItem{}
	rows.Scan(&rowItem.ID, &rowItem.OrderData)
	defer rows.Close()
	strs := []entity.DataItem{}
	for rows.Next() {
		str := entity.DataItem{}
		err := rows.Scan(&str.ID, &str.OrderData)
		if err != nil {
			return strs, err
		}
		strs = append(strs, str)
	}
	return strs, err
}

func (dbService *DBService) GetOrderByID(id string) (*entity.DataItem, error) {
	row := dbService.db.QueryRow("select * from orders where id=$1", id)
	rowData := new(entity.DataItem)
	err := row.Scan(&rowData.ID, &rowData.OrderData)
	return rowData, err
}
