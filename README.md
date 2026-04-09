# Lab 11 - Docker Multi-Platform Build

## H4: Cross-platform build with buildx

### Prerequisites
```bash
docker buildx create --name mybuilder --use
docker buildx inspect --bootstrap

```
(.venv) PS D:\Work\UniStud\6 семместр\МиТП\Лабы\lab_11\task\go-app> docker compose ps
NAME         IMAGE             COMMAND           SERVICE      CREATED          STATUS                    PORTS
go-app       task-go-app       "/app"            go-app       47 seconds ago   Up 47 seconds (healthy)   0.0.0.0:8080->8080/tcp, [::]:8080->8080/tcp
python-app   task-python-app   "python app.py"   python-app   47 seconds ago   Up 41 seconds (healthy)   0.0.0.0:8000->8000/tcp, [::]:8000->8000/tcp
```