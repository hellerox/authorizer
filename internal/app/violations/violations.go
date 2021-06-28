package violations

// Violations of the business rules are not considered to be errors as they are expected to happen}

const ViolationAccountAlreadyExists = "account-already-initialized"
const ViolationCardNotActive = "card-not-active"
const ViolationInsufficientLimit = "insufficient-limit"
const ViolationHighFrequencySmallInterval = "high-frequency-small-interval"
const ViolationDoubledTransaction = "doubled-transaction"
