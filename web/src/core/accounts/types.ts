import { z } from "zod"

export const AccountTypeSchema = z.enum(["cash", "checking", "savings", "credit_card"])

export const AccountSchema = z.object({
	id: z.string().uuid(),
	name: z.string(),
	type: AccountTypeSchema,
	currency_code: z.string(),
	initial_balance_cents: z.number().int(),
	is_active: z.boolean(),
	created_at: z.string().datetime({ offset: true }),
	updated_at: z.string().datetime({ offset: true }),
})

export const AccountSummarySchema = z.object({
	account_id: z.string().uuid(),
	account_name: z.string(),
	balance_cents: z.number().int(),
	reconciled_balance_cents: z.number().int(),
	difference_cents: z.number().int(),
})

export type AccountType = z.infer<typeof AccountTypeSchema>
export type Account = z.infer<typeof AccountSchema>
export type AccountSummary = z.infer<typeof AccountSummarySchema>
