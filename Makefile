.PHONY: nats config scheduler alerter frontend run-all stop-all

# Démarre le serveur NATS (crée le conteneur s'il n'existe pas)
nats:
	@echo "Starting NATS..."
	@docker start nats-server > /dev/null 2>&1 || docker run -d --name nats-server -p 4222:4222 nats -js

# Lance l'API Config
config:
	@echo "Starting Config API..."
	@cd config && go run cmd/main.go

# Lance le Scheduler
scheduler:
	@echo "Starting Scheduler..."
	@cd scheduler && go run cmd/main.go

# Lance l'Alerter
alerter:
	@echo "Starting Alerter..."
	@cd alerter && go run cmd/main.go

# Lance tout en même temps (les logs seront mélangés)
# Utiliser Ctrl+C pour arrêter
run-all: nats
	@echo "Starting all services..."
	@(cd config && go run cmd/main.go) & \
	(cd scheduler && go run cmd/main.go) & \
	(cd alerter && go run cmd/main.go) & \
	(cd frontend && npm run dev) & \
	wait

# Lance le Frontend
frontend:
	@echo "Starting Frontend..."
	@cd frontend && npm run dev
