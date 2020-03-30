package repositories

import (
	"time"

	db "github.com/elnerribeiro/go-mustache-db"
	u "github.com/elnerribeiro/go-ws-db-auth/utils"
)

//ListInserts Retrieve one batch of inserts by id
func (insert *Insert) ListInserts() (*Insert, error) {
	var dados db.Dados
	dados = make(db.Dados)
	dados["id"] = insert.ID
	val, err := db.SelectOne(nil, "getinsid", &Insert{}, dados)
	if err != nil {
		u.Logger.Error("[ListInserts] Error listing inserts: %s", err)
		return nil, err
	}
	insert = val.(*Insert)
	if insert != nil {
		var dados db.Dados
		dados = make(db.Dados)
		dados["id_ins_id"] = insert.ID
		val2, err := db.SelectAll(nil, "getinsbatch", &[]InsertBatch{}, dados)
		if err != nil {
			u.Logger.Error("[ListInserts] Cannot find children: %s", err)
			return insert, nil
		}
		insert.ListVals = *(val2.(*[]InsertBatch))
	}
	return insert, nil
}

//InsertOneBatch Inserts one item of the batch
func (insert *InsertBatch) InsertOneBatch(tx *db.Transacao) error {
	var dados db.Dados
	dados = make(db.Dados)
	dados["id_ins_id"] = insert.ID_Ins_ID
	dados["pos"] = insert.Pos
	_, err := db.Insert(tx, "insert_batch", dados)
	if err != nil {
		u.Logger.Error("[InsertOneBatch] Cannot insert children for id %d: %s", insert.ID_Ins_ID, err)
		return err
	}
	return nil
}

//InsertID Inserts one batch of items
func (insert *Insert) InsertID(tx *db.Transacao) (*Insert, error) {
	var dados db.Dados
	dados = make(db.Dados)
	dados["quantity"] = insert.Quantity
	dados["status"] = "Running"
	dados["type"] = insert.Type
	tstamp := time.Now().Unix()
	dados["tstampinit"] = tstamp
	val, err := db.InsertReturningPostgres(tx, "ins_id", dados, "id", &Insert{})
	if err != nil {
		u.Logger.Error("[InsertID] Cannot insert new batch: %s", err)
		return nil, err
	}
	insertReturn := val.(*Insert)
	insertReturn.Quantity = insert.Quantity
	insertReturn.Status = "Running"
	insertReturn.Type = insert.Type
	insertReturn.Tstampinit = tstamp
	return insertReturn, nil
}

//UpdateInsertID Finishes batch insertion
func (insert *Insert) UpdateInsertID(tx *db.Transacao) (int64, error) {
	var dados db.Dados
	dados = make(db.Dados)
	dados["status"] = insert.Status
	tstamp := time.Now().Unix()
	dados["tstampend"] = tstamp
	var filters db.Dados
	filters = make(db.Dados)
	filters["id"] = insert.ID
	_, err := db.Update(tx, "ins_id", dados, filters)
	if err != nil {
		u.Logger.Error("[UpdateInsertID] Cannot update batch: %s", err)
		return 0, err
	}
	return tstamp, nil
}

//ClearBatches Removes all batches
func (insert *Insert) ClearBatches(tx *db.Transacao) error {
	var dados db.Dados
	dados = make(db.Dados)
	_, err := db.ExecuteSQL(tx, "removeall", dados)
	if err != nil {
		u.Logger.Error("[ClearBatches] Cannot remove children: %s", err)
		return err
	}
	_, err2 := db.ExecuteSQL(tx, "removeallids", dados)
	if err2 != nil {
		u.Logger.Error("[ClearBatches] Cannot remove batches: %s", err)
		return err2
	}
	return nil
}
