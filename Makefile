.PHONY: swag

swag:
	swag init -g cmd/main.go -o docs
