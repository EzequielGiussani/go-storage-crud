package repository_test

import (
	"app/internal"
	"app/internal/repository"
	"database/sql"
	"testing"

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

func TestWarehouse_FindById(t *testing.T) {

	t.Run("success - found by id", func(t *testing.T) {
		db, err := sql.Open("txdb", "test_db")
		require.NoError(t, err)
		defer db.Close()

		defer func(db *sql.DB) {

			_, err := db.Exec("DELETE FROM `warehouses`")

			require.NoError(t, err)

			_, err = db.Exec("ALTER TABLE `warehouses` auto_increment = 1")
			require.NoError(t, err)
		}(db)

		//set up
		func(db *sql.DB) {
			_, err := db.Exec("INSERT INTO `warehouses` (`id`, `name`, `adress`, `telephone`, `capacity`) VALUES (1, 'warehouse 1', 'address 1', 'telephone 1', 100)")
			require.NoError(t, err)
		}(db)

		rp := repository.NewRepositoryWarehouseMySql(db)

		//act
		wh, err := rp.FindById(1)

		//assert
		expectedWh := internal.Warehouse{
			Id: 1,
			WarehouseAttributes: internal.WarehouseAttributes{
				Name:      "warehouse 1",
				Address:   "address 1",
				Telephone: "telephone 1",
				Capacity:  100,
			},
		}
		require.NoError(t, err)
		require.Equal(t, expectedWh, wh)
	})

	t.Run("failure - warehouse not found", func(t *testing.T) {
		db, err := sql.Open("txdb", "test_db")
		require.NoError(t, err)
		defer db.Close()

		rp := repository.NewRepositoryWarehouseMySql(db)

		//act
		wh, err := rp.FindById(1)

		//assert
		expectedWh := internal.Warehouse{}
		expectedErr := internal.ErrRepositoryWarehouseNotFound
		require.Equal(t, expectedWh, wh)
		require.ErrorIs(t, err, expectedErr)
		require.EqualError(t, err, expectedErr.Error())
	})

}

func TestWarehouse_GetAll(t *testing.T) {

	t.Run("success - return 1", func(t *testing.T) {
		db, err := sql.Open("txdb", "test_db")
		require.NoError(t, err)
		defer db.Close()

		//set up
		func(db *sql.DB) {
			_, err := db.Exec("INSERT INTO `warehouses` (`id`, `name`, `adress`, `telephone`, `capacity`) VALUES (1, 'warehouse 1', 'address 1', 'telephone 1', 100)")
			require.NoError(t, err)
		}(db)

		rp := repository.NewRepositoryWarehouseMySql(db)

		//act
		wh, err := rp.GetAll()

		//assert
		expectedWh := []internal.Warehouse{
			{
				Id: 1,
				WarehouseAttributes: internal.WarehouseAttributes{
					Name:      "warehouse 1",
					Address:   "address 1",
					Telephone: "telephone 1",
					Capacity:  100,
				},
			},
		}
		require.NoError(t, err)
		require.Equal(t, expectedWh, wh)
	})

	t.Run("success - return 0", func(t *testing.T) {
		db, err := sql.Open("txdb", "test_db")
		require.NoError(t, err)
		defer db.Close()

		rp := repository.NewRepositoryWarehouseMySql(db)

		//act
		wh, err := rp.GetAll()

		//assert
		expectedWh := []internal.Warehouse(nil)
		require.NoError(t, err)
		require.Equal(t, expectedWh, wh)
	})

}

func TestWarehouse_Save(t *testing.T) {

	t.Run("success - saved", func(t *testing.T) {
		db, err := sql.Open("txdb", "test_db")
		require.NoError(t, err)
		defer db.Close()

		//set up
		wh := internal.Warehouse{
			Id: 1,
			WarehouseAttributes: internal.WarehouseAttributes{
				Name:      "warehouse 1",
				Address:   "address 1",
				Telephone: "telephone 1",
				Capacity:  100,
			},
		}

		rp := repository.NewRepositoryWarehouseMySql(db)

		//act
		err = rp.Save(&wh)

		//assert
		require.NoError(t, err)
	})

}
