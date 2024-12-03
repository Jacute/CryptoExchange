clear-db:
	@docker compose down
	@docker compose -f docker-compose-balance.yml down
	@sudo rm -rf JacuteSQL/storage
run-randombot:
	@python3 TradingRobot/main.py randombot --host 127.0.0.1 --port 8080 --timeout 0.5
run-richbot:
	@python3 TradingRobot/main.py richbot --host 127.0.0.1 --port 8080 --timeout 10