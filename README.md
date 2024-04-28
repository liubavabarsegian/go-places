# A service for finding the closest restaurants by location
> This is my first go app

## Run app
To deploy the app locally you simply run `make start` in your cmd. 
To stop the app run `make stop`

## Loading Data with ElasticSearch
All the restaurants are loaded from CSV file `internal/config/data.csv` and inserted into EalsticSearch storage using Bulk API for a faster performance.

Check info about the index:

```zsh
curl -s -XGET "http://localhost:9200/places"
```

Add new restaurants: (the mapping is in `internal/config/schema.json` file)

```zsh
curl -XPUT "http://localhost:9200/places"
```

You can query the restaurants by ID like this:

```zsh
curl -s -XGET "http://localhost:9200/places/_doc/1"
```

## Interface
To access the view, follow [http://localhost:8888/?page=1](http://localhost:8888/?page=1)

<img width="646" alt="image" src="https://github.com/liubavabarsegian/s21-go-places-app/assets/74152874/5996f28c-93c2-40cb-b096-708c4d6d31af">


## API
```zsh
curl --request GET \
  --url 'http://localhost:8888/api/recommend?lat=55.674&lon=37.666' \
  --header 'Accept: */*' \
  --header 'Content-Type: application/json'
```

By default Elasticsearch doesn't allow you to deal with pagination for more than 10000 entries. To fix that, run `make update_index_settings`

## Closest restaurants
### Get JWT
Before calling API for the closest restaurants, you need go get token.

```zsh
http://127.0.0.1:8888/api/get_token
```
That token is used for the authorization.
### Get recommended restaurants
Get the closest restaurants according to the location (latitude, longitude). Don't forget to use your specific token for authorization.
```zsh
curl --request GET \
  --url 'http://localhost:8888/api/recommend?lat=55.674&lon=37.666' \
  --header 'Accept: */*' \
  --header 'Authorization: Bearer <token>' \
  --header 'Content-Type: application/json'
```
