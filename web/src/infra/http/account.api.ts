import { type Account, AccountSchema } from "#/core/accounts/types";
import { http } from "./client";

const AccountsSchema = AccountSchema.array();

export async function listAccounts(): Promise<Account[]> {
	return http.get("accounts").json(AccountsSchema);
}
