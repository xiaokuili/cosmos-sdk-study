IMAGE_NAME = cosmos

.PHONY: buildchain build run

buildchain:
	rm -rf ./datafactory/build
	./ignite chain build --output ./datafactory/build --path ./datafactory
	./ignite chain init --home ./datafactory/build/.datafactory --path ./datafactory


build:
	docker build -t $(IMAGE_NAME) .

run:
	docker run -it $(IMAGE_NAME) /bin/bash
