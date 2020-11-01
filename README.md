# midnight
Golang middleware Kong plugin for study purposes only.

## Build
```
go build -buildmode plugin midnight.go
```

## Test

In [sample/change-pokemon-type.yml](change-pokemon-type.yml) has the configuration for the kong service / route and plugin that makes the transformation in the api https://docs.pokemontcg.io/#api_v1cards_get. With the following flow:

  - Change service path [:id] v1/cards/:id with path request [.+] /change/.+
  - Make request to api
  - To each response "card.types" random new type
  - Build result on "message"
  - Return to client

```
cd sample
docker-compose up

curl -v localhost:8000/change/xy7-54

# {"message":"Change types for pokemon Gardevoir[xy7-54] from [Fairy] to [Darkness]"}
```

## How it works
  - Configure Kong Service with upstream URL. You can add named parameter on path with ":" + parameter name. 

        EX.
        https://api.pokemontcg.io/v1/cards/:id
  
  - Configure Kong Route

  - Add plugin midnight with following config:

    - [uri] - key (parameter name) value (index on request path) pair where to get named parameters

          EX.
          # config
          {"id":1}
          # request path
          /change/xy7-54
          # kong service path
          /v1/cards/:id
          # become
          /v1/cards/xy7-54

    - [in/out] - golang template to build with request params ([in]) to change request body to upstream and response params ([out]) to change resposen body to client. Midnigth add useful template sprig functions http://masterminds.github.io/sprig/
    
          EX.
          # response body
          {"card":{"name": "Pikachu"}}
          # config [out] template
          {{ {"pokemon": "{{ .card.name }}"} }}
          # become
          {"pokemon": "Pikachu"}

## Etc

  - Kong plugin has some [limitations](https://docs.konghq.com/enterprise/2.1.x/go/#limitations-of-the-go-pdk).There are no header_filter or body_filter phases.

  - Length must be at least 1 on kong plugin configs. Add an empty space in the unused configuration.