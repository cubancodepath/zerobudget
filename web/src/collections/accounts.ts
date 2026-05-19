import {
	createBrowserWASQLitePersistence,
	openBrowserWASQLiteOPFSDatabase,
	persistedCollectionOptions,
} from "@tanstack/browser-db-sqlite-persistence";
import {
	type OfflineExecutor,
	startOfflineExecutor,
} from "@tanstack/offline-transactions";
import { queryCollectionOptions } from "@tanstack/query-db-collection";
import { type Collection, createCollection } from "@tanstack/react-db";
import type { QueryClient } from "@tanstack/react-query";
import { type Account, AccountSchema } from "#/core/accounts/types";
import * as api from "#/infra/http/account.api";
import { queryClient } from "#/infra/tanstack-query/query-client";

export type AccountsCollection = Collection<Account, string, any>;
export type AccountOfflineExecutor = OfflineExecutor;

export const browserPersistence = import.meta.env.SSR
	? null
	: await (async () => {
			const database = await openBrowserWASQLiteOPFSDatabase({
				databaseName: "zerobudget.sqlite",
			});

			return createBrowserWASQLitePersistence({
				database,
			});
		})();

export const accountsCollection = createCollection(
	queryCollectionOptions({
		queryKey: ["accounts"],
		queryClient,
		schema: AccountSchema,
		getKey: (account) => account.id,
		queryFn: api.listAccounts,
	}),
);

export function createAccountsOfflineExecutor(collection: AccountsCollection) {
	const offline = startOfflineExecutor({
		collections: { accounts: collection },
		mutationFns: {
			createAccount: async ({ transaction }) => {
				const mutation = transaction.mutations[0];

				const created = await api.createAccount({
					name: mutation.modified.name,
					type: mutation.modified.type,
					currency_code: mutation.modified.currency_code,
					initial_balance_cents: mutation.modified.initial_balance_cents,
					is_active: mutation.modified.is_active,
				});

				collection.utils.writeUpsert(created);
			},

			updateAccount: async ({ transaction }) => {
				const mutation = transaction.mutations[0];
				const next = mutation.modified as Account;

				const updated = await api.updateAccount(next.id, {
					name: mutation.modified.name,
					type: mutation.modified.type,
					currency_code: mutation.modified.currency_code,
					initial_balance_cents: mutation.modified.initial_balance_cents,
					is_active: mutation.modified.is_active,
				});

				collection.utils.writeUpsert(updated);
			},

			deleteAccount: async ({ transaction }) => {
				const mutation = transaction.mutations[0];
				await api.deactivateAccount(String(mutation.key));
				collection.utils.writeDelete(mutation.key);
			},
		},
	});

	void offline.waitForInit();

	return offline;
}

export function createAccountsCollection(queryClient: QueryClient) {
	const queryOptions = queryCollectionOptions({
		queryKey: ["accounts"],
		queryClient,
		schema: AccountSchema,
		getKey: (account) => account.id,
		queryFn: api.listAccounts,
	});

	if (!browserPersistence) {
		return createCollection(queryOptions);
	}

	const persitenceOptions = persistedCollectionOptions({
		id: "accounts",
		schemaVersion: 1,
		persistence: browserPersistence,
		...queryOptions,
	}) as any;

	return createCollection(persitenceOptions) as unknown as AccountsCollection;
}
