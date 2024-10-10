package main

import (
	"flag"
	"log"
	"net/http"
)

const html = `
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8">
	</head>
	<body>
		<h1>Simple CORS</h1>
		<div id="output"></div>
		<script>
			document.addEventListener('DOMContentLoaded', () => {
				fetch("http://localhost:5000/v1/healthcheck").then(
					function (response) {
						response.text().then((text) => {
							document.getElementById("output").innerHTML = text;
						});
					},
					function(err) {
						document.getElementById("output").innerHTML = err;
					}
				);
			});
		</script>
	</body>
</html>`

func main() {
	addr := flag.String("addr", ":9000", "Server Address")
	flag.Parse()

	log.Printf("starting server on %s", *addr)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	})

	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal(err)
	}
}
