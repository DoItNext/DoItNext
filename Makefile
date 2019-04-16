include make.macro

default: build

run:
	go run cmd/DoItNext/main.go

clean:
	@echo Cleaning project...
	@$(call rmdir,$(BIN_DIR))
	@echo Cleaning project done

build: clean
	@echo Building project...
	@$(call mkdir, $(BIN_DIR)config)
	@$(call cpdir, config $(BIN_DIR))
	@go build -a -o $(call exe, $(BIN_DIR)DoItNext) cmd/DoItNext/main.go
	@echo Building project done

database:
	@echo Starting mariaDB...
	docker run -d --rm \
		-p 3306:3306 \
		-v mariadb_doitnext:/var/lib/mysql \
		-e MYSQL_ROOT_PASSWORD=root \
        -e MYSQL_DATABASE=doitnext \
        -e MYSQL_USER=doitnext \
        -e MYSQL_PASSWORD=doitnext \
		--name mariadb craftdock/mariadb:10.2
	docker run -d --rm \
		-p 3307:80 \
		--link mariadb:db \
		--name myadmin phpmyadmin/phpmyadmin