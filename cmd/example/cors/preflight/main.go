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
		<h1>Preflight CORS</h1>
		<div id="output"></div>
		<script>
			document.addEventListener('DOMContentLoaded', function() {
				fetch("http://localhost:5000/v1/tokens/authentication", {
					method: "POST",
					headers: {'Content-Type': 'application/json'},
					body: JSON.stringify({
						email: 'john@mail.com',
						password: 'password'
					})
				}).then(function (response) {
						response.text().then(function (text) {
							document.getElementById("output").innerHTML = text;
						});
					},
					function (err) {
						document.getElementById("output").innerHTML = err;
					}
				);
			});
		</script>
	</body>
</html>`

func main() {
	addr := flag.String("addr", ":9000", "Server address")
	flag.Parse()

	log.Printf("starting server on %s", *addr)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	})

	http.ListenAndServe(*addr, mux)
}
