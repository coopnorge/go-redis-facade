# db-facade

db-facade is a sharable package for connecting to and interacting with a database. Currently this is package is configured for Redis.

### importing

```go
package main
import "dev.azure.com/coopnorge/Scan-and-pay/db-facade.git"
```
**Note:** This is a private repo. In order to vendor the package you have to [update you git config.](https://dev.azure.com/coopnorge/Engineering/_wiki/wikis/Engineering.wiki/714/Go?anchor=rationale)
### How to use
```go
package main

import dbFacade "dev.azure.com/coopnorge/Scan-and-pay/db-facade.git"

func main(){
  dbConfiguration := dbFacade.Config{
    <your configurations>
  }
  dbFacade := dbFacade.NewRedisFacade(dbConfiguration)
}
```
