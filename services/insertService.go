package services

import (
	db "github.com/elnerribeiro/go-mustache-db"
	repo "github.com/elnerribeiro/go-ws-db-auth/repositories"
)

//ListInserts Lists all inserts for one batch by id
func ListInserts(id int) (*repo.Insert, error) {
	return repo.ListInserts(id)
}

//ClearInserts Clears tables
func ClearInserts() error {
	tx, err := db.GetTransaction()
	defer db.Rollback(tx)
	if err != nil {
		return err
	}
	err2 := repo.ClearBatches(tx)
	if err2 != nil {
		return err2
	}
	db.Commit(tx)
	return nil
}

//InsertBatchSync Inserts a batch of given quantity synchronous
func InsertBatchSync(quantity int) (*repo.Insert, error) {
	tx, err := db.GetTransaction()
	defer db.Rollback(tx)
	if err != nil {
		return nil, err
	}
	ins, err := repo.InsertID(tx, quantity, "sync")
	if err != nil {
		return nil, err
	}
	for i := 1; i < quantity; i++ {
		repo.InsertOneBatch(tx, ins.ID, i)
	}
	repo.UpdateInsertID(tx, "Finished", ins.ID)
	db.Commit(tx)
	return repo.ListInserts(ins.ID)
}

//InsertBatchASync Inserts a batch of given quantity asynchronous
func InsertBatchASync(quantity int) (*repo.Insert, error) {
	tx, err := db.GetTransaction()
	defer db.Rollback(tx)
	if err != nil {
		return nil, err
	}
	ins, err := repo.InsertID(tx, quantity, "async")
	if err != nil {
		return nil, err
	}
	db.Commit(tx)
	go insertBatch(ins.ID, quantity)
	return ins, nil
}

func insertBatch(id int, quantity int) {
	tx, err := db.GetTransaction()
	if err != nil {
		return
	}
	defer db.Rollback(tx)
	for i := 1; i < quantity; i++ {
		err := repo.InsertOneBatch(tx, id, i)
		if err != nil {
			onError(tx, id)
		}
	}
	repo.UpdateInsertID(tx, "Finished", id)
	db.Commit(tx)
}

func onError(tx *db.Transacao, id int) {
	db.Rollback(tx)
	newtx, err := db.GetTransaction()
	if err != nil {
		return
	}
	repo.UpdateInsertID(newtx, "Error", id)
	db.Commit(newtx)
}
