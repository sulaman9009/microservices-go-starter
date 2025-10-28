.PHONY: run-all run-api-gateway run-others

# Run everything: api-gateway (with env) plus all other services found under services/
run-all: run-api-gateway run-others
	@echo "All services started."

# api-gateway must be started with specific env vars
run-api-gateway:
	@echo "Starting api-gateway with TRIP_SERVICE_URL=localhost:9093 and DRIVER_SERVICE_URL=localhost:9092..."
	@cd services/api-gateway && TRIP_SERVICE_URL=localhost:9093 DRIVER_SERVICE_URL=localhost:9092 air &

# Start all other services discovered under services/ using air
run-others:
	@for d in services/*; do \
		name=$$(basename $$d); \
		if [ "$$name" != "api-gateway" ] && [ -d "$$d" ]; then \
			echo "Starting $$name..."; \
			( cd "$$d" && air ) & \
		fi; \
	done; \
	wait
