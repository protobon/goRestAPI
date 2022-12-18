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
	card INT NOT NULL,
	totalPrice INTEGER NOT NULL,
	feeAmount INTEGER NOT NULL,
	fees INTEGER NOT NULL,
	currentFee INTEGER NOT NULL,
	currentFeePaid BOOLEAN DEFAULT FALSE,
	purchaseDate TEXT NOT NULL,
	completed BOOLEAN DEFAULT FALSE,
	createdAt DATE NOT NULL,
	CONSTRAINT products_pkey PRIMARY KEY (id),
	CONSTRAINT card_fk FOREIGN KEY (card) REFERENCES card(id)
);`

var CardTableCreationQuery = `CREATE TABLE IF NOT EXISTS card
(
   id SERIAL,
   type TEXT NOT NULL,
   closeDay INTEGER NOT NULL,
   CONSTRAINT card_pkey PRIMARY KEY (id)
)`
