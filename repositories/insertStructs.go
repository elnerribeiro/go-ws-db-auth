package repositories

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
