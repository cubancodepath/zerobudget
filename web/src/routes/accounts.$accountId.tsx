import { Button, Modal } from "@heroui/react";
import { useLiveQuery } from "@tanstack/react-db";
import { createFileRoute, useRouter } from "@tanstack/react-router";
import { useEffect, useMemo } from "react";
import { EditAccountModal } from "#/components/accounts/edit-account-modal";
import { useModal } from "#/components/modal";
import type { AccountMutationInput } from "#/infra/http/account.api";

function formatCents(cents: number, currencyCode: string) {
	return new Intl.NumberFormat("en-US", {
		style: "currency",
		currency: currencyCode,
		minimumFractionDigits: 2,
		maximumFractionDigits: 2,
	}).format(cents / 100);
}

function formatDate(value: string) {
	return new Intl.DateTimeFormat("en-US", {
		year: "numeric",
		month: "short",
		day: "2-digit",
		hour: "2-digit",
		minute: "2-digit",
	}).format(new Date(value));
}

export const Route = createFileRoute("/accounts/$accountId")({
	component: AccountDetailPage,
});

function AccountDetailPage() {
	const { accountId } = Route.useParams();
	const { accountsCollection, accountsOfflineExecutor } = Route.useRouteContext();
	const modal = useModal();
	const router = useRouter();

	const updateAccountAction = useMemo(
		() =>
			accountsOfflineExecutor.createOfflineAction<{
				accountId: string;
				input: AccountMutationInput;
			}>({
				mutationFnName: "updateAccount",
				onMutate: ({ accountId, input }) => {
					accountsCollection.update(accountId, (draft) => {
						draft.name = input.name;
						draft.type = input.type;
						draft.currency_code = input.currency_code;
						draft.initial_balance_cents = input.initial_balance_cents;
						draft.is_active = input.is_active ?? true;
						draft.updated_at = new Date().toISOString();
					});
				},
			}),
		[accountsCollection, accountsOfflineExecutor],
	);

	const deleteAccountAction = useMemo(
		() =>
			accountsOfflineExecutor.createOfflineAction<{ accountId: string }>({
				mutationFnName: "deleteAccount",
				onMutate: ({ accountId }) => {
					accountsCollection.utils.writeDelete(accountId);
				},
			}),
		[accountsCollection, accountsOfflineExecutor],
	);
	const { data: accounts = [], isLoading, isError } = useLiveQuery((query) =>
		query.from({ accounts: accountsCollection }),
	);

	const account = accounts.find((item) => item.id === accountId);

	useEffect(() => {
		if (!isLoading && !account) {
			void router.navigate({ to: "/" });
		}
	}, [account, isLoading, router]);

	if (isLoading) {
		return <div className="p-4 text-sm text-default-500">Loading account...</div>;
	}

	if (isError) {
		return (
			<div className="p-4 text-sm text-danger">
				Something failed while loading this account.
			</div>
		);
	}

	if (!account) {
		return <div className="p-4 text-sm text-default-500">Redirecting...</div>;
	}

	return (
		<div className="space-y-4 p-4">
			<div className="flex items-start justify-between gap-3">
				<div>
					<h1 className="text-2xl font-semibold">{account.name}</h1>
					<p className="text-sm text-default-500">{account.type}</p>
				</div>
				<div className="flex items-center gap-2">
					<Button
						variant="outline"
						size="sm"
						onPress={() => {
							modal.open({
								render: ({ close }) => (
									<EditAccountModal
										account={account}
										onSubmit={(input) => {
											updateAccountAction({ accountId: account.id, input });
										}}
										onClose={close}
									/>
								),
							});
						}}
					>
						Edit account
					</Button>
					<Button
						variant="danger"
						size="sm"
						onPress={() => {
							modal.open({
								render: ({ close }) => (
									<Modal.Dialog className="sm:max-w-md">
										<Modal.CloseTrigger />
										<Modal.Header>
											<Modal.Heading>Delete account</Modal.Heading>
										</Modal.Header>
										<Modal.Body>
											<p className="text-sm text-default-600">
												Are you sure you want to delete <b>{account.name}</b>? This
												action cannot be undone.
											</p>
										</Modal.Body>
										<Modal.Footer>
											<Button variant="secondary" onPress={close}>
												Cancel
											</Button>
											<Button
												variant="danger"
												onPress={() => {
													deleteAccountAction({ accountId: account.id });
													close();
													void router.navigate({ to: "/" });
												}}
											>
												Delete
											</Button>
										</Modal.Footer>
									</Modal.Dialog>
								),
							});
						}}
					>
						Delete
					</Button>
				</div>
			</div>

			<div className="grid gap-3 sm:grid-cols-2">
				<div className="rounded-lg border border-default-200 p-3">
					<p className="text-xs text-default-500">Currency</p>
					<p className="text-sm font-medium">{account.currency_code}</p>
				</div>
				<div className="rounded-lg border border-default-200 p-3">
					<p className="text-xs text-default-500">Initial balance</p>
					<p className="text-sm font-medium">
						{formatCents(account.initial_balance_cents, account.currency_code)}
					</p>
				</div>
				<div className="rounded-lg border border-default-200 p-3">
					<p className="text-xs text-default-500">Status</p>
					<p className="text-sm font-medium">
						{account.is_active ? "Active" : "Inactive"}
					</p>
				</div>
				<div className="rounded-lg border border-default-200 p-3">
					<p className="text-xs text-default-500">Created at</p>
					<p className="text-sm font-medium">{formatDate(account.created_at)}</p>
				</div>
				<div className="rounded-lg border border-default-200 p-3 sm:col-span-2">
					<p className="text-xs text-default-500">Updated at</p>
					<p className="text-sm font-medium">{formatDate(account.updated_at)}</p>
				</div>
			</div>
		</div>
	);
}
