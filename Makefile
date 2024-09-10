.PHONY: registry cert

registry:
	go run ./main.go -port 50051

cert:
	cd cert; chmod +x gen.sh; ./gen.sh; cd ..