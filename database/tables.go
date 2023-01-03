package database

// DummyTableCreate - Dummy
var DummyTableCreate = `CREATE TABLE IF NOT EXISTS dummy
(
   id SERIAL,
   name TEXT NOT NULL,
   price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
   createdAt DATE NOT NULL,
   updatedAt DATE,
   CONSTRAINT dummy_pkey PRIMARY KEY (id)
)`

// CreditCardTableCreate - Credit card
/*
	Fields:
   - Id int
   - Brand string: 'Visa', 'MasterCard', etc...
   - CloseDay int: Card closing day
   - AmountLimit int: Maximum amount available to spend with this Card
   - CreatedAt time
   - UpdatedAt time
*/
var CreditCardTableCreate = `CREATE TABLE IF NOT EXISTS credit_card
(
   id SERIAL,
   brand TEXT NOT NULL,
   closeDay INTEGER NOT NULL,
   amountLimit INTEGER NOT NULL,
   createdAt DATE NOT NULL,
   updatedAt DATE,
   CONSTRAINT credit_card_pkey PRIMARY KEY (id)
)`

// CreditTableCreate - Credit stores single purchase made with your CC, important fields are: purchaseDate, totalAmount, fees, currentFee
/*
	Fields:
   - Id int
   - Card int: FK credit_card(id)
   - TotalPrice int: Total price spent on this credit
   - FeeAmount int: Partial price (total/fees)
   - ExpiredAmount int: Amount accumulated from previous unpaid fees
   - Fees int: Number of fees (1...n)
   - CurrentFee int: Number of current fee
   - CurrentFeePaid bool: Is this month covered? (True | False)
   - PurchaseDate time: Credit's date of purchase
   - Completed bool: Whole payment complete? (True | False)
   - CreatedAt time
   - UpdatedAt time
*/
var CreditTableCreate = `CREATE TABLE IF NOT EXISTS credit (
	id SERIAL,
	card INT NOT NULL,
	totalPrice INTEGER NOT NULL,
	feeAmount INTEGER NOT NULL,
	expiredAmount INTEGER,
	fees INTEGER NOT NULL,
	currentFee INTEGER NOT NULL,
	currentFeePaid BOOLEAN DEFAULT FALSE,
	purchaseDate TEXT NOT NULL,
	completed BOOLEAN DEFAULT FALSE,
	createdAt DATE NOT NULL,
	updatedAt DATE,
	CONSTRAINT credit_pk PRIMARY KEY (id),
	CONSTRAINT credit_card_fk FOREIGN KEY (card) REFERENCES credit_card(id)
);`

// RentTableCreate - Your Rent every month
/*
	Fields:
   - Id int
   - Amount int: Rent's value
   - DueDate time: Date of payment
   - Active bool: Renews each month (True | False)
   - CreatedAt time
   - UpdatedAt time
*/
var RentTableCreate = `CREATE TABLE IF NOT EXISTS rent
(
   id SERIAL,
   amount INTEGER NOT NULL,
   dueDate DATE NOT NULL,
   active BOOLEAN NOT NULL DEFAULT TRUE,
   createdAt DATE NOT NULL,
   updatedAt DATE,
   CONSTRAINT rent_pkey PRIMARY KEY (id)
)`

// BillTableCreate - A bill that you pay each month. Example: Electricity Bill
/*
	Fields:
   - Id int
   - Type string: 'Electricity', 'Bus Ticket', etc...
   - Amount int: Amount of the bill
   - Currency string: 'UY, 'US', etc...
   - DueDate time: Payment due date
   - CreatedAt time
   - UpdatedAt time
*/
var BillTableCreate = `CREATE TABLE IF NOT EXISTS bill
(
   id SERIAL,
   name TEXT NOT NULL,
   amount INTEGER NOT NULL,
   currency TEXT NOT NULL,
   dueDate DATE NOT NULL,
   paid BOOLEAN NOT NULL DEFAULT FALSE,
   createdAt DATE NOT NULL,
   updatedAt DATE,
   CONSTRAINT bill_pkey PRIMARY KEY (id)
)`

// DebitTableCreate - Immediate payment purchase/expense
/*
	Fields:
   - Id int
   - Total int: The total amount of the expense
   - Currency string: 'UY', 'US'
   - CreatedAt time
   - UpdatedAt time
*/
var DebitTableCreate = `CREATE TABLE IF NOT EXISTS debit
(
    id SERIAL,
    total INTEGER,
    card INT NOT NULL,
    createdAt DATE NOT NULL,
    updatedAt DATE,
    CONSTRAINT debit_pkey PRIMARY KEY (id),
    CONSTRAINT debit_card_fk FOREIGN KEY (card) REFERENCES debit_card(id)
)`

// DebitCardTableCreate - Your Salary
/*
	Fields:
   - Id int
   - Salary int: Your salary each month
   - Currency string: 'UY', 'US'
   - Renew bool: Tells if Salary renews each month (True: Active | False: Inactive)
   - Money int: The total money you currently have
   - CreatedAt time
   - UpdatedAt time
*/
var DebitCardTableCreate = `CREATE TABLE IF NOT EXISTS debit_card
(
    id SERIAL,
    salary INTEGER,
    currency TEXT NOT NULL,
    renew BOOLEAN NOT NULL DEFAULT TRUE,
    money INTEGER,
    createdAt DATE NOT NULL,
    updatedAt DATE,
    CONSTRAINT debit_card_pkey PRIMARY KEY (id)
)`
