include make.macro

default: build

run:
	go run cmd/DoItNext/main.go

build: clean
	@echo Building project...
	@$(call mkdir, $(BIN_DIR)config)
	@$(call cpdir, config $(BIN_DIR))
	@go build -a -o $(call exe, $(BIN_DIR)DoItNext) cmd/DoItNext/main.go
	@echo Building project done

clean:
	@echo Cleaning project...
	@$(call rmdir,$(BIN_DIR))
	@echo Cleaning project done
