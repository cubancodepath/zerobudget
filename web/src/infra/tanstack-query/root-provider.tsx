import { accountsCollection } from "#/collections/accounts";
import { queryClient } from "./query-client";

export function getContext() {
	return {
		queryClient,
		accountsCollection,
	};
}
export default function TanstackQueryProvider() {}
