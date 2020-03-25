select id, email, password as password, role from usuario
where 1 = 1
{{#email}} and email = lower({:email}) {{/email}}
{{#id}} and id = {:id} {{/id}}