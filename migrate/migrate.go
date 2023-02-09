package migrate

import (
	"context"

	"github.com/JongGeonClass/JGC-API/dbmodel"
	"github.com/thak1411/gorn"
	"github.com/thak1411/rnlog"
)

// Drop All FK & INDEX
// Migrate All Table For GORN
// If you Migrate From Gorn Before, You Must Not Use This Function
func MigrateAllReset(db *gorn.DB) error {
	_, tableNames := dbmodel.GetTables()
	err := db.ExecTx(context.Background(), func(txdb *gorn.DB) error {
		// Delete All Foreign Key
		rnlog.Info("Delete All Foreign Key...")
		for _, tableName := range tableNames {
			fkeys, err := txdb.GetForeignKeys(tableName)
			if err != nil {
				rnlog.Error("Get Foreign Key Error: %s", err.Error())
				return err
			}
			for _, fkey := range fkeys {
				if err := txdb.DropForeignKey(tableName, fkey.ConstraintName); err != nil {
					rnlog.Error("Drop Foreign Key Error: %s", err.Error())
					return err
				}
			}
		}
		// Delete All Index
		rnlog.Info("Delete All Index...")
		indexes, err := txdb.GetIndexes()
		if err != nil {
			rnlog.Error("Get Index Error: %s", err.Error())
			return err
		}
		for _, index := range indexes {
			if err := txdb.DropIndex(index); err != nil {
				rnlog.Error("Drop Index Error: %s", err.Error())
				return err
			}
		}
		return nil
	})
	if err != nil {
		rnlog.Error("Migrate Error: %+v\n", err)
		return err
	}
	rnlog.Info("Migrate Success")
	return nil
}

func Migrate(db *gorn.DB) {
	var err error
	// err = MigrateAllReset(db)
	// if err != nil {
	// 	rnlog.Error("Migrate All Reset Error: %+v\n", err)
	// 	return
	// }

	tables, tableNames := dbmodel.GetTables()
	indexes := dbmodel.GetIndexes()
	err = db.ExecTx(context.Background(), func(txdb *gorn.DB) error {
		// Set Foreign Key Check Off
		rnlog.Info("Set Foreign Key Check Off...")
		sql := gorn.NewSql().
			Set("@OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0")
		if result, err := txdb.Exec(context.Background(), sql); err != nil {
			return err
		} else if _, err := result.RowsAffected(); err != nil {
			return err
		}

		// Migrate Table
		rnlog.Info("Migrate Tables...")
		for i, table := range tables {
			rnlog.Info("Migrate Table: %s", tableNames[i])
			tableName := tableNames[i]
			if err := txdb.Migration(tableName, table); err != nil {
				rnlog.Error("Migrate Table Error: %s", err.Error())
				return err
			}
		}
		// Migrate Index
		rnlog.Info("Migrate Indexes...")
		if err := txdb.MigrationIndex(indexes); err != nil {
			rnlog.Error("Migrate Index Error: %s", err.Error())
			return err
		}
		// Set Foreign Key Check On
		rnlog.Info("Set Foreign Key Check On...")
		sql = gorn.NewSql().
			Set("FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS")
		if result, err := txdb.Exec(context.Background(), sql); err != nil {
			return err
		} else if _, err := result.RowsAffected(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		rnlog.Error("Migrate Error: %+v\n", err)
		return
	}
	rnlog.Info("Migrate Success")
}
