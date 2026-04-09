# Lab 11 - Docker Multi-Platform Build

## H4: Cross-platform build with buildx

### Prerequisites
```bash
docker buildx create --name mybuilder --use
docker buildx inspect --bootstrap