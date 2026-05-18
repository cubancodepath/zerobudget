import type { Account, AccountSummary } from "./types"

export const isCreditCard = (account: Account): boolean =>
	account.type === "credit_card"

export const isNegativeBalance = (summary: AccountSummary): boolean =>
	summary.balance_cents < 0
