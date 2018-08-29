<div align="center">
<img src="https://github.com/yudaishimanaka/searcher/blob/master/images/analysis.png" alt="armadillo" width="128" height="128">
</div>

# Stay
In-room management system using arp

## Requirement
- Golang 1.10.x ~
- External package
  - github.com/BurntSushi/toml
  - github.com/go-xorm/xorm
  - github.com/gin-gonic/gin
  - gopkg.in/olahol/melody.v1
  - github.com/go-sql-driver/mysql

## Install

## Usage

### Create user
```
curl -i -X POST \
> -H "Accept: application/json" \
> -H "Content-Type: multipart/form-data" \
> -F "user_name=user_name" \
> -F "hw_addr=XX:XX:XX:XX:XX:XX" \
> -F "icon=@/path/to/example.png"
> localhost:8888/user/register
```

## LICENSE
The MIT License (MIT) -see `LICENSE` for more details.
