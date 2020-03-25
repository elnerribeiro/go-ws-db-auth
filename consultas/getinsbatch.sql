select id, id_ins_id, pos from insert_batch
where 1 = 1
{{#id}} and id = {:id} {{/id}}
{{#id_ins_id}} and id_ins_id = {:id_ins_id} {{/id_ins_id}}