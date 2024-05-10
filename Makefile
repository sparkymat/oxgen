gen: clean
	go run oxgen.go resource --path=webapp --query-field=title --service=blog --skip-git Post title:string:not_null:default=:updateable body:string:not_null:default= photo:attachment
	go run oxgen.go resource --path=webapp --query-field=username --service=blog --parent=post --skip-git Comment username:string:not_null body:string:not_null:default=:updateable

clean:
	rm -f oxgen
	rm -rf webapp/migrations
	rm -rf webapp/internal/service
	rm -rf webapp/internal/*.go
	rm -rf webapp/internal/handler/*.go
	rm -rf webapp/internal/handler/api/*_*.go
	rm -rf webapp/internal/handler/api/presenter/*.go
	rm -rf webapp/frontend/src/models
	rm -rf webapp/frontend/src/slices
	mkdir -p webapp/internal/route
	mkdir -p webapp/internal/service
	mkdir -p webapp/internal/handler
	mkdir -p webapp/frontend/src/models
	mkdir -p webapp/frontend/src/slices
	echo > webapp/internal/database/queries.sql
	echo "package service; type DatabaseProvider interface {}" | goimports > webapp/internal/service/database_iface.go
	echo "package internal; type Services struct {}" | goimports > webapp/internal/services.go
	echo "package route; func registerAPIRoutes(app *echo.Group, cfg internal.ConfigService, services internal.Services) {apiGroup := app.Group(\"api\")}" | goimports > webapp/internal/route/api.go
	psql webapp -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
