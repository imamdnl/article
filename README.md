# article

## Requirements

1. Go 1.18
1. PostgreSQL
1. Redisearch

## Running

Download the required packages:

```
go mod download
```

Run the docker compose

```
docker-compose up -d
```

### API List & Payloads
Here is our API List and its payload:

1. [GET] **/articles?query=harry&author=rowling**  
   `/articles?query=harry&author=rowling`
2. [POST] **/articles**
```javascript
{
    "author": "JK Rowling",
    "title": "Harry Potter",
    "body": "lorem ipsum lorem ipsum lorem ipsum"
}
```