# frisby
API testing framework inspired by frisby-js, written in Go



### Installation

```
go get -U github.com/verdverm/frisby
```

### Basic Usage

First create a frisby object:

```
// create an object with a given name (used in the report)
F := frisby.Create("Test successful user login").
    Get("https://golang.org")
```

Next perform the operation:

```
F.Send()
```

Then inspect the response:

```
F.ExpectStatus(200).
    ExpectContent("The Go Programming Language").
    PrintReport()
```

Check any error(s):

```
err := F.Error()

errs := F.Errors()
for _,e := range errs {
	fmt.Println("Error: ", e)
}

### HTTP Method functions

Your basic HTTP verbs:

- Get(url string)
- Post(url string)
- Put(url string)
- Patch(url string)
- Delete(url string)
- Head(url string)
- Options(url string)

### Pre-flight functions

Functions called before `Send()`

( Most of these come from [github.com/mozillazg/request](https://github.com/mozillazg/request))

- BasicAuth(username,password string)
- Proxy(url string)
- SetHeader(key,value string)
- SetHeaders(map[string]string)
- SetCookies(key,value string)
- SetCookiess(map[string]string)
- SetDate(key,value string)
- SetDates(map[string]string)
- SetParam(key,value string)
- SetParams(map[string]string)
- SetJson(interface{})
- SetFile(filename string)


### Post-flight functions

Functions called after `Send()`

- ExpectStatus(code int)
- ExpectHeader(key, value string)
- ExpectContent(content string)
- ExpectJson(path string, value interface{})
- ExpectJsonType(path string, value_type reflect.Kind)
- AfterContent( func(Frisby,[]byte,error) )
- AfterText( func(Frisby,string,error) )
- AfterJson( func(Frisby,simplejson.Json,error) )
- PrintBody()
- PrintReport()


