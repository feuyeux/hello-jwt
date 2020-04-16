# https://stedolan.github.io/jq/
# jq --version
echo "Generate a JSON Web Token:"
TOKEN=$(curl -s localhost:8080/authenticate -d'{
              "username": "hello_man",
              "password": "password"
            }' -XPOST -H "Content-Type: application/json" | jq --raw-output '.token' )
echo "TOKEN=$TOKEN"
echo
echo "Validate the JSON Web Token(Web):"
echo "localhost:8080/hello:"
curl -H "Authorization: Bearer $TOKEN" "localhost:8080/hello"
echo
echo
echo "Validate the JSON Web Token(Webflux):"
echo "localhost:8080/hello-stream:"
curl -H "Authorization: Bearer $TOKEN" "localhost:8080/hello-stream"