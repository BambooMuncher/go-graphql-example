package main

import (
    "fmt"
    "net/http"
    "context"
    "strings"
    "github.com/graphql-go/handler"
)

func main() {
    h := handler.New(&handler.Config{
        Schema: &schema,
        Pretty: true,
        RootObjectFn: func(ctx context.Context, r *http.Request) map[string]interface{} {
            return map[string]interface{}{
                "token": strings.TrimSpace(strings.Replace(r.Header.Get("Authorization"), "Bearer", "", 1)),
            }
        },
    })

    http.Handle("/graphql", h)
    http.Handle("/sandbox", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write(sandboxHTML)
    }))

    fmt.Println("Server running on port 8080")
    http.ListenAndServe(":8080", nil)
}

var sandboxHTML = []byte(`
<!DOCTYPE html>
<html lang="en">
<body style="margin: 0; overflow-x: hidden; overflow-y: hidden">
<div id="sandbox" style="height:100vh; width:100vw;"></div>
<script src="https://embeddable-sandbox.cdn.apollographql.com/_latest/embeddable-sandbox.umd.production.min.js"></script>
<script>
new window.EmbeddedSandbox({
  target: "#sandbox",
  // Pass through your server href if you are embedding on an endpoint.
  // Otherwise, you can pass whatever endpoint you want Sandbox to start up with here.
  initialEndpoint: "http://localhost:8080/graphql",
});
// advanced options: https://www.apollographql.com/docs/studio/explorer/sandbox#embedding-sandbox
</script>
</body>
 
</html>`)