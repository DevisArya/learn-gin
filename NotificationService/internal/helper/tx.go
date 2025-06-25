package helper

import "gorm.io/gorm"

func CommitOrRollback(tx *gorm.DB) {

	if err := recover(); err != nil {
		tx.Rollback()
		panic(err)
	}

	if errCommit := tx.Commit().Error; errCommit != nil {
		panic(errCommit)
	}

}
