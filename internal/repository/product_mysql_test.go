package repository_test

import (
	"app/internal"
	"app/internal/repository"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func init() {
	cfg := mysql.Config{
		User:      "user1",
		Passwd:    "secret_password",
		Addr:      "localhost:3306",
		Net:       "tcp",
		DBName:    "test_db",
		ParseTime: true,
	}

	txdb.Register("txdb", "mysql", cfg.FormatDSN())
}

func TestProduct_GetAll(t *testing.T) {

	t.Run("success - return 1", func(t *testing.T) {
		db, err := sql.Open("txdb", "test_db")
		require.NoError(t, err)
		defer db.Close()

		//set up
		func(db *sql.DB) {
			_, err := db.Exec("INSERT INTO `products` (`id`, `name`, `quantity`, `code_value`, `is_published`, `expiration`, `price`, `id_warehouse`) VALUES (1, 'product 1', 1, 'code_value 1', true, '2021-01-01', 1, 0)")
			require.NoError(t, err)
		}(db)

		rp := repository.NewRepositoryProductMySql(db)

		//act
		product, err := rp.GetAll()

		date, err := time.Parse("2006-01-02", "2021-01-01")

		require.NoError(t, err)

		//assert
		expectedProduct := []internal.Product{
			{
				Id:          1,
				WarehouseId: 0,
				ProductAttributes: internal.ProductAttributes{

					Name:        "product 1",
					Quantity:    1,
					CodeValue:   "code_value 1",
					IsPublished: true,
					Expiration:  date,
					Price:       1,
				},
			},
		}
		require.NoError(t, err)
		require.Equal(t, expectedProduct, product)
	})

	t.Run("success - return 0", func(t *testing.T) {
		db, err := sql.Open("txdb", "test_db")
		require.NoError(t, err)
		defer db.Close()

		rp := repository.NewRepositoryProductMySql(db)

		//act
		wh, err := rp.GetAll()

		//assert
		expectedWh := []internal.Product(nil)
		require.NoError(t, err)
		require.Equal(t, expectedWh, wh)
	})

}

func TestProduct_Save(t *testing.T) {

	t.Run("success - saved", func(t *testing.T) {
		db, err := sql.Open("txdb", "test_db")
		require.NoError(t, err)
		defer db.Close()

		timeNow := time.Now()

		//set up
		prod := internal.Product{
			Id:          1,
			WarehouseId: 0,
			ProductAttributes: internal.ProductAttributes{
				Name:        "product 1",
				Quantity:    1,
				CodeValue:   "code_value 1",
				IsPublished: true,
				Expiration:  timeNow,
				Price:       1,
			},
		}

		rp := repository.NewRepositoryProductMySql(db)

		//act
		err = rp.Save(&prod)

		//assert
		require.NoError(t, err)
	})

}

func TestProduct_Delete(t *testing.T) {
	t.Run("success - delete by id", func(t *testing.T) {
		db, err := sql.Open("txdb", "test_db")
		require.NoError(t, err)
		defer db.Close()

		// defer func(db *sql.DB) {

		// 	_, err := db.Exec("DELETE FROM `products`")

		// 	require.NoError(t, err)

		// 	_, err = db.Exec("ALTER TABLE `products` auto_increment = 1")
		// 	require.NoError(t, err)
		// }(db)

		//set up
		func(db *sql.DB) {
			_, err := db.Exec("INSERT INTO `products` (`id`, `name`, `quantity`, `code_value`, `is_published`, `expiration`, `price`, `id_warehouse`) VALUES (1, 'product 1', 1, 'code_value 1', true, '2021-01-01', 1, 0)")
			require.NoError(t, err)
		}(db)

		rp := repository.NewRepositoryProductMySql(db)

		//act
		err = rp.Delete(1)

		//assert
		require.NoError(t, err)
	})

	t.Run("fail - not found by id", func(t *testing.T) {
		db, err := sql.Open("txdb", "test_db")
		require.NoError(t, err)
		defer db.Close()

		// defer func(db *sql.DB) {

		// 	_, err := db.Exec("DELETE FROM `products`")

		// 	require.NoError(t, err)

		// 	_, err = db.Exec("ALTER TABLE `products` auto_increment = 1")
		// 	require.NoError(t, err)
		// }(db)

		//set up
		rp := repository.NewRepositoryProductMySql(db)

		//act
		err = rp.Delete(1)

		//assert
		require.Error(t, err)
		require.ErrorIs(t, err, internal.ErrRepositoryProductNotFound)
	})
}

func TestProduct_Update(t *testing.T) {

	t.Run("success - updated", func(t *testing.T) {
		db, err := sql.Open("txdb", "test_db")
		require.NoError(t, err)
		defer db.Close()

		timeNow := time.Now()

		//set up
		func(db *sql.DB) {
			_, err := db.Exec("INSERT INTO `products` (`id`, `name`, `quantity`, `code_value`, `is_published`, `expiration`, `price`, `id_warehouse`) VALUES (1, 'product 1', 1, 'code_value 1', true, '2021-01-01', 1, 0)")
			require.NoError(t, err)
		}(db)

		prod := internal.Product{
			Id:          1,
			WarehouseId: 0,
			ProductAttributes: internal.ProductAttributes{
				Name:        "updated test",
				Quantity:    1,
				CodeValue:   "code_value 1",
				IsPublished: true,
				Expiration:  timeNow,
				Price:       1,
			},
		}

		rp := repository.NewRepositoryProductMySql(db)

		//act
		err = rp.Update(&prod)

		//assert
		require.NoError(t, err)
	})

}
