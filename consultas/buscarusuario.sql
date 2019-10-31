select id, email, newpass as password from produtosclientes.usuario
where 1 = 1
{{#email}} and email = upper({:email}) {{/email}}
{{#id}} and id = {:id} {{/id}}