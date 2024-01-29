package repository

import (
	"app/internal"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

func NewRepositoryProductMySql(db *sql.DB) *ProductMysql {
	return &ProductMysql{
		db: db,
	}
}

type ProductMysql struct {
	db *sql.DB
}

func (r *ProductMysql) FindById(id int) (p internal.Product, err error) {

	row := r.db.QueryRow("SELECT p.`id`, p.`name`, p.`quantity`, p.`code_value`, p.`is_published`, p.`expiration`, p.`price` from `products` `p` where p.`id` = ? ", id)

	err = row.Scan(&p.Id, &p.Name, &p.Quantity, &p.CodeValue, &p.IsPublished, &p.Expiration, &p.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrRepositoryProductNotFound
		}
		return
	}

	return
}

func (r *ProductMysql) Save(p *internal.Product) (err error) {

	var id int
	err = r.db.QueryRow("SELECT COUNT(*) FROM `products`").Scan(&id)
	if err != nil {
		return
	}

	id++

	_, err = r.db.Exec("INSERT INTO `products` (`id`, `name`, `quantity`, `code_value`, `is_published`, `expiration`, `price`) VALUES (?, ?, ?, ?, ?, ?, ?)", id, p.Name, p.Quantity, p.CodeValue, p.IsPublished, p.Expiration, p.Price)
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

	p.Id = id

	return
}

func (r *ProductMysql) UpdateOrSave(p *internal.Product) (err error) {
	res, err := r.db.Exec("UPDATE `products` SET `name` = ?, `quantity` = ?, `code_value` = ?, `is_published` = ?, `expiration` = ?, `price` = ? WHERE `id` = ?", p.Name, p.Quantity, p.CodeValue, p.IsPublished, p.Expiration, p.Price, p.Id)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			//esta de onda esto ahora ya que no hay indexes que chequeen duplicidad de campos en la DB
			switch mysqlErr.Number {
			case 1062:
				err = internal.ErrRepositoryProductDuplicated
			}
			return
		}
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return
	}

	if rowsAffected == 0 {
		err = r.Save(p)
		if err != nil {
			return
		}
	}

	return
}

func (r *ProductMysql) Update(p *internal.Product) (err error) {
	_, err = r.db.Exec("UPDATE `products` SET `name` = ?, `quantity` = ?, `code_value` = ?, `is_published` = ?, `expiration` = ?, `price` = ? WHERE `id` = ?", p.Name, p.Quantity, p.CodeValue, p.IsPublished, p.Expiration, p.Price, p.Id)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			//esta de onda esto ahora ya que no hay indexes que chequeen duplicidad de campos en la DB
			switch mysqlErr.Number {
			case 1062:
				err = internal.ErrRepositoryProductDuplicated
			}
			return
		}
	}

	return
}

func (r *ProductMysql) Delete(id int) (err error) {
	res, err := r.db.Exec("DELETE FROM `products` WHERE `id` = ?", id)
	if err != nil {
		return
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return
	}

	if rowsAffected == 0 {
		err = internal.ErrRepositoryProductNotFound
		return
	}

	return
}
