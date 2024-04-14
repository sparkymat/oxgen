gen: clean
	go run oxgen.go resource -p webapp -q name -s User name:string:default=:updateable username:string:not_null:unique encrypted_password:string:not_null age:int:updateable dob:date:updateable photo:attachment

clean:
	rm -f oxgen
	rm -rf webapp/migrations
	echo > webapp/internal/database/queries.sql
	psql webapp -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
