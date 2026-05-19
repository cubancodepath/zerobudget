import { Avatar, Button } from "@heroui/react";
import { Sidebar } from "@heroui-pro/react";
import { PlusCircle, Settings } from "lucide-react";
import { CreateAccountModal } from "#/components/accounts/create-account-modal";
import { useModal } from "#/components/modal";
import type { Account } from "#/core/accounts/types";

const STATIC_USER = {
	name: "Bryan Valmaseda",
	email: "bjvalmaseda.g@gmail.com",
	initials: "BV",
};

function formatCents(cents: number) {
	return new Intl.NumberFormat("en-US", {
		style: "currency",
		currency: "USD",
		minimumFractionDigits: 0,
		maximumFractionDigits: 0,
	}).format(cents / 100);
}

interface AppSidebarProps {
	accounts: Account[];
	isReady: boolean;
}

export function AppSidebar({ accounts, isReady }: AppSidebarProps) {
	const modal = useModal();

	return (
		<Sidebar>
			<Sidebar.Header>
				<div className="flex items-center gap-2.5 px-1 py-0.5">
					<div className="flex size-8 shrink-0 items-center justify-center rounded-xl bg-[--accent]">
						<span className="text-sm font-bold text-[--accent-foreground]">
							ZB
						</span>
					</div>
					<span className="truncate text-sm font-semibold tracking-tight">
						ZeroBudget
					</span>
				</div>
			</Sidebar.Header>

			<Sidebar.Content>
				<Sidebar.Group>
					<Sidebar.GroupLabel>Accounts</Sidebar.GroupLabel>
					<Sidebar.Menu>
						{accounts.map((account) => (
							<Sidebar.MenuItem
								key={account.id}
								href={`/accounts/${account.id}`}
								tooltip={account.name}
							>
								<Sidebar.MenuLabel>{account.name}</Sidebar.MenuLabel>
								<Sidebar.MenuChip
									className={
										account.initial_balance_cents < 0 ? "text-[--danger]" : ""
									}
								>
									{formatCents(account.initial_balance_cents)}
								</Sidebar.MenuChip>
							</Sidebar.MenuItem>
						))}
					</Sidebar.Menu>

					<div className="px-2 pt-2">
						<Button
							fullWidth
							variant="outline"
							size="sm"
							className="justify-start gap-2"
							onPress={() => {
								modal.open({
									isKeyboardDismissDisabled: true,
									render: ({ close }) => <CreateAccountModal onClose={close} />,
								});
							}}
							isDisabled={!isReady}
						>
							<PlusCircle size={14} />
							Add Account
						</Button>
					</div>
				</Sidebar.Group>
			</Sidebar.Content>

			<Sidebar.Footer>
				<Sidebar.Tooltip content={STATIC_USER.name} placement="right">
					<div className="flex cursor-pointer items-center gap-2.5 rounded-xl p-2 transition-colors hover:bg-[--default]">
						<Avatar
							aria-label={STATIC_USER.name}
							className="shrink-0"
							size="sm"
						>
							<Avatar.Fallback>{STATIC_USER.initials}</Avatar.Fallback>
						</Avatar>
						<div className="flex min-w-0 flex-1 flex-col overflow-hidden">
							<span className="truncate text-sm font-medium leading-tight">
								{STATIC_USER.name}
							</span>
							<span className="truncate text-xs text-[--muted]">
								{STATIC_USER.email}
							</span>
						</div>
						<Button
							isIconOnly
							aria-label="Settings"
							size="sm"
							variant="ghost"
							className="size-6 min-w-0 shrink-0 rounded-lg"
						>
							<Settings size={14} />
						</Button>
					</div>
				</Sidebar.Tooltip>
			</Sidebar.Footer>

			<Sidebar.Rail />
		</Sidebar>
	);
}
