## App launch

### 1. Locally

Create and configure .env file

```bash
   cp .env.example .env
 ```

Make app-launch.sh script executable with:

```bash
   chmod +x ./scripts/app-launch.sh
```

Run an app-launch script with:

```bash
  ./scripts/app-launch.sh
```

### 2. As a container

Create and configure .env file

```bash
   cp .env.example .env
 ```

Build Docker image with Dockerfile

```bash
   docker build -t lo-test-task .
```

Run docker-compose file

```bash
   docker compose up
```

## Launch tests

Make launch-tests.sh script executable with:

```bash
   chmod +x launch-tests.sh
```

Run an launch-tests script with:

```bash
   ./launch-tests.sh
```