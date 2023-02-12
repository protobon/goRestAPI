package database

// DummyTableCreate - Dummy
var DummyTableCreate = `CREATE TABLE IF NOT EXISTS dummy
(
   id SERIAL,
   name TEXT NOT NULL,
   price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
   created_at DATE NOT NULL,
   updated_at DATE,
   CONSTRAINT dummy_pkey PRIMARY KEY (id)
)`

// CreditCardTableCreate - Credit card
/*
	Fields:
   - id int
   - brand string: 'Visa', 'MasterCard', etc...
   - close_day int: Card closing day
   - amount_limit int: Maximum amount available to spend with this Card
   - created_at time
   - updated_at time
*/
var CreditCardTableCreate = `CREATE TABLE IF NOT EXISTS credit_card
(
   id SERIAL,
   brand TEXT NOT NULL,
   close_day INTEGER NOT NULL,
   amount_limit INTEGER NOT NULL,
   created_at DATE NOT NULL,
   updated_at DATE,
   CONSTRAINT credit_card_pkey PRIMARY KEY (id)
)`

// CreditTableCreate - Credit stores single purchase made with your CC, important fields are: purchaseDate, totalAmount, fees, currentFee
/*
	Fields:
   - id int
   - card int: FK credit_card(id)
   - total_price int: Total price spent on this credit
   - fee_amount int: Partial price (total/fees)
   - expired_amount int: Amount accumulated from previous unpaid fees
   - fees int: Number of fees (1...n)
   - current_fee int: Number of current fee
   - current_fee_paid bool: Is this month covered? (True | False)
   - purchase_date time: Credit's date of purchase
   - completed bool: Whole payment complete? (True | False)
   - created_at time
   - updated_at time
*/
var CreditTableCreate = `CREATE TABLE IF NOT EXISTS credit (
	id SERIAL,
	card INT NOT NULL,
	total_price INTEGER NOT NULL,
	fee_amount INTEGER NOT NULL,
	expired_amount INTEGER,
	fees INTEGER NOT NULL,
	current_fee INTEGER NOT NULL,
	current_fee_paid BOOLEAN DEFAULT FALSE,
	purchase_date TEXT NOT NULL,
	completed BOOLEAN DEFAULT FALSE,
	created_at DATE NOT NULL,
	updated_at DATE,
	CONSTRAINT credit_pk PRIMARY KEY (id),
	CONSTRAINT credit_card_fk FOREIGN KEY (card) REFERENCES credit_card(id)
);`

// RentTableCreate - Your Rent every month
/*
	Fields:
   - id int
   - amount int: Rent's value
   - due_date time: Date of payment
   - active bool: Renews each month (True | False)
   - created_at time
   - updated_at time
*/
var RentTableCreate = `CREATE TABLE IF NOT EXISTS rent
(
   id SERIAL,
   amount INTEGER NOT NULL,
   due_date DATE NOT NULL,
   active BOOLEAN NOT NULL DEFAULT TRUE,
   created_at DATE NOT NULL,
   updated_at DATE,
   CONSTRAINT rent_pkey PRIMARY KEY (id)
)`

// BillTableCreate - A bill that you pay each month. Example: Electricity Bill
/*
	Fields:
   - id int
   - type string: 'Electricity', 'Bus Ticket', etc...
   - amount int: Amount of the bill
   - currency string: 'UY, 'US', etc...
   - due_date time: Payment due date
   - created_at time
   - updated_at time
*/
var BillTableCreate = `CREATE TABLE IF NOT EXISTS bill
(
   id SERIAL,
   name TEXT NOT NULL,
   amount INTEGER NOT NULL,
   currency TEXT NOT NULL,
   due_date DATE NOT NULL,
   paid BOOLEAN NOT NULL DEFAULT FALSE,
   created_at DATE NOT NULL,
   updated_at DATE,
   CONSTRAINT bill_pkey PRIMARY KEY (id)
)`

// DebitTableCreate - Immediate payment purchase/expense
/*
	Fields:
   - id int
   - total int: The total amount of the expense
   - currency string: 'UY', 'US'
   - created_at time
   - updated_at time
*/
var DebitTableCreate = `CREATE TABLE IF NOT EXISTS debit
(
    id SERIAL,
    total INTEGER,
    card INT NOT NULL,
    created_at DATE NOT NULL,
    updated_at DATE,
    CONSTRAINT debit_pkey PRIMARY KEY (id),
    CONSTRAINT debit_card_fk FOREIGN KEY (card) REFERENCES debit_card(id)
)`

// DebitCardTableCreate - Your Salary
/*
	Fields:
   - id int
   - salary int: Your salary each month
   - currency string: 'UY', 'US'
   - renew bool: Tells if Salary renews each month (True: Active | False: Inactive)
   - money int: The total money you currently have
   - created_at time
   - updated_at time
*/
var DebitCardTableCreate = `CREATE TABLE IF NOT EXISTS debit_card
(
    id SERIAL,
    salary INTEGER,
    currency TEXT NOT NULL,
    renew BOOLEAN NOT NULL DEFAULT TRUE,
    money INTEGER,
    created_at DATE NOT NULL,
    updated_at DATE,
    CONSTRAINT debit_card_pkey PRIMARY KEY (id)
)`
