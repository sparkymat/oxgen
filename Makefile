gen: clean
	go run oxgen.go resource -p webapp -q name -s blog -g User name:string:not_null:default=:updateable username:string:not_null:unique encrypted_password:string:not_null age:int:updateable dob:date:updateable photo:attachment

clean:
	rm -f oxgen
	rm -rf webapp/migrations
	rm -rf webapp/internal/service
	rm -rf webapp/internal/handler/*.go
	rm -rf webapp/internal/handler/api/*_*.go
	rm -rf webapp/internal/handler/api/presenter/*.go
	mkdir -p webapp/internal/service
	mkdir -p webapp/internal/handler
	echo > webapp/internal/database/queries.sql
	echo "package service; type DatabaseProvider interface {}" | goimports > webapp/internal/service/database_iface.go
	psql webapp -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
