#!/bin/bash

echo "Testing Go In-Memory Database Implementation"
echo "============================================(^_^)"

echo -e "\n1. Testing Basic Operations (Scenario 1):"
echo "Expected: 10, NULL"
echo "SET ex 10
GET ex
UNSET ex
GET ex
END" | ./database

echo -e "\n2. Testing Value Counting (Scenario 2):"
echo "Expected: 2, 0, 1"
echo "SET a 10
SET b 10
NUMEQUALTO 10
NUMEQUALTO 20
SET b 30
NUMEQUALTO 10
END" | ./database

echo -e "\n3. Testing Basic Transactions (Scenario 1):"
echo "Expected: 10, 20, 10"
echo "SET a 10
GET a
BEGIN
SET a 20
GET a
ROLLBACK
GET a
END" | ./database

echo -e "\n4. Testing Nested Transactions (Scenario 2):"
echo "Expected: 40, NO TRANSACTION"
echo "BEGIN
SET a 30
BEGIN
SET a 40
COMMIT
GET a
ROLLBACK
END" | ./database

echo -e "\n5. Testing Complex Nested Transactions (Scenario 3):"
echo "Expected: 50, NULL, 60, 60"
echo "SET a 50
BEGIN
GET a
SET a 60
BEGIN
UNSET a
GET a
ROLLBACK
GET a
COMMIT
GET a
END" | ./database

echo -e "\n6. Testing NUMEQUALTO with Transactions (Scenario 4):"
echo "Expected: 1, 0, 1"
echo "SET a 10
BEGIN
NUMEQUALTO 10
BEGIN
UNSET a
NUMEQUALTO 10
ROLLBACK
NUMEQUALTO 10
COMMIT
END" | ./database

echo -e "\n7. Testing NO TRANSACTION Errors:"
echo "Expected: NO TRANSACTION"
echo "ROLLBACK
COMMIT
END" | ./database

echo -e "\n8. Testing Case Sensitivity:"
echo "Expected: 10, NULL"
echo "SET ex 10
GET ex
GET EX
END" | ./database

echo -e "\n9. Testing EOF Handling:"
echo "Expected: 10"
echo "SET a 10
GET a" | ./database

echo -e "\nAll tests completed!" 