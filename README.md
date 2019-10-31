# gosbxauth
. Criar arquivo consultas/buscarusuario.sql com select retornando campos id, email e password

. Criar application.properties
	printSql=false
    
    dialeto=mysql (postgres)
	
	db.url=postgres://USERNAME:PASSWORD@HOST:PORTA/BANCO (Postgres)
	
ou

	db.url=username:password@host:porta/banco?param=value (MySQL)

. Criar variavel de ambiente PORT com a porta

/api/user/login

/api/user/validate

. Remover rota /api/user/list