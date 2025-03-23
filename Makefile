# Зависимости
dep:
	go mod tidy
.PHONY: dep

# Запуск тестов
test:
	cd test
	docker compose up -d
	go test ./...
	docker compose down
	cd ..
.PHONY: test

# Запуск всей системы в докер контейнере
dk-start:
	docker compose up -d
.PHONY: dk-start

# Остановка всей системы
dk-stop:
	docker compose down
.PHONY: dk-stop

# Help
h:
	@echo "Usage: make [target]"
	@echo "  target is:"
	@echo "       dep	- Обновление зависимостей"
	@echo "    test		- Запуск всех тестов"
	@echo "  dk-start	- Запуск служб в докер контейнерах (окружения)"
	@echo "   dk-stop	- Остановка запущенных служб (окружения)"
.PHONY: h
help: h
.PHONY: help
