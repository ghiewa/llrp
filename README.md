

# llrp
Official golang implement Low Level Reader Protocol(LLRP)
## Document and referent
You can get there in this repo.
## How to run test
structor of envoriment must move main.go to upper folder of llrp package folder e.g.
```
sample-llrp/
    main.go
    llrp/
        llrp.go
        llrp_test.go
        ....
```
then  ```   go run main.go ``` on `sample-llrp` folder
By default it will connect to physical reader (speed way) IP `192.168.33.16` and port `5084`

## Add more functionality
List of function I implement base on we used. So if you want to add more function of this package go to `request.go` ,`respone.go` and `param_parse.go` files.These files are core of this package.
### 
 ## Issue 
 Received oversize incoming message while run long period
>We fixed by reset & reconfigure them (it's take ~ 10 - 20 second).If you have another way to fixed it.Feel free to contact me fixed in this package. 

