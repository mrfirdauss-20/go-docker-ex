run-mem-server:
	docker build -t hex-mathrush-mem-server -f ./build/package/mem_server/Dockerfile .
	docker run -p 9190:9190 hex-mathrush-mem-server

run-sql-server:
	-docker compose -f ./deploy/sql_server/docker-compose.yml down
	docker compose -f ./deploy/sql_server/docker-compose.yml up --build