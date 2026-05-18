import { z } from "zod"

export const AccountTypeSchema = z.enum(["cash", "checking", "savings", "credit_card"])

export const AccountSchema = z.object({
	id: z.string().uuid(),
	name: z.string(),
	type: AccountTypeSchema,
	currencyCode: z.string(),
	initialBalanceCents: z.number().int(),
	isActive: z.boolean(),
	createdAt: z.string().datetime(),
	updatedAt: z.string().datetime(),
})

export const AccountSummarySchema = z.object({
	accountId: z.string().uuid(),
	accountName: z.string(),
	balanceCents: z.number().int(),
	reconciledBalanceCents: z.number().int(),
	differenceCents: z.number().int(),
})

export type AccountType = z.infer<typeof AccountTypeSchema>
export type Account = z.infer<typeof AccountSchema>
export type AccountSummary = z.infer<typeof AccountSummarySchema>
