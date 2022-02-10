# payment_api
Application that will allow a merchant to offer a way for their shoppers to pay for their product.

It simulates the following payment flow:
- A merchant requests an authorisation through the /authorize call. This call contains
the customer credit card data as well as an amount and currency. It will return a
unique ID that will be used in all next API calls.

- The /void call will cancel the whole transaction without billing the customer. No
further action is possible once a transaction is voided.

- The /capture call will capture the money on the customer bank. It can be called
multiple times with an amount that should not be superior to the amount authorised in the first call. For example, a 10£ authorisation can be captured 2 times with a 4£ and 6£ order.

- The /refund call will refund the money taken from the customer bank account. It can be also called multiple times with an amount that can not be superior to the captured amount. For example, a 10£ authorisation, with a 5£ capture, can only be refunded of 5£. Once a refund has occurred, a capture cannot be made on the specific transaction.

## Important assumptions
Main concern for this implementation was to how display a account balance to a user. I believe there should be distinction
between displayed account balance, and actual account balance, because after the transaction is authorized, money are not yet captured.

To complete the implementation I have made following assumptions:
- Any authorized transaction is blocking that amount from account balance. 


**Reason:** Without subtracting authorized transaction from account balance, we could create infinite ammount of transactions. I believe it is not very realistic. I assume authorizing a transaction behaves similarly to blocking that amount of money on the account.

- Actual account balance is modified only by `capture` and `refund` methods.


Those assumptions lead to this simple algorithm to display account balance:
    
    DisplayedBalance = ActualBalance - Transactions - Refunds + Captures

We are adding Captures and subtracting Refunds to nullify their effect on displayed balance as per requirements they are allowed only within authorized transaction range. Actual balance will be effected by Capture and Refund operations.

### **Example**
    ActualBalance: 10000PLN, DisplayedBalance: 10000PLN

    1. Authorize 3000PLN

    ActualBalance: 10000PLN, DisplayedBalance: 7000PLN

    2. Capture 2000PLN

    ActualBalance: 8000PLN, DisplayedBalance: 7000PLN

    3. Refund 1000PLN

    ActualBalance: 9000PLN, DisplayedBalance: 7000PLN
## Running the application

Preffered way to quickly start the app is through the Makefile




    make quicksetup

It will:
1. Run Postgres in docker container hosting on port `5432`
2. Compile and run payment_api server, hosting on `localhost:3000`git s

## Documentation
After creating quicksetup I advice to get familiar with `swagger.yaml`, which contains description of the API

Most convinient way would be to open `localhost:3000/swagger` where Swagger UI is hosted for maximum convinience

## Tests
Next good step to get familiar with the payment_api would be to read created testcases. They are written using Ginkgo
framework which is very descriptive in its nature, therefore informative for the reader.

All tests are located in `/handlers` directory

To run tests execute:

    make tests