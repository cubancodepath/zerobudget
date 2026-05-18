import type { Transaction } from "./types"

export const isExpense = (tx: Transaction): boolean => tx.amountCents < 0

export const isIncome = (tx: Transaction): boolean => tx.amountCents > 0

export const isUncategorized = (tx: Transaction): boolean => tx.categoryId === null
