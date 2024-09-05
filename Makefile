.PHONY: proto_gen proto_clean registry host cert docker_gen docker_clean# make commands work

# If common/pb folder doesn't exists, it is created
proto_gen:
	[ ! -d common/pb ] && mkdir "common/pb" || echo "common/pb already exists"
	protoc --proto_path=common/proto common/proto/*.proto  --go_out=:common/pb --go-grpc_out=:common/pb

proto_clean:
	rm common/pb/*

registry:
	go run ./registry/main.go -port 50051

host:
	go run ./host/main.go -membership_port 50152 -vivaldi_port 50153 -gossip_port 50154

host1:
	go run ./host/main.go -membership_port 50155 -vivaldi_port 50156 -gossip_port 50157

host2:
	go run ./host/main.go -membership_port 50158 -vivaldi_port 50159 -gossip_port 50160

host3:
	go run ./host/main.go -membership_port 50161 -vivaldi_port 50162 -gossip_port 50163

cert:
	cd cert; chmod +x gen.sh; ./gen.sh; cd ..

docker_gen:
	sudo docker compose -f docker-compose.yml up -d

docker_clean:
	sudo docker compose -f docker-compose.yml down # stop and remove all containers
	docker images | grep "sdcc" | awk '{print $3}' | xargs docker rmi # remove all images with name "sdcc*"
