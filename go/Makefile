NAME := checkplease
BUILD_DIR := build
BUILD_OUT := $(BUILD_DIR)/$(NAME)
SSH_HOST := checkplease

.PHONY: pi-build pi-scp pi-run clean

pi-build: $(BUILD_OUT)

pi-scp: $(BUILD_OUT)
	scp $(BUILD_OUT) $(SSH_HOST):'~/$(NAME)'

pi-run: pi-scp
	ssh $(SSH_HOST) -t '~/$(NAME)'

$(BUILD_OUT): *.go
	GOOS=linux GOARCH=arm GOARM=5 go build -mod readonly -o $(BUILD_OUT) .

clean:
	rm $(BUILD_OUT)
	ssh $(SSH_HOST) -t 'rm ~/$(NAME)'
