package database

var CreditNewRecord = `INSERT INTO credit(
                   card,
                   totalPrice,
                   feeAmount,
                   expiredAmount,
                   fees,
                   currentFee,
                   currentFeePaid,
                   purchaseDate,
                   completed,
                   createdAt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

var GetCreditToPay = `SELECT fees, currentFee, currentFeePaid, completed FROM credit WHERE id=$1`
var CreditCompletePayment = `UPDATE credit SET currentFeePaid=true, completed=true WHERE id=$1`
var PayCredit = `UPDATE credit SET currentFeePaid=true WHERE id=$1`

var DueCreditAllTime = `SELECT * FROM credit JOIN card ON credit.card = card.id 
         WHERE credit.completed=false`

var DueCreditThisMonth = `SELECT * FROM credit WHERE currentFeePaid=false
                       AND completed=false AND currentFee!=0`

var CreditDebtThisMonth = `SELECT feeAmount, expiredAmount FROM credit
                                WHERE currentFeePaid!=true AND currentFee!=0`

var AllCredits = `SELECT * FROM credit WHERE completed=false`

var AddToExpired = `UPDATE credit SET expiredAmount=$1 WHERE id=$2`

var NextQuota = `UPDATE credit SET currentFee=currentFee+1, currentFeePaid=false
              WHERE completed!=true AND currentFee<fees`
