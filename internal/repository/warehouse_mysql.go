package repository

import (
	"app/internal"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

func NewRepositoryWarehouseMySql(db *sql.DB) *Warehouse {
	return &Warehouse{
		db: db,
	}
}

type Warehouse struct {
	db *sql.DB
}

func (r *Warehouse) FindById(id int) (w internal.Warehouse, err error) {

	row := r.db.QueryRow("SELECT w.`id`, w.`name`, w.`adress`, w.`telephone`, w.`capacity` from `warehouses` `w` where w.`id` = ? ", id)

	err = row.Scan(&w.Id, &w.Name, &w.Address, &w.Telephone, &w.Capacity)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrRepositoryWarehouseNotFound
		}
		return
	}

	return
}

func (r *Warehouse) Save(w *internal.Warehouse) (err error) {

	res, err := r.db.Exec("INSERT INTO `warehouses` (`id`, `name`, `adress`, `telephone`, `capacity`) VALUES (?, ?, ?, ?, ?)", w.Id, w.Name, w.Address, w.Telephone, w.Capacity)
	if err != nil {
		var mySqlErr *mysql.MySQLError
		if errors.As(err, &mySqlErr) {
			switch mySqlErr.Number {
			//Pongo este case como ejemplo pero no va a tirarlo nunca ya que no hay constraint que chequee duplicidad utilizando un index
			case 1062:
				err = internal.ErrRepositoryProductDuplicated
			}
			return
		}
	}

	id, err := res.LastInsertId()

	if err != nil {
		return
	}

	w.Id = int(id)

	return
}

func (r *Warehouse) ReportProducts(id int) (w []internal.WarehouseProductsCount, err error) {
	var rows *sql.Rows

	if id == 0 {
		rows, err = r.db.Query("SELECT w.`id`, count(p.`id`) from `warehouses` `w` left join `products` `p` on w.`id` = p.`id_warehouse` group by w.`id`, w.`name`")
	} else {
		rows, err = r.db.Query("SELECT w.`id`, count(p.`id`) from `warehouses` `w` left join `products` `p` on w.`id` = p.`id_warehouse` where w.`id` = ? group by w.`id`, w.`name`", id)
	}

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var wh internal.WarehouseProductsCount
		err = rows.Scan(&wh.Name, &wh.Count)
		if err != nil {
			return
		}

		w = append(w, wh)
	}

	return
}

func (r *Warehouse) GetAll() (w []internal.Warehouse, err error) {
	rows, err := r.db.Query("SELECT w.`id`, w.`name`, w.`adress`, w.`telephone`, w.`capacity` from `warehouses` `w`")
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var wh internal.Warehouse
		err = rows.Scan(&wh.Id, &wh.Name, &wh.Address, &wh.Telephone, &wh.Capacity)
		if err != nil {
			return
		}

		w = append(w, wh)
	}

	err = rows.Err()
	if err != nil {
		return
	}

	return
}
