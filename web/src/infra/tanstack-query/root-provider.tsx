import { QueryClient } from "@tanstack/react-query";
import {
	accountsCollection,
	createAccountsCollection,
	createAccountsOfflineExecutor,
} from "#/collections/accounts";

export function getContext() {
	const queryClient = new QueryClient();
	const accountsCollection = createAccountsCollection(queryClient);
	const accountsOfflineExecutor =
		createAccountsOfflineExecutor(accountsCollection);

	return {
		queryClient,
		accountsCollection,
		accountsOfflineExecutor,
	};
}
export default function TanstackQueryProvider() {}
