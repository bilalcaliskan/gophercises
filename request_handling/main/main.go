package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type apiHandler struct {
	
}

func (apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
}

type timeHandler struct {
	format string
}

func (th *timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(th.format)
	w.Write([]byte("The time is: " + tm))
}

func timeHandlerFunction(format string) http.Handler {
	/*
	First it creates fn, an anonymous function which accesses ‐ or closes over – the format variable forming a closure.
	Regardless of what we do with the closure it will always be able to access the variables that are local to the scope
	it was created in – which in this case means it'll always have access to the format variable.
	Secondly our closure has the signature func(http.ResponseWriter, *http.Request). As you may remember from earlier,
	this means that we can convert it into a HandlerFunc type (so that it satisfies the Handler interface). Our timeHandler
	function then returns this converted closure.
	 */
	fn := func(w http.ResponseWriter, r *http.Request) {
		tm := time.Now().Format(format)
		w.Write([]byte("The time is: " + tm))
	}
	// Below line converts fn function to a handler function
	return http.HandlerFunc(fn)
}

func main() {
	/*
	A ServeMux is essentially a HTTP request router (or multiplexor). It compares incoming requests against a list of
	predefined URL paths, and calls the associated handler for the path whenever a match is found.
	Handlers are responsible for writing response headers and bodies. Almost any object can be a handler, so long as it
	satisfies the http.Handler interface. In lay terms, that simply means it must have a ServeHTTP method with the
	following signature:
		ServeHTTP(http.ResponseWriter, *http.Request)
	ServeMux type also has a ServeHTTP method, meaning that it too satisfies the Handler interface. ServeMux as just
	being a special kind of handler
	 */
	mux := http.NewServeMux()
	rh := http.RedirectHandler("http://example.org", 307)
	mux.Handle("/foo", rh)

	/*
	Custom Handlers
	 */
	th1123 := &timeHandler{format:time.RFC1123}
	th3339 := &timeHandler{format:time.RFC3339}
	mux.Handle("/time/rfc1123", th1123)
	mux.Handle("/time/rfc3339", th3339)

	/*
	Functions as Handlers
	Any function which has the signature func(http.ResponseWriter, *http.Request) can be converted into a HandlerFunc
	type. This is useful because HandleFunc objects come with an inbuilt ServeHTTP method which – rather cleverly and
	conveniently – executes the content of the original function.
	In fact, converting a function to a HandlerFunc type and then adding it to a ServeMux like this is so common that
	Go provides a shortcut: the mux.HandleFunc method.
	 */
	// Below is an example of handler function using anonymous function
	mux.HandleFunc("/time0", func(writer http.ResponseWriter, request *http.Request) {
		tm := time.Now().Format(time.RFC3339)
		writer.Write([]byte("The time is: " + tm))
	})
	// Here is an another example
	th0 := timeHandlerFunction(time.RFC3339)
	mux.Handle("/time1", th0)

	/*
	Documentation example
	 */
	mux.Handle("/api/", apiHandler{})
	mux.HandleFunc("/api/v2", func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/v2" {
			http.NotFound(writer, request)
			return
		}
		fmt.Fprintf(writer, "Welcome to the home page!")
	})
	
	/*
	DefaultServerMux
	var DefaultServeMux = NewServeMux()
	Generally you shouldn't use the DefaultServeMux because it poses a security risk.
	Because the DefaultServeMux is stored in a global variable, any package is able to access it and register a route – including
	any third-party packages that your application imports. If one of those third-party packages is compromised, they
	could use the DefaultServeMux to expose a malicious handler to the web.
	So as a rule of thumb it's a good idea to avoid the DefaultServeMux, and instead use your own locally-scoped ServeMux
	 */
	var format string = time.RFC1123
	th1 := timeHandlerFunction(format)
	http.Handle("/time2", th1)

	log.Println("Listening on port 3000 for ServeMux...")
	// run ServeMux on port 3000
	go http.ListenAndServe(":3000", mux)
	log.Println("Listening on port 3001 for DefaultServeMux...")
	// run DefaultServeMux on port 3001
	go http.ListenAndServe(":3001", nil)
	time.Sleep(10 * time.Minute)
}