# Cloud-Computing-och-BigData-A1

curl -X 'POST' 'expense-service:8081/add-expense' -H 'accept: */*' -H 'Content-Type: application/json' -d '{"username": "default_user", "description": "aa", "amount": 110}'

curl -X 'GET' 'http://expense-service:8081/get-expenses?username=default_user' -H 'accept: application/json'