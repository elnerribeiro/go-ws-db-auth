package repositories

import db "github.com/elnerribeiro/go-mustache-db"

//InsertBatch table insert_batch on database
type InsertBatch struct {
	ID        int `json:"id,omitempty"`
	ID_Ins_ID int `json:"id_ins_id,omitempty"`
	Pos       int `json:"pos,omitempty"`
}

//Insert table ins_id on database
type Insert struct {
	ID         int           `json:"id,omitempty"`
	Type       string        `json:"type,omitempty"`
	Quantity   int           `json:"quantity,omitempty"`
	Status     string        `json:"status,omitempty"`
	Tstampinit int64         `json:"tstampinit,omitempty"`
	Tstampend  int64         `json:"tstampend,omitempty"`
	ListVals   []InsertBatch `json:"list,omitempty"`
}

//InsertInterface interface for batch insert tables
type InsertInterface interface {
	UpdateInsertID(tx *db.Transacao) (int64, error)
	ClearBatches(tx *db.Transacao) error
	InsertID(tx *db.Transacao) (*Insert, error)
	InsertOneBatch(tx *db.Transacao) error
	ListInserts() (*Insert, error)
}
