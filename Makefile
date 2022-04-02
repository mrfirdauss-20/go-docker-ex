run:
	docker build -t hex-mathrush-mem-server -f ./build/package/mem_server/Dockerfile .
	docker run -p 9190:9190 hex-mathrush-mem-server