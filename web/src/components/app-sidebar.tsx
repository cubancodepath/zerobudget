import { Sidebar } from "@heroui-pro/react";
import { Avatar, Button } from "@heroui/react";
import { Plus, Settings } from "lucide-react";
import type { AccountType } from "#/core/accounts/types";

const STATIC_USER = {
	name: "Bryan Valmaseda",
	email: "bjvalmaseda.g@gmail.com",
	initials: "BV",
};

const STATIC_ACCOUNTS = [
	{ id: "1", name: "Checking", type: "checking" as AccountType, balanceCents: 245000 },
	{ id: "2", name: "Savings", type: "savings" as AccountType, balanceCents: 820000 },
	{ id: "3", name: "Chase Credit", type: "credit_card" as AccountType, balanceCents: -125000 },
	{ id: "4", name: "Cash", type: "cash" as AccountType, balanceCents: 5000 },
];

function formatCents(cents: number) {
	return new Intl.NumberFormat("en-US", {
		style: "currency",
		currency: "USD",
		minimumFractionDigits: 0,
		maximumFractionDigits: 0,
	}).format(cents / 100);
}

export function AppSidebar() {
	return (
		<Sidebar>
			<Sidebar.Header>
				<div className="flex items-center gap-2.5 px-1 py-0.5">
					<div className="flex size-8 shrink-0 items-center justify-center rounded-xl bg-[--accent]">
						<span className="text-sm font-bold text-[--accent-foreground]">ZB</span>
					</div>
					<span className="truncate text-sm font-semibold tracking-tight">ZeroBudget</span>
				</div>
			</Sidebar.Header>

			<Sidebar.Content>
				<Sidebar.Group>
					<Sidebar.GroupLabel>Accounts</Sidebar.GroupLabel>
					<Sidebar.Menu>
						{STATIC_ACCOUNTS.map((account) => (
							<Sidebar.MenuItem
								key={account.id}
								href={`/accounts/${account.id}`}
								tooltip={account.name}
							>
								<Sidebar.MenuLabel>{account.name}</Sidebar.MenuLabel>
								<Sidebar.MenuChip
									className={account.balanceCents < 0 ? "text-[--danger]" : ""}
								>
									{formatCents(account.balanceCents)}
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
							onPress={() => {}}
						>
							<Plus size={14} />
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
							name={STATIC_USER.initials}
							size="sm"
						/>
						<div className="flex min-w-0 flex-1 flex-col overflow-hidden">
							<span className="truncate text-sm font-medium leading-tight">
								{STATIC_USER.name}
							</span>
							<span className="truncate text-xs text-[--muted]">{STATIC_USER.email}</span>
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
