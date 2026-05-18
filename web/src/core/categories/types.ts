import { z } from "zod"

export const CategorySchema = z.object({
	id: z.string().uuid(),
	name: z.string(),
	createdAt: z.string().datetime(),
	updatedAt: z.string().datetime(),
})

export type Category = z.infer<typeof CategorySchema>
