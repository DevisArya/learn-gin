package helper

import (
	"gorm.io/gorm"
)

type HelperTx interface {
	Begin() *gorm.DB
	CommitOrRollback(tx *gorm.DB)
}

type DBTX struct {
	DB *gorm.DB
}

func NewHelper(DB *gorm.DB) HelperTx {
	return &DBTX{
		DB: DB,
	}
}

func (g *DBTX) Begin() *gorm.DB {
	return g.DB.Begin()
}

func (g *DBTX) CommitOrRollback(tx *gorm.DB) {

	if err := recover(); err != nil {
		tx.Rollback()
		panic(err)
	}

	if errCommit := tx.Commit(); errCommit != nil {
		panic(errCommit)
	}

}
