package database

var ProductTableCreate = `CREATE TABLE IF NOT EXISTS product
(
   id SERIAL,
   name TEXT NOT NULL,
   price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
   CONSTRAINT product_pkey PRIMARY KEY (id)
)`

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
	CONSTRAINT products_pkey PRIMARY KEY (id),
	CONSTRAINT card_fk FOREIGN KEY (card) REFERENCES card(id)
);`

var CardTableCreate = `CREATE TABLE IF NOT EXISTS card
(
   id SERIAL,
   type TEXT NOT NULL,
   closeDay INTEGER NOT NULL,
   CONSTRAINT card_pkey PRIMARY KEY (id)
)`

var FixedServiceTableCreate = `CREATE TABLE IF NOT EXISTS service_fixed
(
   id SERIAL,
   name TEXT NOT NULL,
   amount INTEGER NOT NULL,
   dueDate DATE NOT NULL,
   month INTEGER,
   year INTEGER,
   active BOOLEAN NOT NULL DEFAULT TRUE,
   CONSTRAINT card_pkey PRIMARY KEY (id)
)`

var VariableServiceTableCreate = `CREATE TABLE IF NOT EXISTS service_variable

(
   id SERIAL,
   name TEXT NOT NULL,
   amount INTEGER NOT NULL,
   dueDate DATE NOT NULL,
   month INTEGER,
   year INTEGER,
   CONSTRAINT card_pkey PRIMARY KEY (id)
)`
