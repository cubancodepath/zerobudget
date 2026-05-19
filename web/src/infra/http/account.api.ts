import { z } from "zod";
import {
	type Account,
	AccountSchema,
	type AccountSummary,
	AccountSummarySchema,
	AccountTypeSchema,
} from "#/core/accounts/types";
import { http } from "./client";

const AccountsSchema = AccountSchema.array();
const AccountMutationSchema = z.object({
	name: z.string().min(1),
	type: AccountTypeSchema,
	currency_code: z.string().min(1),
	initial_balance_cents: z.number().int(),
	is_active: z.boolean().optional(),
});

export type AccountMutationInput = z.infer<typeof AccountMutationSchema>;

export async function listAccounts(): Promise<Account[]> {
	return http.get("accounts").json(AccountsSchema);
}

export async function getAccount(accountId: string): Promise<Account> {
	return http.get(`accounts/${accountId}`).json(AccountSchema);
}

export async function createAccount(input: any): Promise<Account> {
	return http
		.post("accounts", {
			json: AccountMutationSchema.parse(input),
		})
		.json(AccountSchema);
}

export async function updateAccount(
	accountId: string,
	input: any,
): Promise<Account> {
	return http
		.put(`accounts/${accountId}`, {
			json: AccountMutationSchema.parse(input),
		})
		.json(AccountSchema);
}

export async function deactivateAccount(accountId: string): Promise<void> {
	await http.delete(`accounts/${accountId}`);
}

export async function getAccountSummary(
	accountId: string,
): Promise<AccountSummary> {
	return http.get(`accounts/${accountId}/summary`).json(AccountSummarySchema);
}
