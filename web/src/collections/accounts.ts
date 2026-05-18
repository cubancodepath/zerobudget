import { queryCollectionOptions } from "@tanstack/query-db-collection";
import { type Collection, createCollection } from "@tanstack/react-db";
import { type Account, AccountSchema } from "../core/accounts/types";
import { listAccounts } from "../infra/http/account.api";
import { queryClient } from "../infra/tanstack-query/query-client";

export type AccountsCollection = Collection<Account, string, any>;

export const accountsCollection = createCollection(
	queryCollectionOptions({
		queryKey: ["accounts"],
		queryClient,
		schema: AccountSchema,
		getKey: (account) => account.id,
		queryFn: listAccounts,
	}),
);
