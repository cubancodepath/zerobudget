import { z } from "zod"

export const PayeeSchema = z.object({
	id: z.string().uuid(),
	name: z.string(),
	createdAt: z.string().datetime(),
	updatedAt: z.string().datetime(),
})

export type Payee = z.infer<typeof PayeeSchema>
