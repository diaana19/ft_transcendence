# Create a user and login to grab the token to be able to connect to the db

curl -k -X POST https://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "123",
    "email": "123test@test.com",
    "password": "Secret1s!2ds3",
    "dateOfBirth": "1995-06-15"
  }'

# 1. Login and grab token
TOKEN=$(curl -k -s -X POST https://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"123test@test.com","password":"Secret1s!2ds3"}' | grep -o '"token":"[^"]*' | cut -d'"' -f4)

# 2. Access protected route — should work
curl -k https://localhost:3000/api/users \
  -H "Authorization: Bearer $TOKEN"

# 3. Logout
curl -k -X POST https://localhost:3000/api/auth/logout \
  -H "Authorization: Bearer $TOKEN"

# 4. Try protected route again — should get 401
curl -k https://localhost:3000/api/users \
  -H "Authorization: Bearer $TOKEN"