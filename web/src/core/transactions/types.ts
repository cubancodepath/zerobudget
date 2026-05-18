import { z } from "zod"

export const TransactionSchema = z.object({
	id: z.string().uuid(),
	accountId: z.string().uuid(),
	categoryId: z.string().uuid().nullable(),
	payeeId: z.string().uuid().nullable(),
	amountCents: z.number().int(),
	transactionDate: z.string(),
	isReconciled: z.boolean(),
	note: z.string().nullable(),
	createdAt: z.string().datetime(),
	updatedAt: z.string().datetime(),
})

export type Transaction = z.infer<typeof TransactionSchema>
