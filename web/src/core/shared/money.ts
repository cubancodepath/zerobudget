import { z } from "zod"

export const MoneySchema = z.number().int()

export type Money = z.infer<typeof MoneySchema>

export const toCents = (amount: number): Money => Math.round(amount * 100)
export const fromCents = (cents: Money): number => cents / 100
