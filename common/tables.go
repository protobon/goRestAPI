package common

var ProductTableCreationQuery = `CREATE TABLE IF NOT EXISTS product
(
   id SERIAL,
   name TEXT NOT NULL,
   price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
   CONSTRAINT product_pkey PRIMARY KEY (id)
)`

var CreditTableCreationQuery = `CREATE TABLE IF NOT EXISTS credit (
	id SERIAL,
	totalPrice INTEGER NOT NULL,
	feeAmount INTEGER NOT NULL,
	fees INTEGER NOT NULL,
	currentFee INTEGER NOT NULL,
	currentFeePaid BOOLEAN DEFAULT FALSE,
	purchaseDate TEXT NOT NULL,
	completed BOOLEAN DEFAULT FALSE,
	createdAt DATE NOT NULL,
	CONSTRAINT products_pkey PRIMARY KEY (id)
);`
