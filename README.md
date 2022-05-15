# sampleSSO


## Steps to run
1. Clone the repository. 
2. `go mod download`
3. Run the following two commands.
```
$ openssl req -x509 -newkey rsa:2048 -keyout myservice.key -out myservice.cert -days 365 -nodes -subj "/CN=myservice.example.com"
```
It should generate two files in root folder.

```
mdpath=saml-test-$USER-$HOST.xml
curl localhost:8000/saml/metadata > $mdpath
```
It should generate an XML file. 
- Navigate to https://samltest.id/upload.php and upload the file you fetched. 
- Copy three generated files to `test` folder.

## Run test (it is using vanilla http)
1. `cd test`
2. `go run main.go`
3. go to your browser and visit `localhost:8080/hello`
4. It will redirect to samltest page. Use credentials that are provided at the bottom of page.

## Run using gin
1. Ensure that you are in root folder.
2. `go run main.go`
3. go to your browser and visit `localhost:8080/login`

## Bug to resolve
When running gin server, it should redirect us to samltest webpage.
