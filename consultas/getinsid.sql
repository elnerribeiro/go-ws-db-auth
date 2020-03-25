select id, type, quantity, status, tstampinit, coalesce(tstampend,0) from ins_id
where 1 = 1
{{#id}} and id = {:id} {{/id}}