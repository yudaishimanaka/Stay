<div align="center">
<img src="https://github.com/yudaishimanaka/Stay/blob/master/assets/images/analysis.png" alt="Stay" width="128" height="128">
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
  - github.com/mdlayher/arp

## Usage  
1. Setting `config.toml`  
  ```toml
  [app]
  interface = "" # e.g.) interface = "ehe0"
  network = "" # e.g.) network = "192.168.1.0/24"
  arp_interval =   # Second designation. e.g.) arp_interval = 60
  arp_timeout =   # Millisecond designation. e.g.) arp_timeout = 10
  [mysql]
  user = "" # Your mysql user
  password = "" # Your mysql password
  database = "" # Using database name
  ```

2. Create database(first only)  
  `$ go run db/migrate.go`  

3. Run server  
  `$ sudo go run *.go`  

4. Access dashboard  
  localhost:8888/dashboard  

5. Register user  
```bash
curl -v -X POST \
> -H "Accept: application/json" \
> -H "Content-Type: multipart/form-data" \
> -F "user_name=example" \
> -F "hw_addr=XX:XX:XX:XX:XX:XX" \
> -F "icon=@/path/to/example.png"
> localhost:8888/user/register
```

## LICENSE  
The MIT License (MIT) -see `LICENSE` for more details.  
