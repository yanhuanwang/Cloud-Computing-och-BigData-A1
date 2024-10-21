# Cloud-Computing-och-BigData-A1

curl -X 'POST' 'localhost:8081/add-expense' -H 'accept: */*' -H 'Content-Type: application/json' -d '{"username": "aa", "description": "aa", "amount": 110}'

curl -X 'GET' 'http://expense-service:8081/get-expenses?username=aa' -H 'accept: application/json'