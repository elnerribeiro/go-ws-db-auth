package repositories

import (
	"time"

	db "github.com/elnerribeiro/go-mustache-db"
)

//ListInserts Retrieve one batch of inserts by id
func ListInserts(id int) (*Insert, error) {
	var dados db.Dados
	dados = make(db.Dados)
	dados["id"] = id
	val, err := db.SelectOne(nil, "getinsid", &Insert{}, dados)
	if err != nil {
		return nil, err
	}
	insert := val.(*Insert)
	if insert != nil {
		var dados db.Dados
		dados = make(db.Dados)
		dados["id_ins_id"] = id
		val2, err := db.SelectAll(nil, "getinsbatch", &[]InsertBatch{}, dados)
		if err != nil {
			return insert, nil
		}
		insert.ListVals = *(val2.(*[]InsertBatch))
	}
	return insert, nil
}

//InsertOneBatch Inserts one item of the batch
func InsertOneBatch(tx *db.Transacao, id int, pos int) error {
	var dados db.Dados
	dados = make(db.Dados)
	dados["id_ins_id"] = id
	dados["pos"] = pos
	_, err := db.Insert(tx, "insert_batch", dados)
	if err != nil {
		return err
	}
	return nil
}

//InsertID Inserts one batch of items
func InsertID(tx *db.Transacao, qtd int, typeCall string) (*Insert, error) {
	var dados db.Dados
	dados = make(db.Dados)
	dados["quantity"] = qtd
	dados["status"] = "Running"
	dados["type"] = typeCall
	tstamp := time.Now().Unix()
	dados["tstampinit"] = tstamp
	val, err := db.InsertReturningPostgres(tx, "ins_id", dados, "id", &Insert{})
	if err != nil {
		return nil, err
	}
	ins := val.(*Insert)
	ins.Quantity = qtd
	ins.Status = "Running"
	ins.Type = typeCall
	ins.Tstampinit = tstamp
	return ins, nil
}

//UpdateInsertID Finishes batch insertion
func UpdateInsertID(tx *db.Transacao, status string, id int) (int64, error) {
	var dados db.Dados
	dados = make(db.Dados)
	dados["status"] = status
	tstamp := time.Now().Unix()
	dados["tstampend"] = tstamp
	var filters db.Dados
	filters = make(db.Dados)
	filters["id"] = id
	_, err := db.Update(tx, "ins_id", dados, filters)
	if err != nil {
		return 0, err
	}
	return tstamp, nil
}

//ClearBatches Removes all batches
func ClearBatches(tx *db.Transacao) error {
	var dados db.Dados
	dados = make(db.Dados)
	_, err := db.ExecuteSQL(tx, "removeall", dados)
	if err != nil {
		return err
	}
	_, err2 := db.ExecuteSQL(tx, "removeallids", dados)
	if err2 != nil {
		return err2
	}
	return nil
}
