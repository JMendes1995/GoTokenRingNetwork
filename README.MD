# how to test

cd token ring

docker exec -it tokenring-server-1 bash
go run main.go server


docker exec -it tokenring-peer1-1 bash
docker exec -it tokenring-peer2-1 bash
docker exec -it tokenring-peer3-1 bash
go run main.go peer



