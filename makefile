.PHONY: build push pb

build:
	docker buildx build --platform=linux/amd64 -t spacecoupe/pricescraper-worker .

push:
	docker push spacecoupe/pricescraper-worker

pb: build push

update ecs:
	aws ecs update-service --cluster default --service priceScraperV3 --force-new-deployment --region us-east-2 --profile rarityshark