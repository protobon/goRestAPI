package database

var CreditNewRecord = `INSERT INTO credit(
                   card,
                   totalPrice,
                   feeAmount,
                   fees,
                   currentFee,
                   currentFeePaid,
                   purchaseDate,
                   completed,
                   createdAt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
