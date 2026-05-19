import { Button, Modal } from "@heroui/react";
import { type ChangeEvent, type FormEvent, useState } from "react";
import type { Account } from "#/core/accounts/types";
import type { AccountMutationInput } from "#/infra/http/account.api";

type EditAccountModalProps = {
	account: Account;
	onSubmit: (input: AccountMutationInput) => void;
	onClose: () => void;
};

export function EditAccountModal({
	account,
	onSubmit,
	onClose,
}: EditAccountModalProps) {
	const [name, setName] = useState(account.name);
	const [type, setType] = useState(account.type);
	const [currencyCode, setCurrencyCode] = useState(account.currency_code);
	const [initialAmount, setInitialAmount] = useState(
		String(account.initial_balance_cents / 100),
	);
	const [isActive, setIsActive] = useState(account.is_active);

	const onInputChange =
		(setter: (value: string) => void) =>
		(event: ChangeEvent<HTMLInputElement>) => {
			setter(event.target.value);
		};

	const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
		event.preventDefault();

		onSubmit({
			name,
			type,
			currency_code: currencyCode,
			initial_balance_cents: Math.round(Number(initialAmount || 0) * 100),
			is_active: isActive,
		});

		onClose();
	};

	return (
		<Modal.Dialog className="sm:max-w-md">
			<Modal.CloseTrigger />
			<Modal.Header>
				<Modal.Heading>Edit account</Modal.Heading>
			</Modal.Header>
			<Modal.Body>
				<form id="edit-account-form" className="space-y-3" onSubmit={handleSubmit}>
					<label className="flex flex-col gap-1 text-sm">
						<span className="text-default-600">Name</span>
						<input
							className="h-10 rounded-medium border border-default-200 bg-transparent px-3"
							value={name}
							onChange={onInputChange(setName)}
							required
						/>
					</label>
					<label className="flex flex-col gap-1 text-sm">
						<span className="text-default-600">Type</span>
						<select
							className="h-10 rounded-medium border border-default-200 bg-transparent px-3"
							value={type}
							onChange={(event) => setType(event.target.value as Account["type"])}
						>
							<option value="cash">Cash</option>
							<option value="checking">Checking</option>
							<option value="savings">Savings</option>
							<option value="credit_card">Credit card</option>
						</select>
					</label>
					<label className="flex flex-col gap-1 text-sm">
						<span className="text-default-600">Currency code</span>
						<input
							className="h-10 rounded-medium border border-default-200 bg-transparent px-3"
							value={currencyCode}
							onChange={onInputChange(setCurrencyCode)}
							required
						/>
					</label>
					<label className="flex flex-col gap-1 text-sm">
						<span className="text-default-600">Initial amount</span>
						<input
							className="h-10 rounded-medium border border-default-200 bg-transparent px-3"
							type="number"
							step="0.01"
							value={initialAmount}
							onChange={onInputChange(setInitialAmount)}
							required
						/>
					</label>
					<label className="flex items-center gap-2 text-sm">
						<input
							type="checkbox"
							checked={isActive}
							onChange={(event) => setIsActive(event.target.checked)}
						/>
						<span>Active account</span>
					</label>
				</form>
			</Modal.Body>
			<Modal.Footer>
				<Button variant="secondary" onPress={onClose}>
					Cancel
				</Button>
				<Button variant="primary" type="submit" form="edit-account-form">
					Save changes
				</Button>
			</Modal.Footer>
		</Modal.Dialog>
	);
}
