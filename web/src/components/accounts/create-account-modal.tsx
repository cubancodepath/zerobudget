import { Button, Modal } from "@heroui/react";
import { type ChangeEvent, type FormEvent, useState } from "react";
import { Route } from "#/routes/__root";

type CreateAccountModalProps = {
	onClose: () => void;
};

export function CreateAccountModal({ onClose }: CreateAccountModalProps) {
	const { accountsOfflineExecutor, accountsCollection } =
		Route.useRouteContext();
	const [name, setName] = useState("");
	const [type, setType] = useState("checking");
	const [currencyCode, setCurrencyCode] = useState("USD");
	const [initialAmount, setInitialAmount] = useState("0");
	const [isActive, setIsActive] = useState(true);

	const onInputChange =
		(setter: (value: string) => void) =>
		(event: ChangeEvent<HTMLInputElement>) => {
			setter(event.target.value);
		};

	const onSubmit = async (event: FormEvent<HTMLFormElement>) => {
		event.preventDefault();

		const tx = accountsOfflineExecutor.createOfflineTransaction({
			mutationFnName: "createAccount",
		});

		tx.mutate(() => {
			accountsCollection.insert({
				id: crypto.randomUUID() as any,
				name,
				type: type as "checking" | "savings",
				initial_balance_cents: Math.round(Number(initialAmount || 0) * 100),
				currency_code: currencyCode,
				is_active: true,
				created_at: new Date().toISOString(),
				updated_at: new Date().toISOString(),
			});
		});

		onClose();
	};

	return (
		<Modal.Dialog className="sm:max-w-md">
			<Modal.CloseTrigger />
			<Modal.Header>
				<Modal.Heading>Create account</Modal.Heading>
			</Modal.Header>
			<Modal.Body>
				<form
					id="create-account-form"
					className="space-y-3"
					onSubmit={onSubmit}
				>
					<label className="flex flex-col gap-1 text-sm">
						<span className="text-default-600">Name</span>
						<input
							className="h-10 rounded-medium border border-default-200 bg-transparent px-3"
							placeholder="Main checking"
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
							onChange={(event) => setType(event.target.value)}
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
				<Button variant="primary" type="submit" form="create-account-form">
					Create
				</Button>
			</Modal.Footer>
		</Modal.Dialog>
	);
}
