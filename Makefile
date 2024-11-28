clear-db:
	@docker compose down
	@sudo rm -rf JacuteSQL/storage
	@docker compose up -d
run-randombot:
	@python3 TradingRobot/main.py randombot --host 127.0.0.1 --port 8080 --timeout 1
run-richbot:
	@python3 TradingRobot/main.py richbot --host 127.0.0.1 --port 8080 --timeout 10