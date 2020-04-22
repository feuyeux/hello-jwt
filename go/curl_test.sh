echo ":::: 1 Test Login ::::"
curl -s :8080/login -d'{
"username":"username",
"password":"password"
}'| jq
echo
read RT AT < <(echo $(curl -s :8080/login -d'{
"username":"username",
"password":"password"
}' | jq -r ".refresh_token,.access_token"))
echo "Access Token  :$AT"
echo "Refresh Token :$RT"
echo
echo ":::: 2 Test Request api ::::"
curl :8080/api -H "authorization:bearer $AT" -d'{"title":"Hello JWT!"}'
echo
echo
echo ":::: 3 Test Logout $ Request api ::::"
curl :8080/logout -H "authorization:bearer $AT" -XPOST
echo
echo
curl :8080/api -H "authorization:bearer $AT" -d'{"title":"Hello JWT!"}'
echo
echo
echo ":::: 4 Test Refresh ::::"
curl -s :8080/refresh -d'{"refresh_token":"'$RT'"}' | jq
echo "Done!"