package services

import (
	db "github.com/elnerribeiro/go-mustache-db"
	repo "github.com/elnerribeiro/go-ws-db-auth/repositories"
	u "github.com/elnerribeiro/go-ws-db-auth/utils"
)

//ListInserts Lists all inserts for one batch by id
func ListInserts(insert *repo.Insert) (*repo.Insert, error) {
	return insert.ListInserts()
}

//ClearInserts Clears tables
func ClearInserts(insert *repo.Insert) error {
	tx, err := db.GetTransaction()
	defer db.Rollback(tx)
	if err != nil {
		u.Logger.Error("[ClearInserts] Error starting transaction: %s", err)
		return err
	}
	if err2 := insert.ClearBatches(tx); err2 != nil {
		u.Logger.Error("[ClearInserts] Error cleaning batches: %s", err2)
		return err2
	}
	db.Commit(tx)
	return nil
}

//InsertBatchSync Inserts a batch of given quantity synchronous
func InsertBatchSync(insert *repo.Insert) (*repo.Insert, error) {
	tx, err := db.GetTransaction()
	defer db.Rollback(tx)
	if err != nil {
		u.Logger.Error("[InsertBatchSync] Error starting transaction: %s", err)
		return nil, err
	}
	insert.Type = "sync"
	ins, err := insert.InsertID(tx)
	if err != nil {
		u.Logger.Error("[InsertBatchSync] Error inserting a batch: %s", err)
		return nil, err
	}

	for i := 1; i < insert.Quantity; i++ {
		insertBatch := &repo.InsertBatch{}
		insertBatch.ID_Ins_ID = ins.ID
		insertBatch.Pos = i
		if err := insertBatch.InsertOneBatch(tx); err != nil {
			u.Logger.Error("[InsertBatchSync] Error inserting one item: %s", err)
		}
	}
	ins.Status = "Finished"
	ins.UpdateInsertID(tx)
	db.Commit(tx)
	return ins.ListInserts()
}

//InsertBatchASync Inserts a batch of given quantity asynchronous
func InsertBatchASync(insert *repo.Insert) (*repo.Insert, error) {
	tx, err := db.GetTransaction()
	defer db.Rollback(tx)
	if err != nil {
		u.Logger.Error("[InsertBatchASync] Error starting transaction: %s", err)
		return nil, err
	}
	insert.Type = "async"
	ins, err := insert.InsertID(tx)
	ins.Quantity = insert.Quantity
	if err != nil {
		u.Logger.Error("[InsertBatchASync] Error inserting a batch: %s", err)
		return nil, err
	}
	db.Commit(tx)
	go insertBatch(ins)
	return ins, nil
}

func insertBatch(insert *repo.Insert) {
	tx, err := db.GetTransaction()
	if err != nil {
		u.Logger.Error("[insertBatch] Error starting transaction: %s", err)
		return
	}
	defer db.Rollback(tx)
	for i := 1; i < insert.Quantity; i++ {
		insertBatch := &repo.InsertBatch{}
		insertBatch.ID_Ins_ID = insert.ID
		insertBatch.Pos = i
		if err := insertBatch.InsertOneBatch(tx); err != nil {
			u.Logger.Error("[insertBatch] Error inserting one item: %s", err)
			onError(tx, insert)
		}
	}
	insert.Status = "Finished"
	insert.UpdateInsertID(tx)
	db.Commit(tx)
}

func onError(tx *db.Transacao, insert *repo.Insert) {
	db.Rollback(tx)
	newtx, err := db.GetTransaction()
	if err != nil {
		u.Logger.Error("[onError] Error starting transaction: %s", err)
		return
	}
	insert.Status = "Error"
	insert.UpdateInsertID(newtx)
	db.Commit(newtx)
}
