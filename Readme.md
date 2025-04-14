# Get started
### Start the Application and Database Containers
```sh
docker compose -f docker-compose.yml up wallet-api-service --build
```
---

### Troubleshooting
If you encounter issues due to a change, run the following command for a clean installation:
```sh
docker-compose down --rmi all --volumes --remove-orphans
```
